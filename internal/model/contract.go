package model

import "time"

type Contract struct {
	ID        int64      `json:"id,omitempty"`
	Number    string     `json:"number,omitempty"`
	Address   string     `json:"address,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
