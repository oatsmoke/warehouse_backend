package dto

type Department struct {
	Title string `json:"title,omitempty" bind:"required"`
}
