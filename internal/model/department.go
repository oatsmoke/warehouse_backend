package model

type Department struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	IsDeleted bool   `json:"is_deleted"`
}
