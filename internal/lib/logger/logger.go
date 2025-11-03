package logger

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
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

//var (
//	NoRowsAffected          = errors.New("no rows affected")
//	WrongUsernameOrPassword = errors.New("wrong username or password")
//	InvalidToken            = errors.New("invalid token")
//	InvalidClaims           = errors.New("invalid claims")
//	TokenIsRevoked          = errors.New("token is revoked")
//)

var mu sync.Mutex

func Init(logLevel string) {
	var level slog.Level

	switch logLevel {
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	case "debug":
		level = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}

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

func ErrResponse(ctx *gin.Context, err error, status int) {
	//slog.Error(err.Error(), slog.String("status", strconv.Itoa(status)), slog.String("fn"))
	template(slog.LevelError, red, err.Error())
	ctx.AbortWithStatusJSON(status, map[string]string{"message": err.Error()})
}

//func ErrInConsole(err error) {
//	template(slog.LevelError, red, err.Error())
//}
//
//func DebugInConsole(message string) {
//	template(slog.LevelDebug, magenta, message)
//}
//
//func WarnInConsole(message string) {
//	template(slog.LevelWarn, yellow, message)
//}

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

var (
	ErrAlreadyExists           = errors.New("already exists")
	ErrNoRowsAffected          = errors.New("no rows affected")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrWrongUsernameOrPassword = errors.New("wrong username or password")
	ErrInvalidToken            = errors.New("invalid token")
	ErrInvalidClaims           = errors.New("invalid claims")
	ErrTokenHasBeenRevoked     = errors.New("token has been revoked")
	ErrWrongPassword           = errors.New("wrong password")
	ErrUserIdNotFound          = errors.New("user id not found")
	ErrInvalidRole             = errors.New("invalid role")
)

const (
	MsgAuthenticationFailed        = "authentication failed"
	MsgAuthorizationDenied         = "authorization denied"
	MsgFailedToInsert              = "failed to insert"
	MsgFailedToSelect              = "failed to select"
	MsgFailedToUpdate              = "failed to update"
	MsgFailedToDelete              = "failed to delete"
	MsgFailedToRestore             = "failed to restore"
	MsgFailedToScan                = "failed to scan"
	MsgFailedToIterateOverRows     = "failed to iterate over rows"
	MsgFailedToGet                 = "failed to get"
	MsgFailedToSet                 = "failed to set"
	MsgFailedToMarshal             = "failed to marshal"
	MsgFailedToUnmarshal           = "failed to unmarshal"
	MsgFailedToConvert             = "failed to convert"
	MsgFailedToParse               = "failed to parse"
	MsgFailedToSigned              = "failed to signed"
	MsgFailedToValidate            = "failed to validate"
	MsgFailedToGenerateHash        = "failed to generate hash"
	MsgFailedToSetSenderAddress    = "failed to set sender address"
	MsgFailedToAddRecipientAddress = "failed to add recipient address"
	MsgFailedToSetBodyText         = "failed to set body text"
	MsgFailedToSetBodyHTML         = "failed to set body html"
	MsgFailedToSetMailClient       = "failed to set mail client"
	MsgFailedToSendMail            = "failed to send mail"
)

func Info(msg string) {
	inConsole("info", msg)
}

func Warn(msg string) {
	inConsole("warn", msg)
}

func Error(msg string, err error) error {
	if fn, file, line, ok := runtime.Caller(1); ok {
		errMsg := err.Error()
		var pgErr *pgconn.PgError

		pgErr, err = pgErrParse(err)
		if pgErr != nil {
			errMsg = pgErr.Message
		}

		inConsole("error",
			fmt.Sprintf("%s: %s, %s:%s:%d",
				msg,
				errMsg,
				filepath.Base(runtime.FuncForPC(fn).Name()),
				filepath.Base(file),
				line,
			),
		)
	}

	return err
}

func ResponseErr(ctx *gin.Context, msg string, err error, status int) {
	inConsole("error",
		fmt.Sprintf("%d: %s: %v",
			status,
			msg,
			err,
		),
	)

	ctx.AbortWithStatusJSON(status, gin.H{"message": msg})
}

func inConsole(outType, msg string) {
	mu.Lock()
	defer mu.Unlock()

	switch outType {
	case "info":
		//fmt.Print(blue)
		slog.Info(msg)
	case "warn":
		//fmt.Print(yellow)
		slog.Warn(msg)
	case "error":
		//fmt.Print(red)
		slog.Error(msg)
	case "debug":
		//fmt.Print(magenta)
		slog.Debug(msg)
	}

	//fmt.Print(reset)
}

func pgErrParse(err error) (*pgconn.PgError, error) {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			err = ErrAlreadyExists
		}
	}

	return pgErr, err
}
