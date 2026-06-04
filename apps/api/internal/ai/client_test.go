package ai

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStreamResponsesForcesStreamingAndProxiesSSE(t *testing.T) {
	var upstreamPayload map[string]any

	client := NewClient(Config{
		BaseURL:      "https://provider.example",
		APIKey:       "test-key",
		DefaultModel: "gpt-5.4",
		Models:       []string{"gpt-5.4"},
	})
	client.http = &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/v1/responses" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Fatalf("unexpected authorization header: %q", got)
		}
		if err := json.NewDecoder(r.Body).Decode(&upstreamPayload); err != nil {
			t.Fatalf("decode upstream payload: %v", err)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"text/event-stream"}},
			Body:       io.NopCloser(bytes.NewBufferString("event: response.output_text.delta\ndata: {\"delta\":\"ok\"}\n\n")),
		}, nil
	})}

	req := httptest.NewRequest(http.MethodPost, "/v1/responses", strings.NewReader(`{"stream":false,"input":"hello"}`))
	rec := httptest.NewRecorder()

	if err := client.StreamResponses(rec, req); err != nil {
		t.Fatalf("stream responses: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); got != "text/event-stream" {
		t.Fatalf("unexpected content type: %q", got)
	}
	if got := rec.Body.String(); !strings.Contains(got, "response.output_text.delta") {
		t.Fatalf("SSE body was not proxied: %q", got)
	}
	if upstreamPayload["stream"] != true {
		t.Fatalf("stream was not forced to true: %#v", upstreamPayload["stream"])
	}
	if upstreamPayload["model"] != "gpt-5.4" {
		t.Fatalf("default model was not applied: %#v", upstreamPayload["model"])
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}
