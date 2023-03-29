package model

type Employee struct {
	Id                int        `json:"id,omitempty" db:"id"`
	Name              string     `json:"name,omitempty" db:"name"`
	Phone             string     `json:"phone,omitempty" db:"phone"`
	Email             string     `json:"email,omitempty" db:"email"`
	Password          string     `json:"password,omitempty" db:"password"`
	Hash              string     `json:"hash,omitempty" db:"hash"`
	RegistrationDate  int64      `json:"registrationDate,omitempty" db:"registration_date"`
	AuthorizationDate int64      `json:"authorizationDate,omitempty" db:"authorization_date"`
	Activate          bool       `json:"activate,omitempty" db:"activate"`
	Hidden            bool       `json:"hidden,omitempty" db:"hidden"`
	Department        Department `json:"department" db:"department"`
	Role              string     `json:"role,omitempty" db:"role"`
	IsDeleted         bool       `json:"isDeleted,omitempty" db:"is_deleted"`
}

type SignInInput struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RequestEmployee struct {
	Ids []int `json:"ids,omitempty"`
	Id  int   `json:"id,omitempty"`
}
