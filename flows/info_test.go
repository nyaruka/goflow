package flows_test

import (
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
)

func TestNewResultSpecs(t *testing.T) {
	assert.Equal(t, []*flows.ResultSpec{}, flows.NewResultSpecs(nil))

	node1 := definition.NewNode(
		flows.NodeUUID("1fb823c3-599a-41e9-b59b-658266af3466"),
		nil,
		nil,
		[]flows.Exit{definition.NewExit(flows.ExitUUID("3c158842-24f3-4a40-bea4-7522952c0131"), "")},
	)
	node2 := definition.NewNode(
		flows.NodeUUID("0ba673a3-63b3-46f9-9246-9c727cf2917f"),
		nil,
		nil,
		[]flows.Exit{definition.NewExit(flows.ExitUUID("434ac29c-abe6-4bd7-b29b-740d517b6bb5"), "")},
	)

	extracted := []flows.ExtractedResult{
		{Node: node1, Info: flows.NewResultInfo("Response 1", []string{"Red", "Green"})},
		{Node: node1, Info: flows.NewResultInfo("Response-1", nil)},
		{Node: node2, Info: flows.NewResultInfo("Response-1", []string{"Green", "Blue"})},
		{Node: node2, Info: flows.NewResultInfo("Favorite Beer", []string{})},
	}

	specs := flows.NewResultSpecs(extracted)
	specsJSON := jsonx.MustMarshal(specs)

	test.AssertEqualJSON(t, []byte(`[
		{
			"key": "response_1",
			"name": "Response 1",
			"categories": [
				"Red",
				"Green",
				"Blue"
			],
			"node_uuids": [
				"1fb823c3-599a-41e9-b59b-658266af3466",
				"0ba673a3-63b3-46f9-9246-9c727cf2917f"
			]
		},
		{
			"key": "favorite_beer",
			"name": "Favorite Beer",
			"categories": [],
			"node_uuids": [
				"0ba673a3-63b3-46f9-9246-9c727cf2917f"
			]
		}
	]`), specsJSON, "result specs JSON mismatch")

	assert.Equal(t, `key=response_1|name=Response 1|categories=Red,Green`, flows.NewResultInfo("Response 1", []string{"Red", "Green"}).String())
}
