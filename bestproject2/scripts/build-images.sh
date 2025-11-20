#!/bin/bash

set -e

echo "Building Docker images..."

# Configuration
AWS_REGION=${AWS_REGION:-us-east-1}
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
ECR_REGISTRY="${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com"
VERSION=${VERSION:-latest}

# Login to ECR
echo "Logging in to Amazon ECR..."
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $ECR_REGISTRY

# Build and push User Management service
echo "Building user-management service..."
cd microservices/user-management
docker build -t user-management:$VERSION .
docker tag user-management:$VERSION $ECR_REGISTRY/user-management:$VERSION
docker tag user-management:$VERSION $ECR_REGISTRY/user-management:latest

# Create ECR repository if it doesn't exist
aws ecr describe-repositories --repository-names user-management --region $AWS_REGION || \
  aws ecr create-repository --repository-name user-management --region $AWS_REGION

docker push $ECR_REGISTRY/user-management:$VERSION
docker push $ECR_REGISTRY/user-management:latest
echo "✓ user-management pushed to ECR"
cd ../..

# Build and push Product Catalog service
echo "Building product-catalog service..."
cd microservices/product-catalog
docker build -t product-catalog:$VERSION .
docker tag product-catalog:$VERSION $ECR_REGISTRY/product-catalog:$VERSION
docker tag product-catalog:$VERSION $ECR_REGISTRY/product-catalog:latest

aws ecr describe-repositories --repository-names product-catalog --region $AWS_REGION || \
  aws ecr create-repository --repository-name product-catalog --region $AWS_REGION

docker push $ECR_REGISTRY/product-catalog:$VERSION
docker push $ECR_REGISTRY/product-catalog:latest
echo "✓ product-catalog pushed to ECR"
cd ../..

echo ""
echo "All images built and pushed successfully!"
echo "Images:"
echo "  $ECR_REGISTRY/user-management:$VERSION"
echo "  $ECR_REGISTRY/product-catalog:$VERSION"
