package dto

type Profile struct {
	Title      string `json:"title,omitempty" binding:"required"`
	CategoryID int64  `json:"category_id,omitempty" binding:"required"`
}
