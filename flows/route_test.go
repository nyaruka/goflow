package flows_test

import (
	"testing"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestRoute(t *testing.T) {
	exitUUID := flows.ExitUUID(utils.NewUUID())
	r := flows.NewRoute(exitUUID, "red", map[string]string{"foo": "bar"})
	assert.Equal(t, exitUUID, r.Exit())
	assert.Equal(t, "red", r.Match())
	assert.Equal(t, map[string]string{"foo": "bar"}, r.MatchExtra())

	assert.Equal(t, flows.ExitUUID(""), flows.NoRoute.Exit())
	assert.Equal(t, "", flows.NoRoute.Match())
}
