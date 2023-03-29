package model

type Category struct {
	Id    int    `json:"id,omitempty" db:"id"`
	Title string `json:"title,omitempty" db:"title"`
}
