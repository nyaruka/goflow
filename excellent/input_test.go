package excellent

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInput(t *testing.T) {
	input := newInput(bufio.NewReader(strings.NewReader("12")))

	assert.Equal(t, '1', input.read())

	input.unread('1')

	assert.Equal(t, '1', input.read())
	assert.Equal(t, '2', input.read())
	assert.Equal(t, eof, input.read())
	assert.Equal(t, eof, input.read())

	input = newInput(bufio.NewReader(strings.NewReader("😊")))
	assert.Equal(t, '😊', input.read())
	assert.Equal(t, eof, input.read())
}
