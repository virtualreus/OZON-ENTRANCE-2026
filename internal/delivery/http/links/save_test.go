package links

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"ozon_entrance/internal/domain/dto"
	"ozon_entrance/internal/errs"
	"ozon_entrance/internal/usecase/mocks"
	"ozon_entrance/pkg/logger"
)

func TestCreateLink(t *testing.T) {
	t.Parallel()
	uc := mocks.NewLinksUseCase(t)
	handler := CreateLink(uc)

	tests := []struct {
		name     string
		body     string
		setup    func()
		wantCode int
	}{
		{
			name: "ok",
			body: `{"original":"https://ya.ru"}`,
			setup: func() {
				uc.On("CreateLink", mock.Anything, dto.OriginalLink{Original: "https://ya.ru"}).
					Return(&dto.ShortLink{Short: "abc123XY_z"}, nil)
			},
			wantCode: http.StatusOK,
		},
		{
			name: "empty url",
			body: `{"original":""}`,
			setup: func() {
				uc.On("CreateLink", mock.Anything, dto.OriginalLink{Original: ""}).Return(nil, errs.ErrEmptyURL)
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "invalid url",
			body: `{"original":"not-a-url"}`,
			setup: func() {
				uc.On("CreateLink", mock.Anything, dto.OriginalLink{Original: "not-a-url"}).Return(nil, errs.ErrInvalidURLFormat)
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "internal error",
			body: `{"original":"https://ok.ru"}`,
			setup: func() {
				uc.On("CreateLink", mock.Anything, dto.OriginalLink{Original: "https://ok.ru"}).Return(nil, errs.ErrInternal)
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "decode error",
			body:     `{invalid`,
			setup:    func() {},
			wantCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc.ExpectedCalls = nil
			tt.setup()

			ctx := logger.ContextWithLogger(context.Background(), logger.New())
			req := httptest.NewRequest(http.MethodPost, "/link", bytes.NewBufferString(tt.body)).WithContext(ctx)
			rec := httptest.NewRecorder()

			handler(rec, req)

			assert.Equal(t, tt.wantCode, rec.Code)
		})
	}
}
