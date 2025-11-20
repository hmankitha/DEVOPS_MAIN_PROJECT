#!/bin/bash

set -e

echo "==================================="
echo "E-Commerce Platform Deployment"
echo "==================================="
echo ""

# Configuration
ENVIRONMENT=${1:-prod}
AWS_REGION=${AWS_REGION:-us-east-1}
EKS_CLUSTER="ecommerce-${ENVIRONMENT}-cluster"

echo "Environment: $ENVIRONMENT"
echo "AWS Region: $AWS_REGION"
echo "EKS Cluster: $EKS_CLUSTER"
echo ""

# Check prerequisites
echo "Checking prerequisites..."
command -v aws >/dev/null 2>&1 || { echo "aws CLI is required but not installed. Aborting." >&2; exit 1; }
command -v kubectl >/dev/null 2>&1 || { echo "kubectl is required but not installed. Aborting." >&2; exit 1; }
command -v terraform >/dev/null 2>&1 || { echo "terraform is required but not installed. Aborting." >&2; exit 1; }
command -v helm >/dev/null 2>&1 || { echo "helm is required but not installed. Aborting." >&2; exit 1; }
echo "✓ All prerequisites met"
echo ""

# Step 1: Provision infrastructure
echo "Step 1: Provisioning infrastructure with Terraform..."
cd infrastructure/terraform
terraform init
terraform plan -var-file=environments/${ENVIRONMENT}.tfvars -out=tfplan
terraform apply tfplan
cd ../..
echo "✓ Infrastructure provisioned"
echo ""

# Step 2: Configure kubectl
echo "Step 2: Configuring kubectl..."
aws eks update-kubeconfig --region $AWS_REGION --name $EKS_CLUSTER
kubectl cluster-info
echo "✓ kubectl configured"
echo ""

# Step 3: Create namespaces
echo "Step 3: Creating Kubernetes namespaces..."
kubectl apply -f infrastructure/kubernetes/namespaces/
echo "✓ Namespaces created"
echo ""

# Step 4: Deploy secrets
echo "Step 4: Deploying secrets..."
# Fetch secrets from AWS Secrets Manager and create K8s secrets
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
kubectl create secret docker-registry ecr-registry \
  --docker-server=${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com \
  --docker-username=AWS \
  --docker-password=$(aws ecr get-login-password --region $AWS_REGION) \
  --namespace=ecommerce-${ENVIRONMENT} \
  --dry-run=client -o yaml | kubectl apply -f -
echo "✓ Secrets deployed"
echo ""

# Step 5: Deploy microservices
echo "Step 5: Deploying microservices..."
kubectl apply -f infrastructure/kubernetes/deployments/user-management.yaml
kubectl apply -f infrastructure/kubernetes/deployments/product-catalog.yaml
echo "✓ Microservices deployed"
echo ""

# Step 6: Deploy ingress
echo "Step 6: Deploying ingress controller..."
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm upgrade --install ingress-nginx ingress-nginx/ingress-nginx \
  --namespace ingress-nginx \
  --create-namespace \
  --set controller.service.type=LoadBalancer
kubectl apply -f infrastructure/kubernetes/ingress/
echo "✓ Ingress deployed"
echo ""

# Step 7: Deploy monitoring
echo "Step 7: Deploying monitoring stack..."
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm upgrade --install prometheus prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  --create-namespace \
  --values monitoring/prometheus/values.yaml
echo "✓ Monitoring deployed"
echo ""

# Step 8: Deploy logging
echo "Step 8: Deploying logging stack..."
helm repo add elastic https://helm.elastic.co
helm repo update
helm upgrade --install elasticsearch elastic/elasticsearch \
  --namespace logging \
  --create-namespace
helm upgrade --install kibana elastic/kibana \
  --namespace logging
echo "✓ Logging deployed"
echo ""

# Step 9: Wait for deployments
echo "Step 9: Waiting for deployments to be ready..."
kubectl wait --for=condition=available --timeout=300s \
  deployment/user-management -n ecommerce-${ENVIRONMENT}
kubectl wait --for=condition=available --timeout=300s \
  deployment/product-catalog -n ecommerce-${ENVIRONMENT}
echo "✓ All deployments ready"
echo ""

# Step 10: Display access information
echo "==================================="
echo "Deployment Complete!"
echo "==================================="
echo ""
echo "Service Endpoints:"
INGRESS_IP=$(kubectl get svc ingress-nginx-controller -n ingress-nginx -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')
echo "  API Gateway: http://${INGRESS_IP}"
echo "  Grafana: http://${INGRESS_IP}/grafana"
echo "  Kibana: http://${INGRESS_IP}/kibana"
echo ""
echo "Kubernetes Dashboard:"
echo "  kubectl proxy"
echo "  http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/"
echo ""
echo "Monitoring:"
kubectl get pods -n monitoring
echo ""
echo "Application Pods:"
kubectl get pods -n ecommerce-${ENVIRONMENT}
echo ""
echo "To view logs:"
echo "  kubectl logs -f deployment/user-management -n ecommerce-${ENVIRONMENT}"
echo "  kubectl logs -f deployment/product-catalog -n ecommerce-${ENVIRONMENT}"
echo ""
