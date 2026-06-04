package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Config struct {
	BaseURL      string
	APIKey       string
	DefaultModel string
	Models       []string
}

type Client struct {
	cfg  Config
	http *http.Client
}

func NewClient(cfg Config) *Client {
	return &Client{
		cfg: cfg,
		http: &http.Client{
			Timeout: 0,
		},
	}
}

func (c *Client) Models() []string {
	out := make([]string, len(c.cfg.Models))
	copy(out, c.cfg.Models)
	return out
}

func (c *Client) DefaultModel() string {
	return c.cfg.DefaultModel
}

func (c *Client) StreamResponses(w http.ResponseWriter, r *http.Request) error {
	if c.cfg.APIKey == "" {
		http.Error(w, "missing TERRAWEAVE_AI_API_KEY", http.StatusServiceUnavailable)
		return nil
	}

	var payload map[string]any
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 32<<20))
	decoder.UseNumber()
	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return nil
	}

	payload["stream"] = true
	if _, ok := payload["model"]; !ok || payload["model"] == "" {
		payload["model"] = c.cfg.DefaultModel
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal upstream payload: %w", err)
	}

	upstreamURL := strings.TrimRight(c.cfg.BaseURL, "/") + "/v1/responses"
	req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, upstreamURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create upstream request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("call upstream responses API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		w.Header().Set("Content-Type", firstHeader(resp.Header, "Content-Type", "application/json"))
		w.WriteHeader(resp.StatusCode)
		_, _ = io.Copy(w, io.LimitReader(resp.Body, 1<<20))
		return nil
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache, no-transform")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)

	flusher, ok := w.(http.Flusher)
	if !ok {
		return errors.New("response writer does not support streaming")
	}

	buf := make([]byte, 32*1024)
	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			if _, err := w.Write(buf[:n]); err != nil {
				return err
			}
			flusher.Flush()
		}
		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				return nil
			}
			return readErr
		}
	}
}

func firstHeader(header http.Header, key string, fallback string) string {
	value := header.Get(key)
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

type Health struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

func NewHealth() Health {
	return Health{Status: "ok", Time: time.Now().UTC().Format(time.RFC3339)}
}
