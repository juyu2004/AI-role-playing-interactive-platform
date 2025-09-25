package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	execpkg "os/exec"
	"time"

	"ai-role-playing-platform/backend-go/internal/config"
)

// Baidu ASR (vop) implementation.
// API: https://ai.baidu.com/ai-doc/SPEECH/ek38lxj1x

var (
	cachedBaiduToken    string
	cachedBaiduTokenExp time.Time
)

func getBaiduAccessToken() (string, error) {
	if cachedBaiduToken != "" && time.Now().Before(cachedBaiduTokenExp) {
		return cachedBaiduToken, nil
	}
	cfg := config.Load()
	if cfg.BaiduAPIKey == "" || cfg.BaiduSecretKey == "" {
		return "", errors.New("missing baidu api keys")
	}
	values := url.Values{}
	values.Set("grant_type", "client_credentials")
	values.Set("client_id", cfg.BaiduAPIKey)
	values.Set("client_secret", cfg.BaiduSecretKey)
	resp, err := http.PostForm("https://aip.baidubce.com/oauth/2.0/token", values)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var out struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if out.AccessToken == "" {
		return "", errors.New("empty access token")
	}
	cachedBaiduToken = out.AccessToken
	cachedBaiduTokenExp = time.Now().Add(time.Duration(out.ExpiresIn-60) * time.Second)
	return cachedBaiduToken, nil
}

// TranscribeAudio sends PCM(16k) audio to Baidu VOP and returns text.
func TranscribeAudio(b []byte) (string, error) {
	return TranscribeAudioWith(b, "pcm", 16000)
}

// TranscribeAudioWith allows specifying format and sample rate explicitly.
func TranscribeAudioWith(b []byte, format string, rate int) (string, error) {
	if len(bytes.TrimSpace(b)) == 0 {
		return "", nil
	}
	// Prefer local SDK CLI if available
	if text, ok := tryLocalSDK(b); ok {
		return text, nil
	}
	// Prefer Baidu JSON base64
	if text, err := baiduASRJSON(b, format, rate); err == nil && text != "" {
		return text, nil
	} else if err != nil {
		return "", err
	}
	// Fallback: multipart
	if text, err := baiduASRMultipart(b, format, rate); err == nil && text != "" {
		return text, nil
	} else if err != nil {
		return "", err
	}
	return "", errors.New("asr failed")
}

func baiduASRJSON(b []byte, format string, rate int) (string, error) {
	token, err := getBaiduAccessToken()
	if err != nil {
		return "", err
	}
	payload := map[string]any{
		"format":  format, // "pcm" per Baidu doc
		"rate":    rate,   // 16000
		"dev_pid": 1537,   // 普通话
		"channel": 1,
		"lan":     "zh",
		"cuid":    "ai-role-playing-platform",
		"len":     len(b),
		"speech":  base64.StdEncoding.EncodeToString(b),
	}
	body, _ := json.Marshal(payload)
	// token 放在 URL 查询串
	url := "http://vop.baidu.com/server_api?dev_pid=1537&cuid=ai-role-playing-platform&token=" + url.QueryEscape(token)
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var v struct {
		ErrNo  int      `json:"err_no"`
		ErrMsg string   `json:"err_msg"`
		Result []string `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return "", err
	}
	if v.ErrNo != 0 || len(v.Result) == 0 {
		return "", fmt.Errorf("baidu asr error: err_no=%d err_msg=%s", v.ErrNo, v.ErrMsg)
	}
	return v.Result[0], nil
}

func baiduASRMultipart(b []byte, format string, rate int) (string, error) {
	token, err := getBaiduAccessToken()
	if err != nil {
		return "", err
	}
	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)
	fw, _ := mw.CreateFormFile("audio", "audio."+format)
	_, _ = io.Copy(fw, bytes.NewReader(b))
	_ = mw.WriteField("dev_pid", "1537")
	_ = mw.WriteField("format", format)
	_ = mw.WriteField("rate", fmt.Sprintf("%d", rate))
	_ = mw.Close()
	req, _ := http.NewRequest("POST", "http://vop.baidu.com/server_api", body)
	q := req.URL.Query()
	q.Set("cuid", "ai-role-playing-platform")
	q.Set("token", token)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", mw.FormDataContentType())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var v struct {
		ErrNo  int      `json:"err_no"`
		ErrMsg string   `json:"err_msg"`
		Result []string `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return "", err
	}
	if v.ErrNo != 0 || len(v.Result) == 0 {
		return "", fmt.Errorf("baidu asr error: err_no=%d err_msg=%s", v.ErrNo, v.ErrMsg)
	}
	return v.Result[0], nil
}

// tryLocalSDK attempts to invoke a local C++ SDK CLI to do ASR when present.
// It writes stdin bytes to the CLI and expects plain text on stdout.
func tryLocalSDK(b []byte) (string, bool) {
	candidates := []string{
		"backend-go/tools/aip-cpp-sdk-4.16.7/bin/asr_cli.exe",
		"backend-go/tools/aip-cpp-sdk-4.16.7/bin/asr_cli",
		"backend-go/tools/aip-cpp-sdk-4.16.6/bin/asr_sdk.exe",
		"backend-go/tools/aip-cpp-sdk-4.16.6/bin/asr_sdk",
	}
	cfg := config.Load()
	for _, path := range candidates {
		cmd := execpkg.Command(path, "--stdin", "--rate", "16000", "--format", "pcm", "--dev_pid", "1537", "--ak", cfg.BaiduAPIKey, "--sk", cfg.BaiduSecretKey, "--appid", cfg.BaiduAppID)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			continue
		}
		go func() { _, _ = stdin.Write(b); _ = stdin.Close() }()
		out, err := cmd.CombinedOutput()
		if err != nil {
			continue
		}
		text := string(bytes.TrimSpace(out))
		if text != "" {
			return text, true
		}
	}
	return "", false
}
