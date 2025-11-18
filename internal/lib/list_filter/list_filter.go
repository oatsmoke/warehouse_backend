package list_filter

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oatsmoke/warehouse_backend/internal/dto"
)

const (
	defaultWithDeleted       = "false"
	defaultSortBy            = "id"
	defaultOrder             = "asc"
	defaultOffset      int32 = 0
	defaultLimit       int32 = 50
)

func ParseQueryParams(c *gin.Context) *dto.QueryParams {
	qp := new(dto.QueryParams)

	qp.WithDeleted = c.DefaultQuery("deleted", defaultWithDeleted) == "true"

	qp.Search = c.Query("search")

	ids := c.QueryArray("ids")
	if len(ids) > 0 {
		for _, id := range ids {
			if parsedId, err := strconv.ParseInt(id, 10, 64); err == nil {
				qp.Ids = append(qp.Ids, parsedId)
			}
		}
	}

	qp.SortColumn = c.DefaultQuery("sort_by", defaultSortBy)
	qp.SortOrder = strings.ToLower(c.DefaultQuery("order", defaultOrder))
	if qp.SortOrder != "asc" && qp.SortOrder != "desc" {
		qp.SortOrder = defaultOrder
	}

	if n, err := strconv.ParseInt(c.Query("limit"), 10, 32); err != nil || n < 0 {
		qp.PaginationLimit = defaultLimit
	}

	if n, err := strconv.ParseInt(c.Query("offset"), 10, 32); err != nil || n < 0 {
		qp.PaginationOffset = defaultOffset
	}

	return qp
}
