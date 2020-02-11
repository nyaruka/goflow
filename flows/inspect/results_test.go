package inspect_test

import (
	"testing"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/inspect"

	"github.com/stretchr/testify/assert"
)

func TestResults(t *testing.T) {
	n := definition.NewNode(flows.NodeUUID("866b06e2-ff54-443e-9d79-2f60074514b5"), []flows.Action{
		actions.NewSetContactName(flows.ActionUUID("52a91ae8-1115-4c17-99a2-58b15ed7de7f"), "Bob"),
		actions.NewSetRunResult(flows.ActionUUID("94790ebc-4f24-4664-a15d-ac758781c720"), "Age", "32", "HasAge"),
	}, nil, []flows.Exit{})

	infos := make([]*flows.ResultInfo, 0)
	inspect.Results(n.Actions(), func(r *flows.ResultInfo) {
		infos = append(infos, r)
	})

	assert.Equal(t, []*flows.ResultInfo{
		flows.NewResultInfo("Age", []string{"HasAge"}),
	}, infos)
}
