BINARY     = bin/api
MIGRATIONS = migrations
GOFLAGS    = -race

-include .env
export

.PHONY: run build test lint swagger migrate-up migrate-down docker-up docker-down install-tools

## run: start the API server (requires postgres to be running)
run:
	go run ./cmd/api

## build: compile the binary to bin/api
build:
	go build -o $(BINARY) ./cmd/api

## test: run all tests with the race detector (requires Docker for integration tests)
test:
	go test $(GOFLAGS) ./...

## test-unit: run only unit tests (no Docker required)
test-unit:
	go test $(GOFLAGS) -short ./internal/service/...

## lint: run golangci-lint (install via install-tools)
lint:
	golangci-lint run ./...

## swagger: regenerate OpenAPI docs from handler annotations
swagger:
	swag init -g cmd/api/main.go -o docs

## migrate-up: apply all pending migrations
migrate-up:
	goose -dir $(MIGRATIONS) postgres "$(DATABASE_URL)" up

## migrate-down: roll back the most recent migration
migrate-down:
	goose -dir $(MIGRATIONS) postgres "$(DATABASE_URL)" down

## docker-up: start the local Postgres container
docker-up:
	docker compose up -d postgres

## docker-down: stop the local Postgres container
docker-down:
	docker compose down

## install-tools: install swag, goose CLI, and golangci-lint
install-tools:
	go install github.com/swaggo/swag/cmd/swag@v1.16.1
	go install github.com/pressly/goose/v3/cmd/goose@v3.24.3
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

## help: list all targets
help:
	@grep -E '^## ' Makefile | sed 's/## //'
