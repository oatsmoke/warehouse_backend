package model

type Company struct {
	Id        int    `json:"id,omitempty" db:"id"`
	Title     string `json:"title,omitempty" db:"title"`
	IsDeleted bool   `json:"isDeleted,omitempty" db:"is_deleted"`
}
