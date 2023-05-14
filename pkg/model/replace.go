package model

type Replace struct {
	Id           int `json:"id,omitempty" db:"id"`
	TransferFrom int `json:"transferFrom,omitempty" db:"transfer_from"`
	TransferTo   int `json:"transferTo,omitempty" db:"transfer_to"`
}
