package flows_test

import (
	"testing"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"

	"github.com/stretchr/testify/assert"
)

func TestLocationPaths(t *testing.T) {
	assert.True(t, flows.IsPossibleLocationPath("Ireland > Antrim"))
	assert.True(t, flows.IsPossibleLocationPath("Ireland>Antrim"))
	assert.False(t, flows.IsPossibleLocationPath("Antrim"))

	assert.Equal(t, "Antrim", flows.LocationPath("Ireland > Antrim").Name())
	assert.Equal(t, "Ireland", flows.LocationPath("Ireland").Name())
	assert.Equal(t, "", flows.LocationPath("").Name())
	assert.Equal(t, "Ireland > Antrim", flows.LocationPath("Ireland > Antrim").String())
	assert.Equal(t, types.NewXText(`"Ireland > Antrim"`), flows.LocationPath("Ireland > Antrim").ToXJSON())
}
