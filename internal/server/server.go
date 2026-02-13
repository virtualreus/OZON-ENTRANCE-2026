package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"ozon_entrance/internal/infrastructure/database/postgres"
	"time"
)

type Server struct {
	database *postgres.Postgres
	router   *chi.Mux
	server   *http.Server
}

func NewServer() (*Server, error) {
	server := &Server{}
	if err := server.init(); err != nil {
		return nil, err
	}
	return server, nil
}

func (s *Server) init() error {
	s.router = chi.NewRouter()
	if err := s.initDB(); err != nil {
		return fmt.Errorf("err on initial database: %w", err)
	}
	if err := postgres.MigrateDB(s.database); err != nil {
		return fmt.Errorf("migration err: %w", err)
	}
	s.initHTTPServer()
	return nil
}

func (s *Server) initDB() error {
	config, err := postgres.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	pg, err := postgres.NewPostgres(config)
	if err != nil {
		log.Fatal(err)
	}
	s.database = pg
	return nil
}

func (s *Server) initHTTPServer() {
	s.server = &http.Server{
		Addr:         ":8080",
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func (s *Server) Run() {
	log.Println("Server started")
	if err := s.server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
