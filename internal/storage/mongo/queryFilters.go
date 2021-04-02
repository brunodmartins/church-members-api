package mongo

type QueryFilters map[string]interface{}


func (qf QueryFilters) AddFilter(key string, value interface{}) {
	qf[key] = value
}
