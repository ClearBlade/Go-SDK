package GoSDK

import ()

type AnalyticsQuery struct {
	Scope    interface{} 		`json:"scope"`
	Filter   AnalyticsFilter	`json:"filter"`
}

type Range struct {
	Start	interface{}  	`json:"start"`
	End		interface{}  	`json:"end"`
}

type AnalyticsFilter struct {
	Module			interface{} 	`json:"module"`
	Id				interface{} 	`json:"id"`
	Action			interface{} 	`json:"action"`
	Interval		interface{} 	`json:"interval"`
	Limit			interface{} 	`json:"limit"`
	QueryRange 		Range 			`json:"range"`
}

func (q *AnalyticsQuery) serialize() map[string]interface{} {
	qMap := make(map[string]interface{})
	if q.Scope != nil {qMap["scope"] = q.Scope}

	if q.Filter.Module != nil 	{ qMap["module"] = 		q.Filter.Module	}
	if q.Filter.Id != nil 		{ qMap["id"] = 			q.Filter.Id }
	if q.Filter.Action != nil 	{ qMap["action"] = 		q.Filter.Action }
	if q.Filter.Interval != nil { qMap["interval"] = 	q.Filter.Interval }
	if q.Filter.Limit != nil	{ qMap["limit"] = 		q.Filter.Limit }
	if q.Filter.QueryRange.Start != nil && q.Filter.QueryRange.End != nil { qMap["range"] = map[string]interface{}{"start":q.Filter.QueryRange.Start,"end":q.Filter.QueryRange.End}}

	return qMap
}