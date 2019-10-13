package filters

import (
	"github.com/go-pg/pg/orm"
	_ "github.com/lib/pq"
)

const DefaultLimit = 25

type Filter func(*orm.Query) (*orm.Query, error)

func PageFilter(page, limit int) Filter {
	if limit == 0 {
		limit = DefaultLimit
	}
	offset := 0
	if page > 1 {
		offset = page * limit
	}

	return func(q *orm.Query) (*orm.Query, error) {
		if limit > 0 {
			q.Limit(limit)
		}
		if offset > 0 {
			q.Offset(offset)
		}
		return q, nil
	}
}

func TrashedFilter(trashed bool) Filter {
	return func(q *orm.Query) (*orm.Query, error) {
		if trashed {
			q.Where("deleted_at is not null")
		} else {
			q.Where("deleted_at is null")
		}
		return q, nil
	}
}
