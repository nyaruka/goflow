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
		{"_Hello World", "_hello_world"},
		{"   Hello World    ", "hello_world"},
		{"hi😀😃😄😁there", "hi_there"},
		{"昨夜のコ", "昨夜のコ"},
		{"this@isn't@email", "this_isn_t_email"},
	}

	for _, test := range snakeTests {
		assert.Equal(t, test.expected, utils.Snakify(test.input), "unexpected result snakifying '%s'", test.input)
	}
}

func TestTokenizeString(t *testing.T) {
	tokenizerTests := []struct {
		text   string
		result []string
	}{
		{" one ", []string{"one"}},
		{"one   two three", []string{"one", "two", "three"}},
		{"one.two.three", []string{"one", "two", "three"}},
		{"O'Grady can't foo_bar", []string{"O'Grady", "can't", "foo_bar"}}, // single quotes and underscores don't split tokens
		{"öne.βήταa.thé", []string{"öne", "βήταa", "thé"}},                 // non-latin letters allowed in tokens
		{"واحد اثنين ثلاثة", []string{"واحد", "اثنين", "ثلاثة"}},           // RTL scripts
		{"  \t\none(two!*@three ", []string{"one", "two", "three"}},        // other punctuation ignored
		{"spend$£€₠₣₪", []string{"spend", "$", "£", "€", "₠", "₣", "₪"}},   // currency symbols treated as individual tokens
		{"math+=×÷√∊", []string{"math", "+", "=", "×", "÷", "√", "∊"}},     // math symbols treated as individual tokens
		{"emoji😄🏥👪👰😟🧟", []string{"emoji", "😄", "🏥", "👪", "👰", "😟", "🧟"}},   // emojis treated as individual tokens
		{"👍🏿 👨🏼", []string{"👍", "🏿", "👨", "🏼"}},                            // tone modifiers treated as individual tokens
		{"ℹ ℹ️", []string{"ℹ", "ℹ️"}},                                      // variation selectors ignored
		{"ยกเลิก sasa", []string{"ยกเลิก", "sasa"}},                        // Thai word means Cancelled
		{"বাতিল sasa", []string{"বাতিল", "sasa"}},                          // Bangla word means Cancel
		{"ထွက်သွား sasa", []string{"ထွက်သွား", "sasa"}},                    // Burmese word means exit
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

func TestPrefixOverlap(t *testing.T) {
	assert.Equal(t, 0, utils.PrefixOverlap("", ""))
	assert.Equal(t, 0, utils.PrefixOverlap("abc", ""))
	assert.Equal(t, 0, utils.PrefixOverlap("", "abc"))
	assert.Equal(t, 0, utils.PrefixOverlap("a", "x"))
	assert.Equal(t, 1, utils.PrefixOverlap("x", "x"))
	assert.Equal(t, 2, utils.PrefixOverlap("xya", "xyz"))
	assert.Equal(t, 2, utils.PrefixOverlap("😄😟👨🏼", "😄😟👰"))
	assert.Equal(t, 4, utils.PrefixOverlap("25078", "25073254252"))
}

func TestStringSlices(t *testing.T) {
	assert.Equal(t, []string{"he", "hello", "world"}, utils.StringSlices("hello world", []int{0, 2, 0, 5, 6, 11}))
}

func TestStringSliceContains(t *testing.T) {
	assert.False(t, utils.StringSliceContains(nil, "a", true))
	assert.False(t, utils.StringSliceContains([]string{}, "a", true))
	assert.False(t, utils.StringSliceContains([]string{"b", "c"}, "a", true))
	assert.True(t, utils.StringSliceContains([]string{"b", "a", "c"}, "a", true))
	assert.False(t, utils.StringSliceContains([]string{"b", "a", "c"}, "A", true))
	assert.True(t, utils.StringSliceContains([]string{"b", "a", "c"}, "A", false))
}

func TestIndent(t *testing.T) {
	assert.Equal(t, "", utils.Indent("", "  "))
	assert.Equal(t, "  x", utils.Indent("x", "  "))
	assert.Equal(t, "  x\n  y", utils.Indent("x\ny", "  "))
	assert.Equal(t, "  x\n\n  y", utils.Indent("x\n\ny", "  "))
	assert.Equal(t, ">>>x", utils.Indent("x", ">>>"))
}

func TestStringSet(t *testing.T) {
	set := utils.NewStringSet(0)
	set.Add("x")
	set.Add("x")
	assert.True(t, set.Contains("x"))
	assert.False(t, set.Contains("y"))
	assert.Equal(t, []string{"x"}, set.Slice())
}
