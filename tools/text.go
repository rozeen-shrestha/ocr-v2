package tools

import (
	"bytes"
	"image"
	"image/jpeg"
	_ "image/png"

	"github.com/corona10/goimagehash"
	"github.com/nfnt/resize"
	"github.com/otiai10/gosseract/v2"
	_ "golang.org/x/image/webp"
)

// smartScale Image adjusts dimensions efficiently:
// - Keeps optimal 1080p images (800px - 1600px) at native resolution.
// - Scales down oversized images (>1600px) to max 1600px width.
// - Upscales small images (<800px) to ~1200px width using fast Bilinear interpolation.
func smartScaleImage(img image.Image) image.Image {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	maxDim := w
	if h > maxDim {
		maxDim = h
	}

	if maxDim > 1600 {
		scale := 1600.0 / float64(maxDim)
		newW := uint(float64(w) * scale)
		newH := uint(float64(h) * scale)
		return resize.Resize(newW, newH, img, resize.Bilinear)
	} else if maxDim < 800 {
		scale := 1200.0 / float64(maxDim)
		newW := uint(float64(w) * scale)
		newH := uint(float64(h) * scale)
		return resize.Resize(newW, newH, img, resize.Bilinear)
	}

	return img
}

// OCR extracts text from imageBytes along with perceptual and difference hashes.
// Robust and fast: supports JPEG, PNG, WebP, applies fast JPEG buffer re-encoding (~10ms),
// and uses optimal Tesseract LSTM settings.
func OCR(imageBytes []byte) (string, string, string, error) {
	// Decode image (supports JPEG, PNG, WebP)
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return "", "", "", err
	}

	// Compute perceptual hashes on native image
	phash, err := goimagehash.ExtPerceptionHash(img, 64, 64)
	if err != nil {
		return "", "", "", err
	}
	dhash, err := goimagehash.ExtDifferenceHash(img, 64, 64)
	if err != nil {
		return "", "", "", err
	}

	// Scale image if needed
	scaled := smartScaleImage(img)

	// Re-encode to high-quality JPEG buffer (takes ~10ms, universally supported by Leptonica)
	var buf bytes.Buffer
	if err = jpeg.Encode(&buf, scaled, &jpeg.Options{Quality: 92}); err != nil {
		return "", "", "", err
	}

	client := gosseract.NewClient()
	defer client.Close()

	// Tesseract LSTM configuration
	client.SetPageSegMode(gosseract.PSM_AUTO)
	client.SetVariable("user_defined_dpi", "300")
	client.SetVariable("preserve_interword_spaces", "1")

	if err = client.SetImageFromBytes(buf.Bytes()); err != nil {
		return "", "", "", err
	}

	text, err := client.Text()
	if err != nil {
		return "", "", "", err
	}

	return text, dhash.ToString(), phash.ToString(), nil
}
