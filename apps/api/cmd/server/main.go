package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Dr-Ai-0018/Terraweave/apps/api/internal/ai"
	"github.com/Dr-Ai-0018/Terraweave/apps/api/internal/cache"
	"github.com/Dr-Ai-0018/Terraweave/apps/api/internal/config"
	"github.com/Dr-Ai-0018/Terraweave/apps/api/internal/modalgpu"
	"github.com/Dr-Ai-0018/Terraweave/apps/api/internal/storage"
	"github.com/Dr-Ai-0018/Terraweave/apps/api/internal/store"
	"github.com/Dr-Ai-0018/Terraweave/apps/api/internal/system"
)

func main() {
	cfg := config.Load()
	client := ai.NewClient(ai.Config{
		BaseURL:      cfg.AIBaseURL,
		APIKey:       cfg.AIAPIKey,
		DefaultModel: cfg.AIDefaultModel,
		Models:       cfg.AIModels,
	})

	dbStore, err := store.Open(cfg.DatabaseURL)
	if err != nil {
		log.Printf("postgis disabled: %v", err)
	}
	defer func() {
		if err := dbStore.Close(); err != nil {
			log.Printf("close postgis: %v", err)
		}
	}()

	redisCache, err := cache.Open(cfg.RedisURL)
	if err != nil {
		log.Printf("redis disabled: %v", err)
	}
	defer func() {
		if err := redisCache.Close(); err != nil {
			log.Printf("close redis: %v", err)
		}
	}()

	objectStorage, err := storage.Open(cfg.S3)
	if err != nil {
		log.Printf("minio disabled: %v", err)
	}
	systemChecker := system.NewChecker(dbStore, redisCache, objectStorage)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", writeJSONHandler(func() any {
		return ai.NewHealth()
	}))
	mux.HandleFunc("GET /v1/system/status", func(w http.ResponseWriter, r *http.Request) {
		status := systemChecker.Check(r.Context())
		if !status.OK {
			writeJSON(w, http.StatusServiceUnavailable, status)
			return
		}
		writeJSON(w, http.StatusOK, status)
	})
	mux.HandleFunc("GET /v1/models", writeJSONHandler(func() any {
		return map[string]any{
			"object":        "list",
			"default_model": client.DefaultModel(),
			"data":          modelObjects(client.Models()),
		}
	}))
	mux.HandleFunc("GET /v1/modal/gpus", writeJSONHandler(func() any {
		return map[string]any{
			"default_gpu": cfg.ModalDefaultGPU,
			"data":        modalgpu.List(),
		}
	}))
	mux.HandleFunc("POST /v1/responses", func(w http.ResponseWriter, r *http.Request) {
		if err := client.StreamResponses(w, r); err != nil {
			log.Printf("stream responses: %v", err)
		}
	})

	log.Printf("terraweave api listening on %s", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, withCommonHeaders(mux)); err != nil {
		log.Fatal(err)
	}
}

func writeJSONHandler(fn func() any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, fn())
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("write json: %v", err)
	}
}

func modelObjects(models []string) []map[string]any {
	out := make([]map[string]any, 0, len(models))
	for _, model := range models {
		out = append(out, map[string]any{
			"id":       model,
			"object":   "model",
			"owned_by": "configured-provider",
		})
	}
	return out
}

func withCommonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
