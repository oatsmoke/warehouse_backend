package model

type Category struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Deleted bool   `json:"deleted"`
}
