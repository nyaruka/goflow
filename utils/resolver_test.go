package utils_test

import (
	"testing"

	"github.com/nyaruka/goflow/utils"
)

func TestSnakify(t *testing.T) {
	var snakeTests = []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello_world"},
		{"hello_world", "hello_world"},
		{"hello-world", "hello_world"},
		{"hi😀😃😄😁there", "hi_there"},
		{"昨夜のコ", "昨夜のコ"},
		{"this@isn't@email", "this_isn_t_email"},
	}

	for _, test := range snakeTests {
		value := utils.Snakify(test.input)

		if value != test.expected {
			t.Errorf("Expected: '%s' Got: '%s' for input: '%s'", test.expected, value, test.input)
		}
	}
}
