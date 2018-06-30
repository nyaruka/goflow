package utils_test

import (
	"testing"

	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
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

func TestTokenizeString(t *testing.T) {
	tokenizerTests := []struct {
		text   string
		result []string
	}{
		{"one   two three", []string{"one", "two", "three"}},
		{"one.two.three", []string{"one", "two", "three"}},
		{"one.βήταa.three", []string{"one", "βήταa", "three"}},
		{"one😄three", []string{"one", "😄", "three"}},
		{"  one.two.*@three ", []string{"one", "two", "three"}},
		{" one ", []string{"one"}},
	}
	for _, test := range tokenizerTests {
		assert.Equal(t, test.result, utils.TokenizeString(test.text), "unexpected result tokenizing '%s'", test.text)
	}
}

func TestTokenizeStringByChars(t *testing.T) {
	tokenizerTests := []struct {
		text   string
		chars  string
		result []string
	}{
		{"one   two three", " ", []string{"one", "two", "three"}},
		{"Jim O'Grady", " ", []string{"Jim", "O'Grady"}},
		{"one.βήταa/three", "./", []string{"one", "βήταa", "three"}},
		{"one😄three", "😄", []string{"one", "three"}},
		{"  one.two.*@three ", " .*@", []string{"one", "two", "three"}},
		{" one ", " ", []string{"one"}},
	}
	for _, test := range tokenizerTests {
		assert.Equal(t, test.result, utils.TokenizeStringByChars(test.text, test.chars), "unexpected result tokenizing '%s'", test.text)
	}
}
