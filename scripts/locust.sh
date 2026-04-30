#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"

PAYLOAD='{"customer_id":"cust-001","items":[{"product_id":"prod-000","category":"electronics","quantity":2,"unit_price":49.99}]}'

WAVES=20
CONCURRENCY=20
INTERVAL=3

echo "Releasing the swarm..."
echo ""

for wave in $(seq 1 "${WAVES}"); do
  for i in $(seq 1 "${CONCURRENCY}"); do
    curl -s -o /dev/null \
      --max-time 5 \
      -X POST "${BASE_URL}/orders" \
      -H "Content-Type: application/json" \
      -d "${PAYLOAD}" &
  done
  printf "\r\033[K  Wave %2d/%d in flight... (%ds elapsed)" "${wave}" "${WAVES}" "${SECONDS}"
  [ "${wave}" -lt "${WAVES}" ] && sleep "${INTERVAL}"
done

wait

echo ""
echo ""
echo "The swarm has landed. Open Grafana to see what happened."
