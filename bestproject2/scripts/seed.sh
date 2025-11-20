#!/usr/bin/env bash
set -euo pipefail

API_PC=${API_PC:-http://localhost:8000}
API_UM=${API_UM:-http://localhost:8080}

echo "Seeding product-catalog..."
curl -sS -X POST "$API_PC/api/v1/products" -H 'Content-Type: application/json' -d '{"name":"Sample Product","price":19.99,"stock":100}' || true

echo "Checking user-management health..."
curl -sS "$API_UM/health" || true

echo "Done."
