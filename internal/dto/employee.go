package dto

type Employee struct {
	LastName     string `json:"last_name,omitempty" binding:"required"`
	FirstName    string `json:"first_name,omitempty" binding:"required"`
	MiddleName   string `json:"middle_name,omitempty"`
	Phone        string `json:"phone,omitempty" binding:"required"`
	Email        string `json:"email,omitempty" binding:"required"`
	DepartmentID int64  `json:"department_id,omitempty"`
}
