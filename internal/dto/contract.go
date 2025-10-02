package dto

type Contract struct {
	Number  string `json:"number,omitempty" binding:"required"`
	Address string `json:"address,omitempty" binding:"required"`
}
