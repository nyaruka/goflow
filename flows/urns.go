package flows

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	validator "gopkg.in/go-playground/validator.v9"
)

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
// To render a URN in a human friendly format, use the [format_urn](#function:format_urn) function.
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
	urns.URN
	channel Channel
}

// NewContactURN creates a new contact URN with associated channel
func NewContactURN(urn urns.URN, channel Channel) *ContactURN {
	return &ContactURN{URN: urn, channel: channel}
}

// Channel gets the channel associated with this URN
func (u *ContactURN) Channel() Channel { return u.channel }

// SetChannel sets the channel associated with this URN
func (u *ContactURN) SetChannel(channel Channel) { u.channel = channel }

// Resolve resolves the given key when this URN is referenced in an expression
func (u *ContactURN) Resolve(env utils.Environment, key string) types.XValue {
	switch key {
	case "scheme":
		return types.NewXText(u.URN.Scheme())
	case "path":
		return types.NewXText(u.URN.Path())
	case "display":
		return types.NewXText(u.URN.Display())
	case "channel":
		return u.Channel()
	}
	return types.NewXResolveError(u, key)
}

// Describe returns a representation of this type for error messages
func (u *ContactURN) Describe() string { return "URN" }

// Reduce is called when this object needs to be reduced to a primitive
func (u *ContactURN) Reduce() types.XPrimitive { return types.NewXText(string(u.URN)) }

// ToXJSON is called when this type is passed to @(json(...))
func (u *ContactURN) ToXJSON(env utils.Environment) types.XText {
	return types.ResolveKeys(env, u, "scheme", "path", "display").ToXJSON(env)
}

var _ types.XValue = (*ContactURN)(nil)
var _ types.XResolvable = (*ContactURN)(nil)

// URNList is the list of a contact's URNs
type URNList []*ContactURN

// ReadURNList parses contact URN list from the given list of raw URNs
func ReadURNList(session Session, rawURNs []urns.URN) (URNList, error) {
	l := make(URNList, len(rawURNs))

	for u := range rawURNs {
		scheme, path, query, display := rawURNs[u].ToParts()

		// re-create the URN without the query component
		queryLess, err := urns.NewURNFromParts(scheme, path, "", display)
		if err != nil {
			return nil, err
		}

		parsedQuery, err := url.ParseQuery(query)
		if err != nil {
			return nil, err
		}

		var channel Channel
		channelUUID := parsedQuery.Get("channel")
		if channelUUID != "" {
			if channel, err = session.Assets().GetChannel(ChannelUUID(channelUUID)); err != nil {
				return nil, err
			}
		}

		l[u] = &ContactURN{URN: queryLess, channel: channel}
	}
	return l, nil
}

// RawURNs returns the raw URNs with or without channel information
func (l URNList) RawURNs(includeChannels bool) []urns.URN {
	raw := make([]urns.URN, len(l))
	for u := range l {
		scheme, path, query, display := l[u].URN.ToParts()

		if includeChannels && l[u].channel != nil {
			query = fmt.Sprintf("channel=%s", l[u].channel.UUID())
		}

		raw[u], _ = urns.NewURNFromParts(scheme, path, query, display)
	}
	return raw
}

// Clone returns a clone of this URN list
func (l URNList) clone() URNList {
	urns := make(URNList, len(l))
	copy(urns, l)
	return urns
}

// WithScheme returns a new URN list containing of only URNs of the given scheme
func (l URNList) WithScheme(scheme string) URNList {
	var matching URNList
	for _, u := range l {
		if u.URN.Scheme() == scheme {
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
func (l URNList) Reduce() types.XPrimitive {
	array := types.NewXArray()
	for _, urn := range l {
		array.Append(urn)
	}
	return array
}

// ToXJSON is called when this type is passed to @(json(...))
func (l URNList) ToXJSON(env utils.Environment) types.XText {
	return l.Reduce().ToXJSON(env)
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
