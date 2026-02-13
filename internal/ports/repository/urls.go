package repository

import "ozon_entrance/internal/domain/entities"

type URLRepository interface {
	SaveURL(originalURL string) (entities.Link, error)
	GetURL(shortUrl string) (entities.Link, error)
}
