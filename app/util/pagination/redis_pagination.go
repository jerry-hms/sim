package pagination

import (
	"github.com/go-redis/redis/v8"
	redis_pagination "github.com/jerry-hms/redis-pagination"
	"sim/app/global/variable"
)

func GetHashDb() *redis_pagination.HashDb {
	return redis_pagination.NewHashDB(&redis.Options{
		Addr:     variable.ConfigYml.GetString("redis.host"),
		Password: variable.ConfigYml.GetString("redis.password"),
		DB:       variable.ConfigYml.GetInt("redis.db"),
	})
}
