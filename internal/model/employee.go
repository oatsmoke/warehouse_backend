package model

import "time"

type Employee struct {
	ID         int64       `json:"id,omitempty"`
	LastName   string      `json:"last_name,omitempty"`
	FirstName  string      `json:"first_name,omitempty"`
	MiddleName string      `json:"middle_name,omitempty"`
	Phone      string      `json:"phone,omitempty"`
	Department *Department `json:"department,omitempty"`
	DeletedAt  *time.Time  `json:"deleted_at,omitempty"`
}

func NewEmployee() *Employee {
	return &Employee{
		Department: &Department{},
	}
}

//type RequestEmployee struct {
//	Ids          []int64 `json:"ids"`
//	DepartmentId int64   `json:"department_id"`
//}
