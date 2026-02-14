package converters

import (
	"ozon_entrance/internal/domain/dto"
	"ozon_entrance/internal/domain/entities"
)

type LinksConverter struct{}

func NewLinksConverter() *LinksConverter {
	return &LinksConverter{}
}

func (c LinksConverter) ToOriginalDTO(link entities.Link) *dto.OriginalLink {
	return &dto.OriginalLink{
		Original: link.Original,
	}
}

func (c LinksConverter) ToShortDTO(link entities.Link) *dto.ShortLink {
	return &dto.ShortLink{
		Short: link.Short,
	}
}
