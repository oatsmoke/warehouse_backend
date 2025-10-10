package dto

import (
	"github.com/oatsmoke/warehouse_backend/internal/lib/role"
)

type UserCreate struct {
	Username   string    `json:"username,omitempty"`
	Email      string    `json:"email,omitempty" binding:"required"`
	Role       role.Role `json:"role,omitempty"`
	EmployeeID int64     `json:"employee_id,omitempty"`
}

type UserUpdate struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty" binding:"required"`
}

type UserPasswordUpdate struct {
	OldPassword string `json:"old_password,omitempty" binding:"required"`
	NewPassword string `json:"new_password,omitempty" binding:"required"`
}

type UserRoleUpdate struct {
	Role role.Role `json:"role,omitempty" binding:"required"`
}

type UserEnabledUpdate struct {
	Enabled bool `json:"enabled,omitempty"`
}

type UserEmployeeUpdate struct {
	EmployeeID int64 `json:"employee_id,omitempty"`
}
