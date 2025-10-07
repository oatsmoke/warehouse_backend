package model

import "time"

type User struct {
	ID           int64      `json:"id,omitempty"`
	Username     string     `json:"username ,omitempty"`
	PasswordHash string     `json:"-"`
	Role         Role       `json:"role,omitempty"`
	Enabled      bool       `json:"enabled,omitempty"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	Employee     *Employee  `json:"employee,omitempty"`
}
