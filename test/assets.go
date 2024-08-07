package test

import (
	"os"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
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

func NewChannel(name string, address string, schemes []string, roles []assets.ChannelRole, features []assets.ChannelFeature) *flows.Channel {
	return flows.NewChannel(static.NewChannel(assets.ChannelUUID(uuids.NewV4()), name, address, schemes, roles, features))
}

func NewTelChannel(name string, address string, roles []assets.ChannelRole, parent *assets.ChannelReference, country i18n.Country, matchPrefixes []string, allowInternational bool) *flows.Channel {
	return flows.NewChannel(static.NewTelChannel(assets.ChannelUUID(uuids.NewV4()), name, address, roles, parent, country, matchPrefixes, allowInternational))
}

func NewClassifier(name, type_ string, intents []string) *flows.Classifier {
	return flows.NewClassifier(static.NewClassifier(assets.ClassifierUUID(uuids.NewV4()), name, type_, intents))
}
