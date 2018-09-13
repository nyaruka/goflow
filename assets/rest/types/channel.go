package types

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
)

// json serializable implementation of a channel asset
type channel struct {
	UUID_          assets.ChannelUUID       `json:"uuid" validate:"required,uuid"`
	Name_          string                   `json:"name"`
	Address_       string                   `json:"address"`
	Schemes_       []string                 `json:"schemes" validate:"min=1"`
	Roles_         []assets.ChannelRole     `json:"roles" validate:"min=1,dive,eq=send|eq=receive|eq=call|eq=answer|eq=ussd"`
	Parent_        *assets.ChannelReference `json:"parent" validate:"omitempty,dive"`
	Country_       string                   `json:"country,omitempty"`
	MatchPrefixes_ []string                 `json:"match_prefixes,omitempty"`
}

// NewChannel creates a new channel
func NewChannel(uuid assets.ChannelUUID, name string, address string, schemes []string, roles []assets.ChannelRole, parent *assets.ChannelReference) assets.Channel {
	return &channel{
		UUID_:    uuid,
		Name_:    name,
		Address_: address,
		Schemes_: schemes,
		Roles_:   roles,
		Parent_:  parent,
	}
}

// NewTelChannel creates a new tel channel
func NewTelChannel(uuid assets.ChannelUUID, name string, address string, roles []assets.ChannelRole, parent *assets.ChannelReference, country string, matchPrefixes []string) assets.Channel {
	return &channel{
		UUID_:          uuid,
		Name_:          name,
		Address_:       address,
		Schemes_:       []string{urns.TelScheme},
		Roles_:         roles,
		Parent_:        parent,
		Country_:       country,
		MatchPrefixes_: matchPrefixes,
	}
}

// UUID returns the UUID of this channel
func (c *channel) UUID() assets.ChannelUUID { return c.UUID_ }

// Name returns the name of this channel
func (c *channel) Name() string { return c.Name_ }

// Address returns the address of this channel
func (c *channel) Address() string { return c.Address_ }

// Schemes returns the supported schemes of this channel
func (c *channel) Schemes() []string { return c.Schemes_ }

// Roles returns the roles of this channel
func (c *channel) Roles() []assets.ChannelRole { return c.Roles_ }

// Parent returns a reference to this channel's parent (if any)
func (c *channel) Parent() *assets.ChannelReference { return c.Parent_ }

// Country returns this channel's associated country code (if any)
func (c *channel) Country() string { return c.Country_ }

// MatchPrefixes returns this channel's match prefixes values used for selecting a channel for a URN (if any)
func (c *channel) MatchPrefixes() []string { return c.MatchPrefixes_ }

// ReadChannel reads a channel from the given JSON
func ReadChannel(data json.RawMessage) (assets.Channel, error) {
	c := &channel{}
	if err := utils.UnmarshalAndValidate(data, c); err != nil {
		return nil, fmt.Errorf("unable to read channel: %s", err)
	}
	return c, nil
}

// ReadChannels reads channels from the given JSON
func ReadChannels(data json.RawMessage) ([]assets.Channel, error) {
	items, err := utils.UnmarshalArray(data)
	if err != nil {
		return nil, err
	}

	channels := make([]assets.Channel, len(items))
	for d := range items {
		if channels[d], err = ReadChannel(items[d]); err != nil {
			return nil, err
		}
	}

	return channels, nil
}
