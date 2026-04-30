#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
PRODUCT_ID="${PRODUCT_ID:-prod-large}"

echo "Slugs on the move..."
echo ""

END=$((SECONDS + 60))
COUNT=0

while [ $SECONDS -lt $END ]; do
  for i in $(seq 1 5); do
    curl -s -o /dev/null \
      --max-time 10 \
      "${BASE_URL}/catalog/products/${PRODUCT_ID}" &
  done
  wait
  COUNT=$((COUNT + 5))
  printf "\r\033[K  %d slugs on the trail... (%ds elapsed)" "${COUNT}" "${SECONDS}"
done

echo ""
echo ""
echo "The trail is complete. Open Grafana to see what happened."
