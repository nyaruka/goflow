package flows

import (
	"encoding/json"
	"fmt"
	"github.com/nyaruka/gocommon/urns"
	"strings"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
)

// ChannelRole is a role that a channel can perform
type ChannelRole string

// different roles that channels can perform
const (
	ChannelRoleSend    ChannelRole = "send"
	ChannelRoleReceive ChannelRole = "receive"
	ChannelRoleCall    ChannelRole = "call"
	ChannelRoleAnswer  ChannelRole = "answer"
	ChannelRoleUSSD    ChannelRole = "ussd"
)

// Channel represents a means for sending and receiving input during a flow run. It renders as its name in a template,
// and has the following properties which can be accessed:
//
//  * `uuid` the UUID of the channel
//  * `name` the name of the channel
//  * `address` the address of the channel
//
// Examples:
//
//   @contact.channel -> My Android Phone
//   @contact.channel.name -> My Android Phone
//   @contact.channel.address -> +12345671111
//   @run.input.channel.uuid -> 57f1078f-88aa-46f4-a59a-948a5739c03d
//   @(json(contact.channel)) -> {"address":"+12345671111","name":"My Android Phone","uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d"}
//
// @context channel
type Channel interface {
	types.XValue
	types.XResolvable

	UUID() ChannelUUID
	Name() string
	Address() string
	Schemes() []string
	SupportsScheme(string) bool
	Roles() []ChannelRole
	HasRole(ChannelRole) bool
	Parent() *ChannelReference

	MatchTelCountry() string
	MatchTelPrefixes() []string
	SetTelMatching(string, []string)

	Reference() *ChannelReference
}

type channel struct {
	uuid    ChannelUUID
	name    string
	address string
	schemes []string
	roles   []ChannelRole
	parent  *ChannelReference

	matchTelCountry  string
	matchTelPrefixes []string
}

// NewChannel creates a new channel
func NewChannel(uuid ChannelUUID, name string, address string, schemes []string, roles []ChannelRole, parent *ChannelReference) Channel {
	return &channel{
		uuid:    uuid,
		name:    name,
		address: address,
		schemes: schemes,
		roles:   roles,
		parent:  parent,
	}
}

// UUID returns the UUID of this channel
func (c *channel) UUID() ChannelUUID { return c.uuid }

// Name returns the name of this channel
func (c *channel) Name() string { return c.name }

// Address returns the address of this channel
func (c *channel) Address() string { return c.address }

// Schemes returns the supported schemes of this channel
func (c *channel) Schemes() []string { return c.schemes }

// Roles returns the roles of this channel
func (c *channel) Roles() []ChannelRole { return c.roles }

// Reference returns a reference to this channel
func (c *channel) Reference() *ChannelReference { return NewChannelReference(c.uuid, c.name) }

// Parent returns a reference to this channel's parent (if any)
func (c *channel) Parent() *ChannelReference { return c.parent }

// SupportsScheme returns whether this channel supports the given URN scheme
func (c *channel) SupportsScheme(scheme string) bool {
	for _, s := range c.schemes {
		if s == scheme {
			return true
		}
	}
	return false
}

// HasRole returns whether this channel has the given role
func (c *channel) HasRole(role ChannelRole) bool {
	for _, r := range c.roles {
		if r == role {
			return true
		}
	}
	return false
}

// MatchTelCountry returns this channel's associated country code (if any)
func (c *channel) MatchTelCountry() string { return c.matchTelCountry }

// MatchTelPrefixes returns this channel's match prefixes values used for selecting a channel for a URN (if any)
func (c *channel) MatchTelPrefixes() []string { return c.matchTelPrefixes }

func (c *channel) SetTelMatching(country string, prefixes []string) {
	c.matchTelCountry = country
	c.matchTelPrefixes = prefixes
}

// Resolve resolves the given key when this channel is referenced in an expression
func (c *channel) Resolve(env utils.Environment, key string) types.XValue {
	switch key {
	case "uuid":
		return types.NewXText(string(c.uuid))
	case "name":
		return types.NewXText(c.name)
	case "address":
		return types.NewXText(c.address)
	}

	return types.NewXResolveError(c, key)
}

// Describe returns a representation of this type for error messages
func (c *channel) Describe() string { return "channel" }

// Reduce is called when this object needs to be reduced to a primitive
func (c *channel) Reduce(env utils.Environment) types.XPrimitive {
	return types.NewXText(c.name)
}

// ToXJSON is called when this type is passed to @(json(...))
func (c *channel) ToXJSON(env utils.Environment) types.XText {
	return types.ResolveKeys(env, c, "uuid", "name", "address").ToXJSON(env)
}

func (c *channel) String() string {
	return c.name
}

var _ Channel = (*channel)(nil)

// ChannelSet defines the unordered set of all channels for a session
type ChannelSet struct {
	channels       []Channel
	channelsByUUID map[ChannelUUID]Channel
}

// NewChannelSet creates a new channel set
func NewChannelSet(channels []Channel) *ChannelSet {
	s := &ChannelSet{channels: channels, channelsByUUID: make(map[ChannelUUID]Channel, len(channels))}
	for _, channel := range s.channels {
		s.channelsByUUID[channel.UUID()] = channel
	}
	return s
}

// GetForURN returns the best channel for the given URN
func (s *ChannelSet) GetForURN(urn *ContactURN, role ChannelRole) Channel {
	// if caller has told us which channel to use for this URN, use that
	if urn.Channel() != nil {
		return s.getDelegate(urn.Channel(), role)
	}

	// tel is a special case because we do number based matching
	if urn.Scheme() == urns.TelScheme {
		countryCode := utils.DeriveCountryFromTel(urn.Path())
		candidates := make([]Channel, 0)

		for _, ch := range s.channels {
			if ch.HasRole(role) && ch.SupportsScheme(urns.TelScheme) && (countryCode == "" || countryCode == ch.MatchTelCountry()) && ch.Parent() == nil {
				candidates = append(candidates, ch)
			}
		}

		var channel Channel
		if len(candidates) > 1 {
			// we don't have a channel for this contact yet, let's try to pick one from the same carrier
			// we need at least one digit to overlap to infer a channel
			contactNumber := strings.TrimPrefix(urn.URN.Path(), "+")
			prefix := 1
			for _, candidate := range candidates {
				candidatePrefixes := candidate.MatchTelPrefixes()
				if len(candidatePrefixes) == 0 {
					candidatePrefixes = []string{strings.TrimPrefix(candidate.Address(), "+")}
				}

				for _, chanPrefix := range candidatePrefixes {
					for idx := prefix; idx <= len(chanPrefix); idx++ {
						if idx >= prefix && chanPrefix[0:idx] == contactNumber[0:idx] {
							prefix = idx
							channel = candidate
						} else {
							break
						}
					}
				}
			}

		} else if len(candidates) == 1 {
			channel = candidates[0]
		}

		if channel != nil {
			return s.getDelegate(channel, role)
		}
	}

	return s.getForSchemeAndRole(urn.Scheme(), role)
}

func (s *ChannelSet) getForSchemeAndRole(scheme string, role ChannelRole) Channel {
	for _, ch := range s.channels {
		if ch.HasRole(ChannelRoleSend) && ch.SupportsScheme(scheme) {
			return s.getDelegate(ch, role)
		}
	}
	return nil
}

// looks for a delegate for the given channel and defaults to the channel itself
func (s *ChannelSet) getDelegate(channel Channel, role ChannelRole) Channel {
	for _, ch := range s.channels {
		if ch.Parent() != nil && ch.Parent().UUID == channel.UUID() && ch.HasRole(role) {
			return ch
		}
	}
	return channel
}

// FindByUUID finds the channel with the given UUID
func (s *ChannelSet) FindByUUID(uuid ChannelUUID) Channel {
	return s.channelsByUUID[uuid]
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type channelEnvelope struct {
	UUID    ChannelUUID       `json:"uuid" validate:"required,uuid"`
	Name    string            `json:"name"`
	Address string            `json:"address"`
	Schemes []string          `json:"schemes" validate:"min=1"`
	Roles   []ChannelRole     `json:"roles" validate:"min=1,dive,eq=send|eq=receive|eq=call|eq=answer|eq=ussd"`
	Parent  *ChannelReference `json:"parent" validate:"omitempty,dive"`

	MatchCountry  string   `json:"match_country"`
	MatchPrefixes []string `json:"match_prefixes"`
}

// ReadChannel decodes a channel from the passed in JSON
func ReadChannel(data json.RawMessage) (Channel, error) {
	ce := channelEnvelope{}
	if err := utils.UnmarshalAndValidate(data, &ce); err != nil {
		return nil, fmt.Errorf("unable to read channel: %s", err)
	}

	return &channel{
		uuid:             ce.UUID,
		name:             ce.Name,
		address:          ce.Address,
		schemes:          ce.Schemes,
		roles:            ce.Roles,
		parent:           ce.Parent,
		matchTelCountry:  ce.MatchCountry,
		matchTelPrefixes: ce.MatchPrefixes,
	}, nil
}

// ReadChannelSet decodes channels from the passed in JSON
func ReadChannelSet(data json.RawMessage) (*ChannelSet, error) {
	items, err := utils.UnmarshalArray(data)
	if err != nil {
		return nil, err
	}

	channels := make([]Channel, len(items))
	for c := range items {
		channels[c], err = ReadChannel(items[c])
		if err != nil {
			return nil, err
		}
	}
	return NewChannelSet(channels), nil
}
