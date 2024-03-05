package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type Config struct {
	Address  string
	Password string
	DB       int
}

type cache struct {
	rdb *redis.Client
}

func NewRedisClient(cfg *Config) Cache {
	return &cache{
		rdb: redis.NewClient(&redis.Options{
			Addr: cfg.Address,
			DB:   cfg.DB,
		}),
	}
}

func (r *cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	val, err := r.rdb.Set(ctx, key, value, expiration).Result()
	if err != nil {
		return err
	}
	if val != "OK" {
		return fmt.Errorf("error setting value in cache: %s", val)
	}
	return nil
}

func (r *cache) Get(ctx context.Context, key string) (string, error) {
	val, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return val, nil
}

func (r *cache) Del(ctx context.Context, key string) error {
	result, err := r.rdb.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	if result == 0 {
		return fmt.Errorf("key %s not found in cache", key)
	}
	return nil
}

func (r *cache) HGet(ctx context.Context, hashKey string, key string) (string, error) {
	value, err := r.rdb.HGet(ctx, hashKey, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func (r *cache) HSet(ctx context.Context, hashKey string, key, value string, expiration time.Duration) error {
	err := r.rdb.HSet(ctx, hashKey, key, value).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *cache) HDel(ctx context.Context, hashKey string, key string) error {
	err := r.rdb.HDel(ctx, hashKey, key).Err()
	if err != nil {
		return err
	}
	return nil
}
