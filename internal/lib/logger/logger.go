package logger

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"log/slog"
	"os"
	"strconv"
)

// Init is logger initialization
func Init(logLevel string) {
	var level slog.Level

	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
		//AddSource: true,
		Level: level,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}

// Err is the error returned
func Err(err error, message, fn string) error {
	slog.Error(err.Error(), slog.String("fn", fn))

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			fmt.Println(pgErr.Where)
			return errors.New(pgErr.Detail)
		}
	}

	if message != "" {
		return errors.New(message)
	}

	return err
}

// ErrResponse is an error for the client
func ErrResponse(ctx *gin.Context, err error, status int, fn string) {
	slog.Error(err.Error(), slog.String("status", strconv.Itoa(status)), slog.String("fn", fn))
	ctx.AbortWithStatusJSON(status, map[string]string{"message": err.Error()})
}

// ErrInConsole is an error in the console
func ErrInConsole(err error, fn string) {
	slog.Error(err.Error(), slog.String("fn", fn))
}

// InfoInConsole is an info in the console
func InfoInConsole(message string, fn string) {
	slog.Info(message, slog.String("fn", fn))
}
