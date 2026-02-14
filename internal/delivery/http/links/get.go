package links

import (
	"net/http"
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
