package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/jwt_auth"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepository repository.Auth
	userRepository repository.User
}

func NewAuthService(authRepository repository.Auth, userRepository repository.User) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		userRepository: userRepository,
	}
}

func (s *AuthService) AuthUser(ctx context.Context, login *dto.UserLogin) (*jwt_auth.Token, error) {
	user, err := s.userRepository.GetByUsername(ctx, login.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, logger.WrongUsernameOrPassword
		} else {
			return nil, err
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(login.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, logger.WrongUsernameOrPassword
		} else {
			return nil, err
		}
	}

	token := &jwt_auth.Token{}
	claims, err := token.New(user.ID)
	if err != nil {
		return nil, err
	}

	if err := s.authRepository.Set(ctx, claims, false); err != nil {
		return nil, err
	}

	logger.InfoInConsole(fmt.Sprintf("User %s logged in", user.Username))
	return token, nil
}

func (s *AuthService) CheckToken(ctx context.Context, token *jwt_auth.Token) (*jwt_auth.Token, error) {
	if claims, err := jwt_auth.CheckToken(token.Access); err != nil {
		logger.WarnInConsole(err.Error())

		claims, err := jwt_auth.CheckToken(token.Refresh)
		if err != nil {
			return nil, err
		}

		revoked, err := s.authRepository.Get(ctx, claims.ID)
		if err != nil {
			return nil, err
		}

		if revoked {
			return nil, logger.TokenIsRevoked
		}

		if err := s.authRepository.Set(ctx, claims, true); err != nil {
			return nil, err
		}

		userId, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			return nil, err
		}

		newClaims, err := token.New(userId)
		if err != nil {
			return nil, err
		}

		if err := s.authRepository.Set(ctx, newClaims, false); err != nil {
			return nil, err
		}
	} else {
		userId, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			return nil, err
		}

		token.UserID = userId
	}

	return token, nil
}
