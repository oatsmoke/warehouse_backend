package model

import "time"

type Employee struct {
	ID                int64       `json:"id,omitempty"`
	Name              string      `json:"name,omitempty"`
	Phone             string      `json:"phone,omitempty"`
	Email             string      `json:"email,omitempty"`
	Password          string      `json:"password,omitempty"`
	Hash              string      `json:"hash,omitempty"`
	RegistrationDate  *time.Time  `json:"registration_date,omitempty"`
	AuthorizationDate *time.Time  `json:"authorization_date,omitempty"`
	Activate          bool        `json:"activate,omitempty"`
	Hidden            bool        `json:"hidden,omitempty"`
	Department        *Department `json:"department,omitempty"`
	Role              string      `json:"role,omitempty"`
	Deleted           bool        `json:"deleted,omitempty"`
}

type RequestEmployee struct {
	Ids          []int64 `json:"ids"`
	DepartmentId int64   `json:"department_id"`
}
