package jsonpath_test

import (
	"testing"

	"github.com/nyaruka/goflow/flows/definition/migrations/jsonpath"
	"github.com/stretchr/testify/assert"
)

func TestParsePath(t *testing.T) {
	assert.Equal(t, []string{}, jsonpath.ParsePath(""))
	assert.Equal(t, []string{"foo"}, jsonpath.ParsePath("$.foo"))
	assert.Equal(t, []string{"*"}, jsonpath.ParsePath("$[*]"))
	assert.Equal(t, []string{"foo", "*"}, jsonpath.ParsePath("$.foo[*]"))
	assert.Equal(t, []string{"foo", "*", "bar"}, jsonpath.ParsePath("$.foo[*].bar"))
}

func TestVisit(t *testing.T) {
	data := map[string]any{
		"foo": "bar",
		"arr": []any{"123", "234"},
		"obj_arr": []any{
			map[string]any{"sub": "345"},
			map[string]any{"sub": "456", "alt": "789"},
		},
	}

	tcs := []struct {
		path     string
		expected []any
	}{
		{"$.foo", []any{"bar"}},
		{"$.arr[*]", []any{"123", "234"}},
		{"$.obj_arr[*]", []any{map[string]any{"sub": "345"}, map[string]any{"sub": "456", "alt": "789"}}},
		{"$.obj_arr[*].sub", []any{"345", "456"}},
		{"$.obj_arr[*].alt", []any{"789"}},
	}

	for _, tc := range tcs {
		var matches []any
		jsonpath.Visit(data, jsonpath.ParsePath(tc.path), func(m any) { matches = append(matches, m) })

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
		jsonpath.Transform(tc.data, jsonpath.ParsePath(tc.path), func(c, k, m any) any { return tc.repl })

		assert.Equal(t, tc.expected, tc.data)
	}
}
