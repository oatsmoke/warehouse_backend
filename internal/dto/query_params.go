package dto

type QueryParams struct {
	WithDeleted bool
	Search      string
	SortBy      string
	Order       string
	Offset      string
	Limit       string
}
