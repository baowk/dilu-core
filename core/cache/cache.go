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
		arr := strings.Split(conf.Addr, ";")
		op := &redis.UniversalOptions{
			Addrs:    arr,
			Password: conf.Password, // no password set
		}
		if conf.DB > 0 {
			op.DB = conf.DB
		}
		if conf.MasterName != "" {
			op.MasterName = conf.MasterName
		}
		rdb := redis.NewUniversalClient(op)

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
