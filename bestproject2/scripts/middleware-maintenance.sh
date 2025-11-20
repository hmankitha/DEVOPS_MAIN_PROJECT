#!/usr/bin/env bash
set -euo pipefail

COMPOSE_FILE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$COMPOSE_FILE_DIR"

usage() {
  cat <<'USAGE'
Middleware maintenance helper

Usage:
  scripts/middleware-maintenance.sh start            # start mysql, redis, kafka and exporters
  scripts/middleware-maintenance.sh stop             # stop mysql, redis, kafka and exporters
  scripts/middleware-maintenance.sh status           # show container status

  scripts/middleware-maintenance.sh backup-mysql <out.sql>
  scripts/middleware-maintenance.sh restore-mysql <in.sql>

  scripts/middleware-maintenance.sh backup-redis <out.rdb>
  scripts/middleware-maintenance.sh restore-redis <in.rdb>

  scripts/middleware-maintenance.sh kafka-topics
  scripts/middleware-maintenance.sh kafka-create-topic <topic> [partitions] [replication]

Notes:
- Requires Docker Desktop running.
- MySQL root password is taken from docker-compose (MYSQL_ROOT_PASSWORD). Default: mysql
USAGE
}

require_compose() {
  if ! command -v docker >/dev/null 2>&1; then
    echo "docker not found on PATH" >&2
    exit 1
  fi
}

start() {
  require_compose
  docker compose up -d mysql redis kafka redis-exporter kafka-exporter mysqld-exporter kafka-ui
}

stop() {
  require_compose
  docker compose stop mysql redis kafka redis-exporter kafka-exporter mysqld-exporter kafka-ui
}

status() {
  require_compose
  docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}' | (grep -E 'mysql|redis|kafka|exporter|kafka-ui' || true)
}

backup_mysql() {
  require_compose
  local out=${1:-}
  [ -n "$out" ] || { echo "Output file required"; exit 2; }
  local pw
  pw=$(awk -F= '/MYSQL_ROOT_PASSWORD/{print $2}' docker-compose.yml | head -n1 | tr -d '" ')
  pw=${pw:-mysql}
  docker compose exec -T mysql sh -c "mysqldump -uroot -p$pw --all-databases" > "$out"
  echo "MySQL backup written to $out"
}

restore_mysql() {
  require_compose
  local in=${1:-}
  [ -f "$in" ] || { echo "Input SQL file not found: $in"; exit 2; }
  local pw
  pw=$(awk -F= '/MYSQL_ROOT_PASSWORD/{print $2}' docker-compose.yml | head -n1 | tr -d '" ')
  pw=${pw:-mysql}
  cat "$in" | docker compose exec -T mysql sh -c "mysql -uroot -p$pw"
  echo "MySQL restore completed from $in"
}

backup_redis() {
  require_compose
  local out=${1:-}
  [ -n "$out" ] || { echo "Output file required"; exit 2; }
  docker compose exec -T redis redis-cli --rdb /data/dump.rdb
  cid=$(docker compose ps -q redis)
  docker cp "$cid:/data/dump.rdb" "$out"
  echo "Redis RDB backup written to $out"
}

restore_redis() {
  require_compose
  local in=${1:-}
  [ -f "$in" ] || { echo "Input RDB file not found: $in"; exit 2; }
  cid=$(docker compose ps -q redis)
  docker cp "$in" "$cid:/data/dump.rdb"
  docker compose restart redis
  echo "Redis restore completed from $in"
}

kafka_topics() {
  require_compose
  docker compose exec -T kafka kafka-topics.sh --bootstrap-server kafka:9092 --list || true
}

kafka_create_topic() {
  require_compose
  local topic=${1:-}
  local partitions=${2:-1}
  local replication=${3:-1}
  [ -n "$topic" ] || { echo "Topic name required"; exit 2; }
  docker compose exec -T kafka kafka-topics.sh --bootstrap-server kafka:9092 \
    --create --if-not-exists --topic "$topic" --partitions "$partitions" --replication-factor "$replication" || true
}

cmd=${1:-}
case "$cmd" in
  start) start ;;
  stop) stop ;;
  status) status ;;
  backup-mysql) shift; backup_mysql "$@" ;;
  restore-mysql) shift; restore_mysql "$@" ;;
  backup-redis) shift; backup_redis "$@" ;;
  restore-redis) shift; restore_redis "$@" ;;
  kafka-topics) kafka_topics ;;
  kafka-create-topic) shift; kafka_create_topic "$@" ;;
  -h|--help|help|"") usage ;;
  *) echo "Unknown command: $cmd"; usage; exit 2 ;;
 esac
