package tools

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"

	"github.com/corona10/goimagehash"
	"github.com/nfnt/resize"
	"github.com/otiai10/gosseract/v2"
)

// smartResizeAndPreprocess adapts resolution and enhances contrast for Tesseract.
// - Images with max dimension > 1800px are scaled down to ~1600px max dimension.
// - Images with min dimension < 800px are scaled up to ~1200px min dimension.
// - Standard 1080p screenshots (800px - 1800px) stay at native resolution.
// - Fast Bilinear interpolation is used instead of heavy Lanczos3.
// - Image is converted to high-contrast grayscale to optimize dark-mode UI text accuracy.
func smartResizeAndPreprocess(img image.Image) image.Image {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	var processed image.Image = img

	maxDim := w
	if h > maxDim {
		maxDim = h
	}
	minDim := w
	if h < minDim {
		minDim = h
	}

	// Adapt resolution
	if maxDim > 1800 {
		scale := 1600.0 / float64(maxDim)
		newW := uint(float64(w) * scale)
		newH := uint(float64(h) * scale)
		processed = resize.Resize(newW, newH, img, resize.Bilinear)
	} else if minDim < 800 {
		scale := 1200.0 / float64(minDim)
		newW := uint(float64(w) * scale)
		newH := uint(float64(h) * scale)
		processed = resize.Resize(newW, newH, img, resize.Bilinear)
	}

	// High-contrast Grayscale conversion for dark mode / light text accuracy
	pBounds := processed.Bounds()
	grayImg := image.NewGray(pBounds)

	for y := pBounds.Min.Y; y < pBounds.Max.Y; y++ {
		for x := pBounds.Min.X; x < pBounds.Max.X; x++ {
			c := processed.At(x, y)
			r, g, b, _ := c.RGBA()
			// Standard luminance calculation (uint8)
			lum := uint8((299*r + 587*g + 114*b) / 1000 >> 8)

			// Light contrast stretch to sharpen dark-mode Instagram text
			var val uint8
			if lum > 200 {
				val = 255
			} else if lum < 45 {
				val = 0
			} else {
				val = uint8(float64(lum-45) * (255.0 / 155.0))
			}

			grayImg.Set(x, y, color.Gray{Y: val})
		}
	}

	return grayImg
}

// OCR extracts text from imageBytes along with perceptual and difference hashes.
// Optimized for AMD EPYC 9355P (4 vCPUs) and tessdata_best neural network models.
func OCR(imageBytes []byte) (string, string, string, error) {
	img, format, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return "", "", "", err
	}

	// Compute perceptual hashes on original full-resolution image
	phash, err := goimagehash.ExtPerceptionHash(img, 64, 64)
	if err != nil {
		return "", "", "", err
	}
	dhash, err := goimagehash.ExtDifferenceHash(img, 64, 64)
	if err != nil {
		return "", "", "", err
	}

	// Preprocess image for maximum Tesseract LSTM accuracy and fast throughput
	processed := smartResizeAndPreprocess(img)

	var buf bytes.Buffer
	if format == "jpeg" {
		if err = jpeg.Encode(&buf, processed, &jpeg.Options{Quality: 90}); err != nil {
			return "", "", "", err
		}
	} else {
		if err = png.Encode(&buf, processed); err != nil {
			return "", "", "", err
		}
	}

	client := gosseract.NewClient()
	defer client.Close()

	// Configuration for high accuracy LSTM recognition
	client.SetPageSegMode(gosseract.PSM_AUTO)
	client.SetVariable("tessedit_ocr_engine_mode", "1") // LSTM engine mode
	client.SetVariable("user_defined_dpi", "300")      // Fixes resolution warning
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
