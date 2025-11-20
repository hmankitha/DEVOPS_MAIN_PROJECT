#!/bin/bash

# Security scanning script using Trivy

set -e

echo "Running security scans with Trivy..."

SEVERITY="HIGH,CRITICAL"
SCAN_DIR=${1:-.}

# Scan filesystem
echo "Scanning filesystem..."
trivy fs --severity $SEVERITY --exit-code 1 $SCAN_DIR

# Scan Docker images if they exist
if command -v docker &> /dev/null; then
    echo "Scanning Docker images..."
    
    # User Management
    if docker images | grep -q "user-management"; then
        echo "Scanning user-management image..."
        trivy image --severity $SEVERITY user-management:latest
    fi
    
    # Product Catalog
    if docker images | grep -q "product-catalog"; then
        echo "Scanning product-catalog image..."
        trivy image --severity $SEVERITY product-catalog:latest
    fi
fi

# Scan Kubernetes manifests
echo "Scanning Kubernetes manifests..."
trivy config --severity $SEVERITY infrastructure/kubernetes/

# Scan Terraform files
echo "Scanning Terraform configurations..."
trivy config --severity $SEVERITY infrastructure/terraform/

echo "Security scan completed!"
