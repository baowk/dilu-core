package cache

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/baowk/dilu-core/config"
	"github.com/redis/go-redis/v9"
)

type ICache interface {
	Type() string
	IsRedis() bool
	Get(key string) (string, error)
	Set(key string, val any, expiration time.Duration) error
	SetNX(key string, val any, expiration time.Duration) error
	Del(key string) error
	HGet(hk, field string) (any, error)
	HDel(hk, fields string) error
	Incr(key string) (int64, error)
	Decr(key string) (int64, error)
	Expire(key string, expiration time.Duration) error
	ExpireAt(key string, t time.Time) error
	Exists(key string) bool
	MGet(keys ...string) ([]any, error)
	MSet(pairs map[string]any) error
	RealKey(key string) string
}

func New(conf config.CacheCfg) (ICache, error) {
	if conf.GetType() == "redis" {
		arr := strings.Split(conf.Addr, ";")
		op := &redis.UniversalOptions{
			Addrs:    arr,
			Password: conf.Password,
		}
		if conf.DB > 0 {
			op.DB = conf.DB
		}
		if conf.MasterName != "" {
			op.MasterName = conf.MasterName
		}
		rdb := redis.NewUniversalClient(op)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, err := rdb.Ping(ctx).Result(); err != nil {
			return nil, fmt.Errorf("redis connect ping failed: %w", err)
		}
		return &RedisCache{redis: rdb, prefix: conf.Prefix}, nil
	}
	return NewMemory(), nil
}
