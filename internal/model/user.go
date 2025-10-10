package model

import (
	"time"

	"github.com/oatsmoke/warehouse_backend/internal/lib/role"
)

type User struct {
	ID           int64      `json:"id,omitempty"`
	Username     string     `json:"username,omitempty"`
	PasswordHash string     `json:"-"`
	Email        string     `json:"email,omitempty"`
	Role         role.Role  `json:"role,omitempty"`
	Enabled      bool       `json:"enabled,omitempty"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	Employee     *Employee  `json:"employee,omitempty"`
}

func NewUser() *User {
	return &User{
		Employee: &Employee{
			Department: &Department{},
		},
	}
}
