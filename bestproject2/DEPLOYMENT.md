# E-Commerce Platform Deployment Guide

## 📋 Prerequisites

### Required Tools
- AWS CLI (v2.x)
- kubectl (v1.28+)
- Terraform (v1.5+)
- Docker (v24+)
- Helm (v3.12+)
- Jenkins (v2.400+)

### AWS Requirements
- AWS Account with appropriate permissions
- IAM roles for EKS, RDS, S3, Secrets Manager
- Route53 hosted zone (optional, for custom domain)

## 🚀 Quick Start

### 1. Local Development Setup

```bash
# Clone the repository
git clone <repository-url>
cd bestproject2

# Start services locally with Docker Compose
docker-compose up -d

# Verify services are running
docker-compose ps

# Access services
# User Management: http://localhost:8080
# Product Catalog: http://localhost:8000
# Grafana: http://localhost:3000
# Kibana: http://localhost:5601
# Prometheus: http://localhost:9090
```

### 2. AWS Infrastructure Deployment

```bash
# Set environment variables
export AWS_REGION=us-east-1
export ENVIRONMENT=prod

# Initialize and apply Terraform
cd infrastructure/terraform
terraform init
terraform plan -var-file=environments/prod.tfvars
terraform apply -var-file=environments/prod.tfvars

# Configure kubectl
aws eks update-kubeconfig --name ecommerce-prod-cluster --region us-east-1
```

### 3. Build and Push Docker Images

```bash
# Build all images
make docker-build

# Or use the script
./scripts/build-images.sh

# Images will be pushed to AWS ECR
```

### 4. Deploy to Kubernetes

```bash
# Deploy using script
./scripts/deploy.sh prod

# Or manually
kubectl apply -f infrastructure/kubernetes/namespaces/
kubectl apply -f infrastructure/kubernetes/deployments/
kubectl apply -f infrastructure/kubernetes/ingress/
```

### 5. Setup CI/CD Pipeline

```bash
# Configure Jenkins
# 1. Install required plugins:
#    - Docker Pipeline
#    - AWS Steps
#    - Kubernetes CLI
#    - SonarQube Scanner

# 2. Add credentials:
#    - AWS credentials
#    - Docker registry credentials
#    - Kubernetes config

# 3. Create multibranch pipeline jobs for each service
#    - User Management: ci-cd/jenkins/Jenkinsfile-user-management
#    - Product Catalog: ci-cd/jenkins/Jenkinsfile-product-catalog
```

## 📊 Monitoring Setup

### Prometheus & Grafana

```bash
# Deploy monitoring stack
helm upgrade --install prometheus prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  --create-namespace \
  --values monitoring/prometheus/values.yaml

# Access Grafana
kubectl port-forward -n monitoring svc/prometheus-grafana 3000:80

# Default credentials: admin/admin
```

### ELK Stack

```bash
# Deploy Elasticsearch
helm upgrade --install elasticsearch elastic/elasticsearch \
  --namespace logging \
  --create-namespace

# Deploy Kibana
helm upgrade --install kibana elastic/kibana \
  --namespace logging

# Access Kibana
kubectl port-forward -n logging svc/kibana 5601:5601
```

## 🧪 Testing

### Unit Tests

```bash
# Run all unit tests
make test-all

# Individual services
cd microservices/user-management && make test
cd microservices/product-catalog && make test
```

### Integration Tests

```bash
# Run integration tests
make integration-test
```

### Load Testing

```bash
# Using Locust
cd testing/locust
locust -f load_test.py --host=https://api.ecommerce.com

# Access Locust UI at http://localhost:8089

# Run headless load test
locust -f load_test.py --headless -u 1000 -r 100 --run-time 10m --host=https://api.ecommerce.com
```

## 🔒 Security

### Container Scanning

```bash
# Scan all images
./security/trivy/scan.sh

# Scan specific image
trivy image user-management:latest
```

### Secrets Management

```bash
# Store secrets in AWS Secrets Manager
aws secretsmanager create-secret \
  --name ecommerce/user-service \
  --secret-string '{"db_password":"xxx","jwt_secret":"xxx"}'

# Retrieve secrets
aws secretsmanager get-secret-value \
  --secret-id ecommerce/user-service
```

## 🔄 Update & Rollback

### Rolling Update

```bash
# Update deployment
kubectl set image deployment/user-management \
  user-management=<ECR_REGISTRY>/user-management:v2.0.0 \
  -n ecommerce-prod

# Check rollout status
kubectl rollout status deployment/user-management -n ecommerce-prod
```

### Rollback

```bash
# Rollback to previous version
kubectl rollout undo deployment/user-management -n ecommerce-prod

# Rollback to specific revision
kubectl rollout undo deployment/user-management --to-revision=2 -n ecommerce-prod
```

## 📈 Scaling

### Manual Scaling

```bash
# Scale deployment
kubectl scale deployment/user-management --replicas=5 -n ecommerce-prod
```

### Auto-scaling

HPA is already configured in deployment manifests:
- Min replicas: 3
- Max replicas: 10
- CPU threshold: 70%
- Memory threshold: 80%

## 🐛 Troubleshooting

### View Logs

```bash
# View pod logs
kubectl logs -f deployment/user-management -n ecommerce-prod

# View logs from all containers
kubectl logs -f deployment/user-management --all-containers=true -n ecommerce-prod

# View previous container logs
kubectl logs deployment/user-management --previous -n ecommerce-prod
```

### Debug Pod Issues

```bash
# Describe pod
kubectl describe pod <pod-name> -n ecommerce-prod

# Execute command in pod
kubectl exec -it <pod-name> -n ecommerce-prod -- /bin/sh

# Check events
kubectl get events -n ecommerce-prod --sort-by='.lastTimestamp'
```

### Database Connection Issues

```bash
# Test database connectivity
kubectl run -it --rm debug --image=postgres:15 --restart=Never -- \
  psql -h <RDS_ENDPOINT> -U postgres -d usermanagement
```

## 📞 Support

For issues and questions:
- Create an issue in the repository
- Contact DevOps team
- Check monitoring dashboards for alerts

## 📝 Additional Resources

- [Terraform AWS Modules](https://registry.terraform.io/namespaces/terraform-aws-modules)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [AWS EKS Best Practices](https://aws.github.io/aws-eks-best-practices/)
- [Prometheus Operator](https://prometheus-operator.dev/)
