package links

import (
	"encoding/json"
	"errors"
	"net/http"
	"ozon_entrance/internal/domain/dto"
	"ozon_entrance/internal/errs"
	"ozon_entrance/internal/usecase"
	httpError "ozon_entrance/pkg/http/error"
	"ozon_entrance/pkg/http/writer"
	"ozon_entrance/pkg/logger"
)

func CreateLink(uc usecase.LinksUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.OriginalLink
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.FromContext(r.Context()).Debug("CreateLink: decode error", "err", err)
			httpError.InternalError(w, err)
			return
		}
		short, err := uc.CreateLink(r.Context(), req)
		if err != nil {
			if errors.Is(err, errs.ErrInvalidURLFormat) || errors.Is(err, errs.ErrEmptyURL) {
				httpError.BadRequest(w, err)
				return
			}
			logger.FromContext(r.Context()).Error("CreateLink handler", "err", err)
			httpError.InternalError(w, err)
			return
		}
		writer.WriteJson(w, short)
	}
}
