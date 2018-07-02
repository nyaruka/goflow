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
		assert.Equal(t, test.expected, utils.Snakify(test.input), "unexpected result snakifying '%s'", test.input)
	}
}

func TestURLEscape(t *testing.T) {
	var urlTests = []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"hello_world-there", "hello_world-there"},
		{"foo: bar ? & some/thing", "foo%3A%20bar%20%3F%20%26%20some%2Fthing"},
	}

	for _, test := range urlTests {
		assert.Equal(t, test.expected, utils.URLEscape(test.input), "unexpected result URL escaping '%s'", test.input)
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
