GO ?= go
OPENAPI := ../openapi.yaml
GEN_PKG := internal/openapigen
OPENAPI_GENERATOR_VERSION := v2.4.1

.DEFAULT_GOAL := help

.PHONY: help setup openapi-gen openapi-check format format-check lint test run build gen-secret

help:
	@echo "Available targets:"
	@echo "  make setup          - download backend Go modules"
	@echo "  make openapi-gen    - generate backend OpenAPI code"
	@echo "  make openapi-check  - fail if generated OpenAPI code is stale"
	@echo "  make format         - format backend Go files"
	@echo "  make format-check   - fail if Go files need formatting"
	@echo "  make lint           - run go vet"
	@echo "  make test           - run backend tests"
	@echo "  make run            - run the backend locally"
	@echo "  make build          - build the backend binary"

setup:
	$(GO) mod download

openapi-gen:
	$(GO) run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OPENAPI_GENERATOR_VERSION) -generate "types,chi-server,spec" -package openapigen -o $(GEN_PKG)/openapi.gen.go $(OPENAPI)

openapi-check:
	@tmp_file=$$(mktemp); \
	trap 'rm -f "$$tmp_file"' EXIT; \
	$(GO) run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OPENAPI_GENERATOR_VERSION) -generate "types,chi-server,spec" -package openapigen -o "$$tmp_file" $(OPENAPI); \
	cmp -s "$$tmp_file" $(GEN_PKG)/openapi.gen.go || { echo "Generated OpenAPI code is stale. Run make openapi-gen."; exit 1; }

format:
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

format-check:
	@test -z "$$(gofmt -l $$(find . -name '*.go' -not -path './vendor/*'))"

lint:
	$(GO) vet ./...

test:
	$(GO) test ./...

run:
	CGO_ENABLED=1 $(GO) run main.go

build:
	mkdir -p bin
	CGO_ENABLED=1 $(GO) build -o bin/mygolangapp main.go

gen-secret:
	$(GO) run ./script/gen-secret/main.go
