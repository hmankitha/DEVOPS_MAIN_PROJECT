.PHONY: build test run clean docker-build docker-run

build: ## Build the application
	go build -o bin/user-management main.go

test: ## Run tests
	go test -v -cover ./...

run: ## Run the application
	go run main.go

clean: ## Clean build artifacts
	rm -rf bin/
	go clean

docker-build: ## Build Docker image
	docker build -t user-management:latest .

docker-run: ## Run Docker container
	docker run -p 8080:8080 --env-file .env user-management:latest

lint: ## Run linter
	golangci-lint run

deps: ## Download dependencies
	go mod download
	go mod tidy

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
