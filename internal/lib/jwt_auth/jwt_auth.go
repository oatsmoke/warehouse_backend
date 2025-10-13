package jwt_auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oatsmoke/warehouse_backend/internal/lib/env"
	"github.com/oatsmoke/warehouse_backend/internal/lib/generate"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
)

type Token struct {
	UserID  int64
	Access  string
	Refresh string
}

func (t *Token) New(userId int64) (*jwt.RegisteredClaims, error) {
	strUserId := strconv.FormatInt(userId, 10)

	if err := t.setAccess(strUserId); err != nil {
		return nil, err
	}

	claims, err := t.setRefresh(strUserId)
	if err != nil {
		return nil, err
	}

	t.UserID = userId

	return claims, nil
}

func (t *Token) setAccess(userId string) error {
	accessTTL, err := strconv.Atoi(env.GetAccessTtl())
	if err != nil {
		return err
	}

	claims := &jwt.RegisteredClaims{
		Subject:   userId,
		Audience:  []string{"user-agent"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(accessTTL) * time.Second)),
		ID:        generate.RandString(10),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(env.GetSigningKey()))
	if err != nil {
		return err
	}

	t.Access = token
	return nil
}

func (t *Token) setRefresh(userId string) (*jwt.RegisteredClaims, error) {
	refreshTTL, err := strconv.Atoi(env.GetRefreshTtl())
	if err != nil {
		return nil, err
	}

	claims := &jwt.RegisteredClaims{
		Subject:   userId,
		Audience:  []string{"user-agent"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(refreshTTL) * time.Second)),
		ID:        generate.RandString(10),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(env.GetSigningKey()))
	if err != nil {
		return nil, err
	}

	t.Refresh = token
	return claims, nil
}

func CheckToken(token string) (*jwt.RegisteredClaims, error) {
	t, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, checkMethod)
	if err != nil {
		return nil, err
	}

	if !t.Valid {
		return nil, logger.InvalidToken
	}

	claims, ok := t.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, logger.InvalidClaims
	}

	return claims, nil
}

func checkMethod(t *jwt.Token) (interface{}, error) {
	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
	}

	return []byte(env.GetSigningKey()), nil
}
