package static_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets/static"
	"github.com/stretchr/testify/assert"
)

func TestGlobal(t *testing.T) {
	global := static.NewGlobal("org_name", "Org Name", "U-Report")
	assert.Equal(t, "org_name", global.Key())
	assert.Equal(t, "Org Name", global.Name())
	assert.Equal(t, "U-Report", global.Value())
}
