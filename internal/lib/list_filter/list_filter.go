package list_filter

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oatsmoke/warehouse_backend/internal/dto"
)

const (
	defaultOrder  = "ASC"
	defaultOffset = "0"
	defaultLimit  = "50"
)

func ParseQueryParams(c *gin.Context) *dto.QueryParams {
	qp := new(dto.QueryParams)

	qp.WithDeleted = c.DefaultQuery("deleted", "false") == "true"

	qp.Search = c.Query("search")

	qp.SortBy = c.Query("sort_by")

	qp.Order = strings.ToUpper(c.DefaultQuery("order", defaultOrder))
	if qp.Order != "ASC" && qp.Order != "DESC" {
		qp.Order = defaultOrder
	}

	qp.Offset = c.DefaultQuery("offset", defaultOffset)

	qp.Limit = c.DefaultQuery("limit", defaultLimit)

	return qp
}

func BuildQuery(qp *dto.QueryParams, fields []string, table string) (string, []interface{}) {
	i := 1
	var args []interface{}

	withDeleted := fmt.Sprintf("WHERE ($%d OR %s.deleted_at IS NULL)", i, table)
	args = append(args, qp.WithDeleted)
	i++

	var search string
	if qp.Search != "" {
		first := true
		qp.Search = fmt.Sprintf("%%%s%%", qp.Search)
		for _, field := range fields {
			if first {
				search = fmt.Sprintf(" AND (%s ILIKE $%d", field, i)
				args = append(args, qp.Search)
				i++
				first = false
				continue
			}
			search += fmt.Sprintf(" OR %s ILIKE $%d", field, i)
			args = append(args, qp.Search)
			i++
		}
		search += ")"
	}

	var exist bool
	defaultSortBy := fmt.Sprintf("%s.id", table)
	for _, field := range fields {
		if qp.SortBy == field || qp.SortBy == defaultSortBy {
			exist = true
			break
		}
	}
	if !exist {
		qp.SortBy = defaultSortBy
	}
	order := fmt.Sprintf(" ORDER BY %s %s", qp.SortBy, qp.Order)

	limit := fmt.Sprintf(" LIMIT $%d OFFSET $%d", i, i+1)
	args = append(args, qp.Limit, qp.Offset)

	return fmt.Sprintf("%s%s%s%s;", withDeleted, search, order, limit), args
}
