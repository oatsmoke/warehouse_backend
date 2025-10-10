package role

type Role string

const (
	AdminRole     Role = "admin"
	EmployeeRole  Role = "employee"
	GoverningRole Role = "governing"
)

func (r Role) IsValid() bool {
	return r == AdminRole || r == EmployeeRole || r == GoverningRole
}
