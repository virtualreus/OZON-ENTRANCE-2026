package links_usecase

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/singleflight"
	"net/url"
	"ozon_entrance/internal/converters"
	"ozon_entrance/internal/domain/dto"
	"ozon_entrance/internal/domain/entities"
	"ozon_entrance/internal/domain/ports/generator"
	"ozon_entrance/internal/domain/ports/repository"
	"ozon_entrance/internal/errs"
	"ozon_entrance/internal/usecase"
	"strings"
)

const maxAttempts = 10

type linksUseCase struct {
	linksRepository repository.LinksRepository
	linksConverter  *converters.LinksConverter
	shortGenerator  generator.ShortLinkGenerator
	sf              singleflight.Group
}

func NewLinksUseCase(linksRepository repository.LinksRepository, shortGenerator generator.ShortLinkGenerator) usecase.LinksUseCase {
	return &linksUseCase{
		linksRepository: linksRepository,
		linksConverter:  converters.NewLinksConverter(),
		shortGenerator:  shortGenerator,
		sf:              singleflight.Group{},
	}
}

func (u *linksUseCase) CreateLink(ctx context.Context, originalLink dto.OriginalLink) (*dto.ShortLink, error) {
	original := strings.TrimSpace(originalLink.Original)
	if original == "" {
		return nil, errs.ErrEmptyURL
	}

	_, err := url.ParseRequestURI(original)
	if err != nil {
		return nil, errs.ErrInvalidURLFormat
	}

	key := original

	v, err, _ := u.sf.Do(key, func() (interface{}, error) {
		var lastErr error

		for attempt := 1; attempt <= maxAttempts; attempt++ {
			short := u.shortGenerator.GenerateShortLink()

			saved, err := u.linksRepository.SaveLink(ctx, original, short)
			if err == nil {
				return saved, nil
			}

			lastErr = err

			if errors.Is(err, errs.ErrDuplicate) {
				continue
			}
			return nil, fmt.Errorf("create link: %w", err)
		}

		return nil, fmt.Errorf("failed to generate unique short after %d attempts: %w", maxAttempts, lastErr)
	})

	if err != nil {
		return nil, err
	}

	savedLink := v.(entities.Link)
	return u.linksConverter.ToShortDTO(savedLink), nil
}

func (u *linksUseCase) GetLink(ctx context.Context, shortLink string) (*dto.OriginalLink, error) {
	if len(shortLink) != 10 {
		return nil, errs.ErrInvalidShortLink
	}

	v, err, _ := u.sf.Do(shortLink, func() (interface{}, error) {
		link, err := u.linksRepository.GetLink(ctx, shortLink)
		if err != nil {
			return nil, err
		}
		return link, nil
	})

	if err != nil {
		return nil, err
	}

	return u.linksConverter.ToOriginalDTO(v.(entities.Link)), nil
}
