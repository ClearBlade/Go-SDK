package GoSDK

import ()

// Filter is the atomic structure inside a query it contains
// A field a value and an operator
type Filter struct {
	Field    string
	Value    interface{}
	Operator string
}

type Ordering struct {
	SortOrder bool
	OrderKey  string
}

type Query struct {
	Filters    [][]Filter
	PageSize   int
	PageNumber int
	Order      []Ordering
}

func NewQuery() *Query {
	query := &Query{
		Filters: [][]Filter{[]Filter{}},
		Order:   []Ordering{},
	}
	return query
}

func (q *Query) EqualTo(field string, value interface{}) {
	f := Filter{
		Field:    field,
		Value:    value,
		Operator: "=",
	}
	q.Filters[0] = append(q.Filters[0], f)
}

func (q *Query) GreaterThan(field string, value interface{}) {
	f := Filter{
		Field:    field,
		Value:    value,
		Operator: ">",
	}
	q.Filters[0] = append(q.Filters[0], f)
}

func (q *Query) GreaterThanEqualTo(field string, value interface{}) {
	f := Filter{
		Field:    field,
		Value:    value,
		Operator: ">=",
	}
	q.Filters[0] = append(q.Filters[0], f)
}

func (q *Query) LessThan(field string, value interface{}) {
	f := Filter{
		Field:    field,
		Value:    value,
		Operator: "<",
	}
	q.Filters[0] = append(q.Filters[0], f)
}

func (q *Query) LessThanEqualTo(field string, value interface{}) {
	f := Filter{
		Field:    field,
		Value:    value,
		Operator: "<=",
	}
	q.Filters[0] = append(q.Filters[0], f)
}

func (q *Query) NotEqualTo(field string, value interface{}) {
	f := Filter{
		Field:    field,
		Value:    value,
		Operator: "!=",
	}
	q.Filters[0] = append(q.Filters[0], f)
}

func (q *Query) Matches(field, regex string) {
	f := Filter{
		Field:    field,
		Value:    regex,
		Operator: "~",
	}
	q.Filters[0] = append(q.Filters[0], f)
}

func (q *Query) Or(orQuery *Query) {
	q.Filters = append(q.Filters, orQuery.Filters...)
}

// Map will produce the kind of thing that is sent as a query
// either as the body of a request or as a queryString
func (q *Query) serialize() map[string]interface{} {
	qrMap := make(map[string]interface{})
	qrMap["PAGENUM"] = q.PageNumber
	qrMap["PAGESIZE"] = q.PageSize
	sortMap := make([]map[string]interface{}, len(q.Order))
	for i, ordering := range q.Order {
		sortMap[i] = make(map[string]interface{})
		if ordering.SortOrder {
			sortMap[i]["ASC"] = ordering.OrderKey
		} else {
			sortMap[i]["DESC"] = ordering.OrderKey
		}
	}
	qrMap["SORT"] = sortMap
	filterSlice := make([][]map[string]interface{}, len(q.Filters))
	for i, querySlice := range q.Filters {
		qm := make([]map[string]interface{}, len(querySlice))
		for j, query := range querySlice {
			mapForQuery := make(map[string]interface{})
			var op string
			switch query.Operator {
			case "=":
				op = "EQ"
			case ">":
				op = "GT"
			case "<":
				op = "LT"
			case ">=":
				op = "GTE"
			case "<=":
				op = "LTE"
			case "/=":
				op = "NEQ"
			default:
				op = "EQ"
			}
			mapForQuery[op] = []map[string]interface{}{map[string]interface{}{query.Field: query.Value}}
			qm[j] = mapForQuery
		}
		filterSlice[i] = qm
	}
	qrMap["FILTERS"] = filterSlice
	return qrMap
}
