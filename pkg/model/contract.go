package model

type Contract struct {
	Id        int    `json:"id,omitempty" db:"id"`
	Number    string `json:"number,omitempty" db:"number"`
	Address   string `json:"address,omitempty" db:"address"`
	IsDeleted bool   `json:"isDeleted,omitempty" db:"is_deleted"`
}
