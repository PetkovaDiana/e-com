package cache

import (
	"context"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
	HGet(ctx context.Context, hashKey string, key string) (string, error)
	HSet(ctx context.Context, hashKey string, key, value string, expiration time.Duration) error
	HDel(ctx context.Context, hashKey string, key string) error
}
