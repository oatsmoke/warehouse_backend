package dto

type QueryParams struct {
	WithDeleted      bool
	Search           string
	IDs              []int64
	SortColumn       string
	SortOrder        string
	PaginationLimit  int32
	PaginationOffset int32
	Param            string
	ParamID          int64
}
