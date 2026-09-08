.DEFAULT_GOAL := help

BACKEND_DIR := backend
FRONTEND_DIR := frontend
NPM ?= npm

.PHONY: help setup run run-backend run-frontend generate generate-check format format-check lint test build frontend-typecheck

help:
	@echo "Available targets:"
	@echo "  make setup          - download backend and frontend dependencies"
	@echo "  make run            - run backend and frontend development servers"
	@echo "  make run-backend    - run only the backend development server"
	@echo "  make run-frontend   - run only the frontend development server"
	@echo "  make generate       - regenerate backend and frontend OpenAPI code"
	@echo "  make generate-check - fail if generated OpenAPI code is stale"
	@echo "  make format         - format backend Go files"
	@echo "  make format-check   - fail if backend Go files need formatting"
	@echo "  make lint           - run backend lint and frontend typecheck"
	@echo "  make test           - run backend and frontend tests"
	@echo "  make build          - build backend and frontend"
	@echo "  make frontend-typecheck - typecheck the frontend"

setup:
	$(MAKE) -C $(BACKEND_DIR) setup
	$(NPM) --prefix $(FRONTEND_DIR) install

run:
	$(MAKE) -j2 run-backend run-frontend

run-backend:
	$(MAKE) -C $(BACKEND_DIR) run

run-frontend:
	$(NPM) --prefix $(FRONTEND_DIR) run dev

generate:
	$(MAKE) -C $(BACKEND_DIR) openapi-gen
	$(NPM) --prefix $(FRONTEND_DIR) run generate:api

generate-check:
	$(MAKE) -C $(BACKEND_DIR) openapi-check
	$(NPM) --prefix $(FRONTEND_DIR) run generate:api:check

format:
	$(MAKE) -C $(BACKEND_DIR) format

format-check:
	$(MAKE) -C $(BACKEND_DIR) format-check

lint:
	$(MAKE) -C $(BACKEND_DIR) lint
	$(NPM) --prefix $(FRONTEND_DIR) run typecheck

test:
	$(MAKE) -C $(BACKEND_DIR) test
	$(NPM) --prefix $(FRONTEND_DIR) test

build:
	$(MAKE) -C $(BACKEND_DIR) build
	$(NPM) --prefix $(FRONTEND_DIR) run build

frontend-typecheck:
	$(NPM) --prefix $(FRONTEND_DIR) run typecheck
