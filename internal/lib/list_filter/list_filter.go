package list_filter

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oatsmoke/warehouse_backend/internal/dto"
)

const (
	defaultWithDeleted = "false"
	defaultSortBy      = "id"
	defaultOrder       = "asc"
	defaultOffset      = "0"
	defaultLimit       = "50"
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

	qp.SortBy = c.DefaultQuery("sort_by", defaultSortBy)
	qp.Order = strings.ToLower(c.DefaultQuery("order", defaultOrder))
	if qp.Order != "asc" && qp.Order != "desc" {
		qp.Order = defaultOrder
	}

	qp.Limit = c.DefaultQuery("limit", defaultLimit)
	if n, err := strconv.ParseInt(qp.Limit, 10, 64); err != nil || n < 0 {
		qp.Limit = defaultLimit
	}

	qp.Offset = c.DefaultQuery("offset", defaultOffset)
	if n, err := strconv.ParseInt(qp.Offset, 10, 64); err != nil || n < 0 {
		qp.Offset = defaultOffset
	}

	return qp
}
