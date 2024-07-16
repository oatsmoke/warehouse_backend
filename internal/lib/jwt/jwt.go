package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strconv"
	"time"
)

// GenerateToken is token generation
func GenerateToken(id int64) (string, error) {
	num, err := strconv.Atoi(os.Getenv("tokenTTL"))
	if err != nil {
		return "", err
	}

	tokenTTL := time.Duration(num) * time.Minute

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
		"exp": jwt.NewNumericDate(time.Unix(time.Now().Add(tokenTTL).Unix(), 0)),
		"iat": jwt.NewNumericDate(time.Unix(time.Now().Unix(), 0)),
	})

	signedString, err := token.SignedString([]byte(os.Getenv("signingKey")))
	if err != nil {
		return "", err
	}

	return signedString, nil
}

// ParseToken is token parsing
func ParseToken(accessToken string) (int64, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing")
		}
		return []byte(os.Getenv("signingKey")), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid type token")
	}

	userId := int64(claims["sub"].(float64))

	return userId, nil
}
