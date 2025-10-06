package model

import "time"

type Department struct {
	ID        int64      `json:"id,omitempty"`
	Title     string     `json:"title,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
