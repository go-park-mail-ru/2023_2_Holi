package redis_connector

import (
	"context"
	"os"

	logs "2023_2_Holi/logger"

	"github.com/redis/go-redis/v9"
)

func RedisConnector() *redis.Client {
	redis := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	_, err := redis.Ping(context.Background()).Result()
	if err != nil {
		logs.LogFatal(logs.Logger, "redis_connector", "redisConnector", err, "Failed to connect to redis")
	}
	logs.Logger.Info("Connected to redis")
	logs.Logger.Debug("redis client :", redis)

	return redis
}
