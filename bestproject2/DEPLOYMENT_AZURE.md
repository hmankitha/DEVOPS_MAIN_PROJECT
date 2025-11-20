# Azure Deployment (Azure Container Apps + ACR)

This project includes Azure deployment via Bicep that provisions:
- Azure Container Registry (ACR)
- Azure Log Analytics Workspace
- Azure Container Apps Environment (ACA)
- Three Container Apps (product-catalog, user-management, order-management)
- Azure Key Vault (RBAC enabled) and a user-assigned managed identity for secure secrets & ACR access

## Prerequisites
- Azure CLI logged in: `az login`
- Docker running (to build images)
- Permissions to create RG, ACR, Container Apps
- Optional: `jq` for parsing outputs (script uses it)

## One-shot deploy
```zsh
# From repo root
export RG_NAME=bestproject2-rg
export LOCATION=eastus
# Optionally set ACR name (must be globally unique, 5-50 alphanumeric)
# export ACR_NAME=bp2acr123abc

bash scripts/azure-deploy.sh
```
The script will:
1) Create or reuse the resource group.
2) Deploy Bicep infra (ACR [admin disabled], Key Vault [RBAC], Managed Identity, Log Analytics, Container Apps Env).
3) Build and push Docker images for all three services to ACR.
4) Re-deploy templates with the pushed image tags so Container Apps run the latest images.

## Manual (advanced)
Deploy infra:
```zsh
az group create -n $RG_NAME -l $LOCATION
az deployment group create -g $RG_NAME \
  -f infrastructure/azure/bicep/main.bicep \
  -p location=$LOCATION acrName=$ACR_NAME acaEnvName=bp2-aca-env
```
Build and push images:
```zsh
ACR_LOGIN_SERVER=$(az acr show -n $ACR_NAME --query loginServer -o tsv)
az acr login -n $ACR_NAME

docker build -t $ACR_LOGIN_SERVER/product-catalog:latest microservices/product-catalog
 docker push $ACR_LOGIN_SERVER/product-catalog:latest

docker build -t $ACR_LOGIN_SERVER/user-management:latest microservices/user-management
 docker push $ACR_LOGIN_SERVER/user-management:latest

docker build -t $ACR_LOGIN_SERVER/order-management:latest microservices/order-management
 docker push $ACR_LOGIN_SERVER/order-management:latest
```
Update apps to latest images:
```zsh
az deployment group create -g $RG_NAME \
  -f infrastructure/azure/bicep/main.bicep \
  -p productCatalogImage=$ACR_LOGIN_SERVER/product-catalog:latest \
     userManagementImage=$ACR_LOGIN_SERVER/user-management:latest \
     orderManagementImage=$ACR_LOGIN_SERVER/order-management:latest
```

## Configuration and environment variables
- Default container ports: product-catalog 8000, user-management 8080, order-management 8090.
- Secrets are read from Key Vault using a user-assigned managed identity. By default, an env var `JWT_SECRET` is sourced from the Key Vault secret name `jwt-secret`.
- Set the secret value (example):
```zsh
KV_NAME=$(az deployment group show -g $RG_NAME -n apps --query properties.outputs.keyVaultUri.value -o tsv | sed -E 's#https://([^\.]+).*#\1#')
az keyvault secret set --vault-name "$KV_NAME" --name jwt-secret --value "change-me-secure"
```
- To add more secrets, extend `infrastructure/azure/bicep/apps.bicep` (add entries under `configuration.secrets` and map to container `env` with `secretRef`).
- For managed databases (e.g., Azure Database for PostgreSQL Flexible Server) and Redis (Azure Cache for Redis), add resources and wire connection strings as secrets.

## Observability
- Container Apps emit logs to the provisioned Log Analytics Workspace.
- Consider Azure Managed Grafana and Azure Monitor for production monitoring.

## Security hardening summary
- ACR admin user disabled; images pulled via RBAC with a user-assigned managed identity (role: AcrPull).
- Secrets fetched from Key Vault using RBAC (role: Key Vault Secrets User). No secrets are stored in templates.
- You can lock down ingress and use VNET injection for Container Apps in advanced setups.

## Cleanup
```zsh
az group delete -n $RG_NAME --yes --no-wait
```
