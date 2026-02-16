package in_memory_repo

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"ozon_entrance/internal/domain/entities"
	"ozon_entrance/internal/errs"
	"ozon_entrance/internal/infrastructure/database/in_memory"
)

func setUpTestLinksClient() *in_memory.InMemory {
	store := in_memory.NewInMemory()
	now := time.Now()
	store.ByShort["1111111111"] = entities.Link{
		Short:     "1111111111",
		Original:  "https://fijma.com",
		CreatedAt: now,
	}
	store.ByOriginal["https://fijma.com"] = "1111111111"

	store.ByShort["2222222222"] = entities.Link{
		Short:     "2222222222",
		Original:  "https://twitch.tv",
		CreatedAt: now,
	}
	store.ByOriginal["https://twitch.tv"] = "2222222222"

	store.ByShort["3333333333"] = entities.Link{
		Short:     "3333333333",
		Original:  "https://t.me",
		CreatedAt: now,
	}
	store.ByOriginal["https://t.me"] = "3333333333"

	store.ByShort["4444444444"] = entities.Link{
		Short:     "4444444444",
		Original:  "https://www.ozon.ru",
		CreatedAt: now,
	}
	store.ByOriginal["https://www.ozon.ru"] = "4444444444"
	return store
}

func setUpTestLinksClientInconsistent() *in_memory.InMemory {
	store := in_memory.NewInMemory()
	store.ByOriginal["https://max.ru"] = "onparkovka"
	return store
}

func Test_SaveLink(t *testing.T) {
	client := setUpTestLinksClient()
	repo := NewLinksRepository(client)

	tests := []struct {
		name     string
		original string
		short    string
		wantErr  error
	}{
		{
			name:     "valid create",
			original: "https://link.ru",
			short:    "5555555555",
			wantErr:  nil,
		},
		{
			name:     "existing link",
			original: "https://t.me",
			short:    "9999999999",
			wantErr:  nil,
		},
		{
			name:     "duplicate short",
			original: "https://other.com",
			short:    "1111111111",
			wantErr:  errs.ErrDuplicate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			got, err := repo.SaveLink(ctx, tt.original, tt.short)

			assert.Equal(t, tt.wantErr, err)
			if tt.wantErr == nil {
				assert.NotEmpty(t, got.Short)
				assert.Equal(t, tt.original, got.Original)
				assert.False(t, got.CreatedAt.IsZero())
				found, getErr := repo.GetLink(ctx, got.Short)
				assert.NoError(t, getErr)
				assert.Equal(t, tt.original, found.Original)
			} else {
				assert.Equal(t, entities.Link{}, got)
			}
		})
	}
}

func Test_SaveLink_InconsistentStorage(t *testing.T) {
	client := setUpTestLinksClientInconsistent()
	repo := NewLinksRepository(client)

	tests := []struct {
		name     string
		original string
		short    string
		wantErr  error
	}{
		{
			name:     "inconsistent link",
			original: "https://max.ru",
			short:    "5555555555",
			wantErr:  errs.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			got, err := repo.SaveLink(ctx, tt.original, tt.short)

			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, entities.Link{}, got)
		})
	}
}

func Test_GetLink(t *testing.T) {
	client := setUpTestLinksClient()
	repo := NewLinksRepository(client)

	tests := []struct {
		name     string
		short    string
		wantErr  error
		wantOrig string
	}{
		{
			name:     "first",
			short:    "1111111111",
			wantOrig: "https://fijma.com",
			wantErr:  nil,
		},
		{
			name:     "second",
			short:    "2222222222",
			wantOrig: "https://twitch.tv",
			wantErr:  nil,
		},
		{
			name:     "third",
			short:    "3333333333",
			wantOrig: "https://t.me",
			wantErr:  nil,
		},
		{
			name:     "fourth",
			short:    "4444444444",
			wantOrig: "https://www.ozon.ru",
			wantErr:  nil,
		},
		{
			name:    "not found",
			short:   "asdasdasdd",
			wantErr: errs.ErrNotFound,
		},
		{
			name:    "empty short",
			short:   "",
			wantErr: errs.ErrNotFound,
		},
		{
			name:    "wrong length",
			short:   "ab",
			wantErr: errs.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			got, err := repo.GetLink(ctx, tt.short)

			assert.Equal(t, tt.wantErr, err)
			if tt.wantErr == nil {
				assert.Equal(t, tt.short, got.Short)
				assert.Equal(t, tt.wantOrig, got.Original)
			} else {
				assert.Equal(t, entities.Link{}, got)
			}
		})
	}
}
