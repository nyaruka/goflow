package jsonpath

// Very barebones jsonpath implementation, only supports the path types we need for now.. no expressions, only *
// wildcards.
//
// But allows us to transform values, with access to the container which may be a localizable object.

import (
	"strconv"
	"strings"
)

// ParsePath parses a jsonpath into a slice of path parts
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

func Visit(j any, path []string, on func(any)) {
	visit(nil, j, path, on, nil)
}

// Transform applies a transformation function to the value at the given path. The transformation function takes 3
// parameters:
//  1. the container object (will be either map[string]any or []any)
//  2. the key within the container (will be either string or int depending on the container type)
//  3. the value itself
//
// The transformation function should return the new value to be set at the same path.
func Transform(j any, path []string, tx func(any, any, any) any) {
	visit(nil, j, path, nil, tx)
}

func visit(container, j any, path []string, on func(any), tx func(any, any, any) any) {
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
						typed[k] = tx(j, k, val)
					} else {
						on(val)
					}
				} else {
					visit(j, val, rem, on, tx)
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
						typed[i] = tx(j, i, val)
					} else {
						on(val)
					}
				} else {
					visit(j, val, rem, on, tx)
				}
			}
		}
	}
}
