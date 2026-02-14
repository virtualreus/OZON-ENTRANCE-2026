package usecase

import (
	"context"
	"ozon_entrance/internal/domain/dto"
)

type LinksUseCase interface {
	CreateLink(ctx context.Context, originalLink string) (*dto.ShortLink, error)
	GetLink(ctx context.Context, shortLink string) (*dto.OriginalLink, error)
}
