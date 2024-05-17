package es

func Any(queries ...map[string]any) map[string]any {
	return map[string]any{"bool": map[string]any{"should": queries}}
}

func All(queries ...map[string]any) map[string]any {
	return map[string]any{"bool": map[string]any{"must": queries}}
}

func Not(query map[string]any) map[string]any {
	return map[string]any{"bool": map[string]any{"must_not": query}}
}

func Ids(values ...string) map[string]any {
	return map[string]any{"ids": map[string]any{"values": values}}
}

func Term(field string, value any) map[string]any {
	return map[string]any{"term": map[string]any{field: value}}
}

func Exists(field string) map[string]any {
	return map[string]any{"exists": map[string]any{"field": field}}
}

func Nested(path string, query map[string]any) map[string]any {
	return map[string]any{"nested": map[string]any{"path": path, "query": query}}
}

func Match(field string, value any) map[string]any {
	return map[string]any{"match": map[string]any{field: map[string]any{"query": value}}}
}

func MatchPhrase(field, value string) map[string]any {
	return map[string]any{"match_phrase": map[string]any{field: map[string]any{"query": value}}}
}

func GreaterThan(field string, value any) map[string]any {
	return map[string]any{"range": map[string]any{field: map[string]any{
		"from":          value,
		"include_lower": false,
		"include_upper": true,
		"to":            nil,
	}}}
}

func GreaterThanOrEqual(field string, value any) map[string]any {
	return map[string]any{"range": map[string]any{field: map[string]any{
		"from":          value,
		"include_lower": true,
		"include_upper": true,
		"to":            nil,
	}}}
}

func LessThan(field string, value any) map[string]any {
	return map[string]any{"range": map[string]any{field: map[string]any{
		"from":          nil,
		"include_lower": true,
		"include_upper": false,
		"to":            value,
	}}}
}

func LessThanOrEqual(field string, value any) map[string]any {
	return map[string]any{"range": map[string]any{field: map[string]any{
		"from":          nil,
		"include_lower": true,
		"include_upper": true,
		"to":            value,
	}}}
}

func Between(field string, from, to any) map[string]any {
	return map[string]any{"range": map[string]any{field: map[string]any{
		"from":          from,
		"include_lower": true,
		"include_upper": false,
		"to":            to,
	}}}
}
