package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"

	"ozon_entrance/internal/adapter/repository/in_memory_repo"
	"ozon_entrance/internal/adapter/repository/postgres_repo"
	"ozon_entrance/internal/domain/ports/repository"
	"ozon_entrance/internal/infrastructure/database/in_memory"
	"ozon_entrance/internal/infrastructure/database/postgres"
	"ozon_entrance/internal/infrastructure/generator"
	"ozon_entrance/internal/usecase"
	"ozon_entrance/internal/usecase/links_usecase"
)

type Server struct {
	// storages
	database *postgres.Postgres
	inMemory *in_memory.InMemory
	// repo
	linksRepository repository.LinksRepository

	// uc
	linksUseCase usecase.LinksUseCase
	router       *chi.Mux
	server       *http.Server
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
	if s.database != nil {
		if err := postgres.MigrateDB(s.database); err != nil {
			return fmt.Errorf("migration err: %w", err)
		}
	}
	s.initRepo()
	s.initUseCases()
	s.initRoutes()
	s.initHTTPServer()
	return nil
}

func (s *Server) initDB() error {
	if s.storageInMemory() {
		s.inMemory = in_memory.NewInMemory()
		return nil
	}
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

func (s *Server) storageInMemory() bool {
	return os.Getenv("IN_MEM") == "true"
}

func (s *Server) initRepo() {
	if s.storageInMemory() {
		s.linksRepository = in_memory_repo.NewLinksRepository(s.inMemory)
		return
	}
	s.linksRepository = postgres_repo.NewLinksRepository(s.database)
}

func (s *Server) initUseCases() {
	shortGenerator := generator.NewShortGenerator()
	s.linksUseCase = links_usecase.NewLinksUseCase(s.linksRepository, shortGenerator)
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
