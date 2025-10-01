package model

import "time"

type Profile struct {
	ID        int64      `json:"id,omitempty"`
	Title     string     `json:"title,omitempty"`
	Category  *Category  `json:"category,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
