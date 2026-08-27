package test

import (
	"os"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/core"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition/migrations"
	"github.com/nyaruka/goflow/flows/engine"
)

// LoadSessionAssets loads a session assets instance from a static JSON file
func LoadSessionAssets(env envs.Environment, path string) (flows.SessionAssets, error) {
	assetsJSON, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	source, err := static.NewSource(assetsJSON)
	if err != nil {
		return nil, err
	}

	mconfig := &migrations.Config{BaseMediaURL: "http://temba.io/"}

	return engine.NewSessionAssets(env, source, mconfig)
}

func LoadFlowFromAssets(env envs.Environment, path string, uuid assets.FlowUUID) (flows.Flow, error) {
	sa, err := LoadSessionAssets(env, path)
	if err != nil {
		return nil, err
	}

	return sa.Flows().Get(uuid)
}

// Test asset files should use these channels rather than inventing new ones, and reuse an existing definition
// verbatim - same UUID, name, address, schemes and roles - whenever one of them is needed:
//
//	57f1078f-88aa-46f4-a59a-948a5739c03d  Android Channel   +17035550111   tel        send, receive (country US)
//	3a05eaf5-cb1b-4246-bef1-f277419c83a7  Vonage Channel    +17035550112   tel        send, receive
//	a78930fe-6a40-4aa8-99c3-e61b02f45ca1  Twilio Channel    +17035550113   tel        send, receive, call, answer
//	8e21f093-99aa-413b-b55b-758b54308fcb  Telegram Channel  345765375445   telegram   send, receive
//	eb9fee95-d762-4679-a7d5-91532e400c54  Receive Only      56789          ext        receive
//
// Telegram is the preferred non-tel type. Channels created directly with the functions below are for exercising
// channel selection logic and can use whatever addresses and countries that requires.
func NewChannel(name string, address string, schemes []string, roles []assets.ChannelRole) *core.Channel {
	return core.NewChannel(static.NewChannel(assets.ChannelUUID(uuids.NewV4()), name, address, schemes, roles))
}

func NewTelChannel(name string, address string, roles []assets.ChannelRole, parent *assets.ChannelReference, country i18n.Country, matchPrefixes []string, allowInternational bool) *core.Channel {
	return core.NewChannel(static.NewTelChannel(assets.ChannelUUID(uuids.NewV4()), name, address, roles, parent, country, matchPrefixes, allowInternational))
}
