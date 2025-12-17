package dto

import (
	"github.com/oatsmoke/warehouse_backend/internal/lib/role"
)

type UserRequest struct {
	Username   string    `json:"username,omitempty" binding:"required"`
	Email      string    `json:"email,omitempty" binding:"required"`
	Role       role.Role `json:"role,omitempty" binding:"required"`
	EmployeeID int64     `json:"employee_id,omitempty"`
}

type UserResponse struct {
	ID           int64  `json:"id,omitempty"`
	Username     string `json:"username,omitempty"`
	Email        string `json:"email,omitempty"`
	Role         string `json:"role,omitempty"`
	Enabled      bool   `json:"enabled,omitempty"`
	LastLoginAt  string `json:"last_login_at,omitempty"`
	EmployeeName string `json:"employee_name,omitempty"`
}

type UserPasswordUpdate struct {
	OldPassword string `json:"old_password,omitempty" binding:"required"`
	NewPassword string `json:"new_password,omitempty" binding:"required"`
}

type UserEnabledUpdate struct {
	Enabled bool `json:"enabled,omitempty"`
}

type UserLogin struct {
	Username string `json:"username,omitempty" binding:"required"`
	Password string `json:"password,omitempty" binding:"required"`
}
