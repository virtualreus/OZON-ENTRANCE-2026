package links_usecase

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"ozon_entrance/internal/domain/dto"
	"ozon_entrance/internal/domain/entities"
	"ozon_entrance/internal/domain/ports/repository/mocks"
	"ozon_entrance/internal/errs"
	"ozon_entrance/internal/infrastructure/generator"
	"ozon_entrance/pkg/logger"
)

func ctxWithLogger(ctx context.Context) context.Context {
	return logger.ContextWithLogger(ctx, logger.New())
}

func Test_GetLink(t *testing.T) {
	repo := mocks.NewLinksRepository(t)
	uc := NewLinksUseCase(repo, generator.NewShortGenerator())

	tests := []struct {
		name    string
		short   string
		setup   func()
		wantErr error
		want    string
	}{
		{
			name:  "valid find",
			short: "abc123XY_z",
			setup: func() {
				repo.On("GetLink", mock.Anything, "abc123XY_z").
					Return(entities.Link{Short: "abc123XY_z", Original: "https://ya.ru", CreatedAt: time.Now()}, nil)
			},
			want: "https://ya.ru",
		},
		{
			name:  "not found",
			short: "bebrochkaz",
			setup: func() {
				repo.On("GetLink", mock.Anything, "bebrochkaz").Return(entities.Link{}, errs.ErrNotFound)
			},
			wantErr: errs.ErrNotFound,
		},
		{
			name:    "invalid length",
			short:   "krauninriver",
			setup:   func() {},
			wantErr: errs.ErrInvalidShortLink,
		},
		{
			name:    "empty short",
			short:   "",
			setup:   func() {},
			wantErr: errs.ErrInvalidShortLink,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.ExpectedCalls = nil
			tt.setup()

			got, err := uc.GetLink(ctxWithLogger(context.Background()), tt.short)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want, got.Original)
			}
		})
	}
}

func Test_CreateLink(t *testing.T) {
	repo := mocks.NewLinksRepository(t)
	uc := NewLinksUseCase(repo, generator.NewShortGenerator())

	tests := []struct {
		name    string
		input   dto.OriginalLink
		setup   func()
		wantErr error
	}{
		{
			name:  "ok",
			input: dto.OriginalLink{Original: "https://max.ru"},
			setup: func() {
				repo.On("SaveLink", mock.Anything, "https://max.ru", mock.Anything).
					Return(entities.Link{Short: "shotnblock", Original: "https://max.ru", CreatedAt: time.Now()}, nil)
			},
		},
		{
			name:    "empty url",
			input:   dto.OriginalLink{Original: ""},
			setup:   func() {},
			wantErr: errs.ErrEmptyURL,
		},
		{
			name:    "invalid link",
			input:   dto.OriginalLink{Original: "danilkolbasenko"},
			setup:   func() {},
			wantErr: errs.ErrInvalidURLFormat,
		},
		{
			name:  "not found",
			input: dto.OriginalLink{Original: "https://example.com"},
			setup: func() {
				repo.On("SaveLink", mock.Anything, "https://example.com", mock.Anything).
					Return(entities.Link{}, errs.ErrNotFound)
			},
			wantErr: errs.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.ExpectedCalls = nil
			tt.setup()

			got, err := uc.CreateLink(ctxWithLogger(context.Background()), tt.input)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Len(t, got.Short, 10)
			}
		})
	}
}
