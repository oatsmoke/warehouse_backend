package service

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/smtp"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/oatsmoke/warehouse_backend/internal/lib/jwt_auth"
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
func (s *AuthService) AuthUser(ctx context.Context, login, password string) (*jwt_auth.Token, error) {
	user, err := s.AuthRepository.FindByPhone(ctx, &model.Employee{Phone: login})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, logger.Err(err, "wrong login or password")
		} else {
			return nil, logger.Err(err, "something wrong")
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, logger.Err(err, "wrong login or password")
		} else {
			return nil, logger.Err(err, "something wrong")
		}
	}
	t := &jwt_auth.Token{}
	claims, err := t.New(user.ID)
	if err != nil {
		return nil, logger.Err(err, "")
	}

	if err := s.AuthRepository.Set(ctx, claims, false); err != nil {
		return nil, logger.Err(err, "")
	}

	return t, nil
	//return user.ID, nil
}

func (s *AuthService) Check(ctx context.Context, token *jwt_auth.Token) (*jwt_auth.Token, error) {
	if claims, err := jwt_auth.CheckToken(token.Access); err != nil {
		logger.WarnInConsole(err.Error())

		claims, err := jwt_auth.CheckToken(token.Refresh)
		if err != nil {
			return nil, logger.Err(err, "")
		}

		revoked, err := s.AuthRepository.Get(ctx, claims.ID)
		if err != nil {
			return nil, logger.Err(err, "")
		}

		if revoked {
			return nil, logger.Err(errors.New("token is revoked"), "")
		}

		if err := s.AuthRepository.Set(ctx, claims, true); err != nil {
			return nil, logger.Err(err, "")
		}

		userId, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			return nil, logger.Err(err, "")
		}

		newClaims, err := token.New(userId)
		if err != nil {
			return nil, logger.Err(err, "")
		}

		if err := s.AuthRepository.Set(ctx, newClaims, false); err != nil {
			return nil, logger.Err(err, "")
		}
	} else {
		userId, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			return nil, logger.Err(err, "")
		}

		token.UserID = userId
	}

	return token, nil
}

// GenerateHash is to generate hash
//func (s *AuthService) GenerateHash(ctx context.Context, id int64) (string, error) {
//	str := generate.RandString(10)
//
//	if err := s.AuthRepository.SetHash(ctx, id, str); err != nil {
//		return "", logger.Err(err, "")
//	}
//
//	return str, nil
//}

// FindByHash is to find by hash
//func (s *AuthService) FindByHash(ctx context.Context, hash string) (int64, error) {
//	user, err := s.AuthRepository.FindByHash(ctx, &model.Employee{Hash: hash})
//	if err != nil {
//		return 0, logger.Err(err, "")
//	}
//
//	return user.ID, nil
//}

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
