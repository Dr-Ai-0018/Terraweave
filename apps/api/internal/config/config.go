package config

import (
	"os"
	"strings"
)

type Config struct {
	Addr            string
	AIBaseURL       string
	AIAPIKey        string
	AIDefaultModel  string
	AIModels        []string
	ModalDefaultGPU string
}

func Load() Config {
	models := splitCSV(getenv("TERRAWEAVE_AI_MODELS", "gpt-5.4"))
	defaultModel := getenv("TERRAWEAVE_AI_DEFAULT_MODEL", "gpt-5.4")
	if len(models) == 0 {
		models = []string{defaultModel}
	}

	return Config{
		Addr:            getenv("TERRAWEAVE_API_ADDR", ":8080"),
		AIBaseURL:       strings.TrimRight(getenv("TERRAWEAVE_AI_BASE_URL", "https://api.openai.com"), "/"),
		AIAPIKey:        os.Getenv("TERRAWEAVE_AI_API_KEY"),
		AIDefaultModel:  defaultModel,
		AIModels:        models,
		ModalDefaultGPU: getenv("TERRAWEAVE_MODAL_DEFAULT_GPU", "L4"),
	}
}

func getenv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	seen := map[string]struct{}{}
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	return out
}
