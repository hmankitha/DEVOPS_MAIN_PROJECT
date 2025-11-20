# Middleware Maintenance (MySQL, Redis, Kafka)

Local middleware services and Prometheus exporters are available via Docker Compose and Make targets.

## Requirements
- Docker Desktop running
- `make` and `bash`

## Lifecycle
```bash
# Start / Stop / Status
make middleware-start
make middleware-status
make middleware-stop
```

## Backups / Restore
```bash
# MySQL
make middleware-backup-mysql OUT=backups/mysql-$(date +%F).sql
make middleware-restore-mysql IN=backups/mysql.sql

# Redis
make middleware-backup-redis OUT=backups/redis-$(date +%F).rdb
make middleware-restore-redis IN=backups/redis.rdb
```

## Kafka Utilities
```bash
make kafka-topics
make kafka-create-topic NAME=orders PARTITIONS=3 REPL=1
```

## Monitoring
Prometheus scrapes Redis, Kafka, and MySQL exporters for local development. See `monitoring/prometheus/prometheus.yml`.
