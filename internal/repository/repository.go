package repository

import (
	"strconv"
	"strings"
)

type QueryOptions struct {
	Distinct  []string
	Join      []JoinOption
	Where     []string
	Arguments []interface{}
	Order     []string
	Limit     *int
	Offset    *int
}

func (q *QueryOptions) IsEmpty() bool {
	return len(q.Distinct) == 0 &&
		len(q.Join) == 0 &&
		len(q.Where) == 0 &&
		len(q.Order) == 0 &&
		q.Limit == nil &&
		q.Offset == nil
}

func (q *QueryOptions) AddDistinct(options ...DistinctOption) {
	if len(options) == 0 {
		q.Distinct = append(q.Distinct, "")
	}

	for _, o := range options {
		q.Distinct = append(q.Distinct, string(o))
	}
}

func (q *QueryOptions) AddWhere(options ...WhereOption) {
	for _, o := range options {
		start := len(q.Arguments) + 1
		stmt, args := o.GetWhere(start)

		q.Where = append(q.Where, stmt)
		q.Arguments = append(q.Arguments, args...)

		if owj, ok := o.(WhereOptionWithJoin); ok {
			q.Join = append(q.Join, owj.GetJoin()...)
		}
	}
}

func (q *QueryOptions) AddOrder(options ...OrderOption) {
	for _, o := range options {
		q.Order = append(q.Order, string(o))
	}
}

func (q *QueryOptions) SetLimit(limit int) {
	q.Limit = &limit
}

func (q *QueryOptions) SetOffset(offset int) {
	q.Offset = &offset
}

func (q *QueryOptions) BuildDistinct() string {
	if len(q.Distinct) == 0 {
		return ""
	}

	if q.Distinct[0] == "" {
		return "DISTINCT"
	}

	return "DISTINCT ON (" + strings.Join(q.Distinct, ",") + ")"
}

func (q *QueryOptions) BuildJoin() string {
	stmt := make([]string, 0, len(q.Join))
	tables := make(map[string]struct{}, len(q.Join))

	for _, j := range q.Join {
		if _, ok := tables[j.Table]; ok {
			continue
		}

		stmt = append(stmt, j.Build())
		tables[j.Table] = struct{}{}
	}

	return strings.Join(stmt, "\n")
}

func (q *QueryOptions) BuildWhere() string {
	if len(q.Where) == 0 {
		return ""
	}

	return "WHERE (" + strings.Join(q.Where, ") AND (") + ")"
}

func (q *QueryOptions) BuildOrder() string {
	if len(q.Order) == 0 {
		return ""
	}

	return "ORDER BY " + strings.Join(q.Order, ",")
}

func (q *QueryOptions) BuildLimit() string {
	stmt := ""

	if q.Limit != nil {
		stmt += "LIMIT " + strconv.Itoa(*q.Limit)
	}

	if q.Offset != nil {
		stmt += " OFFSET " + strconv.Itoa(*q.Offset)
	}

	return stmt
}

type JoinOption struct {
	Type  string
	Table string
	On    string
}

func (j *JoinOption) Build() string {
	return j.Type + " JOIN " + j.Table + " ON " + j.On
}

type WhereOption interface {
	GetWhere(start int) (stmt string, args []interface{})
}

type WhereOptionWithJoin interface {
	WhereOption
	GetJoin() []JoinOption
}

type DistinctOption string

type OrderOption string
