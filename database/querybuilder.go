package database

import (
	"errors"
	. "github.com/volix-dev/leopard/helpers"
	"strings"
)

type QueryBuilder struct {
	table     string
	operation Operation
	wheres    []where
	selects   []string
	order     *orderBy
	limit     *int
	groupBy   *string
}

type orderBy struct {
	column string
	order  string
}

type where struct {
	or       bool
	field    string
	operator string
	value    interface{}
}

func (qb *QueryBuilder) Where(field string, operator string, value interface{}) *QueryBuilder {
	qb.wheres = append(qb.wheres, where{false, field, operator, value})
	return qb
}

func (qb *QueryBuilder) OrWhere(field string, operator string, value interface{}) *QueryBuilder {
	qb.wheres = append(qb.wheres, where{true, field, operator, value})
	return qb
}

func (qb *QueryBuilder) Select(selects ...string) *QueryBuilder {
	qb.selects = append(qb.selects, selects...)
	return qb
}

func (qb *QueryBuilder) OverwriteSelect(selects ...string) *QueryBuilder {
	qb.selects = selects
	return qb
}

func (qb *QueryBuilder) OrderBy(column string, order string) *QueryBuilder {
	qb.order = &orderBy{column, order}
	return qb
}

func (qb QueryBuilder) OrderByAsc(column string) *QueryBuilder {
	return qb.OrderBy(column, "ASC")
}

func (qb QueryBuilder) OrderByDesc(column string) *QueryBuilder {
	return qb.OrderBy(column, "DESC")
}

func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limit = &limit
	return qb
}

func (qb *QueryBuilder) GroupBy(column string) *QueryBuilder {
	qb.groupBy = &column
	return qb
}

// todo: Change this into different drivers

func (qb *QueryBuilder) Build() (string, []interface{}, error) {
	switch qb.operation {
	case Select:
		return qb.buildSelect()
	}
	return "", nil, errors.New("invalid operation")
}

func (qb *QueryBuilder) buildSelect() (string, []interface{}, error) {
	builder := strings.Builder{}

	builder.WriteString("SELECT ")

	if len(qb.selects) > 0 {
		builder.WriteString(strings.Join(qb.selects, ", "))
	} else {
		builder.WriteString("*")
	}

	builder.WriteString(" FROM ")
	builder.WriteString(qb.table)
	builder.WriteString(" ")

	var whereValues []interface{}
	if len(qb.wheres) == 0 {
		builder.WriteString("WHERE ")
		wheres, whereValues2 := qb.buildWheres()
		whereValues = whereValues2
		builder.WriteString(wheres)
	}

	builder.WriteString(" ")

	if qb.order != nil {
		builder.WriteString("ORDER BY ")
		builder.WriteString(qb.order.column)
		builder.WriteString(" ")
		builder.WriteString(qb.order.order)
	}

	var limitValues []interface{}
	if qb.limit != nil {
		builder.WriteString(" ")
		limit, values := qb.buildLimit()
		builder.WriteString(limit)

		limitValues = values
	}

	var groupByValues []interface{}
	if qb.limit != nil {
		builder.WriteString(" ")
		groupBy, values := qb.buildGroupBy()
		builder.WriteString(groupBy)

		groupByValues = values
	}

	return builder.String(), append(whereValues, append(limitValues, groupByValues...)...), nil
}

func (qb *QueryBuilder) buildWheres() (string, []interface{}) {
	wheres := strings.Builder{}
	var values []interface{}

	for i, where := range qb.wheres {
		if i > 0 {
			wheres.WriteString(" ")
			wheres.WriteString(Ternary(where.or, "OR ", "AND "))
		}

		wheres.WriteString(where.field)
		wheres.WriteString(" ")
		wheres.WriteString(where.operator)
		wheres.WriteString(" ?")

		values = append(values, where.value)
	}

	return wheres.String(), values
}

func (qb *QueryBuilder) buildLimit() (string, []interface{}) {
	if qb.limit == nil {
		return "", []interface{}{}
	}

	return "LIMIT ?", []interface{}{*qb.limit}
}

func (qb *QueryBuilder) buildGroupBy() (string, []interface{}) {
	if qb.groupBy == nil {
		return "", []interface{}{}
	}

	return "GROUP BY ?", []interface{}{*qb.groupBy}
}

func NewQueryBuilder(table string) *QueryBuilder {
	return &QueryBuilder{table: table, operation: Select}
}
