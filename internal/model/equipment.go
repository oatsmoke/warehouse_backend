package model

type Equipment struct {
	ID           int64    `json:"id"`
	SerialNumber string   `json:"serial_number"`
	Profile      *Profile `json:"profile"`
	IsDeleted    bool     `json:"is_deleted"`
}
