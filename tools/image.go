package tools

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

// Shared HTTP client with connection pooling & keep-alives for high throughput.
var httpClient = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 5 * time.Second,
	},
}

func GetBytesFromURL(imageURL string) ([]byte, error) {
	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		return nil, err
	}

	// Discord CDN URLs (cdn.discordapp.com / media.discordapp.net) require
	// ex, is, and hm query params (URL signature).
	q := parsedURL.Query()
	cleaned := url.Values{}
	for _, key := range []string{"ex", "is", "hm"} {
		if v := q.Get(key); v != "" {
			cleaned.Set(key, v)
		}
	}
	parsedURL.RawQuery = cleaned.Encode()

	cleanedURL := parsedURL.String()

	req, err := http.NewRequest(http.MethodGet, cleanedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Quotient-OCR-Bot/2.0")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status %d", resp.StatusCode)
	}

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return imageBytes, nil
}
