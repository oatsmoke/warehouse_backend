package model

type Profile struct {
	Id       int      `json:"id,omitempty" db:"id"`
	Title    string   `json:"title,omitempty" db:"title"`
	Category Category `json:"category" db:"category"`
}
