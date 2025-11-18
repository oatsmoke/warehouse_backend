package dto

type ListResponse[T any] struct {
	List  T     `json:"list,omitempty"`
	Total int64 `json:"total,omitempty"`
}
