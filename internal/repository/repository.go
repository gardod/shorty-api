package repository

import (
	"strconv"
	"strings"
)

type QueryOptions struct {
	Distinct  []string
	Join      []Join
	Where     []string
	Group     []string
	Having    []string
	Order     []string
	Limit     *int
	Offset    *int
	Arguments []interface{}
}

func (q *QueryOptions) AddDistinct(options ...DistinctOption) {
	if len(options) == 0 {
		q.Distinct = append(q.Distinct, "")
	}

	for _, o := range options {
		q.addOption(o)
	}
}

func (q *QueryOptions) AddWhere(options ...WhereOption) {
	for _, o := range options {
		q.addOption(o)
	}
}

func (q *QueryOptions) AddOrder(options ...OrderOption) {
	for _, o := range options {
		q.addOption(o)
	}
}

func (q *QueryOptions) addOption(option interface{}) {
	if o, ok := option.(DistinctOption); ok {
		q.Distinct = append(q.Distinct, o.GetDistinct())
	}

	if o, ok := option.(JoinOption); ok {
		q.Join = append(q.Join, o.GetJoin()...)
	}

	if o, ok := option.(WhereOption); ok {
		start := len(q.Arguments) + 1
		stmt, args := o.GetWhere(start)

		if stmt != "" {
			q.Where = append(q.Where, stmt)
			q.Arguments = append(q.Arguments, args...)
		}
	}

	if o, ok := option.(GroupOption); ok {
		q.Group = append(q.Group, o.GetGroup())
	}

	if o, ok := option.(HavingOption); ok {
		start := len(q.Arguments) + 1
		stmt, args := o.GetHaving(start)

		if stmt != "" {
			q.Having = append(q.Having, stmt)
			q.Arguments = append(q.Arguments, args...)
		}
	}

	if o, ok := option.(OrderOption); ok {
		q.Order = append(q.Order, o.GetOrder())
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
	if len(q.Join) == 0 {
		return ""
	}

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

func (q *QueryOptions) BuildGroup() string {
	if len(q.Group) == 0 {
		return ""
	}

	fields := make([]string, 0, len(q.Group))
	unique := make(map[string]struct{}, len(q.Group))
	for _, s := range q.Group {
		if _, ok := unique[s]; ok {
			continue
		}

		fields = append(fields, s)
		unique[s] = struct{}{}
	}

	return "GROUP BY " + strings.Join(fields, ",")
}

func (q *QueryOptions) BuildHaving() string {
	if len(q.Having) == 0 {
		return ""
	}

	return "HAVING (" + strings.Join(q.Having, ") AND (") + ")"
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

type Join struct {
	Type  string
	Table string
	On    string
}

func (j *Join) Build() string {
	return j.Type + " JOIN " + j.Table + " ON " + j.On
}

type DistinctOption interface {
	GetDistinct() string
}

type JoinOption interface {
	GetJoin() []Join
}

type WhereOption interface {
	GetWhere(start int) (stmt string, args []interface{})
}

type HavingOption interface {
	WhereOption
	GetHaving(start int) (stmt string, args []interface{})
}

type GroupOption interface {
	GetGroup() string
}

type OrderOption interface {
	GetOrder() string
}

type OrderDirection string

const (
	OrderDirectionAsc           = OrderDirection(`ASC`)
	OrderDirectionAscNullsFirst = OrderDirection(`ASC NULLS FIRST`)
	OrderDirectionDesc          = OrderDirection(`DESC`)
	OrderDirectionDescNullsLast = OrderDirection(`DESC NULLS LAST`)
)
