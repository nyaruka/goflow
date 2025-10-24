package flows

import (
	"fmt"
	"net/url"
	"slices"

	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
)

var redacted = "********"

func init() {
	utils.RegisterValidatorTag("urn", ValidateURN, func(validator.FieldError) string {
		return "is not a valid URN"
	})
	utils.RegisterValidatorTag("urnscheme", ValidateURNScheme, func(validator.FieldError) string {
		return "is not a valid URN scheme"
	})
}

// ValidateURN validates whether the field value is a valid URN
func ValidateURN(fl validator.FieldLevel) bool {
	err := urns.URN(fl.Field().String()).Validate()
	return err == nil
}

// ValidateURNScheme validates whether the field value is a valid URN scheme
func ValidateURNScheme(fl validator.FieldLevel) bool {
	return urns.IsValidScheme(fl.Field().String())
}

// Route represents a URN and a channel. It can model a URN and its channel affinity, or a destination for a message.
type Route struct {
	urn     urns.URN
	channel *Channel
}

func NewRoute(urn urns.URN, channel *Channel) *Route {
	return &Route{urn: urn, channel: channel}
}

// URN gets the underlying URN
func (r *Route) URN() urns.URN { return r.urn }

// Channel gets the channel associated with this URN
func (r *Route) Channel() *Channel { return r.channel }

// SetChannel sets the channel associated with this URN
func (r *Route) SetChannel(channel *Channel) {
	r.channel = channel
}

// Encode encodes the identity URN and the channel as a URN with query parameters
func (r *Route) Encode() urns.URN {
	scheme, path, query, display := r.urn.ToParts()

	parsedQuery, _ := url.ParseQuery(query)

	if r.channel != nil {
		parsedQuery.Set("channel", string(r.channel.UUID()))
	} else {
		parsedQuery.Del("channel")
	}

	urn, _ := urns.NewFromParts(scheme, path, parsedQuery, display)
	return urn
}

// Equal determines if this route is equal to another
func (r *Route) Equal(other *Route) bool {
	return other != nil && r.urn == other.urn && r.channel == other.channel
}

// ToXValue returns a representation of this object for use in expressions
func (r *Route) ToXValue(env envs.Environment) types.XValue {
	if env.RedactionPolicy() == envs.RedactionPolicyURNs {
		scheme, _, _, _ := r.urn.ToParts()

		return types.NewXText(fmt.Sprintf("%s:%s", scheme, redacted))
	}

	return types.NewXText(string(r.urn))
}

type RouteEnvelope struct {
	URN     urns.URN                 `json:"urn"               validate:"required,urn"`
	Channel *assets.ChannelReference `json:"channel,omitempty"`
}

func (e *RouteEnvelope) Unmarshal(sa SessionAssets, missing assets.MissingCallback) *Route {
	var channel *Channel
	if e.Channel != nil {
		if channel = sa.Channels().Get(e.Channel.UUID); channel == nil {
			missing(e.Channel, nil)
		}
	}

	return &Route{urn: e.URN, channel: channel}
}

// RouteList is the list of a contact's routes
type RouteList []*Route

func NewRouteList(sa SessionAssets, envelopes []*RouteEnvelope, missing assets.MissingCallback) RouteList {
	routes := make(RouteList, len(envelopes))
	for i, e := range envelopes {
		routes[i] = e.Unmarshal(sa, missing)
	}
	return routes
}

// Encode returns the encoded URNs in this list
func (l RouteList) Encode() []urns.URN {
	encoded := make([]urns.URN, len(l))
	for i, r := range l {
		encoded[i] = r.Encode()
	}
	return encoded
}

// Equal returns whether this list of URNs is equal to another
func (l RouteList) Equal(other RouteList) bool {
	if len(l) != len(other) {
		return false
	}

	for i := range l {
		if !l[i].Equal(other[i]) {
			return false
		}
	}
	return true
}

// Clone returns a clone of this route list
func (l RouteList) clone() RouteList {
	urns := make(RouteList, len(l))
	for i, r := range l {
		urns[i] = &Route{urn: r.urn, channel: r.channel}
	}
	return urns
}

func (l RouteList) marshal() []*RouteEnvelope {
	es := make([]*RouteEnvelope, len(l))
	for i, r := range l {
		es[i] = &RouteEnvelope{URN: r.urn, Channel: r.channel.Reference()}
	}
	return es
}

// WithScheme returns a new URN list containing of only URNs of the given schemes
func (l RouteList) WithScheme(schemes ...string) RouteList {
	var matching RouteList
	for _, u := range l {
		if slices.Contains(schemes, u.urn.Scheme()) {
			matching = append(matching, u)
		}
	}
	return matching
}

// ToXValue returns a representation of this object for use in expressions
func (l RouteList) ToXValue(env envs.Environment) types.XValue {
	return types.NewXLazyArray(func() []types.XValue {
		array := make([]types.XValue, len(l))
		for i, urn := range l {
			array[i] = urn.ToXValue(env)
		}
		return array
	})
}

// MapContext returns a map of the highest priority URN for each scheme - exposed in expressions as @urns
func (l RouteList) MapContext(env envs.Environment) map[string]types.XValue {
	byScheme := make(map[string]types.XValue)

	for _, u := range l {
		scheme := u.URN().Scheme()
		if _, seen := byScheme[scheme]; !seen {
			byScheme[scheme] = u.ToXValue(env)
		}
	}

	// and add nils for all other schemes
	for _, scheme := range urns.Schemes {
		if _, seen := byScheme[scheme.Prefix]; !seen {
			byScheme[scheme.Prefix] = nil
		}
	}

	return byScheme
}
