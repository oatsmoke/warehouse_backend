package model

type Equipment struct {
	ID           int64    `json:"id,omitempty"`
	SerialNumber string   `json:"serial_number,omitempty"`
	Profile      *Profile `json:"profile,omitempty"`
	Deleted      bool     `json:"deleted,omitempty"`
}
