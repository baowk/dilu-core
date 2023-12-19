package locker

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/bsm/redislock"
)

// NewRedis 初始化locker
func NewRedis(c *redis.Client) *Redis {
	return &Redis{
		client: c,
	}
}

type Redis struct {
	client *redis.Client
	mutex  *redislock.Client
}

func (Redis) String() string {
	return "redis"
}

func (r *Redis) Lock(key string, ttl time.Duration, options *redislock.Options) (*redislock.Lock, error) {
	if r.client == nil {
		return nil, errors.New("redis client is nil")
	}
	if r.mutex == nil {
		r.mutex = redislock.New(r.client)
	}
	return r.mutex.Obtain(context.TODO(), key, ttl, options)
}
