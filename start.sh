#!/bin/bash
set -e

echo "=================================================="
echo "   Starting Quotient OCR v2 Service               "
echo "=================================================="

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

# Parse PORT from .env (defaults to 8080 if not set)
PORT=$(grep -v '^#' .env 2>/dev/null | grep -i '^PORT=' | cut -d '=' -f2 | tr -d ' \r\n"')
if [ -z "$PORT" ]; then
    PORT=8080
fi

# Stop and remove existing container if running
docker stop ocr_v2 >/dev/null 2>&1 || true
docker rm ocr_v2 >/dev/null 2>&1 || true

# Launch container in background (-d) with automatic restart on reboot
echo "Launching container in background on port ${PORT}..."
docker run -d \
    --name ocr_v2 \
    --restart unless-stopped \
    -p ${PORT}:${PORT} \
    --env-file .env \
    quotientbot/ocr:v2

echo ""
echo "=================================================="
echo " 🎉 OCR v2 is running in background on port ${PORT}"
echo "    It will stay online even after closing your terminal."
echo ""
echo " To view live logs : docker logs -f ocr_v2"
echo " To stop service   : bash stop.sh"
echo "=================================================="
