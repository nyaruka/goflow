package jsonpath

import (
	"strconv"
	"strings"
)

func ParsePath(path string) []string {
	path = strings.ReplaceAll(path, "[*]", ".*")
	split := strings.Split(path, ".")
	parts := make([]string, 0, len(split))

	for _, p := range split {
		if p != "" && p != "$" {
			parts = append(parts, p)
		}
	}
	return parts
}

func Visit(j any, path []string, callback func(any)) error {
	return visit(nil, j, path, callback, nil)
}

func Transform(j any, path []string, callback func(any, any) any) error {
	return visit(nil, j, path, nil, callback)
}

func visit(container, j any, path []string, on func(any), tx func(any, any) any) error {
	selector := path[0]
	rem := path[1:]

	switch typed := j.(type) {
	case map[string]any:
		filter := func(k string, v any) bool {
			return k == selector || selector == "*"
		}

		for k, val := range typed {
			if filter(k, val) {
				if len(rem) == 0 {
					if tx != nil {
						typed[k] = tx(j, val)
					} else {
						on(val)
					}
				} else {
					if err := visit(j, val, rem, on, tx); err != nil {
						return err
					}
				}
			}
		}
	case []any:
		index, err := strconv.Atoi(selector)
		filter := func(i int, v any) bool {
			return (i == index && err == nil) || selector == "*"
		}

		for i, val := range typed {
			if filter(i, val) {
				if len(rem) == 0 {
					if tx != nil {
						typed[i] = tx(j, val)
					} else {
						on(val)
					}
				} else {
					if err := visit(j, val, rem, on, tx); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
