# Architecture Document

## System Architecture

### High-Level Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                          API Gateway / Ingress                   │
│                        (AWS ALB / NGINX Ingress)                │
└────────────────┬────────────────────────────────────────────────┘
                 │
        ┌────────┴────────┐
        │                 │
┌───────▼──────┐  ┌──────▼────────┐  ┌──────────────┐
│     User     │  │    Product    │  │    Order     │
│  Management  │  │   Catalog     │  │  Management  │
│   (Golang)   │  │   (Python)    │  │    (Java)    │
└──────┬───────┘  └───────┬───────┘  └──────┬───────┘
       │                  │                  │
       │         ┌────────┴─────────┐        │
       │         │                  │        │
┌──────▼─────┐  ┌▼──────────┐  ┌───▼────┐  ┌▼────────┐
│ PostgreSQL │  │ PostgreSQL│  │  Redis │  │   S3    │
│  (Users)   │  │ (Products)│  │ (Cache)│  │(Storage)│
└────────────┘  └───────────┘  └────────┘  └─────────┘
```

## Microservices Architecture

### 1. User Management Service
**Technology**: Golang + Gin Framework
**Database**: PostgreSQL
**Port**: 8080

**Responsibilities**:
- User registration and authentication
- JWT token generation and validation
- User profile management
- Password reset functionality
- Role-based access control (RBAC)

**API Endpoints**:
- POST /api/v1/auth/register
- POST /api/v1/auth/login
- POST /api/v1/auth/refresh
- GET /api/v1/users/me
- PUT /api/v1/users/me

### 2. Product Catalog Service
**Technology**: Python + FastAPI
**Database**: PostgreSQL
**Cache**: Redis
**Search**: Elasticsearch
**Port**: 8000

**Responsibilities**:
- Product management (CRUD)
- Category management
- Inventory tracking
- Product search and filtering
- Reviews and ratings

**API Endpoints**:
- GET /api/v1/products
- POST /api/v1/products
- GET /api/v1/products/{id}
- PUT /api/v1/products/{id}
- GET /api/v1/categories
- GET /api/v1/inventory/{product_id}

### 3. Order Management Service
**Technology**: Java + Spring Boot
**Database**: PostgreSQL
**Port**: 8081

**Responsibilities**:
- Order creation and processing
- Order status tracking
- Payment integration
- Order history
- Shipping management

## Infrastructure Components

### AWS Services

#### 1. VPC Architecture
- **CIDR**: 10.0.0.0/16
- **Availability Zones**: 3 (us-east-1a, us-east-1b, us-east-1c)
- **Public Subnets**: 3 (for load balancers)
- **Private Subnets**: 3 (for application workloads)
- **NAT Gateways**: 3 (one per AZ for high availability)

#### 2. EKS (Elastic Kubernetes Service)
- **Version**: 1.28
- **Node Groups**:
  - General: 3-10 nodes (t3.large, ON_DEMAND)
  - Spot: 1-5 nodes (t3.large/t3a.large, SPOT)
- **Add-ons**:
  - VPC CNI
  - CoreDNS
  - kube-proxy
  - AWS Load Balancer Controller

#### 3. RDS (Relational Database Service)
- **Engine**: PostgreSQL 15.4
- **Instance Types**: 
  - Dev: db.t3.medium
  - Prod: db.r6g.xlarge
- **Storage**: 100-500 GB SSD
- **Multi-AZ**: Enabled in production
- **Automated Backups**: 7-30 days retention
- **Encryption**: At rest and in transit

#### 4. ElastiCache (Redis)
- **Engine**: Redis 7.x
- **Node Type**:
  - Dev: cache.t3.medium
  - Prod: cache.r6g.large
- **Cluster Mode**: Disabled (single node in dev, cluster in prod)

#### 5. S3 Buckets
- **ecommerce-users**: User avatars and documents
- **ecommerce-products**: Product images
- **ecommerce-orders**: Order documents
- **ecommerce-terraform-state**: Terraform state files
- **Versioning**: Enabled
- **Encryption**: AES-256

#### 6. ECR (Elastic Container Registry)
- Stores Docker images for all microservices
- Lifecycle policies for automatic cleanup
- Image scanning enabled

#### 7. Secrets Manager
- Stores sensitive configuration
- Automatic rotation enabled
- Secrets:
  - Database credentials
  - JWT secrets
  - API keys
  - Third-party service credentials

#### 8. IAM Roles
- EKS Cluster Role
- EKS Node Group Role
- Service Account Roles (IRSA)
- CI/CD Pipeline Roles

## CI/CD Pipeline

### Jenkins Pipeline Stages

1. **Checkout**: Pull code from repository
2. **Build**: Compile application
3. **Unit Tests**: Run automated tests
4. **Code Quality**: SonarQube analysis
5. **Quality Gate**: Enforce quality standards
6. **Docker Build**: Create container image
7. **Security Scan**: Trivy vulnerability scanning
8. **Push to ECR**: Upload image to registry
9. **Deploy to Staging**: Automatic deployment
10. **Integration Tests**: End-to-end tests
11. **Deploy to Production**: Manual approval required
12. **Smoke Tests**: Verify production deployment

### Deployment Strategy
- **Strategy**: Rolling Update
- **Max Surge**: 1
- **Max Unavailable**: 0
- **Zero-downtime deployments**

## Monitoring & Observability

### Prometheus
- **Metrics Collection**: Every 15 seconds
- **Retention**: 15 days
- **Monitored Metrics**:
  - HTTP request rate and duration
  - Error rates
  - CPU and memory usage
  - Database connections
  - Cache hit rates

### Grafana
- **Dashboards**:
  - E-Commerce Overview
  - Service Health
  - Infrastructure Metrics
  - Business Metrics
- **Alerts**: Integrated with Prometheus AlertManager

### ELK Stack (Elasticsearch, Logstash, Kibana)
- **Log Aggregation**: Centralized logging
- **Retention**: 30 days
- **Indices**: Per namespace per day
- **Log Levels**: INFO, WARN, ERROR, CRITICAL

## Security Architecture

### Network Security
- **Security Groups**: Restrictive ingress/egress rules
- **Network Policies**: Kubernetes network isolation
- **TLS/SSL**: End-to-end encryption
- **WAF**: AWS Web Application Firewall

### Application Security
- **Authentication**: JWT tokens
- **Authorization**: Role-based access control
- **Input Validation**: All API inputs validated
- **Rate Limiting**: Protection against abuse
- **CORS**: Configured for allowed origins

### Container Security
- **Image Scanning**: Trivy scans on every build
- **Non-root Users**: Containers run as non-root
- **Read-only Filesystems**: Where possible
- **Security Contexts**: Applied to all pods

### Data Security
- **Encryption at Rest**: All databases and S3 buckets
- **Encryption in Transit**: TLS 1.2+
- **Secrets Management**: AWS Secrets Manager
- **Backup**: Automated daily backups

## High Availability & Disaster Recovery

### High Availability
- **Multi-AZ Deployment**: Services across 3 AZs
- **Load Balancing**: AWS ALB for traffic distribution
- **Auto-scaling**: HPA for pods, ASG for nodes
- **Health Checks**: Liveness and readiness probes

### Disaster Recovery
- **RTO (Recovery Time Objective)**: 1 hour
- **RPO (Recovery Point Objective)**: 15 minutes
- **Backup Strategy**:
  - Database: Automated daily snapshots
  - Configuration: Version controlled in Git
  - Secrets: Replicated in Secrets Manager
- **Recovery Procedures**: Documented runbooks

## Scaling Strategy

### Horizontal Pod Autoscaling (HPA)
- **Metrics**: CPU (70%), Memory (80%)
- **Min Replicas**: 3
- **Max Replicas**: 10
- **Scale Up**: Aggressive (quick response to load)
- **Scale Down**: Conservative (avoid thrashing)

### Cluster Autoscaling
- **Min Nodes**: 3
- **Max Nodes**: 20
- **Scale Up**: When pods can't be scheduled
- **Scale Down**: After 10 minutes of low utilization

## Performance Targets

### API Response Times
- **P50**: < 100ms
- **P95**: < 500ms
- **P99**: < 1000ms

### Throughput
- **Requests per Second**: 10,000+
- **Concurrent Users**: 50,000+

### Availability
- **SLA**: 99.9% uptime
- **Planned Downtime**: < 4 hours/year

## Cost Optimization

### Strategies
- **Spot Instances**: For non-critical workloads
- **Auto-scaling**: Scale down during low traffic
- **Reserved Instances**: For baseline capacity
- **S3 Lifecycle Policies**: Archive old data
- **Right-sizing**: Regular review of resource allocation

### Estimated Monthly Costs (Production)
- **EKS**: $500-1000
- **EC2 (Nodes)**: $1000-2000
- **RDS**: $500-1000
- **ElastiCache**: $200-400
- **S3**: $100-200
- **Data Transfer**: $200-500
- **Other Services**: $200-400
- **Total**: $2700-5500/month

## Compliance & Governance

### Standards
- **PCI DSS**: For payment processing
- **GDPR**: For user data protection
- **SOC 2**: Security and availability

### Auditing
- **CloudTrail**: All API calls logged
- **Audit Logs**: Application-level audit trails
- **Compliance Reports**: Quarterly reviews

## Future Enhancements

### Phase 2
- Service mesh (Istio/Linkerd)
- Multi-region deployment
- Advanced caching strategies
- Real-time analytics

### Phase 3
- Machine learning recommendations
- Event-driven architecture
- Serverless components
- Edge computing integration
