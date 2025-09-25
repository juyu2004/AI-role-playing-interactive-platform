package handlers

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"ai-role-playing-platform/backend-go/internal/services"
)

type asrResponse struct {
	Text string `json:"text"`
}

// POST /api/asr
// 支持以下方式：
// 1) multipart/form-data: 字段 "file" 或 "audio"，可附带 "format" 与 "rate"
// 2) application/octet-stream: 直接音频字节流，可用查询参数 ?format=pcm&rate=16000
func HandleASR(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var (
		audio  []byte
		format = strings.ToLower(r.URL.Query().Get("format"))
		rate   = 16000
	)
	if v := r.URL.Query().Get("rate"); v != "" {
		// simple parse, ignore error -> keep default on failure
		if n, err := strconv.Atoi(v); err == nil {
			rate = n
		}
	}

	ctype := r.Header.Get("Content-Type")
	if strings.HasPrefix(ctype, "multipart/form-data") {
		if err := r.ParseMultipartForm(50 << 20); err != nil { // 50MB
			badRequest(w, "invalid multipart form")
			return
		}
		var file multipart.File
		var header *multipart.FileHeader
		if f, h, e := r.FormFile("file"); e == nil {
			file, header = f, h
		} else if f, h, e2 := r.FormFile("audio"); e2 == nil {
			file, header = f, h
		} else {
			badRequest(w, "missing file field")
			return
		}
		defer func() { _ = file.Close() }()
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, file); err != nil {
			badRequest(w, "failed to read file")
			return
		}
		audio = buf.Bytes()
		if format == "" {
			format = strings.TrimPrefix(strings.ToLower(filepath.Ext(header.Filename)), ".")
			if format == "wav" || format == "mp3" || format == "amr" || format == "m4a" || format == "pcm" {
				// ok
			} else {
				format = "pcm"
			}
		}
		if f := r.FormValue("format"); f != "" {
			format = strings.ToLower(f)
		}
		if v := r.FormValue("rate"); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				rate = n
			}
		}
	} else {
		// raw bytes
		b, err := io.ReadAll(r.Body)
		if err != nil {
			badRequest(w, "failed to read body")
			return
		}
		audio = b
		if format == "" {
			if strings.Contains(strings.ToLower(ctype), "wav") {
				format = "wav"
			} else if strings.Contains(strings.ToLower(ctype), "amr") {
				format = "amr"
			} else if strings.Contains(strings.ToLower(ctype), "mp3") {
				format = "mp3"
			} else if strings.Contains(strings.ToLower(ctype), "m4a") {
				format = "m4a"
			} else {
				format = "pcm"
			}
		}
	}

	text, err := services.TranscribeAudioWith(audio, format, rate)
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	ok(w, asrResponse{Text: text})
}
