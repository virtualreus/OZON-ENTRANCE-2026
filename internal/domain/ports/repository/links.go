package repository

import (
	"context"
	"ozon_entrance/internal/domain/entities"
)

type LinksRepository interface {
	SaveLink(ctx context.Context, originalLink string) (entities.Link, error)
	GetLink(ctx context.Context, shortLink string) (entities.Link, error)
}
