package dto

type Company struct {
	Title string `json:"title,omitempty" binding:"required"`
}
