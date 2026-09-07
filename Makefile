.DEFAULT_GOAL := help

BACKEND_DIR := backend

.PHONY: help setup generate generate-check format format-check lint test build

help:
	@echo "Available targets:"
	@echo "  make setup          - download backend dependencies"
	@echo "  make generate       - regenerate backend OpenAPI code"
	@echo "  make generate-check - fail if generated backend OpenAPI code is stale"
	@echo "  make format         - format backend Go files"
	@echo "  make format-check   - fail if backend Go files need formatting"
	@echo "  make lint           - run backend static checks"
	@echo "  make test           - run backend tests"
	@echo "  make build          - build the backend"

setup:
	$(MAKE) -C $(BACKEND_DIR) setup

generate:
	$(MAKE) -C $(BACKEND_DIR) openapi-gen

generate-check:
	$(MAKE) -C $(BACKEND_DIR) openapi-check

format:
	$(MAKE) -C $(BACKEND_DIR) format

format-check:
	$(MAKE) -C $(BACKEND_DIR) format-check

lint:
	$(MAKE) -C $(BACKEND_DIR) lint

test:
	$(MAKE) -C $(BACKEND_DIR) test

build:
	$(MAKE) -C $(BACKEND_DIR) build
