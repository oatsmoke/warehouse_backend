package redis

import (
	"github.com/oatsmoke/warehouse_backend/internal/lib/env"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Conn *redis.Client
}

func Connect() *Redis {
	dsn := env.GetRedisDsn()
	return &Redis{Conn: redis.NewClient(&redis.Options{Addr: dsn})}
}

func (r *Redis) Disconnect() {
	if err := r.Conn.Close(); err != nil {
		logger.Warn(err.Error())
	}
}
