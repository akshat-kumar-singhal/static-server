package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gofr.dev/pkg/gofr/datasource/file"
)

type staticFileHandler struct {
	fs               file.FileSystem
	staticFilePath   string
	spaMode          bool
	defaultExtension string
	next             http.Handler
}

func (h *staticFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "/.well-known/") {
		h.next.ServeHTTP(w, r)

		return
	}

	filePath, hasExtension := h.resolveFilePath(r.URL.Path)

	if _, err := h.fs.Stat(filePath); os.IsNotExist(err) {
		if h.spaMode && !hasExtension {
			http.ServeFile(w, r, filepath.Join(h.staticFilePath, indexHTML))
			return
		}

		w.WriteHeader(http.StatusNotFound)

		http.ServeFile(w, r, filepath.Join(h.staticFilePath, "404.html"))

		return
	}

	http.ServeFile(w, r, filePath)
}

func (h *staticFileHandler) resolveFilePath(urlPath string) (string, bool) {
	filePath := filepath.Join(h.staticFilePath, urlPath)

	hasExtension := filepath.Ext(filePath) != ""

	if urlPath == rootPath {
		filePath += indexHTML
	} else if !hasExtension {
		if _, err := h.fs.Stat(filePath + h.defaultExtension); err == nil {
			filePath += h.defaultExtension
		} else if info, err := h.fs.Stat(filePath); err == nil && info.IsDir() {
			filePath += indexHTML
		}
	}

	return filePath, hasExtension
}
