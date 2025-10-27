package dto

type ListResponse[T any] struct {
	List  T   `json:"list,omitempty"`
	Total int `json:"total,omitempty"`
}
