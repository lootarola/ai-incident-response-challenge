#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"

echo "Seeding the colony..."
echo ""

END=$((SECONDS + 90))
COUNT=0

while [ $SECONDS -lt $END ]; do
  for i in $(seq 1 10); do
    curl -s -o /dev/null --max-time 10 "${BASE_URL}/orders/report" &
  done
  wait
  COUNT=$((COUNT + 10))
  printf "\r\033[K  Colony at %d entries... (%ds elapsed)" "${COUNT}" "${SECONDS}"
done

echo ""
echo ""
echo "The infestation is underway. Open Grafana to see what happened."
