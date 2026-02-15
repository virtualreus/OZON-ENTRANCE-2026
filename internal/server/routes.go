package server

import (
	"ozon_entrance/internal/delivery/http/links"
	"ozon_entrance/pkg/middleware"

	"github.com/go-chi/chi/v5"
)

func (s *Server) initRoutes() {
	s.router.Route("/api/v0", func(r chi.Router) {
		r.Use(middleware.LoggerContext(s.logger))
		r.Use(middleware.RequestLog)
		r.Post("/link", links.CreateLink(s.linksUseCase))
		r.Get("/link/{short}", links.GetLinkByShort(s.linksUseCase))
	})
}
