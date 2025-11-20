.PHONY: install run test clean docker-build docker-run

install: ## Install dependencies
	pip install -r requirements.txt

run: ## Run the application
	uvicorn main:app --reload --host 0.0.0.0 --port 8000

test: ## Run tests
	pytest tests/ -v --cov=app --cov-report=html

clean: ## Clean cache and temp files
	find . -type d -name __pycache__ -exec rm -rf {} +
	find . -type f -name '*.pyc' -delete
	rm -rf .pytest_cache
	rm -rf htmlcov
	rm -rf .coverage

docker-build: ## Build Docker image
	docker build -t product-catalog:latest .

docker-run: ## Run Docker container
	docker run -p 8000:8000 --env-file .env product-catalog:latest

migrate: ## Run database migrations
	alembic upgrade head

lint: ## Run linter
	flake8 app/
	black --check app/

format: ## Format code
	black app/

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
