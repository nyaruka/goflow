package flows_test

import (
	"testing"

	"github.com/nyaruka/goflow/flows"

	"github.com/stretchr/testify/assert"
)

func TestMergeResultInfos(t *testing.T) {
	assert.Equal(t, []*flows.ResultInfo{}, flows.MergeResultInfos(nil))

	assert.Equal(t, []*flows.ResultInfo{
		{Key: "response_1", Name: "Response 1", Categories: []string{"Red", "Green", "Blue"}},
		{Key: "favorite_beer", Name: "Favorite Beer", Categories: []string{}},
	}, flows.MergeResultInfos([]*flows.ResultInfo{
		flows.NewResultInfo("Response 1", []string{"Red", "Green"}),
		flows.NewResultInfo("Response-1", nil),
		flows.NewResultInfo("Response-1", []string{"Green", "Blue"}),
		flows.NewResultInfo("Favorite Beer", []string{}),
	}))
}
