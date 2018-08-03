package flows_test

import (
	"testing"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestRoute(t *testing.T) {
	exitUUID := flows.ExitUUID(utils.NewUUID())
	r := flows.NewRoute(exitUUID, "red")
	assert.Equal(t, exitUUID, r.Exit())
	assert.Equal(t, "red", r.Match())

	assert.Equal(t, flows.ExitUUID(""), flows.NoRoute.Exit())
	assert.Equal(t, "", flows.NoRoute.Match())
}
