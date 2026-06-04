package config

import (
	"os"
	"strings"
)

type Config struct {
	Addr                string
	DatabaseURL         string
	DatabaseAdminURL    string
	DatabaseReadonlyURL string
	RedisURL            string
	S3                  S3Config
	AIBaseURL           string
	AIAPIKey            string
	AIDefaultModel      string
	AIModels            []string
	ModalDefaultGPU     string
}

type S3Config struct {
	Endpoint           string
	AccessKeyID        string
	SecretAccessKey    string
	BucketRaw          string
	BucketIntermediate string
	BucketOutputs      string
	BucketModels       string
}

func (cfg S3Config) Buckets() []string {
	return []string{
		cfg.BucketRaw,
		cfg.BucketIntermediate,
		cfg.BucketOutputs,
		cfg.BucketModels,
	}
}

func Load() Config {
	models := splitCSV(getenv("TERRAWEAVE_AI_MODELS", "gpt-5.4"))
	defaultModel := getenv("TERRAWEAVE_AI_DEFAULT_MODEL", "gpt-5.4")
	if len(models) == 0 {
		models = []string{defaultModel}
	}

	return Config{
		Addr:                getenv("TERRAWEAVE_API_ADDR", ":8080"),
		DatabaseURL:         os.Getenv("TERRAWEAVE_DATABASE_URL"),
		DatabaseAdminURL:    os.Getenv("TERRAWEAVE_DATABASE_ADMIN_URL"),
		DatabaseReadonlyURL: os.Getenv("TERRAWEAVE_DATABASE_READONLY_URL"),
		RedisURL:            os.Getenv("TERRAWEAVE_REDIS_URL"),
		S3: S3Config{
			Endpoint:           strings.TrimRight(os.Getenv("TERRAWEAVE_S3_ENDPOINT"), "/"),
			AccessKeyID:        os.Getenv("TERRAWEAVE_S3_ACCESS_KEY_ID"),
			SecretAccessKey:    os.Getenv("TERRAWEAVE_S3_SECRET_ACCESS_KEY"),
			BucketRaw:          getenv("TERRAWEAVE_S3_BUCKET_RAW", "terraweave-raw"),
			BucketIntermediate: getenv("TERRAWEAVE_S3_BUCKET_INTERMEDIATE", "terraweave-intermediate"),
			BucketOutputs:      getenv("TERRAWEAVE_S3_BUCKET_OUTPUTS", "terraweave-outputs"),
			BucketModels:       getenv("TERRAWEAVE_S3_BUCKET_MODELS", "terraweave-models"),
		},
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
