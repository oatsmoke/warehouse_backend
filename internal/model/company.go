package model

type Company struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Deleted bool   `json:"deleted"`
}
