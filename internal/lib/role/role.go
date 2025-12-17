package role

type Role int32

const (
	RootRole Role = iota + 1
	AdminRole
	GoverningRole
	EmployeeRole
	UserRole
)

const UndefinedRole = "undefined"

var roles = [...]string{
	"",
	"root",
	"Администратор",
	"Управляющий",
	"Сотрудник",
	"Пользователь",
}

func (r Role) IsValid() bool {
	return RootRole <= r && r <= UserRole
}

func (r Role) CanAccess(required Role) bool {
	if !r.IsValid() {
		return false
	}

	return r <= required
}

func (r Role) String() string {
	if !r.IsValid() {
		return UndefinedRole
	}

	return roles[r]
}

type FullRole struct {
	ID   int    `json:"id"`
	Role string `json:"role"`
}

func AllRole() []*FullRole {
	r := make([]*FullRole, 0, len(roles)-2)
	for i, role := range roles {
		if i < 2 {
			continue
		}

		r = append(r, &FullRole{i, role})
	}
	return r
}
