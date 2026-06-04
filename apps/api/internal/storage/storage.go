package storage

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Dr-Ai-0018/Terraweave/apps/api/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage struct {
	client  *minio.Client
	buckets []string
}

type BucketStatus struct {
	Name  string `json:"name"`
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

type Status struct {
	OK      bool           `json:"ok"`
	Buckets []BucketStatus `json:"buckets"`
	Error   string         `json:"error,omitempty"`
}

func Open(cfg config.S3Config) (*Storage, error) {
	if cfg.Endpoint == "" {
		return nil, fmt.Errorf("missing TERRAWEAVE_S3_ENDPOINT")
	}
	if cfg.AccessKeyID == "" {
		return nil, fmt.Errorf("missing TERRAWEAVE_S3_ACCESS_KEY_ID")
	}
	if cfg.SecretAccessKey == "" {
		return nil, fmt.Errorf("missing TERRAWEAVE_S3_SECRET_ACCESS_KEY")
	}

	endpoint, secure, err := parseEndpoint(cfg.Endpoint)
	if err != nil {
		return nil, err
	}
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: secure,
	})
	if err != nil {
		return nil, err
	}

	return &Storage{client: client, buckets: cfg.Buckets()}, nil
}

func (s *Storage) Check(ctx context.Context) Status {
	if s == nil || s.client == nil {
		return Status{OK: false, Error: "storage not configured"}
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	status := Status{OK: true, Buckets: make([]BucketStatus, 0, len(s.buckets))}
	for _, bucket := range s.buckets {
		bucketStatus := BucketStatus{Name: bucket}
		if strings.TrimSpace(bucket) == "" {
			bucketStatus.Error = "empty bucket name"
			status.OK = false
			status.Buckets = append(status.Buckets, bucketStatus)
			continue
		}
		exists, err := s.client.BucketExists(ctx, bucket)
		if err != nil {
			bucketStatus.Error = err.Error()
			status.OK = false
		} else if !exists {
			bucketStatus.Error = "bucket does not exist"
			status.OK = false
		} else {
			bucketStatus.OK = true
		}
		status.Buckets = append(status.Buckets, bucketStatus)
	}
	return status
}

func parseEndpoint(raw string) (string, bool, error) {
	if !strings.Contains(raw, "://") {
		return raw, false, nil
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", false, err
	}
	if parsed.Host == "" {
		return "", false, fmt.Errorf("invalid S3 endpoint: %s", raw)
	}
	return parsed.Host, parsed.Scheme == "https", nil
}
