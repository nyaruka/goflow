package types_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets/static/types"

	"github.com/stretchr/testify/assert"
)

func TestGlobal(t *testing.T) {
	global := types.NewGlobal("org_name", "Org Name", "U-Report")
	assert.Equal(t, "org_name", global.Key())
	assert.Equal(t, "Org Name", global.Name())
	assert.Equal(t, "U-Report", global.Value())
}
