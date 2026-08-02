package routers

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"

	"github.com/quotientbot/ocr_v2/tools"
)

type Images struct {
	ImageURLs []string `json:"urls"`
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "online",
		"service": "Quotient OCR API v2 (AMD EPYC Optimized)",
	})
}

func OCRHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	secretKey := os.Getenv("SECRET_KEY")
	authHeader := r.Header.Get("Authorization")
	if secretKey != "" && authHeader != "Bearer "+secretKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var imgs Images
	if err := json.NewDecoder(r.Body).Decode(&imgs); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result := make([]map[string]string, len(imgs.ImageURLs))
	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)

	for i, imageURL := range imgs.ImageURLs {
		wg.Add(1)
		go func(idx int, url string) {
			defer wg.Done()

			imageBytes, err := tools.GetBytesFromURL(url)
			if err != nil {
				return
			}

			text, dhash, phash, err := tools.OCR(imageBytes)
			if err != nil {
				return
			}

			var dhashStr, phashStr string
			if len(dhash) >= 2 {
				dhashStr = dhash[2:]
			} else {
				dhashStr = dhash
			}
			if len(phash) >= 2 {
				phashStr = phash[2:]
			} else {
				phashStr = phash
			}

			res := map[string]string{
				"url":   url,
				"dhash": dhashStr,
				"phash": phashStr,
				"text":  text,
			}

			mu.Lock()
			result[idx] = res
			mu.Unlock()
		}(i, imageURL)
	}

	wg.Wait()

	filtered := make([]map[string]string, 0, len(result))
	for _, r := range result {
		if r != nil {
			filtered = append(filtered, r)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(filtered)
}
