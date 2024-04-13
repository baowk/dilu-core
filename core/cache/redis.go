package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	mode          int8 //1 单机 2 cluster
	redis         *redis.Client
	clusterClient *redis.ClusterClient
	prefix        string
}

func (c *RedisCache) Type() string {
	return "redis"
}

func (c *RedisCache) Get(key string) (string, error) {
	if c.prefix != "" {
		key = c.prefix + ":" + key
	}
	if c.mode == 1 {
		return c.redis.Get(context.TODO(), key).Result()
	} else {
		return c.clusterClient.Get(context.TODO(), key).Result()
	}
}

func (c *RedisCache) Set(key string, val any, expiration time.Duration) error {
	if c.prefix != "" {
		key = c.prefix + ":" + key
	}
	if c.mode == 1 {
		return c.redis.Set(context.TODO(), key, val, expiration).Err()
	} else {
		return c.clusterClient.Set(context.TODO(), key, val, expiration).Err()
	}
}

func (c *RedisCache) Del(key string) error {
	if c.prefix != "" {
		key = c.prefix + ":" + key
	}

	if c.mode == 1 {
		return c.redis.Del(context.TODO(), key).Err()
	} else {
		return c.clusterClient.Del(context.TODO(), key).Err()
	}
}

func (c *RedisCache) HGet(hk, field string) (string, error) {
	if c.prefix != "" {
		hk = c.prefix + ":" + hk
	}

	if c.mode == 1 {
		return c.redis.HGet(context.TODO(), hk, field).Result()
	} else {
		return c.clusterClient.HGet(context.TODO(), hk, field).Result()
	}
}

func (c *RedisCache) HDel(hk, fields string) error {
	if c.prefix != "" {
		hk = c.prefix + ":" + hk
	}

	if c.mode == 1 {
		return c.redis.HDel(context.TODO(), hk, fields).Err()
	} else {
		return c.clusterClient.HDel(context.TODO(), hk, fields).Err()
	}
}

func (c *RedisCache) Incr(key string) error {
	if c.prefix != "" {
		key = c.prefix + ":" + key
	}

	if c.mode == 1 {
		return c.redis.Incr(context.TODO(), key).Err()
	} else {
		return c.clusterClient.Incr(context.TODO(), key).Err()
	}
}

func (c *RedisCache) Decr(key string) error {
	if c.prefix != "" {
		key = c.prefix + ":" + key
	}

	if c.mode == 1 {
		return c.redis.Decr(context.TODO(), key).Err()
	} else {
		return c.clusterClient.Decr(context.TODO(), key).Err()
	}
}

func (c *RedisCache) Expire(key string, expiration time.Duration) error {
	if c.prefix != "" {
		key = c.prefix + ":" + key
	}

	if c.mode == 1 {
		return c.redis.Expire(context.TODO(), key, expiration).Err()
	} else {
		return c.clusterClient.Expire(context.TODO(), key, expiration).Err()
	}
}

func (c *RedisCache) GetClient() *redis.Client {
	return c.redis
}

func (c *RedisCache) GetClusterClient() *redis.ClusterClient {
	return c.clusterClient
}
