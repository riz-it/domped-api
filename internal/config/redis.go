package config

import (
	"strconv"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(conf *Config) *redis.Client {
	database, _ := strconv.Atoi(conf.Redis.DB)
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Address,
		Username: conf.Redis.User,
		Password: conf.Redis.Pass,
		DB:       database,
	})

	return client
}
