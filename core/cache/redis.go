package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
)

// 默认 Redis 操作超时时间
const defaultRedisTimeout = 3 * time.Second

type RedisCache struct {
	redis  redis.UniversalClient
	prefix string
}

// ctx 创建带超时的 context，避免 Redis 操作无限挂起
func (c *RedisCache) ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), defaultRedisTimeout)
}

func (c *RedisCache) Type() string {
	return "redis"
}

func (c *RedisCache) IsRedis() bool {
	return true
}

func (c *RedisCache) RealKey(key string) string {
	if c.prefix != "" {
		key = c.prefix + ":" + key
	}
	return key
}

func (c *RedisCache) Get(key string) (string, error) {
	if c.prefix != "" {
		key = c.prefix + ":" + key
	}
	ctx, cancel := c.ctx()
	defer cancel()
	return c.redis.Get(ctx, key).Result()
}

func (c *RedisCache) Set(key string, val any, expiration time.Duration) error {
	if c.prefix != "" {
		key = c.prefix + ":" + key
	}

	s, err := cast.ToStringE(val)
	if err != nil {
		bs, err := json.Marshal(val)
		if err != nil {
			return err
		}
		s = string(bs)
	}

	ctx, cancel := c.ctx()
	defer cancel()
	return c.redis.Set(ctx, key, s, expiration).Err()
}

func (c *RedisCache) SetNX(key string, val any, expiration time.Duration) error {
	if c.prefix != "" {
		key = c.prefix + ":" + key
	}

	s, err := cast.ToStringE(val)
	if err != nil {
		bs, err := json.Marshal(val)
		if err != nil {
			return err
		}
		s = string(bs)
	}

	ctx, cancel := c.ctx()
	defer cancel()
	ok, err := c.redis.SetNX(ctx, key, s, expiration).Result()
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	return errors.New("set err")
}

func (c *RedisCache) Del(key string) error {
	if c.prefix != "" {
		key = c.prefix + ":" + key
	}
	ctx, cancel := c.ctx()
	defer cancel()
	return c.redis.Del(ctx, key).Err()
}

func (c *RedisCache) HGet(hk, field string) (any, error) {
	if c.prefix != "" {
		hk = c.prefix + ":" + hk
	}
	ctx, cancel := c.ctx()
	defer cancel()
	return c.redis.HGet(ctx, hk, field).Result()
}

func (c *RedisCache) HDel(hk, fields string) error {
	if c.prefix != "" {
		hk = c.prefix + ":" + hk
	}
	ctx, cancel := c.ctx()
	defer cancel()
	return c.redis.HDel(ctx, hk, fields).Err()
}

func (c *RedisCache) Incr(key string) (int64, error) {
	if c.prefix != "" {
		key = c.prefix + ":" + key
	}
	ctx, cancel := c.ctx()
	defer cancel()
	return c.redis.Incr(ctx, key).Result()
}

func (c *RedisCache) Decr(key string) (int64, error) {
	if c.prefix != "" {
		key = c.prefix + ":" + key
	}
	ctx, cancel := c.ctx()
	defer cancel()
	return c.redis.Decr(ctx, key).Result()
}

func (c *RedisCache) Expire(key string, expiration time.Duration) error {
	if c.prefix != "" {
		key = c.prefix + ":" + key
	}
	ctx, cancel := c.ctx()
	defer cancel()
	return c.redis.Expire(ctx, key, expiration).Err()
}

func (c *RedisCache) ExpireAt(key string, tm time.Time) error {
	if c.prefix != "" {
		key = c.prefix + ":" + key
	}
	ctx, cancel := c.ctx()
	defer cancel()
	return c.redis.ExpireAt(ctx, key, tm).Err()
}

func (c *RedisCache) Exists(key string) bool {
	if c.prefix != "" {
		key = c.prefix + ":" + key
	}
	ctx, cancel := c.ctx()
	defer cancel()
	exists, err := c.redis.Exists(ctx, key).Result()
	if err != nil {
		return false
	}
	return exists > 0
}

func (c *RedisCache) MGet(keys ...string) ([]any, error) {
	if c.prefix != "" {
		prefixed := make([]string, len(keys))
		for i, key := range keys {
			prefixed[i] = c.prefix + ":" + key
		}
		keys = prefixed
	}
	ctx, cancel := c.ctx()
	defer cancel()
	return c.redis.MGet(ctx, keys...).Result()
}

func (c *RedisCache) MSet(pairs map[string]any) error {
	ctx, cancel := c.ctx()
	defer cancel()
	if c.prefix == "" {
		return c.redis.MSet(ctx, pairs).Err()
	}

	withPrefix := make(map[string]any, len(pairs))
	for key, val := range pairs {
		withPrefix[c.prefix+":"+key] = val
	}
	return c.redis.MSet(ctx, withPrefix).Err()
}

func (c *RedisCache) GetClient() redis.UniversalClient {
	return c.redis
}

func GetRedisClient(c ICache) (redis.UniversalClient, error) {
	if c != nil && c.Type() == "redis" {
		return c.(*RedisCache).redis, nil
	}
	return nil, errors.ErrUnsupported
}
