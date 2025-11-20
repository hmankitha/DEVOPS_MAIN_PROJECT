#!/usr/bin/env bash
set -euo pipefail

# Azure deployment helper for Azure Container Apps + ACR (Bicep)
# Prereqs: az CLI logged in (az login), Docker, JDK/Maven/Go/Python for builds

RG_NAME=${RG_NAME:-bestproject2-rg}
LOCATION=${LOCATION:-eastus}
ACA_ENV_NAME=${ACA_ENV_NAME:-bp2-aca-env}
ACR_NAME=${ACR_NAME:-}

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
INFRA_DIR="$ROOT_DIR/infrastructure/azure/bicep"

rand_suffix() { LC_ALL=C tr -dc 'a-z0-9' </dev/urandom | head -c 6; echo; }

if [[ -z "${ACR_NAME}" ]]; then
  ACR_NAME="bp2acr$(rand_suffix)" # must be globally unique and 5-50 chars
fi

echo "==> Using resource group: $RG_NAME in $LOCATION"
az group create -n "$RG_NAME" -l "$LOCATION" >/dev/null

# 1) Deploy infra (ACR + Log Analytics + Container Apps Environment)
echo "==> Deploying Bicep (infra)"
az deployment group create \
  -g "$RG_NAME" \
  -f "$INFRA_DIR/main.bicep" \
  -p location="$LOCATION" acrName="$ACR_NAME" acaEnvName="$ACA_ENV_NAME" >/tmp/azure-deploy-infra.json

ACR_LOGIN_SERVER=$(jq -r '.properties.outputs.acrLoginServer.value' /tmp/azure-deploy-infra.json)
CONTAINERAPPS_ENV_ID=$(jq -r '.properties.outputs.containerAppsEnvironmentId.value' /tmp/azure-deploy-infra.json)

echo "==> ACR: $ACR_LOGIN_SERVER"

# 2) Build and push images to ACR (docker login required)
echo "==> Logging into ACR"
az acr login -n "$ACR_NAME"

pushd "$ROOT_DIR" >/dev/null

# product-catalog
IMG_PC="$ACR_LOGIN_SERVER/product-catalog:latest"
docker build -t "$IMG_PC" microservices/product-catalog

docker push "$IMG_PC"

# user-management
IMG_UM="$ACR_LOGIN_SERVER/user-management:latest"
docker build -t "$IMG_UM" microservices/user-management

docker push "$IMG_UM"

# order-management
IMG_OM="$ACR_LOGIN_SERVER/order-management:latest"
docker build -t "$IMG_OM" microservices/order-management

docker push "$IMG_OM"

popd >/dev/null

# 3) Update Container Apps images (idempotent via bicep or CLI)
echo "==> Updating Container Apps images"
az deployment group create \
  -g "$RG_NAME" \
  -f "$INFRA_DIR/main.bicep" \
  -p location="$LOCATION" acrName="$ACR_NAME" acaEnvName="$ACA_ENV_NAME" \
     productCatalogImage="$IMG_PC" userManagementImage="$IMG_UM" orderManagementImage="$IMG_OM" >/dev/null

echo "==> Done"
echo "Container Apps Environment: $CONTAINERAPPS_ENV_ID"
echo "Images:" 
printf ' - %s\n' "$IMG_PC" "$IMG_UM" "$IMG_OM"
