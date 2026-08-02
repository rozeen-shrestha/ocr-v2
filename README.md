# Quotient OCR API v2 (AMD EPYC Optimized)

High-performance, high-accuracy OCR microservice optimized for **AMD EPYC 9355P (4 vCPUs @ 3.545GHz)** servers and specialized for esports verification screenshots (Instagram profiles, YouTube channels, etc.).

---

## ⚡ 1-Command Setup & Run

On your VPS, run **one command** to automatically create configuration, build the Docker container, run it, and verify health:

```bash
cd ocr_v2
bash setup.sh
```

---

## 🔄 1-Command Update (After Code Changes)

Whenever you edit any code in `ocr_v2/`, rebuild and restart with a single command:

```bash
cd ocr_v2
bash update.sh
```

---

## Key Improvements over v1

1. **`tessdata_best` Neural Network Models**: Full float32 neural network parameters instead of `tessdata_fast`, significantly increasing text, handle (`@`), number, and username accuracy.
2. **Smart Resolution Scaling**: Dynamically adjusts image dimensions without forcing unnecessary 2x upscaling on 1080p images.
3. **High-Contrast Grayscale Preprocessing**: Enhances dark-mode UI text (Instagram/Discord/YouTube) to eliminate missed handles and "Following" status buttons.
4. **4-vCPU Parallel Execution**: Fully utilizes OpenMP multithreading on AMD EPYC 9355P processor (`OMP_THREAD_LIMIT=4`).
5. **Pooled Network Layer**: Reuses HTTP connections with Keep-Alive transport to download Discord attachment URLs with minimal latency.
