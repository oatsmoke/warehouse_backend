package env

import (
	"fmt"
	"os"

	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
)

const (
	LogLevel    = "LOG_LEVEL"
	HttpPort    = "HTTP_PORT"
	AccessTtl   = "ACCESS_TTL"
	RefreshTtl  = "REFRESH_TTL"
	SigningKey  = "SIGNING_KEY"
	PostgresDsn = "POSTGRES_DSN"
	RedisDsn    = "REDIS_DSN"
	ClientUrl   = "CLIENT_URL"
)

func GetLogLevel() string {
	return get(LogLevel)
}

func GetHttpPort() string {
	return get(HttpPort)
}

func GetAccessTtl() string {
	return get(AccessTtl)
}

func GetRefreshTtl() string {
	return get(RefreshTtl)
}

func GetSigningKey() string {
	return get(SigningKey)
}

func GetPostgresDsn() string {
	return get(PostgresDsn)
}

func GetRedisDsn() string {
	return get(RedisDsn)
}

func GetClientUrl() string {
	return get(ClientUrl)
}

func get(key string) string {
	val, ok := os.LookupEnv(key)
	if ok {
		return val
	} else {
		switch key {
		case LogLevel:
			//	message(LogLevel)
			return "debug"
		case HttpPort:
			message(HttpPort)
			return ":8080"
		case AccessTtl:
			message(AccessTtl)
			return "300"
		case RefreshTtl:
			message(RefreshTtl)
			return "3600"
		case SigningKey:
			message(SigningKey)
			return "secret"
		case PostgresDsn:
			message(PostgresDsn)
			return "postgres://root:password@localhost:5432/postgres?sslmode=disable"
		case RedisDsn:
			message(RedisDsn)
			return "localhost:6379"
		case ClientUrl:
			message(ClientUrl)
			return "http://localhost"
		default:
			logger.InfoInConsole(fmt.Sprintf("%s not found", key))
			return ""
		}
	}
}

func message(key string) {
	logger.InfoInConsole(fmt.Sprintf("%s not set, set default value", key))
}
