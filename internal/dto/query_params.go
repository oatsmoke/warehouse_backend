package dto

type QueryParams struct {
	WithDeleted      bool
	Search           string
	Ids              []int64
	SortColumn       string
	SortOrder        string
	PaginationLimit  int32
	PaginationOffset int32
}
