#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"

ORDER_PAYLOAD='{"customer_id":"cust-002","items":[{"product_id":"prod-001","category":"clothing","quantity":1,"unit_price":29.99}]}'
ORDER_ID=$(curl -s -X POST "${BASE_URL}/orders" \
  -H "Content-Type: application/json" \
  -d "${ORDER_PAYLOAD}" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ -z "${ORDER_ID}" ]; then
  echo "Could not reach the API. Is the stack running?" >&2
  exit 1
fi

echo "Leaving the light on..."
echo ""

NOTIFY_PAYLOAD='{"event":"order.confirmed"}'
END=$((SECONDS + 90))
BATCH=10
COUNT=0

while [ $SECONDS -lt $END ]; do
  for i in $(seq 1 $BATCH); do
    curl -s -o /dev/null \
      --max-time 3 \
      -X POST "${BASE_URL}/orders/${ORDER_ID}/notify" \
      -H "Content-Type: application/json" \
      -d "${NOTIFY_PAYLOAD}" &
  done
  wait
  COUNT=$((COUNT + BATCH))
  printf "\r\033[K  %d moths gathering... (%ds elapsed)" "${COUNT}" "${SECONDS}"
  sleep 1
done

echo ""
echo ""
echo "The trap has been set. Open Grafana to see what happened."
