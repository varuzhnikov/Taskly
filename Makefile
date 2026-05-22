SHELL := /bin/bash

.PHONY: dev backend-db-up backend-dev frontend-dev check-env help

## dev: start postgres, backend API, and frontend dev server
dev: check-env backend-db-up
	@set -euo pipefail; \
	trap 'kill $$backend_pid $$frontend_pid 2>/dev/null || true' EXIT INT TERM; \
	$(MAKE) -C backend run & backend_pid=$$!; \
	$(MAKE) -C frontend dev & frontend_pid=$$!; \
	wait

## backend-db-up: start local postgres for the backend
backend-db-up:
	@$(MAKE) -C backend docker-up

## backend-dev: start only the backend API
backend-dev: check-env
	@$(MAKE) -C backend run

## frontend-dev: start only the frontend dev server
frontend-dev: check-env
	@$(MAKE) -C frontend dev

## check-env: ensure local env files exist before starting services
check-env:
	@test -f backend/.env || (echo "Missing backend/.env. Copy backend/.env.example to backend/.env first."; exit 1)
	@test -f backend/.env.compose || (echo "Missing backend/.env.compose. Copy backend/.env.compose.example to backend/.env.compose first."; exit 1)
	@test -f frontend/.env.local || (echo "Missing frontend/.env.local. Copy frontend/.env.local.example to frontend/.env.local first."; exit 1)

## help: list available targets
help:
	@grep -E '^## ' Makefile | sed 's/## //'
