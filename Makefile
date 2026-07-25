APP_NAME := go-architecture-template
MAIN_PATH := ./src/main.go
GOBIN := $(shell go env GOPATH)/bin
export PATH := $(GOBIN):$(PATH)

.PHONY: run run-dev build tidy openapi docker-up docker-down test vet

run: ## Run with APP_ENV=local
	APP_ENV=local go run $(MAIN_PATH)

run-dev: ## Run with APP_ENV=development
	APP_ENV=development go run $(MAIN_PATH)

build: ## Build the single binary
	go build -o bin/$(APP_NAME) $(MAIN_PATH)

tidy: ## Tidy go.mod/go.sum
	go mod tidy

openapi: ## Regenerate src/docs from controller annotations (commit the result)
	@if [ ! -x "$(GOBIN)/swag" ]; then \
		echo "installing swag..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	"$(GOBIN)/swag" init -g src/main.go -o src/docs --parseDependency --parseInternal --outputTypes go,json

docker-up: ## Start Postgres and Redis for local development
	docker compose up -d

docker-down: ## Stop local infrastructure containers
	docker compose down

test: ## Run the test suite
	go test ./...

vet: ## Static analysis
	go vet ./...
