package dto

import (
	"github.com/oatsmoke/warehouse_backend/internal/lib/role"
)

type User struct {
	Username   string    `json:"username,omitempty"`
	Email      string    `json:"email,omitempty" binding:"required"`
	Role       role.Role `json:"role,omitempty"`
	EmployeeID int64     `json:"employee_id,omitempty"`
}
