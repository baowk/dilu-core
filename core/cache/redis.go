package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	redis *redis.Client
}

func (c *RedisCache) Type() string {
	return "redis"
}

func (c *RedisCache) Get(key string) (string, error) {
	return c.redis.Get(context.TODO(), key).Result()
}

func (c *RedisCache) Set(key string, val any, expiration time.Duration) error {
	return c.redis.Set(context.TODO(), key, val, expiration).Err()
}

func (c *RedisCache) Del(key string) error {
	return c.redis.Del(context.TODO(), key).Err()
}

func (c *RedisCache) HGet(hk, field string) (string, error) {
	return c.redis.HGet(context.TODO(), hk, field).Result()
}

func (c *RedisCache) HDel(hk, fields string) error {
	return c.redis.HDel(context.TODO(), hk, fields).Err()
}

func (c *RedisCache) Incr(key string) error {
	return c.redis.Incr(context.TODO(), key).Err()
}

func (c *RedisCache) Decr(key string) error {
	return c.redis.Decr(context.TODO(), key).Err()
}

func (c *RedisCache) Expire(key string, expiration time.Duration) error {
	return c.redis.Expire(context.TODO(), key, expiration).Err()
}

func (c *RedisCache) GetClient() *redis.Client {
	return c.redis
}
