package links

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"ozon_entrance/internal/domain/dto"
	"ozon_entrance/internal/errs"
	"ozon_entrance/internal/usecase/mocks"
	"ozon_entrance/pkg/logger"
)

func reqWithShortParam(short string) *http.Request {
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("short", short)
	return httptest.NewRequest(http.MethodGet, "/link/"+short, nil).
		WithContext(context.WithValue(logger.ContextWithLogger(context.Background(), logger.New()), chi.RouteCtxKey, ctx))
}

func TestGetLinkByShort(t *testing.T) {
	t.Parallel()
	uc := mocks.NewLinksUseCase(t)
	handler := GetLinkByShort(uc)

	tests := []struct {
		name     string
		short    string
		setup    func()
		wantCode int
	}{
		{
			name:  "ok",
			short: "1111111111",
			setup: func() {
				uc.On("GetLink", mock.Anything, "1111111111").
					Return(&dto.OriginalLink{Original: "https://ozon.ru"}, nil)
			},
			wantCode: http.StatusOK,
		},
		{
			name:     "not found",
			short:    "0000000000",
			setup:    func() { uc.On("GetLink", mock.Anything, "0000000000").Return(nil, errs.ErrNotFound) },
			wantCode: http.StatusNotFound,
		},
		{
			name:     "invalid short",
			short:    "ab",
			setup:    func() { uc.On("GetLink", mock.Anything, "ab").Return(nil, errs.ErrInvalidShortLink) },
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "empty short",
			short:    "",
			setup:    func() {},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "internal error",
			short:    "3333333333",
			setup:    func() { uc.On("GetLink", mock.Anything, "3333333333").Return(nil, errs.ErrInternal) },
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "link nil",
			short:    "4444444444",
			setup:    func() { uc.On("GetLink", mock.Anything, "4444444444").Return(nil, nil) },
			wantCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc.ExpectedCalls = nil
			tt.setup()

			req := reqWithShortParam(tt.short)
			rec := httptest.NewRecorder()

			handler(rec, req)

			assert.Equal(t, tt.wantCode, rec.Code)
		})
	}
}
