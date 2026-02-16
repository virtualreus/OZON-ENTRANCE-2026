package usecase

import (
	"context"
	"ozon_entrance/internal/domain/dto"
)

//go:generate mockery --name=LinksUseCase
type LinksUseCase interface {
	CreateLink(ctx context.Context, originalLink dto.OriginalLink) (*dto.ShortLink, error)
	GetLink(ctx context.Context, shortLink string) (*dto.OriginalLink, error)
}
