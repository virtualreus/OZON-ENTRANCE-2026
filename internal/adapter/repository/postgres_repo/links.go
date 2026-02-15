package postgres_repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ozon_entrance/internal/domain/entities"
	"ozon_entrance/internal/domain/ports/repository"
	"ozon_entrance/internal/errs"
	"ozon_entrance/internal/infrastructure/database/postgres"
	"ozon_entrance/pkg/logger"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
)

type linksRepository struct {
	db *postgres.Postgres
}

func NewLinksRepository(db *postgres.Postgres) repository.LinksRepository {
	return &linksRepository{
		db: db,
	}
}

func (l linksRepository) SaveLink(ctx context.Context, originalLink, shortUrl string) (entities.Link, error) {
	prefix := "SaveLink: "
	var link entities.Link

	insertQb := l.db.Builder.
		Insert("urls").
		Columns("short_url", "original_url").
		Values(shortUrl, originalLink).
		Suffix("ON CONFLICT (original_url) DO NOTHING RETURNING short_url, original_url, created_at")

	query, args, err := insertQb.ToSql()
	if err != nil {
		return entities.Link{}, fmt.Errorf(prefix+": %w", err)
	}

	err = l.db.SqlxDB().QueryRowxContext(ctx, query, args...).StructScan(&link)
	if err == nil {
		return link, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return entities.Link{}, errs.ErrDuplicate
		}
		logger.FromContext(ctx).Error("SaveLink: insert failed", "err", err)
		return entities.Link{}, fmt.Errorf(prefix+"insert failed: %w", err)
	}

	selectQb := l.db.Builder.
		Select("short_url", "original_url", "created_at").
		From("urls").
		Where(squirrel.Eq{"original_url": originalLink})

	query, args, err = selectQb.ToSql()
	if err != nil {
		return entities.Link{}, fmt.Errorf(prefix+": %w", err)
	}

	err = l.db.SqlxDB().GetContext(ctx, &link, query, args...)
	if err != nil {
		return entities.Link{}, fmt.Errorf(prefix+"get existing: %w", err)
	}
	return link, nil
}

func (l linksRepository) GetLink(ctx context.Context, shortLink string) (entities.Link, error) {
	prefix := "GetLinks: "
	var link entities.Link

	qb := l.db.Builder.
		Select("short_url", "original_url", "created_at").
		From("urls").
		Where(squirrel.Eq{"short_url": shortLink})

	query, args, err := qb.ToSql()
	if err != nil {
		return entities.Link{}, fmt.Errorf(prefix+"convert to sql: %w", err)
	}

	if err = l.db.SqlxDB().GetContext(ctx, &link, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Link{}, errs.ErrNotFound
		}
		logger.FromContext(ctx).Error("GetLink: query failed", "err", err)
		return entities.Link{}, fmt.Errorf(prefix+"query err: %w", err)
	}

	return link, nil
}
