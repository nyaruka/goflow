package engine

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/routers"
	"io/ioutil"
	"testing"
	"time"

	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/stretchr/testify/assert"
)

func TestSession(t *testing.T) {
	assetURLs := map[AssetItemType]string{
		"channel": "http://testserver/assets/channel",
		"field":   "http://testserver/assets/field",
		"flow":    "http://testserver/assets/flow",
		"group":   "http://testserver/assets/group",
		"label":   "http://testserver/assets/label",
	}

	assetsJSON, err := ioutil.ReadFile("testdata/assets.json")
	assert.NoError(t, err)

	// build our session
	assetCache := NewAssetCache(100, 5)
	err = assetCache.Include(assetsJSON)
	assert.NoError(t, err)

	session := NewSession(assetCache, assetURLs)

	flow, err := session.Assets().GetFlow("76f0a02f-3b75-4b86-9064-e9195e1b3a02")
	assert.NoError(t, err)

	// break the add label action so references an invalid label
	addLabelAction := flow.Nodes()[0].Actions()[0].(*actions.AddLabelAction)
	addLabelAction.Labels[0].UUID = "xyx"

	// check that start fails with validation error
	err = session.Start(triggers.NewManualTrigger(flow, time.Now()), []flows.Event{})
	assert.EqualError(t, err, "validation failed for flow[uuid=76f0a02f-3b75-4b86-9064-e9195e1b3a02]: validation failed for action[uuid=ad154980-7bf7-4ab8-8728-545fd6378912, type=add_label]: no such label with uuid 'xyx'")

	// fix the add_label action
	addLabelAction.Labels[0].UUID = "3f65d88a-95dc-4140-9451-943e94e06fea"

	// and break the router so it references an invalid exit
	switchRouter := flow.Nodes()[0].Router().(*routers.SwitchRouter)
	switchRouter.Cases[0].ExitUUID = "xyx"

	// check that start fails with validation error
	err = session.Start(triggers.NewManualTrigger(flow, time.Now()), []flows.Event{})
	assert.EqualError(t, err, "validation failed for flow[uuid=76f0a02f-3b75-4b86-9064-e9195e1b3a02]: validation of router failed for node[uuid=a58be63b-907d-4a1a-856b-0bb5579d7507]: exit 'xyx' missing from node")

	// fix the router
	switchRouter.Cases[0].ExitUUID = "37d8813f-1402-4ad2-9cc2-e9054a96525b"

	// and break an exit so it references an invalid destination
	exit := flow.Nodes()[0].Exits()[0]
	flow.Nodes()[0].Exits()[0] = definition.NewExit(exit.UUID(), flows.NodeUUID("xyz"), "")

	// check that start fails with validation error
	err = session.Start(triggers.NewManualTrigger(flow, time.Now()), []flows.Event{})
	assert.EqualError(t, err, "validation failed for flow[uuid=76f0a02f-3b75-4b86-9064-e9195e1b3a02]: validation failed for exit[uuid=37d8813f-1402-4ad2-9cc2-e9054a96525b]: no such destination xyz")

	// fix the exit
	flow.Nodes()[0].Exits()[0] = exit
}
