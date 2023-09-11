package static

import (
	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
)

// Channel is a JSON serializable implementation of a channel asset
type Channel struct {
	UUID_               assets.ChannelUUID   `json:"uuid" validate:"required,uuid"`
	Name_               string               `json:"name"`
	Address_            string               `json:"address"`
	Schemes_            []string             `json:"schemes" validate:"min=1"`
	Roles_              []assets.ChannelRole `json:"roles" validate:"min=1,dive,eq=send|eq=receive|eq=call|eq=answer|eq=ussd"`
	Country_            i18n.Country         `json:"country,omitempty"`
	MatchPrefixes_      []string             `json:"match_prefixes,omitempty"`
	AllowInternational_ bool                 `json:"allow_international,omitempty"`
}

// NewChannel creates a new channel
func NewChannel(uuid assets.ChannelUUID, name string, address string, schemes []string, roles []assets.ChannelRole, parent *assets.ChannelReference) assets.Channel {
	return &Channel{
		UUID_:               uuid,
		Name_:               name,
		Address_:            address,
		Schemes_:            schemes,
		Roles_:              roles,
		AllowInternational_: true,
	}
}

// NewTelChannel creates a new tel channel
func NewTelChannel(uuid assets.ChannelUUID, name string, address string, roles []assets.ChannelRole, parent *assets.ChannelReference, country i18n.Country, matchPrefixes []string, allowInternational bool) assets.Channel {
	return &Channel{
		UUID_:               uuid,
		Name_:               name,
		Address_:            address,
		Schemes_:            []string{urns.TelScheme},
		Roles_:              roles,
		Country_:            country,
		MatchPrefixes_:      matchPrefixes,
		AllowInternational_: allowInternational,
	}
}

// UUID returns the UUID of this channel
func (c *Channel) UUID() assets.ChannelUUID { return c.UUID_ }

// Name returns the name of this channel
func (c *Channel) Name() string { return c.Name_ }

// Address returns the address of this channel
func (c *Channel) Address() string { return c.Address_ }

// Schemes returns the supported schemes of this channel
func (c *Channel) Schemes() []string { return c.Schemes_ }

// Roles returns the roles of this channel
func (c *Channel) Roles() []assets.ChannelRole { return c.Roles_ }

// Country returns this channel's associated country code (if any)
func (c *Channel) Country() i18n.Country { return c.Country_ }

// MatchPrefixes returns this channel's match prefixes values used for selecting a channel for a URN (if any)
func (c *Channel) MatchPrefixes() []string { return c.MatchPrefixes_ }

// AllowInternational returns whether this channel allows sending internationally (only applies to TEL schemes)
func (c *Channel) AllowInternational() bool { return c.AllowInternational_ }
