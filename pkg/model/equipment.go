package model

type Equipment struct {
	Id           int     `json:"id,omitempty" db:"id"`
	SerialNumber string  `json:"serialNumber,omitempty" db:"serial_number"`
	Profile      Profile `json:"profile" db:"profile"`
	IsDeleted    bool    `json:"isDeleted,omitempty" db:"is_deleted"`
}
