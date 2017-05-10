package utils

import (
	"reflect"
	"testing"
)

var tokenizerTests = []struct {
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

func TestTokenizer(t *testing.T) {
	for _, test := range tokenizerTests {
		result := TokenizeString(test.text)
		if !reflect.DeepEqual(result, test.result) {
			t.Errorf("Unexpected result tokenizing '%s', got: %s expected: %v", test.text, result, test.result)
		}
	}
}
