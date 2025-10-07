package model

import "time"

type Equipment struct {
	ID           int64      `json:"id,omitempty"`
	SerialNumber string     `json:"serial_number,omitempty"`
	Profile      *Profile   `json:"profile,omitempty"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

func NewEquipment() *Equipment {
	return &Equipment{
		Profile: &Profile{
			Category: &Category{},
		},
	}
}
