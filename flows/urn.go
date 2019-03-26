package flows

import (
	"net/url"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	validator "gopkg.in/go-playground/validator.v9"
)

var redactedURN = types.NewXText("********")

func init() {
	utils.Validator.RegisterValidation("urn", ValidateURN)
	utils.Validator.RegisterValidation("urnscheme", ValidateURNScheme)
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

// ContactURN represents a destination for an outgoing message or a source of an incoming message. It is string composed of 3
// components: scheme, path, and display (optional). For example:
//
//  - _tel:+16303524567_
//  - _twitterid:54784326227#nyaruka_
//  - _telegram:34642632786#bobby_
//
// To render a URN in a human friendly format, use the [function:format_urn] function.
//
// Examples:
//
//   @(urns.tel) -> tel:+12065551212
//   @(urn_parts(urns.tel).scheme) -> tel
//   @(format_urn(urns.tel)) -> (206) 555-1212
//   @(json(contact.urns[0])) -> "tel:+12065551212"
//
// @context urn
type ContactURN struct {
	urn     urns.URN
	channel *Channel
}

// NewContactURN creates a new contact URN with associated channel
func NewContactURN(urn urns.URN, channel *Channel) *ContactURN {
	return &ContactURN{urn: urn, channel: channel}
}

// ParseRawURN converts a raw URN to a ContactURN by extracting it's channel reference
func ParseRawURN(ca *ChannelAssets, rawURN urns.URN, missing assets.MissingCallback) (*ContactURN, error) {
	_, _, query, _ := rawURN.ToParts()

	parsedQuery, err := url.ParseQuery(query)
	if err != nil {
		return nil, err
	}

	var channel *Channel
	channelUUID := assets.ChannelUUID(parsedQuery.Get("channel"))
	if channelUUID != "" {
		if channel = ca.Get(channelUUID); channel == nil {
			missing(assets.NewChannelReference(channelUUID, ""))
		}
	}

	return NewContactURN(rawURN, channel), nil
}

// URN gets the underlying URN
func (u *ContactURN) URN() urns.URN { return u.urn }

// Channel gets the channel associated with this URN
func (u *ContactURN) Channel() *Channel { return u.channel }

// SetChannel sets the channel associated with this URN
func (u *ContactURN) SetChannel(channel *Channel) {
	u.channel = channel

	scheme, path, query, display := u.urn.ToParts()

	parsedQuery, _ := url.ParseQuery(query)

	if channel != nil {
		parsedQuery.Set("channel", string(channel.UUID()))
	} else {
		parsedQuery.Del("channel")
	}

	urn, _ := urns.NewURNFromParts(scheme, path, parsedQuery.Encode(), display)
	u.urn = urn
}

func (u *ContactURN) String() string {
	return string(u.urn)
}

// Equal determines if this URN is equal to another
func (u *ContactURN) Equal(other *ContactURN) bool {
	return other != nil && u.String() == other.String()
}

// returns this URN as a raw URN without the query portion (i.e. only scheme, path, display)
func (u *ContactURN) withoutQuery() urns.URN {
	scheme, path, _, display := u.urn.ToParts()
	urn, _ := urns.NewURNFromParts(scheme, path, "", display)
	return urn
}

// Describe returns a representation of this type for error messages
func (u *ContactURN) Describe() string { return "URN" }

// Reduce is called when this object needs to be reduced to a primitive
func (u *ContactURN) Reduce(env utils.Environment) types.XPrimitive {
	if env.RedactionPolicy() == utils.RedactionPolicyURNs {
		return redactedURN
	}
	return types.NewXText(string(u.withoutQuery()))
}

// ToXJSON is called when this type is passed to @(json(...))
func (u *ContactURN) ToXJSON(env utils.Environment) types.XText {
	return u.Reduce(env).ToXJSON(env)
}

var _ types.XValue = (*ContactURN)(nil)

// URNList is the list of a contact's URNs
type URNList []*ContactURN

// ReadURNList parses contact URN list from the given list of raw URNs
func ReadURNList(a SessionAssets, rawURNs []urns.URN, missing assets.MissingCallback) (URNList, error) {
	l := make(URNList, len(rawURNs))

	for u := range rawURNs {
		parsed, err := ParseRawURN(a.Channels(), rawURNs[u], missing)
		if err != nil {
			return nil, err
		}
		l[u] = parsed
	}

	return l, nil
}

// RawURNs returns the raw URNs
func (l URNList) RawURNs() []urns.URN {
	raw := make([]urns.URN, len(l))
	for u := range l {
		raw[u] = l[u].urn
	}
	return raw
}

// Equal returns whether this list of URNs is equal to another
func (l URNList) Equal(other URNList) bool {
	if len(l) != len(other) {
		return false
	}

	for u := range l {
		if !l[u].Equal(other[u]) {
			return false
		}
	}
	return true
}

// Clone returns a clone of this URN list
func (l URNList) clone() URNList {
	urns := make(URNList, len(l))
	for u := range l {
		urns[u] = NewContactURN(l[u].urn, l[u].channel)
	}
	return urns
}

// WithScheme returns a new URN list containing of only URNs of the given scheme
func (l URNList) WithScheme(scheme string) URNList {
	var matching URNList
	for _, u := range l {
		if u.urn.Scheme() == scheme {
			matching = append(matching, u)
		}
	}
	return matching
}

// Context returns this as an XArray - exposed in expressions as @contact.urns, @parent.contact.urns etc
func (l URNList) Context() types.XValue {
	array := types.NewXArray()
	for _, urn := range l {
		array.Append(urn)
	}
	return array
}

// MapContext returns a map of the highest priority URN for each scheme - exposed in expressions as @urns
func (l URNList) MapContext() types.XValue {
	byScheme := make(map[string]types.XValue)

	for _, u := range l {
		scheme := u.URN().Scheme()
		if _, seen := byScheme[scheme]; !seen {
			byScheme[scheme] = u
		}
	}

	// and add nils for all other schemes
	for scheme := range urns.ValidSchemes {
		if _, seen := byScheme[scheme]; !seen {
			byScheme[scheme] = nil
		}
	}

	return types.NewXMap(byScheme)
}
