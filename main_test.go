package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"gofr.dev/pkg/gofr/datasource/file"
	"gofr.dev/pkg/gofr/logging"
)

func TestServer(t *testing.T) {
	dir := setupTestDir(t)
	fs := file.New(logging.NewMockLogger(logging.ERROR))

	handler := &staticFileHandler{
		fs:               fs,
		staticFilePath:   dir,
		defaultExtension: ".html",
		next:             http.NotFoundHandler(),
	}

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	tests := []struct {
		path       string
		statusCode int
	}{
		{"/", http.StatusOK},
		{"/docs", http.StatusOK},
		{"/index", http.StatusOK},
		{"/index/", http.StatusOK},
		{filepath.Join(dir, "index.html"), http.StatusNotFound},
		{"/index.html", http.StatusOK},
		{"/nonexistent", http.StatusNotFound},
	}

	for _, test := range tests {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL+test.path, http.NoBody)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to perform request: %v", err)
		}

		if resp.StatusCode != test.statusCode {
			t.Errorf("Expected status code %v, got %v for path %v", test.statusCode, resp.StatusCode, test.path)
		}

		_ = resp.Body.Close()
	}
}
