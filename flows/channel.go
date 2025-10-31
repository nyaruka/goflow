package flows

import (
	"fmt"
	"slices"
	"strings"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
)

// Channel represents a means for sending and receiving input during a flow run
type Channel struct {
	assets.Channel
}

// NewChannel creates a new channenl
func NewChannel(asset assets.Channel) *Channel {
	return &Channel{Channel: asset}
}

// Asset returns the underlying asset
func (c *Channel) Asset() assets.Channel { return c.Channel }

// Reference returns a reference to this channel
func (c *Channel) Reference() *assets.ChannelReference {
	if c == nil {
		return nil
	}
	return assets.NewChannelReference(c.UUID(), c.Name())
}

// SupportsScheme returns whether this channel supports the given URN scheme
func (c *Channel) SupportsScheme(scheme string) bool {
	return slices.Contains(c.Schemes(), scheme)
}

// HasRole returns whether this channel has the given role
func (c *Channel) HasRole(role assets.ChannelRole) bool {
	return slices.Contains(c.Roles(), role)
}

// HasFeature returns whether this channel has the given feature
func (c *Channel) HasFeature(feat assets.ChannelFeature) bool {
	return slices.Contains(c.Features(), feat)
}

// Context returns the properties available in expressions
//
//	__default__:text -> the name
//	uuid:text -> the UUID of the channel
//	name:text -> the name of the channel
//	address:text -> the address of the channel
//
// @context channel
func (c *Channel) Context(env envs.Environment) map[string]types.XValue {
	return map[string]types.XValue{
		"__default__": types.NewXText(c.Name()),
		"uuid":        types.NewXText(string(c.UUID())),
		"name":        types.NewXText(c.Name()),
		"address":     types.NewXText(c.Address()),
	}
}

func (c *Channel) String() string {
	return fmt.Sprintf("%s (%s)", c.Address(), c.Name())
}

// ChannelAssets provides access to all channel assets
type ChannelAssets struct {
	all    []*Channel
	byUUID map[assets.ChannelUUID]*Channel
}

// NewChannelAssets creates a new set of channel assets
func NewChannelAssets(channels []assets.Channel) *ChannelAssets {
	s := &ChannelAssets{
		all:    make([]*Channel, len(channels)),
		byUUID: make(map[assets.ChannelUUID]*Channel, len(channels)),
	}
	for i, asset := range channels {
		channel := NewChannel(asset)
		s.all[i] = channel
		s.byUUID[channel.UUID()] = channel
	}
	return s
}

// Get returns the channel with the given UUID
func (s *ChannelAssets) Get(uuid assets.ChannelUUID) *Channel {
	return s.byUUID[uuid]
}

// GetForURN returns the best channel for the given URN
func (s *ChannelAssets) GetForURN(urn *URN, role assets.ChannelRole) *Channel {
	// if caller has told us which channel to use for this URN, use that
	if urn.Channel != nil && urn.Channel.HasRole(role) {
		return urn.Channel
	}

	// tel is a special case because we do number based matching
	if urn.Scheme == urns.Phone.Prefix {
		countryCode := i18n.DeriveCountryFromTel(urn.Path)
		candidates := make([]*Channel, 0)

		for _, ch := range s.all {
			// skip if not tel and not sendable
			if !ch.SupportsScheme(urns.Phone.Prefix) || !ch.HasRole(role) {
				continue
			}
			// skip if international and channel doesn't allow that
			if ch.Country() != "" && countryCode != "" && countryCode != ch.Country() && !ch.AllowInternational() {
				continue
			}

			candidates = append(candidates, ch)
		}

		var channel *Channel
		if len(candidates) > 1 {
			// we don't have a channel for this contact yet, let's try to pick one from the same carrier
			// we need at least one digit to overlap to infer a channel
			contactNumber := strings.TrimPrefix(urn.Path, "+")
			maxOverlap := 0
			for _, candidate := range candidates {
				candidatePrefixes := candidate.MatchPrefixes()
				if len(candidatePrefixes) == 0 {
					candidatePrefixes = []string{strings.TrimPrefix(candidate.Address(), "+")}
				}

				for _, prefix := range candidatePrefixes {
					overlap := utils.PrefixOverlap(prefix, contactNumber)
					if overlap >= maxOverlap {
						maxOverlap = overlap
						channel = candidate
					}
				}
			}

		} else if len(candidates) == 1 {
			channel = candidates[0]
		}

		if channel != nil {
			return channel
		}

		return nil
	}

	return s.getForSchemeAndRole(urn.Scheme, role)
}

func (s *ChannelAssets) getForSchemeAndRole(scheme string, role assets.ChannelRole) *Channel {
	for _, ch := range s.all {
		if ch.HasRole(role) && ch.SupportsScheme(scheme) {
			return ch
		}
	}
	return nil
}
