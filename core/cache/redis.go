package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
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
	return c.redis.Get(context.TODO(), key).Result()
}

func (c *RedisCache) Set(key string, val any, expiration time.Duration) error {
	if c.prefix != "" {
		key = c.prefix + ":" + key
	}

	s, err := cast.ToStringE(val)
	if err != nil {
		bs, err := json.Marshal(val)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		s = string(bs)
	}

	return c.redis.Set(context.TODO(), key, s, expiration).Err()
}

func (c *RedisCache) SetNX(key string, val any, expiration time.Duration) error {
	if c.prefix != "" {
		key = c.prefix + ":" + key
	}

	s, err := cast.ToStringE(val)
	if err != nil {
		bs, err := json.Marshal(val)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		s = string(bs)
	}

	ok, err := c.redis.SetNX(context.TODO(), key, s, expiration).Result()
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
	return c.redis.Del(context.TODO(), key).Err()
}

func (c *RedisCache) HGet(hk, field string) (any, error) {
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

func (c *RedisCache) Incr(key string) (int64, error) {
	if c.prefix != "" {
		key = c.prefix + ":" + key
	}
	return c.redis.Incr(context.TODO(), key).Result()
}

func (c *RedisCache) Decr(key string) (int64, error) {
	if c.prefix != "" {
		key = c.prefix + ":" + key
	}
	return c.redis.Decr(context.TODO(), key).Result()
}

func (c *RedisCache) Expire(key string, expiration time.Duration) error {
	if c.prefix != "" {
		key = c.prefix + ":" + key
	}
	return c.redis.Expire(context.TODO(), key, expiration).Err()
}

func (c *RedisCache) ExpireAt(key string, tm time.Time) error {
	if c.prefix != "" {
		key = c.prefix + ":" + key
	}
	return c.redis.ExpireAt(context.TODO(), key, tm).Err()
}

func (c *RedisCache) Exists(key string) bool {
	if c.prefix != "" {
		key = c.prefix + ":" + key
	}
	exists, err := c.redis.Exists(context.TODO(), key).Result()
	if err != nil {
		return false
	} else {
		return exists > 0
	}
}

func (c *RedisCache) MGet(keys ...string) ([]any, error) {
	if c.prefix != "" {
		for i, key := range keys {
			keys[i] = c.prefix + ":" + key
		}
	}
	return c.redis.MGet(context.TODO(), keys...).Result()
}

func (c *RedisCache) MSet(pairs map[string]any) error {
	if c.prefix != "" {
		for i, key := range pairs {
			pairs[i] = c.prefix + ":" + key.(string)
		}
	}
	return c.redis.MSet(context.TODO(), pairs).Err()
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
