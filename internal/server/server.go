package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	"ozon_entrance/pkg/logger"
)

type Server struct {
	logger *slog.Logger

	database *postgres.Postgres
	inMemory *in_memory.InMemory

	linksRepository repository.LinksRepository

	linksUseCase usecase.LinksUseCase

	router *chi.Mux
	server *http.Server
}

func NewServer() (*Server, error) {
	server := &Server{}
	if err := server.init(); err != nil {
		return nil, err
	}
	return server, nil
}

func (s *Server) init() error {
	s.logger = logger.New()
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
		s.logger.Info("storage: in-memory")
		return nil
	}
	config, err := postgres.NewConfig()
	if err != nil {
		return err
	}
	pg, err := postgres.NewPostgres(config)
	if err != nil {
		return err
	}
	s.database = pg
	s.logger.Info("storage: postgres")
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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func (s *Server) Run() {
	go func() {
		s.logger.Info("server started", "addr", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("listen failed", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("shutdown", "err", err)
		return
	}
	s.logger.Info("server stopped")
}
