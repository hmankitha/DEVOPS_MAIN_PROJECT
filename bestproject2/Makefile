.PHONY: help build test deploy clean ansible-setup ansible-start ansible-stop ansible-deploy ansible-health ansible-install \
	middleware-start middleware-stop middleware-status middleware-backup-mysql middleware-restore-mysql \
	middleware-backup-redis middleware-restore-redis kafka-topics kafka-create-topic

# Colors for output
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# ============================================
# Ansible Commands
# ============================================

ansible-install: ## Install Ansible and required collections
	@echo "$(BLUE)Installing Ansible dependencies...$(NC)"
	brew install ansible || echo "Ansible already installed"
	ansible-galaxy collection install community.docker
	ansible-galaxy collection install community.general
	@echo "$(GREEN)✅ Dependencies installed$(NC)"

ansible-setup: ## Initial infrastructure setup (databases, ELK)
	@echo "$(BLUE)Setting up infrastructure with Ansible...$(NC)"
	cd ansible && ansible-playbook playbooks/setup.yml
	@echo "$(GREEN)✅ Infrastructure setup complete$(NC)"

ansible-start: ## Start all infrastructure services (docker-compose)
	@echo "$(BLUE)Starting all services with Ansible...$(NC)"
	cd ansible && ansible-playbook playbooks/start-services.yml
	@echo "$(GREEN)✅ Services started$(NC)"

ansible-stop: ## Stop all services
	@echo "$(YELLOW)Stopping all services with Ansible...$(NC)"
	cd ansible && ansible-playbook playbooks/stop-services.yml
	@echo "$(GREEN)✅ Services stopped$(NC)"

ansible-deploy: ## Deploy all microservices with Ansible
	@echo "$(BLUE)Deploying microservices with Ansible...$(NC)"
	cd ansible && ansible-playbook playbooks/deploy.yml
	@echo "$(GREEN)✅ Deployment complete$(NC)"

ansible-deploy-python: ## Deploy only Python microservices
	@echo "$(BLUE)Deploying Python services...$(NC)"
	cd ansible && ansible-playbook playbooks/deploy.yml --tags "python,product-catalog"
	@echo "$(GREEN)✅ Python services deployed$(NC)"

ansible-deploy-go: ## Deploy only Go microservices
	@echo "$(BLUE)Deploying Go services...$(NC)"
	cd ansible && ansible-playbook playbooks/deploy.yml --tags "go,user-management"
	@echo "$(GREEN)✅ Go services deployed$(NC)"

ansible-health: ## Check health of all services
	@echo "$(BLUE)Checking service health...$(NC)"
	cd ansible && ansible-playbook playbooks/deploy.yml --tags healthcheck
	@echo "$(GREEN)✅ Health check complete$(NC)"

ansible-restart: ansible-stop ansible-start ansible-deploy ansible-health ## Restart all services with Ansible

ansible-status: ## Show status of all services
	@echo "$(BLUE)Service Status:$(NC)"
	@echo ""
	@echo "$(GREEN)Docker Containers:$(NC)"
	@docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" | grep bestproject2 || echo "No containers running"
	@echo ""
	@echo "$(GREEN)Microservices:$(NC)"
	@curl -s http://localhost:8000/api/v1/products > /dev/null && echo "✅ Product Catalog (8000): Running" || echo "❌ Product Catalog (8000): Not responding"
	@curl -s http://localhost:8080/health > /dev/null && echo "✅ User Management (8080): Running" || echo "❌ User Management (8080): Not responding"
	@echo ""
	@echo "$(GREEN)Infrastructure:$(NC)"
	@curl -s http://localhost:9200 > /dev/null && echo "✅ Elasticsearch (9200): Running" || echo "❌ Elasticsearch (9200): Not responding"
	@curl -s http://localhost:5601 > /dev/null && echo "✅ Kibana (5601): Running" || echo "❌ Kibana (5601): Not responding"

ansible-logs-python: ## View Product Catalog logs
	@tail -f /tmp/product-catalog.log

ansible-logs-go: ## View User Management logs
	@tail -f /tmp/user-management.log

# ============================================
# ArgoCD GitOps Commands
# ============================================

argocd-install: ## Install ArgoCD on Kubernetes cluster
	@echo "$(BLUE)Installing ArgoCD...$(NC)"
	chmod +x argocd/install/install.sh
	./argocd/install/install.sh
	@echo "$(GREEN)✅ ArgoCD installed$(NC)"

argocd-deploy: ## Deploy ArgoCD using Ansible
	@echo "$(BLUE)Deploying ArgoCD with Ansible...$(NC)"
	cd ansible && ansible-playbook playbooks/deploy-argocd.yml
	@echo "$(GREEN)✅ ArgoCD deployed$(NC)"

argocd-portforward: ## Start port forwarding to ArgoCD UI
	@echo "$(BLUE)Starting port forward to ArgoCD UI...$(NC)"
	@echo "Access ArgoCD at: https://localhost:8081"
	kubectl port-forward svc/argocd-server -n argocd 8081:443

argocd-password: ## Get ArgoCD admin password
	@echo "$(GREEN)ArgoCD Admin Password:$(NC)"
	@kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
	@echo ""

argocd-login: ## Login to ArgoCD CLI
	@echo "$(BLUE)Logging into ArgoCD...$(NC)"
	@argocd login localhost:8081 --username admin --password $$(kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d) --insecure
	@echo "$(GREEN)✅ Logged in to ArgoCD$(NC)"

argocd-project: ## Apply ArgoCD AppProject
	@echo "$(BLUE)Applying ArgoCD AppProject...$(NC)"
	kubectl apply -f argocd/projects/microservices-project.yaml
	@echo "$(GREEN)✅ AppProject applied$(NC)"

argocd-apps: ## Deploy all applications to ArgoCD
	@echo "$(BLUE)Deploying applications to ArgoCD...$(NC)"
	kubectl apply -f argocd/applications/product-catalog.yaml
	kubectl apply -f argocd/applications/user-management.yaml
	@echo "$(GREEN)✅ Applications deployed$(NC)"

argocd-app-of-apps: ## Deploy using App of Apps pattern
	@echo "$(BLUE)Deploying with App of Apps pattern...$(NC)"
	kubectl apply -f argocd/applications/app-of-apps.yaml
	@echo "$(GREEN)✅ App of Apps deployed$(NC)"

argocd-list: ## List all ArgoCD applications
	@argocd app list

argocd-sync: ## Sync all applications
	@echo "$(BLUE)Syncing all applications...$(NC)"
	argocd app sync --all
	@echo "$(GREEN)✅ All applications synced$(NC)"

argocd-sync-product: ## Sync product-catalog application
	@echo "$(BLUE)Syncing product-catalog...$(NC)"
	argocd app sync product-catalog
	@echo "$(GREEN)✅ Product Catalog synced$(NC)"

argocd-sync-user: ## Sync user-management application
	@echo "$(BLUE)Syncing user-management...$(NC)"
	argocd app sync user-management
	@echo "$(GREEN)✅ User Management synced$(NC)"

argocd-status: ## Show status of all ArgoCD applications
	@echo "$(BLUE)ArgoCD Application Status:$(NC)"
	@argocd app get product-catalog --refresh || echo "❌ Product Catalog not deployed"
	@echo ""
	@argocd app get user-management --refresh || echo "❌ User Management not deployed"

argocd-ui: ## Open ArgoCD UI in browser
	@echo "$(BLUE)Opening ArgoCD UI...$(NC)"
	@echo "Username: admin"
	@echo -n "Password: "
	@kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
	@echo ""
	@open https://localhost:8081 || xdg-open https://localhost:8081

argocd-uninstall: ## Uninstall ArgoCD from cluster
	@echo "$(RED)⚠️  This will remove ArgoCD and all applications!$(NC)"
	@read -p "Are you sure? (yes/no): " confirm && [ "$$confirm" = "yes" ] || exit 1
	kubectl delete -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/v2.9.3/manifests/install.yaml
	kubectl delete namespace argocd
	@echo "$(GREEN)✅ ArgoCD uninstalled$(NC)"

argocd-full-setup: argocd-install argocd-project argocd-apps ## Complete ArgoCD setup with applications

# ============================================
# Original Build & Test Commands
# ============================================

build-all: ## Build all microservices
	@echo "Building all microservices..."
	cd microservices/user-management && make build
	cd microservices/product-catalog && make build
	cd microservices/order-management && make build

test-all: ## Run tests for all services
	@echo "Running tests for all services..."
	cd microservices/user-management && make test
	cd microservices/product-catalog && make test
	cd microservices/order-management && make test

docker-build: ## Build Docker images for all services
	@echo "Building Docker images..."
	cd microservices/user-management && docker build -t user-management:latest .
	cd microservices/product-catalog && docker build -t product-catalog:latest .
	cd microservices/order-management && docker build -t order-management:latest .

docker-push: ## Push Docker images to registry
	@echo "Pushing Docker images..."
	docker tag user-management:latest $(DOCKER_REGISTRY)/user-management:$(VERSION)
	docker tag product-catalog:latest $(DOCKER_REGISTRY)/product-catalog:$(VERSION)
	docker tag order-management:latest $(DOCKER_REGISTRY)/order-management:$(VERSION)
	docker push $(DOCKER_REGISTRY)/user-management:$(VERSION)
	docker push $(DOCKER_REGISTRY)/product-catalog:$(VERSION)
	docker push $(DOCKER_REGISTRY)/order-management:$(VERSION)

terraform-init: ## Initialize Terraform
	cd infrastructure/terraform && terraform init

terraform-plan: ## Plan Terraform changes
	cd infrastructure/terraform && terraform plan -var-file=environments/$(ENV).tfvars

terraform-apply: ## Apply Terraform changes
	cd infrastructure/terraform && terraform apply -var-file=environments/$(ENV).tfvars

k8s-deploy: ## Deploy to Kubernetes
	kubectl apply -f infrastructure/kubernetes/namespaces/
	kubectl apply -f infrastructure/kubernetes/configmaps/
	kubectl apply -f infrastructure/kubernetes/secrets/
	helm upgrade --install user-mgmt infrastructure/kubernetes/helm-charts/user-management
	helm upgrade --install product-catalog infrastructure/kubernetes/helm-charts/product-catalog
	helm upgrade --install order-mgmt infrastructure/kubernetes/helm-charts/order-management

monitoring-deploy: ## Deploy monitoring stack
	kubectl apply -f monitoring/prometheus/
	kubectl apply -f monitoring/grafana/

logging-deploy: ## Deploy logging stack
	kubectl apply -f logging/elk/

integration-test: ## Run integration tests
	@echo "Running integration tests..."
	cd testing && ./run-integration-tests.sh

load-test: ## Run load tests
	cd testing/locust && locust -f load_test.py --headless -u 1000 -r 100 --run-time 5m

security-scan: ## Run security scans
	@echo "Scanning containers with Trivy..."
	trivy image user-management:latest
	trivy image product-catalog:latest
	trivy image order-management:latest

clean: ## Clean build artifacts
	@echo "Cleaning up..."
	cd microservices/user-management && make clean
	cd microservices/product-catalog && make clean
	cd microservices/order-management && make clean

full-deploy: terraform-apply k8s-deploy monitoring-deploy logging-deploy ## Full deployment

.DEFAULT_GOAL := help

# ============================================
# Middleware Maintenance
# ============================================

middleware-start: ## Start MySQL, Redis, Kafka and exporters
	@bash scripts/middleware-maintenance.sh start

middleware-stop: ## Stop MySQL, Redis, Kafka and exporters
	@bash scripts/middleware-maintenance.sh stop

middleware-status: ## Show middleware container status
	@bash scripts/middleware-maintenance.sh status

middleware-backup-mysql: ## Backup MySQL to file (usage: make middleware-backup-mysql OUT=backup.sql)
	@[ -n "$(OUT)" ] || (echo "OUT=<file.sql> required" && exit 2)
	@bash scripts/middleware-maintenance.sh backup-mysql $(OUT)

middleware-restore-mysql: ## Restore MySQL from file (usage: make middleware-restore-mysql IN=backup.sql)
	@[ -n "$(IN)" ] || (echo "IN=<file.sql> required" && exit 2)
	@bash scripts/middleware-maintenance.sh restore-mysql $(IN)

middleware-backup-redis: ## Backup Redis RDB to file (usage: make middleware-backup-redis OUT=dump.rdb)
	@[ -n "$(OUT)" ] || (echo "OUT=<file.rdb> required" && exit 2)
	@bash scripts/middleware-maintenance.sh backup-redis $(OUT)

middleware-restore-redis: ## Restore Redis from RDB file (usage: make middleware-restore-redis IN=dump.rdb)
	@[ -n "$(IN)" ] || (echo "IN=<file.rdb> required" && exit 2)
	@bash scripts/middleware-maintenance.sh restore-redis $(IN)

kafka-topics: ## List Kafka topics
	@bash scripts/middleware-maintenance.sh kafka-topics

kafka-create-topic: ## Create Kafka topic (usage: make kafka-create-topic NAME=orders PARTITIONS=1 REPL=1)
	@[ -n "$(NAME)" ] || (echo "NAME=<topic> required" && exit 2)
	@bash scripts/middleware-maintenance.sh kafka-create-topic $(NAME) $(PARTITIONS) $(REPL)

# ============================================
# Demo
# ============================================

demo: ## Start core services and seed sample data
	@echo "Starting core services (dbs + services)..."
	cd "$(PWD)" && docker compose up -d postgres-user postgres-product redis elasticsearch kibana prometheus grafana user-management product-catalog
	@echo "Seeding sample data..."
	bash scripts/seed.sh
	@echo "Open Postman collection: docs/postman/ecommerce.postman_collection.json"
