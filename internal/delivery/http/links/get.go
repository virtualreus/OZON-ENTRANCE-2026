package links

import (
	"errors"
	"net/http"
	"ozon_entrance/internal/errs"
	"ozon_entrance/internal/usecase"
	httpError "ozon_entrance/pkg/http/error"
	"ozon_entrance/pkg/http/writer"

	"github.com/go-chi/chi/v5"
)

func GetLinkByShort(uc usecase.LinksUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		short := chi.URLParam(r, "short")
		if short == "" {
			httpError.InternalError(w, nil)
			return
		}
		link, err := uc.GetLink(r.Context(), short)
		if err != nil {
			if errors.Is(err, errs.ErrNotFound) || errors.Is(err, errs.ErrInvalidShortLink) {
				httpError.BadRequest(w, err)
				return
			}
			httpError.InternalError(w, err)
			return
		}
		if link == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		writer.WriteJson(w, link)
	}
}
