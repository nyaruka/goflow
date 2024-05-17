package elastic

type Query map[string]any

// Any is a shortcut for a bool query with a should clause
func Any(queries ...Query) Query {
	return Query{"bool": Query{"should": queries}}
}

// All is a shortcut for a bool query with a must clause
func All(queries ...Query) Query {
	return Query{"bool": Query{"must": queries}}
}

// Not is a shortcut for a bool query with a must_not clause
func Not(query Query) Query {
	return Query{"bool": Query{"must_not": query}}
}

// Not is a shortcut for an ids query
func Ids(values ...string) Query {
	return Query{"ids": Query{"values": values}}
}

// Term is a shortcut for a term query
func Term(field string, value any) Query {
	return Query{"term": Query{field: value}}
}

// Exists is a shortcut for an exists query
func Exists(field string) Query {
	return Query{"exists": Query{"field": field}}
}

// Nested is a shortcut for a nested query
func Nested(path string, query Query) Query {
	return Query{"nested": Query{"path": path, "query": query}}
}

// Match is a shortcut for a match query
func Match(field string, value any) Query {
	return Query{"match": Query{field: Query{"query": value}}}
}

// MatchPhrase is a shortcut for a match_phrase query
func MatchPhrase(field, value string) Query {
	return Query{"match_phrase": Query{field: Query{"query": value}}}
}

// GreaterThan is a shortcut for a range query where x > value
func GreaterThan(field string, value any) Query {
	return Query{"range": Query{field: Query{
		"from":          value,
		"include_lower": false,
		"include_upper": true,
		"to":            nil,
	}}}
}

// GreaterThanOrEqual is a shortcut for a range query where x >= value
func GreaterThanOrEqual(field string, value any) Query {
	return Query{"range": Query{field: Query{
		"from":          value,
		"include_lower": true,
		"include_upper": true,
		"to":            nil,
	}}}
}

// LessThan is a shortcut for a range query where x < value
func LessThan(field string, value any) Query {
	return Query{"range": Query{field: Query{
		"from":          nil,
		"include_lower": true,
		"include_upper": false,
		"to":            value,
	}}}
}

// LessThanOrEqual is a shortcut for a range query where x <= value
func LessThanOrEqual(field string, value any) Query {
	return Query{"range": Query{field: Query{
		"from":          nil,
		"include_lower": true,
		"include_upper": true,
		"to":            value,
	}}}
}

// Between is a shortcut for a range query where from <= x < to
func Between(field string, from, to any) Query {
	return Query{"range": Query{field: Query{
		"from":          from,
		"include_lower": true,
		"include_upper": false,
		"to":            to,
	}}}
}
