package links

import (
	"errors"
	"net/http"
	"ozon_entrance/internal/errs"
	"ozon_entrance/internal/usecase"
	httpError "ozon_entrance/pkg/http/error"
	"ozon_entrance/pkg/http/writer"
	"ozon_entrance/pkg/logger"

	"github.com/go-chi/chi/v5"
)

func GetLinkByShort(uc usecase.LinksUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		short := chi.URLParam(r, "short")
		if short == "" {
			logger.FromContext(r.Context()).Debug("GetLink: empty short param")
			httpError.BadRequest(w, errs.ErrInvalidShortLink)
			return
		}
		link, err := uc.GetLink(r.Context(), short)
		if err != nil {
			if errors.Is(err, errs.ErrInvalidShortLink) {
				httpError.BadRequest(w, err)
				return
			}
			if errors.Is(err, errs.ErrNotFound) {
				httpError.NotFound(w, err)
				return
			}
			logger.FromContext(r.Context()).Error("GetLink handler", "err", err)
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
