package model

import "time"

type Equipment struct {
	ID           int64      `json:"id,omitempty"`
	Company      *Company   `json:"company,omitempty"`
	Profile      *Profile   `json:"profile,omitempty"`
	SerialNumber string     `json:"serial_number,omitempty"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}
