package links_usecase

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"ozon_entrance/internal/converters"
	"ozon_entrance/internal/domain/dto"
	"ozon_entrance/internal/domain/entities"
	"ozon_entrance/internal/domain/ports/generator"
	"ozon_entrance/internal/domain/ports/repository"
	"ozon_entrance/internal/errs"
	"ozon_entrance/internal/usecase"
	"ozon_entrance/pkg/logger"
	"strings"

	"golang.org/x/sync/singleflight"
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
		logger.FromContext(ctx).Debug("CreateLink: empty URL")
		return nil, errs.ErrEmptyURL
	}

	_, err := url.ParseRequestURI(original)
	if err != nil {
		logger.FromContext(ctx).Debug("CreateLink: invalid URL", "url", original, "err", err)
		return nil, errs.ErrInvalidURLFormat
	}

	logger.FromContext(ctx).Debug("CreateLink", "original", original)

	key := original

	v, err, _ := u.sf.Do(key, func() (interface{}, error) {
		var lastErr error

		for attempt := 1; attempt <= maxAttempts; attempt++ {
			short, errg := u.shortGenerator.GenerateShortLink()
			if errg != nil {
				return nil, errg
			}
			saved, err := u.linksRepository.SaveLink(ctx, original, short)
			if err == nil {
				return saved, nil
			}

			lastErr = err

			if errors.Is(err, errs.ErrDuplicate) {
				continue
			}
			logger.FromContext(ctx).Error("CreateLink: save failed", "err", err)
			return nil, fmt.Errorf("create link: %w", err)
		}

		logger.FromContext(ctx).Error("CreateLink: max attempts exceeded", "attempts", maxAttempts)
		return nil, fmt.Errorf("failed to generate unique short after %d attempts: %w", maxAttempts, lastErr)
	})

	if err != nil {
		return nil, err
	}

	savedLink := v.(entities.Link)
	logger.FromContext(ctx).Debug("CreateLink: ok", "short", savedLink.Short)
	return u.linksConverter.ToShortDTO(savedLink), nil
}

func (u *linksUseCase) GetLink(ctx context.Context, shortLink string) (*dto.OriginalLink, error) {
	if len(shortLink) != 10 {
		logger.FromContext(ctx).Debug("GetLink: invalid short length", "short", shortLink)
		return nil, errs.ErrInvalidShortLink
	}
	link, err := u.linksRepository.GetLink(ctx, shortLink)
	if err != nil {
		logger.FromContext(ctx).Debug("GetLink: not found or error", "short", shortLink, "err", err)
		return nil, err
	}
	return u.linksConverter.ToOriginalDTO(link), nil
}
