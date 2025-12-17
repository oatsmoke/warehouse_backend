package model

import (
	"github.com/oatsmoke/warehouse_backend/internal/lib/jwt_auth"
)

type AuthClaims struct {
	RegisteredClaims *jwt_auth.CustomClaims `json:"registered_claims"`
	Revoked          bool                   `json:"revoked"`
}
