#!/bin/bash
set -e

echo "=================================================="
echo "   Quotient OCR v2 Setup & Deployment Script      "
echo "   Optimized for AMD EPYC 9355P (4 vCPUs)        "
echo "=================================================="

# Navigate to script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

# 0. Check and install Docker on the host if not present
if ! command -v docker >/dev/null 2>&1; then
    echo "[0/4] Docker is not installed on this host. Installing Docker..."
    if command -v apt-get >/dev/null 2>&1; then
        sudo apt-get update -qq
        sudo apt-get install -y -qq docker.io docker-buildx-plugin docker-compose-plugin
        sudo systemctl enable --now docker || true
    else
        curl -fsSL https://get.docker.com | sh
        sudo systemctl enable --now docker || true
    fi
    echo "✔ Docker installed successfully."
else
    echo "[0/4] Docker is already installed."
fi

# 1. Ensure .env file exists
if [ ! -f .env ]; then
    echo "[1/4] Creating default .env file from .example.env..."
    if [ -f .example.env ]; then
        cp .example.env .env
    else
        echo "SECRET_KEY=change_me_secret" > .env
        echo "PORT=8080" >> .env
    fi
    echo "✔ Created .env file."
else
    echo "[1/4] Found existing .env file."
fi

# 2. Build Docker container (installs all Tesseract, C++, Go dependencies inside container)
echo "[2/4] Building Docker container image (quotientbot/ocr:v2)..."
docker build -t quotientbot/ocr:v2 .
echo "✔ Docker image built successfully."

# 3. Restart Docker container
echo "[3/4] Launching container..."
docker stop ocr_v2 >/dev/null 2>&1 || true
docker rm ocr_v2 >/dev/null 2>&1 || true

docker run -d \
    --name ocr_v2 \
    --restart unless-stopped \
    -p 8080:8080 \
    --env-file .env \
    quotientbot/ocr:v2

# 4. Health Check
echo "[4/4] Verifying OCR service health..."
sleep 2

if command -v curl >/dev/null 2>&1; then
    HEALTH=$(curl -s http://localhost:8080/ || true)
    echo "Service response: $HEALTH"
fi

echo ""
echo "=================================================="
echo " 🎉 SUCCESS! OCR v2 is running on http://localhost:8080"
echo "=================================================="
