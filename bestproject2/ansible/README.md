# Ansible for E-Commerce Microservices Platform

This directory contains Ansible playbooks and roles for automating the deployment and management of the microservices platform.

## Directory Structure

```
ansible/
├── ansible.cfg              # Ansible configuration
├── inventory/
│   └── hosts.yml           # Inventory file with all hosts
├── playbooks/
│   ├── deploy.yml          # Main deployment playbook
│   ├── start-services.yml  # Start all services
│   └── stop-services.yml   # Stop all services
├── roles/
│   ├── docker/             # Docker setup and verification
│   ├── databases/          # PostgreSQL and Redis management
│   ├── elk-stack/          # Elasticsearch and Kibana setup
│   ├── microservices/      # Microservices deployment
│   ├── prometheus/         # Prometheus monitoring
│   └── grafana/            # Grafana dashboards
└── group_vars/             # Group-specific variables
```

## Prerequisites

```bash
# Install Ansible on macOS
brew install ansible

# Install required collections
ansible-galaxy collection install community.docker
ansible-galaxy collection install community.general
```

## Usage

### Start All Services

```bash
cd ansible
ansible-playbook playbooks/start-services.yml
```

### Deploy Microservices

```bash
# Deploy all services
ansible-playbook playbooks/deploy.yml

# Deploy specific service
ansible-playbook playbooks/deploy.yml --tags "python,product-catalog"
ansible-playbook playbooks/deploy.yml --tags "go,user-management"

# Deploy only infrastructure
ansible-playbook playbooks/deploy.yml --tags "infrastructure"

# Deploy only monitoring
ansible-playbook playbooks/deploy.yml --tags "monitoring"
```

### Stop All Services

```bash
ansible-playbook playbooks/stop-services.yml
```

### Health Check

```bash
ansible-playbook playbooks/deploy.yml --tags "healthcheck"
```

## Available Tags

- `infrastructure` - Setup databases, docker, ELK
- `docker` - Docker verification
- `databases` - PostgreSQL and Redis
- `elk` - Elasticsearch and Kibana
- `microservices` - Deploy microservices
- `python` - Python services
- `go` - Go services
- `product-catalog` - Product Catalog service
- `user-management` - User Management service
- `monitoring` - Prometheus and Grafana
- `healthcheck` - Service health checks

## Inventory

The inventory file (`inventory/hosts.yml`) defines:

- **local**: Local development environment
- **database_servers**: PostgreSQL and Redis
- **app_servers**: Python and Go microservices
- **monitoring**: ELK, Prometheus, Grafana
- **cache_servers**: Redis

## Variables

Global variables are defined in `inventory/hosts.yml`:

- `project_root`: Project directory path
- `environment`: Deployment environment (production/development)
- `jwt_secret`: JWT secret key
- `db_user`, `db_password`: Database credentials

## Examples

### Deploy Everything

```bash
ansible-playbook playbooks/deploy.yml
```

### Deploy Only Python Service

```bash
ansible-playbook playbooks/deploy.yml --limit python_services
```

### Check Service Health

```bash
ansible-playbook playbooks/deploy.yml --tags healthcheck --limit app_servers
```

### Restart Specific Service

```bash
# Stop service
ansible-playbook playbooks/stop-services.yml

# Deploy specific service
ansible-playbook playbooks/deploy.yml --tags "python,product-catalog"
```

## Logs

Service logs are stored in:
- Product Catalog: `/tmp/product-catalog.log`
- User Management: `/tmp/user-management.log`
- Elasticsearch logs: Via Kibana at http://localhost:5601

## Troubleshooting

### Check Service Status

```bash
# Check if services are running
docker ps

# Check logs
tail -f /tmp/product-catalog.log
tail -f /tmp/user-management.log
```

### Verify Connectivity

```bash
curl http://localhost:8000/api/v1/products
curl http://localhost:8080/health
```

### Run Ansible in Verbose Mode

```bash
ansible-playbook playbooks/deploy.yml -v   # verbose
ansible-playbook playbooks/deploy.yml -vv  # more verbose
ansible-playbook playbooks/deploy.yml -vvv # very verbose
```

## Security Notes

**Important**: The current configuration uses plaintext passwords in the inventory file. For production:

1. Use Ansible Vault to encrypt sensitive data:
```bash
ansible-vault create group_vars/all/vault.yml
ansible-vault encrypt_string 'your-secret' --name 'jwt_secret'
```

2. Run playbooks with vault password:
```bash
ansible-playbook playbooks/deploy.yml --ask-vault-pass
```

## Maintenance

### Update Dependencies

```bash
# Update Python dependencies
ansible-playbook playbooks/deploy.yml --tags python

# Rebuild Go binary
ansible-playbook playbooks/deploy.yml --tags go
```

### Cleanup

```bash
# Stop all services
ansible-playbook playbooks/stop-services.yml

# Remove docker volumes (destructive)
docker-compose down -v
```
