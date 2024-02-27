package jsonpath

// Very barebones jsonpath implementation that only supports the following syntax:
//
//  - $.foo
//  - $.foo.*
//  - $.foo[0]
//  - $.foo[*]
//  - $.foo[*].bar
//
// This exists because goflow has some very specific requirements for making transformations in JSON flow definitions
// knowing the parent container of the thing being transformed, and the name of thing in the container.

import (
	"errors"
	"strconv"
	"strings"
)

// parses a jsonpath into a slice of path parts
func parsePath(path string) ([]string, error) {
	runes := []rune(path)
	steps := make([]string, 0, 5)

	if len(runes) == 0 || runes[0] != '$' {
		return nil, errors.New("path must begin with $")
	}
	i := 1

	for {
		if i == len(runes) {
			break
		}

		if runes[i] == '.' {
			i++

			var b strings.Builder
			for i < len(runes) && runes[i] != '.' && runes[i] != '[' {
				b.WriteRune(runes[i])
				i++
			}
			steps = append(steps, b.String())
		} else if runes[i] == '[' {
			i++

			var b strings.Builder
			for i < len(runes) && runes[i] != ']' {
				b.WriteRune(runes[i])
				i++
			}
			if i < len(runes) && runes[i] == ']' {
				i++
			}
			s := b.String()
			if len(s) == 0 {
				return nil, errors.New("subscript value can't be empty")
			}
			steps = append(steps, s)
		}
	}

	return steps, nil
}

func Visit(j any, path string, on func(any)) error {
	p, err := parsePath(path)
	if err != nil {
		return err
	}
	visit(nil, j, p, on, nil)
	return nil
}

// Transform applies a transformation function to the value at the given path. The transformation function takes 3
// parameters:
//  1. the container object (will be either map[string]any or []any)
//  2. the key within the container (will be either string or int depending on the container type)
//  3. the value itself
//
// The transformation function should return the new value to be set at the same path.
func Transform(j any, path string, tx func(any, any, any) any) error {
	p, err := parsePath(path)
	if err != nil {
		return err
	}
	visit(nil, j, p, nil, tx)
	return nil
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
