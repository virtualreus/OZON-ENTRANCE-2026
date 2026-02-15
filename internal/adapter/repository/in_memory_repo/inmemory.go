package in_memory_repo

import (
	"context"
	"ozon_entrance/internal/domain/entities"
	"ozon_entrance/internal/domain/ports/repository"
	"ozon_entrance/internal/errs"
	"ozon_entrance/internal/infrastructure/database/in_memory"
	"time"
)

type linksRepository struct {
	db *in_memory.InMemory
}

func NewLinksRepository(db *in_memory.InMemory) repository.LinksRepository {
	return &linksRepository{db: db}
}

func (l *linksRepository) SaveLink(ctx context.Context, originalLink, shortLink string) (entities.Link, error) {
	l.db.Mutex.Lock()
	defer l.db.Mutex.Unlock()

	if existingShort, ok := l.db.ByOriginal[originalLink]; ok {
		if link, ok := l.db.ByShort[existingShort]; ok {
			return link, nil
		}
		return entities.Link{}, errs.ErrNotFound
	}

	if _, exists := l.db.ByShort[shortLink]; exists {
		return entities.Link{}, errs.ErrDuplicate
	}
	now := time.Now()
	link := entities.Link{
		Short:     shortLink,
		Original:  originalLink,
		CreatedAt: now,
	}

	l.db.ByShort[shortLink] = link
	l.db.ByOriginal[originalLink] = shortLink

	return link, nil
}

func (l *linksRepository) GetLink(ctx context.Context, shortLink string) (entities.Link, error) {
	l.db.Mutex.RLock()
	defer l.db.Mutex.RUnlock()

	link, exists := l.db.ByShort[shortLink]
	if !exists {
		return entities.Link{}, errs.ErrNotFound
	}

	return link, nil
}
