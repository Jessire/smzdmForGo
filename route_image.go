package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const maxProxyImageBytes = 5 << 20

func ImageProxyHandler(w http.ResponseWriter, r *http.Request) {
	rawURL := strings.TrimSpace(r.URL.Query().Get("url"))
	if rawURL == "" {
		http.NotFound(w, r)
		return
	}

	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		http.NotFound(w, r)
		return
	}

	req, err := http.NewRequest(http.MethodGet, parsed.String(), nil)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0 Safari/537.36")
	req.Header.Set("Referer", "https://www.smzdm.com/")
	req.Header.Set("Accept", "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "image fetch failed", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		http.Error(w, "image fetch failed", http.StatusBadGateway)
		return
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = imageContentTypeFromPath(parsed.Path)
	}
	if !strings.HasPrefix(strings.ToLower(contentType), "image/") {
		http.Error(w, "not an image", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if resp.ContentLength > 0 && resp.ContentLength <= maxProxyImageBytes {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", resp.ContentLength))
	}
	_, _ = io.Copy(w, io.LimitReader(resp.Body, maxProxyImageBytes))
}

func imageContentTypeFromPath(path string) string {
	path = strings.ToLower(path)
	switch {
	case strings.HasSuffix(path, ".jpg"), strings.HasSuffix(path, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(path, ".png"):
		return "image/png"
	case strings.HasSuffix(path, ".webp"):
		return "image/webp"
	case strings.HasSuffix(path, ".gif"):
		return "image/gif"
	case strings.HasSuffix(path, ".svg"):
		return "image/svg+xml"
	default:
		return ""
	}
}
