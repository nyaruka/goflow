package jsonpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePath(t *testing.T) {
	tcs := []struct {
		path     string
		expected []string
		err      string
	}{
		{"", []string{}, "path must begin with $"},
		{"$.foo", []string{"foo"}, ""},
		{"$[*]", []string{"*"}, ""},
		{"$[2]", []string{"2"}, ""},
		{"$[]", []string{"2"}, "subscript value can't be empty"},
		{"$.foo[*]", []string{"foo", "*"}, ""},
		{"$.foo[*].bar", []string{"foo", "*", "bar"}, ""},
		{"$.foo[*].bar[5]", []string{"foo", "*", "bar", "5"}, ""},
	}

	for _, tc := range tcs {
		actual, err := parsePath(tc.path)
		if tc.err != "" {
			assert.EqualError(t, err, tc.err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		}
	}
}

func TestVisit(t *testing.T) {
	data := map[string]any{
		"foo": "bar",
		"arr": []any{"123", "234", "345"},
		"objs": []any{
			map[string]any{"sub": "345"},
			map[string]any{"sub": "456", "alt": "789"},
		},
	}

	tcs := []struct {
		path     string
		expected []any
	}{
		{"$.foo", []any{"bar"}},
		{"$.arr[*]", []any{"123", "234", "345"}},
		{"$.arr[0]", []any{"123"}},
		{"$.arr[2]", []any{"345"}},
		{"$.objs[*]", []any{map[string]any{"sub": "345"}, map[string]any{"sub": "456", "alt": "789"}}},
		{"$.objs[*].sub", []any{"345", "456"}},
		{"$.objs[1].sub", []any{"456"}},
		{"$.objs[*].alt", []any{"789"}},
	}

	for _, tc := range tcs {
		var matches []any
		Visit(data, tc.path, func(m any) { matches = append(matches, m) })

		assert.Equal(t, tc.expected, matches)
	}
}

func TestTransform(t *testing.T) {
	tcs := []struct {
		data     any
		path     string
		repl     any
		expected any
	}{
		{[]any{"foo", "bar"}, "$[*]", "baz", []any{"baz", "baz"}},
		{map[string]any{"foo": "bar"}, "$.foo", "baz", map[string]any{"foo": "baz"}},
		{map[string]any{"foo": []any{"1", "2"}, "bar": []any{"3"}}, "$.foo[*]", "baz", map[string]any{"foo": []any{"baz", "baz"}, "bar": []any{"3"}}},
	}

	for _, tc := range tcs {
		Transform(tc.data, tc.path, func(c, k, m any) any { return tc.repl })

		assert.Equal(t, tc.expected, tc.data)
	}
}
