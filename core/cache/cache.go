package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/baowk/dilu-core/config"
	"github.com/redis/go-redis/v9"
)

type ICache interface {
	Type() string
	Get(key string) (string, error)
	Set(key string, val any, expiration time.Duration) error
	Del(key string) error
	HGet(hk, field string) (string, error)
	HDel(hk, fields string) error
	Incr(key string) error
	Decr(key string) error
	Expire(key string, expiration time.Duration) error
}

func New(conf config.CacheCfg) ICache {
	if conf.GetType() == "redis" {
		rdb := redis.NewClient(&redis.Options{
			Addr:     conf.Addr,
			Password: conf.Password, // no password set
			DB:       conf.DB,       // use default DB
		})
		pong, err := rdb.Ping(context.Background()).Result()
		if err != nil {
			panic("redis connect ping failed, err:" + err.Error())
		} else {
			fmt.Println("redis connect ping response:", "pong", pong)
			r := RedisCache{
				redis:  rdb,
				prefix: conf.Prefix,
			}
			return &r
		}
	} else {
		return NewMemory()
	}
}
