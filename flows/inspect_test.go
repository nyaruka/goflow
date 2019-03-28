package flows_test

import (
	"testing"

	"github.com/nyaruka/goflow/flows"

	"github.com/stretchr/testify/assert"
)

func TestMergeResultSpecs(t *testing.T) {
	assert.Equal(t, []*flows.ResultSpec{}, flows.MergeResultSpecs(nil))

	assert.Equal(t, []*flows.ResultSpec{
		{Key: "response_1", Name: "Response 1", Categories: []string{"Red", "Green", "Blue"}},
		{Key: "favorite_beer", Name: "Favorite Beer", Categories: []string{}},
	}, flows.MergeResultSpecs([]*flows.ResultSpec{
		flows.NewResultSpec("Response 1", []string{"Red", "Green"}),
		flows.NewResultSpec("Response-1", nil),
		flows.NewResultSpec("Response-1", []string{"Green", "Blue"}),
		flows.NewResultSpec("Favorite Beer", []string{}),
	}))
}
