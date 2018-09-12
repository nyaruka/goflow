package test

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/server/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func NewGroup(name string, query string) *flows.Group {
	return flows.NewGroup(types.NewGroup(assets.GroupUUID(utils.NewUUID()), name, query))
}

func NewChannel(name string, address string, schemes []string, roles []assets.ChannelRole, parentUUID assets.ChannelUUID) *flows.Channel {
	return flows.NewChannel(types.NewChannel(assets.ChannelUUID(utils.NewUUID()), name, address, schemes, roles, parentUUID))
}

func NewTelChannel(name string, address string, roles []assets.ChannelRole, parentUUID assets.ChannelUUID, country string, matchPrefixes []string) *flows.Channel {
	return flows.NewChannel(types.NewTelChannel(assets.ChannelUUID(utils.NewUUID()), name, address, roles, parentUUID, country, matchPrefixes))
}
