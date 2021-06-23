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
