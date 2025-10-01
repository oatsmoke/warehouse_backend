package dto

type Category struct {
	Title string `json:"title,omitempty" binding:"required"`
}
