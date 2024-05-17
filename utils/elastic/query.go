package elastic

// Any is a shortcut for a bool query with a should clause
func Any(queries ...map[string]any) map[string]any {
	return map[string]any{"bool": map[string]any{"should": queries}}
}

// All is a shortcut for a bool query with a must clause
func All(queries ...map[string]any) map[string]any {
	return map[string]any{"bool": map[string]any{"must": queries}}
}

// Not is a shortcut for a bool query with a must_not clause
func Not(query map[string]any) map[string]any {
	return map[string]any{"bool": map[string]any{"must_not": query}}
}

// Not is a shortcut for an ids query
func Ids(values ...string) map[string]any {
	return map[string]any{"ids": map[string]any{"values": values}}
}

// Term is a shortcut for a term query
func Term(field string, value any) map[string]any {
	return map[string]any{"term": map[string]any{field: value}}
}

// Exists is a shortcut for an exists query
func Exists(field string) map[string]any {
	return map[string]any{"exists": map[string]any{"field": field}}
}

// Nested is a shortcut for a nested query
func Nested(path string, query map[string]any) map[string]any {
	return map[string]any{"nested": map[string]any{"path": path, "query": query}}
}

// Match is a shortcut for a match query
func Match(field string, value any) map[string]any {
	return map[string]any{"match": map[string]any{field: map[string]any{"query": value}}}
}

// MatchPhrase is a shortcut for a match_phrase query
func MatchPhrase(field, value string) map[string]any {
	return map[string]any{"match_phrase": map[string]any{field: map[string]any{"query": value}}}
}

// GreaterThan is a shortcut for a range query where x > value
func GreaterThan(field string, value any) map[string]any {
	return map[string]any{"range": map[string]any{field: map[string]any{
		"from":          value,
		"include_lower": false,
		"include_upper": true,
		"to":            nil,
	}}}
}

// GreaterThanOrEqual is a shortcut for a range query where x >= value
func GreaterThanOrEqual(field string, value any) map[string]any {
	return map[string]any{"range": map[string]any{field: map[string]any{
		"from":          value,
		"include_lower": true,
		"include_upper": true,
		"to":            nil,
	}}}
}

// LessThan is a shortcut for a range query where x < value
func LessThan(field string, value any) map[string]any {
	return map[string]any{"range": map[string]any{field: map[string]any{
		"from":          nil,
		"include_lower": true,
		"include_upper": false,
		"to":            value,
	}}}
}

// LessThanOrEqual is a shortcut for a range query where x <= value
func LessThanOrEqual(field string, value any) map[string]any {
	return map[string]any{"range": map[string]any{field: map[string]any{
		"from":          nil,
		"include_lower": true,
		"include_upper": true,
		"to":            value,
	}}}
}

// Between is a shortcut for a range query where from <= x < to
func Between(field string, from, to any) map[string]any {
	return map[string]any{"range": map[string]any{field: map[string]any{
		"from":          from,
		"include_lower": true,
		"include_upper": false,
		"to":            to,
	}}}
}
