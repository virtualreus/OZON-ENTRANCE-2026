package server

import (
	"fmt"
	"log"
	"net/http"
	"ozon_entrance/internal/adapter/repository/postgres_repo"
	"ozon_entrance/internal/domain/ports/repository"
	"ozon_entrance/internal/infrastructure/database/postgres"
	"ozon_entrance/internal/usecase"
	"ozon_entrance/internal/usecase/links_usecase"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	database *postgres.Postgres
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
	if err := postgres.MigrateDB(s.database); err != nil {
		return fmt.Errorf("migration err: %w", err)
	}
	s.initRepo()
	s.initUseCases()
	s.initRoutes()
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

func (s *Server) initRepo() {
	//TODO: сделать проверку на переменную окружения и в зависимости от нее присваивать конкретную имплементацию репы
	// пока пг по дефолту. потом 3 строчки добавить. во всю использую прелести гекс архи :D
	s.linksRepository = postgres_repo.NewLinksRepository(s.database)
}

func (s *Server) initUseCases() {
	s.linksUseCase = links_usecase.NewLinksUseCase(s.linksRepository)
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
