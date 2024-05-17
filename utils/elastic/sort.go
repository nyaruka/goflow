package esq

func Sort(field string, ascending bool) map[string]any {
	return map[string]any{field: map[string]any{"order": order(ascending)}}
}

func SortNested(field string, filter map[string]any, path string, ascending bool) map[string]any {
	return map[string]any{field: map[string]any{
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
