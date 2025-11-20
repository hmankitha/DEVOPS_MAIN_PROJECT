@description('Azure location for all resources')
param location string = resourceGroup().location

@description('Container Apps environment name')
param acaEnvName string = 'bp2-aca-env'

@description('Azure Container Registry name (must be globally unique, 5-50 alphanumeric)')
param acrName string = toLower('bp2acr${uniqueString(resourceGroup().id)}')

@description('Log Analytics Workspace name')
param logAnalyticsName string = 'bp2-law'

@description('Key Vault name')
param keyVaultName string = toLower('bp2-kv-${uniqueString(resourceGroup().id)}')

@description('User-assigned managed identity name for apps')
param appIdentityName string = 'bp2-app-mi'

@description('Image for product-catalog')
param productCatalogImage string = ''

@description('Image for user-management')
param userManagementImage string = ''

@description('Image for order-management')
param orderManagementImage string = ''

// ACR
resource acr 'Microsoft.ContainerRegistry/registries@2023-01-01-preview' = {
  name: acrName
  location: location
  sku: {
    name: 'Basic'
  }
  properties: {
    // Use RBAC (AcrPull) with a user-assigned identity instead of admin creds
    adminUserEnabled: false
  }
}

// Log Analytics
resource law 'Microsoft.OperationalInsights/workspaces@2022-10-01' = {
  name: logAnalyticsName
  location: location
  properties: {
    retentionInDays: 30
  }
}

var lawKeys = listKeys(resourceId('Microsoft.OperationalInsights/workspaces', law.name), '2020-08-01')

// Container Apps Environment
resource acaEnv 'Microsoft.App/managedEnvironments@2024-03-01' = {
  name: acaEnvName
  location: location
  properties: {
    appLogsConfiguration: {
      destination: 'log-analytics'
      logAnalyticsConfiguration: {
        customerId: law.properties.customerId
        sharedKey: lawKeys.primarySharedKey
      }
    }
  }
}

// Key Vault (RBAC enabled)
resource kv 'Microsoft.KeyVault/vaults@2023-02-01' = {
  name: keyVaultName
  location: location
  properties: {
    enableRbacAuthorization: true
    tenantId: subscription().tenantId
    sku: {
      name: 'standard'
      family: 'A'
    }
    softDeleteRetentionInDays: 90
    enablePurgeProtection: true
    publicNetworkAccess: 'Enabled'
  }
}

// User-assigned managed identity for Container Apps (pull images + read KV secrets)
resource appMI 'Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31' = {
  name: appIdentityName
  location: location
}

// RBAC: Allow MI to pull from ACR (AcrPull)
resource acrPullAssignment 'Microsoft.Authorization/roleAssignments@2022-04-01' = {
  name: guid(acr.id, appMI.id, 'AcrPull')
  scope: acr
  properties: {
    principalId: appMI.properties.principalId
    roleDefinitionId: subscriptionResourceId('Microsoft.Authorization/roleDefinitions', '7f951dda-4ed3-4680-a7ca-43fe172d538d')
    principalType: 'ServicePrincipal'
  }
}

// RBAC: Allow MI to read Key Vault secrets (Key Vault Secrets User)
resource kvSecretsUserAssignment 'Microsoft.Authorization/roleAssignments@2022-04-01' = {
  name: guid(kv.id, appMI.id, 'KVSecretsUser')
  scope: kv
  properties: {
    principalId: appMI.properties.principalId
    roleDefinitionId: subscriptionResourceId('Microsoft.Authorization/roleDefinitions', '4633458b-17de-408a-b874-0445c86b69e6')
    principalType: 'ServicePrincipal'
  }
}

module apps './apps.bicep' = {
  name: 'apps'
  params: {
    location: location
    managedEnvId: acaEnv.id
    registryServer: acr.properties.loginServer
    // use user-assigned identity for ACR auth and Key Vault secret references
    userAssignedIdentityId: appMI.id
    keyVaultUri: kv.properties.vaultUri
    productCatalogImage: empty(productCatalogImage) ? '${acr.properties.loginServer}/product-catalog:latest' : productCatalogImage
    userManagementImage: empty(userManagementImage) ? '${acr.properties.loginServer}/user-management:latest' : userManagementImage
    orderManagementImage: empty(orderManagementImage) ? '${acr.properties.loginServer}/order-management:latest' : orderManagementImage
  }
}

output containerAppsEnvironmentId string = acaEnv.id
output acrLoginServer string = acr.properties.loginServer
output acrNameOut string = acr.name
output keyVaultUri string = kv.properties.vaultUri
output appIdentityResourceId string = appMI.id
