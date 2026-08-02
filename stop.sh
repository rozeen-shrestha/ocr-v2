#!/bin/bash
echo "Stopping Quotient OCR v2 Service..."
docker stop ocr_v2 >/dev/null 2>&1 || true
docker rm ocr_v2 >/dev/null 2>&1 || true
echo "✔ OCR v2 service stopped."
