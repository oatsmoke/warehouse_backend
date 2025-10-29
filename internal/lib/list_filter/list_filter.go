package list_filter

import (
	"fmt"
	"strconv"
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

	ids := c.QueryArray("ids")
	for _, id := range ids {
		if parsedId, err := strconv.ParseInt(id, 10, 64); err == nil {
			qp.Ids = append(qp.Ids, parsedId)
		}
	}

	qp.SortBy = c.Query("sort_by")
	qp.Order = strings.ToUpper(c.DefaultQuery("order", defaultOrder))

	qp.Offset = c.Query("offset")
	qp.Limit = c.Query("limit")

	return qp
}

func BuildQuery(qp *dto.QueryParams, fields []string, table string) (string, []interface{}) {
	i := 1
	var (
		args   []interface{}
		search string
		ids    string
		exist  bool
	)

	withDeleted := fmt.Sprintf("WHERE ($%d OR %s.deleted_at IS NULL)", i, table)
	args = append(args, qp.WithDeleted)
	i++

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

	if len(qp.Ids) > 0 {
		ids = fmt.Sprintf(" AND %s.id = ANY($%d)", table, i)
		args = append(args, qp.Ids)
		i++
	}

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
	if qp.Order != "ASC" && qp.Order != "DESC" {
		qp.Order = defaultOrder
	}
	order := fmt.Sprintf(" ORDER BY %s %s", qp.SortBy, qp.Order)

	limit := fmt.Sprintf(" LIMIT $%d OFFSET $%d", i, i+1)
	if n, err := strconv.ParseInt(qp.Limit, 10, 64); err != nil || n < 1 {
		qp.Limit = defaultLimit
	} else {
		qp.Limit = strconv.FormatInt(n, 10)
	}
	if n, err := strconv.ParseInt(qp.Offset, 10, 64); err != nil || n < 1 {
		qp.Offset = defaultOffset
	} else {
		qp.Offset = strconv.FormatInt(n, 10)
	}
	args = append(args, qp.Limit, qp.Offset)

	return fmt.Sprintf("%s%s%s%s%s;", withDeleted, search, ids, order, limit), args
}
