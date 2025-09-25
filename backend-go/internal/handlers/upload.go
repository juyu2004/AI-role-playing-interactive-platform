package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// POST /api/upload/audio  multipart/form-data: field "file"
func HandleUploadAudio(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseMultipartForm(100 << 20); err != nil { // 100MB
		badRequest(w, "invalid multipart form")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		badRequest(w, "missing file")
		return
	}
	defer func() { _ = file.Close() }()

	// ensure uploads dir
	_ = os.MkdirAll("uploads/audio", 0o755)
	filename := time.Now().Format("20060102_150405") + "_" + header.Filename
	dstPath := filepath.Join("uploads", "audio", filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		serverError(w, "failed to save file")
		return
	}
	defer func() { _ = dst.Close() }()
	if _, err := io.Copy(dst, file); err != nil {
		serverError(w, "failed to write file")
		return
	}
	// Static URL (assuming server serves cwd). In production, map /static/ to ./uploads.
	url := "/" + dstPath
	ok(w, map[string]string{"url": url})
}
