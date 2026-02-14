package links_usecase

import (
	"context"
	"ozon_entrance/internal/converters"
	"ozon_entrance/internal/domain/dto"
	"ozon_entrance/internal/domain/ports/repository"
	"ozon_entrance/internal/usecase"
)

type linksUseCase struct {
	linksRepository repository.LinksRepository
	linksConverter  *converters.LinksConverter
}

func NewLinksUseCase(linksRepository repository.LinksRepository) usecase.LinksUseCase {
	return &linksUseCase{
		linksRepository: linksRepository,
	}
}

func (u *linksUseCase) CreateLink(ctx context.Context, originalLink string) (*dto.ShortLink, error) {
	link, err := u.linksRepository.SaveLink(ctx, originalLink)
	if err != nil {
		return nil, err
	}

	return u.linksConverter.ToShortDTO(link), nil
}

func (u *linksUseCase) GetLink(ctx context.Context, shortLink string) (*dto.OriginalLink, error) {
	link, err := u.linksRepository.GetLink(ctx, shortLink)
	if err != nil {
		return nil, err
	}

	return u.linksConverter.ToOriginalDTO(link), nil
}
