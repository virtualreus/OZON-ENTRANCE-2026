.PHONY: up down logs tests

up:
	docker-compose --env-file .env -f docker-compose.yml up -d --build

down:
	docker-compose --env-file .env -f docker-compose.yml down

logs:
	docker-compose --env-file .env -f docker-compose.yml logs -f

tests:
	go clean --testcache
	go test --cover ./internal/adapter/repository/in_memory_repo \
		./internal/usecase/links_usecase \
		./internal/converters \
		./internal/infrastructure/generator \
		./internal/delivery/http/links \
		--coverprofile=coverage.out
	go tool cover --html=coverage.out

.DEFAULT_GOAL := up
