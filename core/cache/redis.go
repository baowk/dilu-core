package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	redis  redis.UniversalClient
	prefix string
	//mode  int8 //1 单机 2 cluster
	//clusterClient *redis.ClusterClient

}

func (c *RedisCache) Type() string {
	return "redis"
}

func (c *RedisCache) Get(key string) (string, error) {
	if c.prefix != "" {
		key = c.prefix + ":" + key
	}
	return c.redis.Get(context.TODO(), key).Result()
}

func (c *RedisCache) Set(key string, val any, expiration time.Duration) error {
	if c.prefix != "" {
		key = c.prefix + ":" + key
	}
	return c.redis.Set(context.TODO(), key, val, expiration).Err()
}

func (c *RedisCache) Del(key string) error {
	if c.prefix != "" {
		key = c.prefix + ":" + key
	}
	return c.redis.Del(context.TODO(), key).Err()
}

func (c *RedisCache) HGet(hk, field string) (string, error) {
	if c.prefix != "" {
		hk = c.prefix + ":" + hk
	}
	return c.redis.HGet(context.TODO(), hk, field).Result()
}

func (c *RedisCache) HDel(hk, fields string) error {
	if c.prefix != "" {
		hk = c.prefix + ":" + hk
	}
	return c.redis.HDel(context.TODO(), hk, fields).Err()
}

func (c *RedisCache) Incr(key string) error {
	if c.prefix != "" {
		key = c.prefix + ":" + key
	}
	return c.redis.Incr(context.TODO(), key).Err()
}

func (c *RedisCache) Decr(key string) error {
	if c.prefix != "" {
		key = c.prefix + ":" + key
	}
	return c.redis.Decr(context.TODO(), key).Err()
}

func (c *RedisCache) Expire(key string, expiration time.Duration) error {
	if c.prefix != "" {
		key = c.prefix + ":" + key
	}
	return c.redis.Expire(context.TODO(), key, expiration).Err()
}

func (c *RedisCache) GetClient() redis.UniversalClient {
	return c.redis
}
