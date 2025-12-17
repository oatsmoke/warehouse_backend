package redis

import (
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/redis/go-redis/v9"
)

func Connect(dsn string) *redis.Client {
	return redis.NewClient(&redis.Options{Addr: dsn})
}

func Disconnect(client *redis.Client) {
	if err := client.Close(); err != nil {
		logger.Warn(err.Error())
	}
}
