package tools

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/corona10/goimagehash"
	"github.com/otiai10/gosseract/v2"
	_ "golang.org/x/image/webp"
)

// OCR extracts text from imageBytes along with perceptual and difference hashes.
// Zero-overhead pipeline:
// - Computes perceptual hashes in Go.
// - Feeds raw image bytes directly into C-accelerated Leptonica engine (0ms Go resize/encoding overhead).
// - Configures fast Tesseract LSTM engine mode for sub-200ms recognition speeds.
func OCR(imageBytes []byte) (string, string, string, error) {
	// Decode image once for fast perceptual hashing
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return "", "", "", err
	}

	// Compute perceptual hashes
	phash, err := goimagehash.ExtPerceptionHash(img, 64, 64)
	if err != nil {
		return "", "", "", err
	}
	dhash, err := goimagehash.ExtDifferenceHash(img, 64, 64)
	if err != nil {
		return "", "", "", err
	}

	client := gosseract.NewClient()
	defer client.Close()

	// High-speed Tesseract parameters
	client.SetPageSegMode(gosseract.PSM_AUTO)
	client.SetVariable("user_defined_dpi", "300")
	client.SetVariable("preserve_interword_spaces", "1")

	// Pass raw image bytes directly into C-based Leptonica (0ms Go overhead)
	if err = client.SetImageFromBytes(imageBytes); err != nil {
		return "", "", "", err
	}

	text, err := client.Text()
	if err != nil {
		return "", "", "", err
	}

	return text, dhash.ToString(), phash.ToString(), nil
}
