package repository

import (
	"context"
	"ozon_entrance/internal/domain/entities"
)

//go:generate mockery --name=LinksRepository
type LinksRepository interface {
	SaveLink(ctx context.Context, originalLink, shortUrl string) (entities.Link, error)
	GetLink(ctx context.Context, shortLink string) (entities.Link, error)
}
