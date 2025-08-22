package service

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/smtp"

	"github.com/jackc/pgx/v5"
	"github.com/oatsmoke/warehouse_backend/internal/lib/generate"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/oatsmoke/warehouse_backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	AuthRepository repository.Auth
}

func NewAuthService(authRepository repository.Auth) *AuthService {
	return &AuthService{
		AuthRepository: authRepository,
	}
}

// AuthUser is user authentication for login
func (s *AuthService) AuthUser(ctx context.Context, login, password string) (int64, error) {
	user, err := s.AuthRepository.FindByPhone(ctx, &model.Employee{Phone: login})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, logger.Err(err, "wrong login or password")
		} else {
			return 0, logger.Err(err, "something wrong")
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, logger.Err(err, "wrong login or password")
		} else {
			return 0, logger.Err(err, "something wrong")
		}
	}

	return user.ID, nil
}

// GenerateHash is to generate hash
func (s *AuthService) GenerateHash(ctx context.Context, id int64) (string, error) {
	str := generate.RandString(10)

	if err := s.AuthRepository.SetHash(ctx, id, str); err != nil {
		return "", logger.Err(err, "")
	}

	return str, nil
}

// FindByHash is to find by hash
func (s *AuthService) FindByHash(ctx context.Context, hash string) (int64, error) {
	user, err := s.AuthRepository.FindByHash(ctx, &model.Employee{Hash: hash})
	if err != nil {
		return 0, logger.Err(err, "")
	}

	return user.ID, nil
}

func sendMail(recipient, phone, password string) error {
	authEmail := "oatsmoke@yandex.ru"
	authPassword := "kbbaojsmxlnboajk"
	host := "smtp.yandex.ru"
	port := "465"
	subject := "Authorization data"
	body := fmt.Sprintf("Login: %s\nPassword: %s", phone, password)
	sendString := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", authEmail, recipient, subject, body)
	msg := []byte(sendString)
	auth := smtp.PlainAuth("", authEmail, authPassword, host)
	conf := &tls.Config{ServerName: host}
	conn, err := tls.Dial("tcp", host+":"+port, conf)
	if err != nil {
		return err
	}
	cl, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	if err := cl.Auth(auth); err != nil {
		return err
	}
	if err := cl.Mail(authEmail); err != nil {
		return err
	}
	if err := cl.Rcpt(recipient); err != nil {
		return err
	}
	w, err := cl.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(msg); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	if err := cl.Quit(); err != nil {
		return err
	}
	return nil
}
