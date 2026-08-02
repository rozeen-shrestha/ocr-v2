#!/bin/bash
set -e

echo "=== Updating Quotient OCR v2 Container ==="

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

# Parse PORT from .env (defaults to 8080 if not set)
PORT=$(grep -v '^#' .env 2>/dev/null | grep -i '^PORT=' | cut -d '=' -f2 | tr -d ' \r\n"')
if [ -z "$PORT" ]; then
    PORT=8080
fi

echo "Rebuilding Docker image..."
docker build -t quotientbot/ocr:v2 .

echo "Restarting container on port ${PORT}..."
docker stop ocr_v2 >/dev/null 2>&1 || true
docker rm ocr_v2 >/dev/null 2>&1 || true

docker run -d \
    --name ocr_v2 \
    --restart unless-stopped \
    -p ${PORT}:${PORT} \
    --env-file .env \
    quotientbot/ocr:v2

echo "✔ Update complete! Quotient OCR v2 is up and running on http://localhost:${PORT}"
