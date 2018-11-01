package flows

import (
	"net/url"
	"strings"

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
// It has several properties which can be accessed in expressions:
//
//  * `scheme` the scheme of the URN, e.g. "tel", "twitter"
//  * `path` the path of the URN, e.g. "+16303524567"
//  * `display` the display portion of the URN, e.g. "+16303524567"
//  * `channel` the preferred [channel](#context:channel) of the URN
//
// To render a URN in a human friendly format, use the [function:format_urn] function.
//
// Examples:
//
//   @contact.urns.0 -> tel:+12065551212
//   @contact.urns.0.scheme -> tel
//   @contact.urns.0.path -> +12065551212
//   @contact.urns.1.display -> nyaruka
//   @(format_urn(contact.urns.0)) -> (206) 555-1212
//   @(json(contact.urns.0)) -> {"display":"","path":"+12065551212","scheme":"tel"}
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
func ParseRawURN(a SessionAssets, rawURN urns.URN) (*ContactURN, error) {
	_, _, query, _ := rawURN.ToParts()

	parsedQuery, err := url.ParseQuery(query)
	if err != nil {
		return nil, err
	}

	var channel *Channel
	channelUUID := parsedQuery.Get("channel")
	if channelUUID != "" {
		if channel, err = a.Channels().Get(assets.ChannelUUID(channelUUID)); err != nil {
			return nil, err
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

// Resolve resolves the given key when this URN is referenced in an expression
func (u *ContactURN) Resolve(env utils.Environment, key string) types.XValue {
	switch strings.ToLower(key) {
	case "scheme":
		return types.NewXText(u.urn.Scheme())
	case "path":
		if env.RedactionPolicy() == utils.RedactionPolicyURNs {
			return redactedURN
		}
		return types.NewXText(u.urn.Path())
	case "display":
		if env.RedactionPolicy() == utils.RedactionPolicyURNs {
			return redactedURN
		}
		return types.NewXText(u.urn.Display())
	case "channel":
		return u.Channel()
	}
	return types.NewXResolveError(u, key)
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
	return types.ResolveKeys(env, u, "scheme", "path", "display").ToXJSON(env)
}

var _ types.XValue = (*ContactURN)(nil)
var _ types.XResolvable = (*ContactURN)(nil)

// URNList is the list of a contact's URNs
type URNList []*ContactURN

// ReadURNList parses contact URN list from the given list of raw URNs
func ReadURNList(a SessionAssets, rawURNs []urns.URN) (URNList, error) {
	l := make(URNList, len(rawURNs))

	for u := range rawURNs {
		parsed, err := ParseRawURN(a, rawURNs[u])
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

// Resolve resolves the given key when this URN list is referenced in an expression
func (l URNList) Resolve(env utils.Environment, key string) types.XValue {
	scheme := strings.ToLower(key)

	// if this isn't a valid scheme, bail
	if !urns.IsValidScheme(scheme) {
		return types.NewXErrorf("no such URN scheme '%s'", key)
	}

	return l.WithScheme(scheme)
}

// Describe returns a representation of this type for error messages
func (l URNList) Describe() string { return "URNs" }

// Reduce is called when this object needs to be reduced to a primitive
func (l URNList) Reduce(env utils.Environment) types.XPrimitive {
	array := types.NewXArray()
	for _, urn := range l {
		array.Append(urn)
	}
	return array
}

// ToXJSON is called when this type is passed to @(json(...))
func (l URNList) ToXJSON(env utils.Environment) types.XText {
	return l.Reduce(env).ToXJSON(env)
}

// Index is called when this object is indexed into in an expression
func (l URNList) Index(index int) types.XValue {
	return l[index]
}

// Length is called when the length of this object is requested in an expression
func (l URNList) Length() int {
	return len(l)
}

var _ types.XValue = (URNList)(nil)
var _ types.XIndexable = (URNList)(nil)
var _ types.XResolvable = (URNList)(nil)

// URNShortcuts provides a simpler way to access single URNs
type URNShortcuts struct {
	urns URNList
}

// NewURNShortcuts creates a new URN shortcuts
func NewURNShortcuts(urns URNList) *URNShortcuts {
	return &URNShortcuts{urns: urns}
}

// Resolve resolves the given key when this is referenced in an expression
func (s *URNShortcuts) Resolve(env utils.Environment, key string) types.XValue {
	scheme := strings.ToLower(key)

	// a scheme means find the first URN with that scheme
	if urns.IsValidScheme(scheme) {
		filtered := s.urns.WithScheme(scheme)
		if len(filtered) > 0 {
			return filtered[0]
		}
		return types.XTextEmpty
	}

	if len(s.urns) > 0 {
		return s.urns[0].Resolve(env, key)
	}

	return types.NewXResolveError(s, key)
}

// Describe returns a representation of this type for error messages
func (s *URNShortcuts) Describe() string { return "URN" }

// Reduce is called when this object needs to be reduced to a primitive
func (s *URNShortcuts) Reduce(env utils.Environment) types.XPrimitive {
	if len(s.urns) > 0 {
		return s.urns[0].Reduce(env)
	}
	return nil
}

// ToXJSON is called when this type is passed to @(json(...))
func (s *URNShortcuts) ToXJSON(env utils.Environment) types.XText {
	if len(s.urns) > 0 {
		return s.urns[0].ToXJSON(env)
	}
	return types.XTextEmpty
}

var _ types.XValue = (*URNShortcuts)(nil)
var _ types.XResolvable = (*URNShortcuts)(nil)
