#!/bin/bash

# ArgoCD Installation Script
# This script installs ArgoCD on your Kubernetes cluster

set -e

ARGOCD_VERSION="v2.9.3"
NAMESPACE="argocd"

echo "================================================"
echo "🚀 Installing ArgoCD ${ARGOCD_VERSION}"
echo "================================================"

# Create namespace
echo "📦 Creating namespace: ${NAMESPACE}"
kubectl create namespace ${NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -

# Install ArgoCD
echo "📥 Installing ArgoCD components..."
kubectl apply -n ${NAMESPACE} -f https://raw.githubusercontent.com/argoproj/argo-cd/${ARGOCD_VERSION}/manifests/install.yaml

# Wait for ArgoCD to be ready
echo "⏳ Waiting for ArgoCD pods to be ready..."
kubectl wait --for=condition=Ready pods --all -n ${NAMESPACE} --timeout=300s

# Get initial admin password
echo ""
echo "================================================"
echo "✅ ArgoCD Installation Complete!"
echo "================================================"
echo ""
echo "📝 Access Information:"
echo ""
echo "1️⃣  Port Forward ArgoCD Server:"
echo "   kubectl port-forward svc/argocd-server -n argocd 8081:443"
echo ""
echo "2️⃣  Access ArgoCD UI:"
echo "   https://localhost:8081"
echo ""
echo "3️⃣  Login Credentials:"
echo "   Username: admin"
echo -n "   Password: "
kubectl -n ${NAMESPACE} get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
echo ""
echo ""
echo "4️⃣  Login via CLI:"
echo "   argocd login localhost:8081 --username admin --password \$(kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath=\"{.data.password}\" | base64 -d)"
echo ""
echo "================================================"
echo "🎯 Next Steps:"
echo "   1. Apply AppProject: kubectl apply -f argocd/projects/microservices-project.yaml"
echo "   2. Deploy Apps: kubectl apply -f argocd/applications/"
echo "================================================"
