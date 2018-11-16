package flows_test

import (
	"testing"

	"github.com/nyaruka/goflow/flows"

	"github.com/stretchr/testify/assert"
)

func TestIsVersionSupported(t *testing.T) {
	assert.False(t, flows.IsVersionSupported("x"))
	assert.False(t, flows.IsVersionSupported("11.9"))
	assert.True(t, flows.IsVersionSupported("12"))
	assert.True(t, flows.IsVersionSupported("12.0"))
	assert.True(t, flows.IsVersionSupported("12.99"))
	assert.False(t, flows.IsVersionSupported("13"))
	assert.False(t, flows.IsVersionSupported("13.0"))
}
