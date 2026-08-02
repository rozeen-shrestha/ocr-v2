package tools

import (
	"bytes"
	"image"
	"image/jpeg"
	_ "image/png"
	"log"
	"sync"

	"github.com/corona10/goimagehash"
	"github.com/nfnt/resize"
	"github.com/otiai10/gosseract/v2"
	_ "golang.org/x/image/webp"
)

// clientPool holds pre-warmed, reusable Tesseract clients.
// Re-using an already-initialized client eliminates the 200-500ms startup
// overhead paid every time gosseract.NewClient() is called from scratch.
var (
	poolSize    = 4 // matches the number of vCPUs
	clientPool  chan *gosseract.Client
	poolOnce    sync.Once
)

func initPool() {
	clientPool = make(chan *gosseract.Client, poolSize)
	for i := 0; i < poolSize; i++ {
		c := gosseract.NewClient()
		applyTesseractConfig(c)
		clientPool <- c
	}
	log.Printf("[OCR] Warmed up %d Tesseract client(s) in pool", poolSize)
}

func applyTesseractConfig(c *gosseract.Client) {
	// PSM 11: sparse text — instantly pulls all text fragments without building layout trees.
	// Perfect for ssverify which only looks for a handful of keywords.
	c.SetPageSegMode(gosseract.PSM_SPARSE_TEXT)

	// Only recognise alphanumeric + common punctuation — massively reduces neural net search space.
	c.SetVariable("tessedit_char_whitelist", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 @._-/|")

	// Disable heavy post-processing not needed for keyword matching
	c.SetVariable("tessedit_enable_dict_correction", "0")
	c.SetVariable("tessedit_enable_doc_dict", "0")
	c.SetVariable("user_defined_dpi", "300")
}

// WarmPool explicitly pre-warms the Tesseract client pool at server startup.
// Call this from main() so all clients are ready before the first request arrives.
func WarmPool() {
	poolOnce.Do(initPool)
}

func acquireClient() *gosseract.Client {
	poolOnce.Do(initPool)
	return <-clientPool
}

func releaseClient(c *gosseract.Client) {
	clientPool <- c
}

// fastScaleImage downscales large phone screenshots to max 900px on the long edge
// using fast Bilinear interpolation — cuts pixel count 4-6x for ~10ms resize.
// For ssverify keywords this resolution is more than enough.
func fastScaleImage(img image.Image) image.Image {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	maxDim := w
	if h > maxDim {
		maxDim = h
	}
	if maxDim > 900 {
		scale := 900.0 / float64(maxDim)
		newW := uint(float64(w) * scale)
		newH := uint(float64(h) * scale)
		return resize.Resize(newW, newH, img, resize.Bilinear)
	}
	return img
}

// OCR extracts text from imageBytes with perceptual hashes.
// Uses a pre-warmed client pool to eliminate init overhead, a character whitelist
// to shrink the recognition search space, and PSM 11 sparse-text mode for
// instant keyword extraction — targeting <200ms total per image.
func OCR(imageBytes []byte) (string, string, string, error) {
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return "", "", "", err
	}

	// Compute perceptual hashes (fast, pure-Go)
	phash, err := goimagehash.ExtPerceptionHash(img, 64, 64)
	if err != nil {
		return "", "", "", err
	}
	dhash, err := goimagehash.ExtDifferenceHash(img, 64, 64)
	if err != nil {
		return "", "", "", err
	}

	// Downscale large screenshot for faster processing
	scaled := fastScaleImage(img)

	// Re-encode to JPEG for universal Leptonica compatibility (~5ms)
	var buf bytes.Buffer
	if encErr := jpeg.Encode(&buf, scaled, &jpeg.Options{Quality: 80}); encErr != nil {
		buf.Reset()
		buf.Write(imageBytes)
	}

	// Borrow a pre-warmed Tesseract client from the pool
	client := acquireClient()
	defer releaseClient(client)

	if err = client.SetImageFromBytes(buf.Bytes()); err != nil {
		return "", "", "", err
	}

	text, err := client.Text()
	if err != nil {
		return "", "", "", err
	}

	return text, dhash.ToString(), phash.ToString(), nil
}
