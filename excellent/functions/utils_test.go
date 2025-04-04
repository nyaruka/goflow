package functions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractWords(t *testing.T) {
	assert.Equal(t, []string(nil), extractWords("", ""))
	assert.Equal(t, []string{"foo"}, extractWords("foo", ""))
	assert.Equal(t, []string{"foo", "bar"}, extractWords("foo  bar  ", ""))
	assert.Equal(t, []string{"foo", "foo"}, extractWords("foo foo", ""))
	assert.Equal(t, []string{"foo.bar", "zed | doh"}, extractWords("foo.bar$zed | doh", "$"))
	assert.Equal(t, []string{"foo.bar", "zed ", " doh"}, extractWords("foo.bar$zed | doh", "$|"))
	assert.Equal(t, []string{"foo", "bar", "zed", "doh"}, extractWords("foo.bar$zed | doh", "$| ."))
	assert.Equal(t, []string{"😁", "😃"}, extractWords("😁🎟️😃", "🎟️"))
}

func TestSlice(t *testing.T) {
	assert.Equal(t, []int{1, 2}, slice([]int{1, 2, 3, 4}, 0, 2))
	assert.Equal(t, []int{1, 2, 3, 4}, slice([]int{1, 2, 3, 4}, 0, 4))
	assert.Equal(t, []int{1, 2, 3, 4}, slice([]int{1, 2, 3, 4}, 0, 10)) // end can exceed length
	assert.Equal(t, []int{2, 3}, slice([]int{1, 2, 3, 4}, 1, 3))
	assert.Equal(t, []int{2, 3}, slice([]int{1, 2, 3, 4}, -3, -1))
	assert.Equal(t, []int{1, 2, 3}, slice([]int{1, 2, 3, 4}, -10, -1))
	assert.Equal(t, []int{}, slice([]int{1, 2, 3, 4}, 0, 0))
	assert.Equal(t, []int{}, slice([]int{1, 2, 3, 4}, 3, 2)) // end is greater than start
}
