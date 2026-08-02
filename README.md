# Quotient OCR API v2 (AMD EPYC Optimized)

High-performance, high-accuracy OCR microservice optimized for **AMD EPYC 9355P (4 vCPUs @ 3.545GHz)** servers and specialized for esports verification screenshots (Instagram profiles, YouTube channels, etc.).

---

## ⚡ Quick Start Scripts

### 1-Click Initial Setup & Run
Run once on a fresh VPS to install Docker, build the container, dynamically configure the port from `.env`, and start the service:
```bash
cd ocr_v2
bash setup.sh
```

### Persistent Start (Background / Non-closing)
Starts the service persistently in background mode (`-d` and `--restart unless-stopped`). Closing your SSH terminal window will **not** stop the service:
```bash
cd ocr_v2
bash start.sh
```

### Stop Service
```bash
bash stop.sh
```

### 1-Click Update (After Code Changes)
Rebuilds and restarts the service on the configured `.env` port:
```bash
bash update.sh
```

---

## Key Improvements over v1

1. **`tessdata_best` Neural Network Models**: Full float32 neural network parameters instead of `tessdata_fast`, significantly increasing text, handle (`@`), number, and username accuracy.
2. **Smart Resolution Scaling**: Dynamically adjusts image dimensions without forcing unnecessary 2x upscaling on 1080p images.
3. **High-Contrast Grayscale Preprocessing**: Enhances dark-mode UI text (Instagram/Discord/YouTube) to eliminate missed handles and "Following" status buttons.
4. **4-vCPU Parallel Execution**: Fully utilizes OpenMP multithreading on AMD EPYC 9355P processor (`OMP_THREAD_LIMIT=4`).
5. **Pooled Network Layer**: Reuses HTTP connections with Keep-Alive transport to download Discord attachment URLs with minimal latency.
