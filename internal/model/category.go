package model

// Category represents a product category.
// Contains the identifier, title, and soft delete flag.
type Category struct {
	ID      int64  `json:"id,omitempty"`
	Title   string `json:"title,omitempty"`
	Deleted bool   `json:"deleted,omitempty"`
}
