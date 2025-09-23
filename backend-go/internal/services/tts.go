package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"math"
	"net/http"
	"net/url"

	"ai-role-playing-platform/backend-go/internal/config"
)

// SynthesizeSpeech creates a tiny sine-wave WAV and returns a data URL.
// This is a self-contained placeholder so the frontend can play audio without external services.
func SynthesizeSpeech(text string) (string, error) {
	// Try Baidu TTS if keys present; fallback to local tone
	if url, ok := tryBaiduTTS(text); ok {
		return url, nil
	}
	if text == "" {
		text = "hi"
	}
	// Duration heuristic: 0.6s + 40ms per char, capped [0.6s, 5s]
	base := 0.6
	extra := 0.04 * float64(len([]rune(text)))
	dur := base + extra
	if dur < 0.6 {
		dur = 0.6
	}
	if dur > 5.0 {
		dur = 5.0
	}
	sampleRate := 16000
	totalSamples := int(dur * float64(sampleRate))
	freq := 440.0
	// Build 16-bit mono PCM WAV
	pcm := make([]int16, totalSamples)
	for i := 0; i < totalSamples; i++ {
		// Simple attack/decay envelope to avoid clicks
		t := float64(i) / float64(sampleRate)
		env := envelope(t, dur)
		s := math.Sin(2*math.Pi*freq*t) * 0.4 * env
		if s > 1 {
			s = 1
		} else if s < -1 {
			s = -1
		}
		pcm[i] = int16(s * math.MaxInt16)
	}
	wav := buildWAVFromPCM16Mono(sampleRate, pcm)
	b64 := base64.StdEncoding.EncodeToString(wav)
	return "data:audio/wav;base64," + b64, nil
}

func envelope(t float64, dur float64) float64 {
	attack := 0.02
	release := 0.08
	if t < attack {
		return t / attack
	}
	if t > dur-release {
		return math.Max((dur-t)/release, 0)
	}
	return 1
}

// buildWAVFromPCM16Mono constructs a minimal RIFF/WAVE header with PCM16 mono data.
func buildWAVFromPCM16Mono(sampleRate int, pcm []int16) []byte {
	byteRate := sampleRate * 2 // mono * 16-bit
	blockAlign := 2
	dataSize := len(pcm) * 2
	riffSize := 36 + dataSize
	buf := &bytes.Buffer{}
	// RIFF header
	buf.WriteString("RIFF")
	writeLE32(buf, uint32(riffSize))
	buf.WriteString("WAVE")
	// fmt chunk
	buf.WriteString("fmt ")
	writeLE32(buf, 16) // PCM fmt chunk size
	writeLE16(buf, 1)  // audio format = PCM
	writeLE16(buf, 1)  // channels = 1
	writeLE32(buf, uint32(sampleRate))
	writeLE32(buf, uint32(byteRate))
	writeLE16(buf, uint16(blockAlign))
	writeLE16(buf, 16) // bits per sample
	// data chunk
	buf.WriteString("data")
	writeLE32(buf, uint32(dataSize))
	// samples
	for _, s := range pcm {
		writeLE16(buf, uint16(s))
	}
	return buf.Bytes()
}

func writeLE16(buf *bytes.Buffer, v uint16) {
	buf.WriteByte(byte(v))
	buf.WriteByte(byte(v >> 8))
}

func writeLE32(buf *bytes.Buffer, v uint32) {
	buf.WriteByte(byte(v))
	buf.WriteByte(byte(v >> 8))
	buf.WriteByte(byte(v >> 16))
	buf.WriteByte(byte(v >> 24))
}

func tryBaiduTTS(text string) (string, bool) {
	cfg := config.Load()
	if cfg.BaiduAPIKey == "" || cfg.BaiduSecretKey == "" {
		return "", false
	}
	token, err := getBaiduAccessToken()
	if err != nil {
		return "", false
	}
	v := url.Values{}
	v.Set("tex", text)
	v.Set("tok", token)
	v.Set("cuid", "ai-role-playing-platform")
	v.Set("ctp", "1")
	v.Set("lan", "zh")
	v.Set("per", "0")
	v.Set("spd", "5")
	v.Set("pit", "5")
	v.Set("vol", "7")
	resp, err := http.PostForm("https://tsn.baidu.com/text2audio", v)
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()
	// If success, content-type is audio/mp3 or audio/wav; error returns JSON
	ctype := resp.Header.Get("Content-Type")
	if ctype == "application/json" {
		var e map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&e)
		return "", false
	}
	b, _ := io.ReadAll(resp.Body)
	if len(b) == 0 {
		return "", false
	}
	mime := ctype
	if mime == "" {
		mime = "audio/mp3"
	}
	return "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(b), true
}
