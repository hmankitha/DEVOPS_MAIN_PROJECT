# E-Commerce Platform - Production DevOps Infrastructure

## 🏗️ Architecture Overview

This is a production-grade, cloud-native e-commerce platform built with microservices architecture and fully automated DevOps workflows.

### Microservices:
- **User Management Service** (Golang) - Authentication, authorization, user profiles
- **Product Catalog Service** (Python/FastAPI) - Product listings, inventory, search
- **Order Management Service** (Java/Spring Boot) - Order processing, payment integration

### Infrastructure:
- **Cloud Provider**: AWS (EKS, RDS, S3, IAM)
- **Container Orchestration**: Kubernetes (EKS)
- **Infrastructure as Code**: Terraform
- **CI/CD**: Jenkins (Multi-branch pipelines)
- **Monitoring**: Prometheus + Grafana
- **Logging**: ELK Stack (Elasticsearch, Logstash, Kibana)
- **Security**: Trivy, AWS Secrets Manager, Vault
- **Load Testing**: JMeter, Locust

## 📁 Project Structure

```
bestproject2/
├── microservices/
│   ├── user-management/        # Golang service
│   ├── product-catalog/        # Python service
│   └── order-management/       # Java service
├── infrastructure/
│   ├── terraform/              # IaC for AWS resources
│   └── kubernetes/             # K8s manifests & Helm charts
├── ci-cd/
│   ├── jenkins/                # Jenkins pipelines
│   └── docker/                 # Dockerfiles & docker-compose
├── monitoring/
│   ├── prometheus/             # Prometheus configs
│   └── grafana/                # Grafana dashboards
├── logging/
│   └── elk/                    # ELK stack configs
├── security/
│   ├── trivy/                  # Container scanning configs
│   └── secrets/                # Vault/AWS Secrets setup
└── testing/
    ├── jmeter/                 # Load test scenarios
    └── locust/                 # Python load tests
```

## 🚀 Quick Start

### Prerequisites
- AWS CLI configured
- kubectl installed
- Terraform >= 1.5
- Docker & Docker Compose
- Jenkins
- helm

### Deployment Steps

1. **Provision Infrastructure**
```bash
cd infrastructure/terraform
terraform init
terraform plan -var-file=environments/prod.tfvars
terraform apply -var-file=environments/prod.tfvars
```

2. **Configure kubectl**
```bash
aws eks update-kubeconfig --name ecommerce-prod-cluster --region us-east-1
```

3. **Deploy Microservices**
```bash
cd infrastructure/kubernetes
kubectl apply -f namespaces/
kubectl apply -f secrets/
helm install user-mgmt ./helm-charts/user-management
helm install product-catalog ./helm-charts/product-catalog
helm install order-mgmt ./helm-charts/order-management
```

4. **Setup Monitoring**
```bash
kubectl apply -f monitoring/prometheus/
kubectl apply -f monitoring/grafana/
```

5. **Setup Logging**
```bash
kubectl apply -f logging/elk/
```

## ☁️ Azure Deployment (Preview)

Azure Container Apps + ACR deployment is available via Bicep and helper scripts.

```zsh
# From repo root
export RG_NAME=bestproject2-rg
export LOCATION=eastus
bash scripts/azure-deploy.sh
```

Details and advanced usage: see `DEPLOYMENT_AZURE.md`.

## 🔒 Security

- All containers scanned with Trivy
- Secrets managed via AWS Secrets Manager
- Network policies enforced
- RBAC configured
- TLS/SSL encryption enabled

## 📊 Monitoring

- Prometheus metrics on port 9090
- Grafana dashboards on port 3000
- Default credentials in secrets manager

## 📝 Logging

- Centralized logs via ELK stack
- Kibana dashboard on port 5601
- Log retention: 30 days (configurable)

## 🧪 Testing

```bash
# Unit tests
make test-all

# Integration tests
make integration-test

# Load tests
cd testing/locust
locust -f load_test.py --host=https://api.ecommerce.com
```

## 📈 Environments

- **Dev**: Auto-deploy on feature branch push
- **Staging**: Auto-deploy on merge to develop
- **Production**: Manual approval required

## 🤝 Contributing

See CONTRIBUTING.md for development workflow.

## 📄 License

MIT License - see LICENSE file
