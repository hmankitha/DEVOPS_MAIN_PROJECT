@description('Azure region')
param location string

@description('Container Apps Environment resource ID')
param managedEnvId string

@description('ACR login server, e.g. myacr.azurecr.io')
param registryServer string

@description('User-assigned managed identity resource ID for app auth (ACR pull, Key Vault)')
param userAssignedIdentityId string

@description('Key Vault URI, e.g. https://mykv.vault.azure.net/')
param keyVaultUri string

@description('Key Vault secret name for JWT secret (referenced by all apps)')
param jwtSecretName string = 'jwt-secret'

@description('Image for product-catalog')
param productCatalogImage string

@description('Image for user-management')
param userManagementImage string

@description('Image for order-management')
param orderManagementImage string


resource productCatalog 'Microsoft.App/containerApps@2024-03-01' = {
  name: 'product-catalog'
  location: location
  identity: {
    type: 'UserAssigned'
    userAssignedIdentities: {
      '${userAssignedIdentityId}': {}
    }
  }
  properties: {
    managedEnvironmentId: managedEnvId
    configuration: {
      ingress: {
        external: true
        targetPort: 8000
      }
      registries: [
        {
          server: registryServer
          identity: userAssignedIdentityId
        }
      ]
      secrets: [
        {
          name: 'jwt-secret'
          identity: userAssignedIdentityId
          keyVaultUrl: '${keyVaultUri}secrets/${jwtSecretName}'
        }
      ]
    }
    template: {
      containers: [
        {
          name: 'app'
          image: productCatalogImage
          resources: {
            cpu: 0.5
            memory: '1Gi'
          }
          env: [
            {
              name: 'JWT_SECRET'
              secretRef: 'jwt-secret'
            }
          ]
        }
      ]
    }
  }
}

resource userManagement 'Microsoft.App/containerApps@2024-03-01' = {
  name: 'user-management'
  location: location
  identity: {
    type: 'UserAssigned'
    userAssignedIdentities: {
      '${userAssignedIdentityId}': {}
    }
  }
  properties: {
    managedEnvironmentId: managedEnvId
    configuration: {
      ingress: {
        external: true
        targetPort: 8080
      }
      registries: [
        {
          server: registryServer
          identity: userAssignedIdentityId
        }
      ]
      secrets: [
        {
          name: 'jwt-secret'
          identity: userAssignedIdentityId
          keyVaultUrl: '${keyVaultUri}secrets/${jwtSecretName}'
        }
      ]
    }
    template: {
      containers: [
        {
          name: 'app'
          image: userManagementImage
          resources: {
            cpu: 0.5
            memory: '1Gi'
          }
          env: [
            {
              name: 'JWT_SECRET'
              secretRef: 'jwt-secret'
            }
          ]
        }
      ]
    }
  }
}

resource orderManagement 'Microsoft.App/containerApps@2024-03-01' = {
  name: 'order-management'
  location: location
  identity: {
    type: 'UserAssigned'
    userAssignedIdentities: {
      '${userAssignedIdentityId}': {}
    }
  }
  properties: {
    managedEnvironmentId: managedEnvId
    configuration: {
      ingress: {
        external: true
        targetPort: 8090
      }
      registries: [
        {
          server: registryServer
          identity: userAssignedIdentityId
        }
      ]
      secrets: [
        {
          name: 'jwt-secret'
          identity: userAssignedIdentityId
          keyVaultUrl: '${keyVaultUri}secrets/${jwtSecretName}'
        }
      ]
    }
    template: {
      containers: [
        {
          name: 'app'
          image: orderManagementImage
          resources: {
            cpu: 0.5
            memory: '1Gi'
          }
          env: [
            {
              name: 'JWT_SECRET'
              secretRef: 'jwt-secret'
            }
          ]
        }
      ]
    }
  }
}
