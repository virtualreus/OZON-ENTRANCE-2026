package converters

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"ozon_entrance/internal/domain/dto"
	"ozon_entrance/internal/domain/entities"
)

func TestLinksConverter_ToOriginalDTO(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name string
		link entities.Link
		want *dto.OriginalLink
	}{
		{
			name: "valid link",
			link: entities.Link{
				Short:     "1111111111",
				Original:  "https://ozon.ru",
				CreatedAt: now,
			},
			want: &dto.OriginalLink{Original: "https://ozon.ru"},
		},
	}
	c := NewLinksConverter()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, c.ToOriginalDTO(tt.link))
		})
	}
}

func TestLinksConverter_ToShortDTO(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name string
		link entities.Link
		want *dto.ShortLink
	}{
		{
			name: "valid link",
			link: entities.Link{
				Short:     "2222222222",
				Original:  "https://t.me",
				CreatedAt: now,
			},
			want: &dto.ShortLink{Short: "2222222222"},
		},
	}
	c := NewLinksConverter()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, c.ToShortDTO(tt.link))
		})
	}
}
