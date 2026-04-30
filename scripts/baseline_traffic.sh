#!/bin/sh
set -u
BASE_URL="${BASE_URL:-http://localhost:8080}"

until curl -sf "$BASE_URL/healthz" >/dev/null 2>&1; do
  sleep 2
done

ORDER_PAYLOAD='{"customer_id":"cust-baseline","items":[{"product_id":"prod-000","category":"electronics","quantity":1,"unit_price":9.99}]}'

while true; do
  curl -sf -m 5 "$BASE_URL/healthz" >/dev/null 2>&1 || true
  sleep 1
  curl -sf -m 5 -G "$BASE_URL/catalog/products" --data-urlencode "search=widget" >/dev/null 2>&1 || true
  sleep 1
  curl -sf -m 5 "$BASE_URL/catalog/products/prod-000" >/dev/null 2>&1 || true
  sleep 1
  curl -sf -m 5 -X POST "$BASE_URL/orders" -H "Content-Type: application/json" -d "$ORDER_PAYLOAD" >/dev/null 2>&1 || true
  sleep 1
  curl -sf -m 5 "$BASE_URL/orders/order-0000" >/dev/null 2>&1 || true
  sleep 1
  curl -sf -m 5 "$BASE_URL/orders/report" >/dev/null 2>&1 || true
  sleep 1
  curl -sf -m 5 -X POST "$BASE_URL/orders/order-0000/notify" -H "Content-Type: application/json" -d '{"event":"confirmed"}' >/dev/null 2>&1 || true
  sleep 1
done
