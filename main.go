package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/quotientbot/ocr_v2/routers"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Note: No .env file loaded, reading from environment variables")
	} else {
		log.Println("Loaded .env successfully")
	}

	// Tune OpenMP multi-threading to 2 threads per image to eliminate thread lock contention
	os.Setenv("OMP_THREAD_LIMIT", "2")
	os.Setenv("OMP_DYNAMIC", "FALSE")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", routers.IndexHandler)
	mux.HandleFunc("/ocr", routers.OCRHandler)

	s := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: mux,
	}

	log.Printf("Quotient OCR v2 API listening on 0.0.0.0:%s (EPYC 4-vCPU ultra-fast mode)\n", port)
	log.Fatal(s.ListenAndServe())
}
