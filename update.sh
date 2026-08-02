#!/bin/bash
set -e

echo "=== Updating Quotient OCR v2 Container ==="

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

echo "Rebuilding Docker image..."
docker build -t quotientbot/ocr:v2 .

echo "Restarting container..."
docker stop ocr_v2 >/dev/null 2>&1 || true
docker rm ocr_v2 >/dev/null 2>&1 || true

docker run -d \
    --name ocr_v2 \
    --restart unless-stopped \
    -p 8080:8080 \
    --env-file .env \
    quotientbot/ocr:v2

echo "✔ Update complete! Quotient OCR v2 is up and running."
