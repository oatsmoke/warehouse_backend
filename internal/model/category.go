package model

type Category struct {
	ID      int64  `json:"id,omitempty"`
	Title   string `json:"title,omitempty"`
	Deleted bool   `json:"deleted,omitempty"`
}
