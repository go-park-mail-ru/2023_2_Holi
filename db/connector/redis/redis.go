package redis

import (
	"context"
	"os"

	logs "2023_2_Holi/logger"

	"github.com/redis/go-redis/v9"
)

func Connect() *redis.Client {
	r := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	_, err := r.Ping(context.Background()).Result()
	if err != nil {
		logs.LogError(logs.Logger, "redis", "Connect", err, err.Error())
	}
	logs.Logger.Info("Connected to redis")
	logs.Logger.Debug("Redis client :", r)

	return r
}
