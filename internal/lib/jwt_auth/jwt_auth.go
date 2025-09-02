package jwt_auth

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oatsmoke/warehouse_backend/internal/lib/env"
	"github.com/oatsmoke/warehouse_backend/internal/lib/generate"
)

//const (
//	AccessTTL  = time.Minute * 1
//	RefreshTTL = time.Hour * 24 * 30
//)

// var sign = RandString(10)
//var sign = "1234567890"

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
	tokenTTL, err := strconv.Atoi(env.GetAccessTtl())
	if err != nil {
		return err
	}

	claims := &jwt.RegisteredClaims{
		Subject:   userId,
		Audience:  []string{"user-agent"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(tokenTTL) * time.Second)),
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
		return nil, errors.New("invalid token")
	}

	claims, ok := t.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}

func checkMethod(t *jwt.Token) (interface{}, error) {
	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
	}

	return []byte(env.GetSigningKey()), nil
}

//// GenerateToken is token generation
//func GenerateToken(id int64) (string, error) {
//	num, err := strconv.Atoi(env.GetTokenTtl())
//	if err != nil {
//		return "", err
//	}
//
//	tokenTTL := time.Duration(num) * time.Minute
//
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
//		"sub": id,
//		"exp": jwt.NewNumericDate(time.Unix(time.Now().Add(tokenTTL).Unix(), 0)),
//		"iat": jwt.NewNumericDate(time.Unix(time.Now().Unix(), 0)),
//	})
//
//	signedString, err := token.SignedString([]byte(env.GetSigningKey()))
//	if err != nil {
//		return "", err
//	}
//
//	return signedString, nil
//}
//
//// ParseToken is token parsing
//func ParseToken(accessToken string) (int64, error) {
//	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, errors.New("invalid signing")
//		}
//		return []byte(env.GetSigningKey()), nil
//	})
//
//	if err != nil {
//		return 0, err
//	}
//
//	claims, ok := token.Claims.(jwt.MapClaims)
//	if !ok || !token.Valid {
//		return 0, errors.New("invalid type token")
//	}
//
//	userId := int64(claims["sub"].(float64))
//
//	return userId, nil
//}
