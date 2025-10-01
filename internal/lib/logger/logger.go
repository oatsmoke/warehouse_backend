package logger

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	magenta = "\033[35m"
	blue    = "\033[34m"
	yellow  = "\033[33m"
	red     = "\033[31m"
	reset   = "\033[0m"
)

var (
	RequestBodyIsEmpty  = errors.New("request body is empty")
	NoRowsAffected      = errors.New("no rows affected")
	MethodNotAllowed    = errors.New("method not allowed")
	LaterDate           = errors.New("StartDate is later than EndDate")
	EmptyRequiredValues = errors.New("empty required values")
)

var mu sync.Mutex

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
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}

// Err is the error returned
func Err(err error, message string) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			fmt.Println(pgErr.Where)
			err = errors.New(pgErr.Detail)
		}
	}

	if message != "" {
		err = fmt.Errorf("%s: %v", message, err)
	}

	template(slog.LevelError, red, err.Error())

	return err
}

// ErrResponse is an error for the client
func ErrResponse(ctx *gin.Context, err error, status int) {
	//slog.Error(err.Error(), slog.String("status", strconv.Itoa(status)), slog.String("fn"))
	template(slog.LevelError, red, err.Error())
	ctx.AbortWithStatusJSON(status, map[string]string{"message": err.Error()})
}

// ErrInConsole is an error in the console
func ErrInConsole(err error) {
	template(slog.LevelError, red, err.Error())
}

func DebugInConsole(message string) {
	template(slog.LevelDebug, magenta, message)
}

func WarnInConsole(message string) {
	template(slog.LevelWarn, yellow, message)
}

// InfoInConsole is an info in the console
func InfoInConsole(message string) {
	template(slog.LevelInfo, blue, message)
}

func template(level slog.Level, color, message string) {
	mu.Lock()
	defer mu.Unlock()

	_, file, line, _ := runtime.Caller(2)
	fmt.Print(color)

	switch level {
	case slog.LevelInfo:
		slog.LogAttrs(nil,
			level,
			message,
		)
	default:
		slog.LogAttrs(nil,
			level,
			message,
			slog.String("file", file),
			slog.Int("line", line),
		)
	}

	fmt.Print(reset)
}
