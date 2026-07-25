APP_NAME := go-architecture-template
MAIN_PATH := ./src/main.go

.PHONY: run run-dev build clean tidy test vet

run: ## Run with APP_ENV=local
	APP_ENV=local go run $(MAIN_PATH)

run-dev: ## Run with APP_ENV=development
	APP_ENV=development go run $(MAIN_PATH)

build: ## Build the single binary
	go build -o bin/$(APP_NAME) $(MAIN_PATH)

clean: ## Remove build artifacts
	rm -rf bin/

tidy: ## Tidy go.mod/go.sum
	go mod tidy

test: ## Run the test suite
	go test ./...

vet: ## Static analysis
	go vet ./...
