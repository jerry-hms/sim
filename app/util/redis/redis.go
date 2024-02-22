package redis

import (
	"github.com/go-redis/redis/v8"
	"sim/app/global/variable"
	"sync"
)

var RedisIns *redis.Client
var RedisOnce sync.Once

func ConnRedis() *redis.Client {
	RedisOnce.Do(func() {
		RedisIns = redis.NewClient(&redis.Options{
			Addr:     variable.ConfigYml.GetString("redis.host"),
			Password: variable.ConfigYml.GetString("redis.password"),
			DB:       variable.ConfigYml.GetInt("redis.db"),
		})
	})
	return RedisIns
}
