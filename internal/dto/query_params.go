package dto

type QueryParams struct {
	WithDeleted bool
	Search      string
	Ids         []int64
	SortBy      string
	Order       string
	Limit       string
	Offset      string
}
