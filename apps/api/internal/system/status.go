package system

import (
	"context"
	"time"

	"github.com/Dr-Ai-0018/Terraweave/apps/api/internal/cache"
	"github.com/Dr-Ai-0018/Terraweave/apps/api/internal/storage"
	"github.com/Dr-Ai-0018/Terraweave/apps/api/internal/store"
)

type Checker struct {
	store   *store.Store
	cache   *cache.Cache
	storage *storage.Storage
}

type Status struct {
	OK        bool           `json:"ok"`
	CheckedAt string         `json:"checked_at"`
	API       Component      `json:"api"`
	PostGIS   store.Status   `json:"postgis"`
	Redis     cache.Status   `json:"redis"`
	MinIO     storage.Status `json:"minio"`
}

type Component struct {
	OK bool `json:"ok"`
}

func NewChecker(store *store.Store, cache *cache.Cache, storage *storage.Storage) *Checker {
	return &Checker{store: store, cache: cache, storage: storage}
}

func (c *Checker) Check(ctx context.Context) Status {
	status := Status{
		OK:        true,
		CheckedAt: time.Now().UTC().Format(time.RFC3339),
		API:       Component{OK: true},
		PostGIS:   c.store.Check(ctx),
		Redis:     c.cache.Check(ctx),
		MinIO:     c.storage.Check(ctx),
	}
	status.OK = status.API.OK && status.PostGIS.OK && status.Redis.OK && status.MinIO.OK
	return status
}
