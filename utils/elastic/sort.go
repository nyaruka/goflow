package elastic

type Sort map[string]any

// SortBy is a shortcut for a simple field sort
func SortBy(field string, ascending bool) Sort {
	return Sort{field: map[string]any{"order": order(ascending)}}
}

// SortNested is a shortcut for a nested field sort
func SortNested(field string, filter Query, path string, ascending bool) Sort {
	return Sort{field: map[string]any{
		"nested": map[string]any{"filter": filter, "path": path},
		"order":  order(ascending),
	}}
}

func order(asc bool) string {
	if asc {
		return "asc"
	}
	return "desc"
}
