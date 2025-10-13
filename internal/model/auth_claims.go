package model

import "github.com/golang-jwt/jwt/v5"

type AuthClaims struct {
	RegisteredClaims *jwt.RegisteredClaims `json:"registered_claims"`
	Revoked          bool                  `json:"revoked"`
}
