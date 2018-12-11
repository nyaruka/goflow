package mobile_test

import (
	"testing"

	"github.com/nyaruka/goflow/mobile"

	"github.com/stretchr/testify/assert"
)

func TestStringSlice(t *testing.T) {
	s := mobile.NewStringSlice(5)
	s.Add("Foo")
	s.Add("Bar")

	assert.Equal(t, 2, s.Length())
	assert.Equal(t, "Foo", s.Get(0))
	assert.Equal(t, "Bar", s.Get(1))
}
