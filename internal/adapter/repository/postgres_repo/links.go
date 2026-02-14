package postgres_repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ozon_entrance/internal/domain/entities"
	"ozon_entrance/internal/domain/ports/repository"
	"ozon_entrance/internal/infrastructure/database/postgres"

	"github.com/Masterminds/squirrel"
)

type linksRepository struct {
	db *postgres.Postgres
}

func NewLinksRepository(db *postgres.Postgres) repository.LinksRepository {
	return &linksRepository{
		db: db,
	}
}

// Если пытаемся вставить и уже есть такая - возвращаем короткую уже готовую (+метод мб)
// Если пытаемся вставить и еще нет - делаем и возвращаем ссылочку. пока хз че
func (l linksRepository) SaveLink(ctx context.Context, originalLink string) (entities.Link, error) {
	return entities.Link{}, nil
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
			return entities.Link{}, fmt.Errorf(prefix+"link not found: %w", err)
		}
		return entities.Link{}, fmt.Errorf(prefix+"query err: %w", err)
	}
	return link, nil
}
