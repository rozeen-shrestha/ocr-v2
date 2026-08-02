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

	// Strict single-thread limit + PASSIVE wait policy to stop CPU busy-spinning/hogging
	os.Setenv("OMP_NUM_THREADS", "1")
	os.Setenv("OMP_THREAD_LIMIT", "1")
	os.Setenv("OMP_WAIT_POLICY", "PASSIVE")
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

	log.Printf("Quotient OCR v2 API listening on 0.0.0.0:%s (Low-CPU mode)\n", port)
	log.Fatal(s.ListenAndServe())
}
