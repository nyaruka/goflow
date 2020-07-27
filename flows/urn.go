package flows

import (
	"fmt"
	"net/url"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	validator "gopkg.in/go-playground/validator.v9"
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

// ContactURN represents a destination for an outgoing message or a source of an incoming message. It is string composed of 3
// components: scheme, path, and display (optional). For example:
//
//  - _tel:+16303524567_
//  - _twitterid:54784326227#nyaruka_
//  - _telegram:34642632786#bobby_
//
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
			missing(assets.NewChannelReference(channelUUID, ""), nil)
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
func (u *ContactURN) withoutQuery(redact bool) urns.URN {
	scheme, path, _, display := u.urn.ToParts()

	if redact {
		return urns.URN(fmt.Sprintf("%s:%s", scheme, redacted))
	}

	urn, _ := urns.NewURNFromParts(scheme, path, "", display)

	return urn
}

// ToXValue returns a representation of this object for use in expressions
func (u *ContactURN) ToXValue(env envs.Environment) types.XValue {
	redact := env.RedactionPolicy() == envs.RedactionPolicyURNs

	return types.NewXText(string(u.withoutQuery(redact)))
}

// URNList is the list of a contact's URNs
type URNList []*ContactURN

// ReadURNList parses contact URN list from the given list of raw URNs
func ReadURNList(a SessionAssets, rawURNs []urns.URN, missing assets.MissingCallback) (URNList, error) {
	l := make(URNList, len(rawURNs))

	for i := range rawURNs {
		parsed, err := ParseRawURN(a.Channels(), rawURNs[i], missing)
		if err != nil {
			return nil, err
		}
		l[i] = parsed
	}

	return l, nil
}

// RawURNs returns the raw URNs
func (l URNList) RawURNs() []urns.URN {
	raw := make([]urns.URN, len(l))
	for i := range l {
		raw[i] = l[i].urn
	}
	return raw
}

// Equal returns whether this list of URNs is equal to another
func (l URNList) Equal(other URNList) bool {
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

// Clone returns a clone of this URN list
func (l URNList) clone() URNList {
	urns := make(URNList, len(l))
	for i := range l {
		urns[i] = NewContactURN(l[i].urn, l[i].channel)
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

// ToXValue returns a representation of this object for use in expressions
func (l URNList) ToXValue(env envs.Environment) types.XValue {
	return types.NewXLazyArray(func() []types.XValue {
		array := make([]types.XValue, len(l))
		for i, urn := range l {
			array[i] = urn.ToXValue(env)
		}
		return array
	})
}

// MapContext returns a map of the highest priority URN for each scheme - exposed in expressions as @urns
func (l URNList) MapContext(env envs.Environment) map[string]types.XValue {
	byScheme := make(map[string]types.XValue)

	for _, u := range l {
		scheme := u.URN().Scheme()
		if _, seen := byScheme[scheme]; !seen {
			byScheme[scheme] = u.ToXValue(env)
		}
	}

	// and add nils for all other schemes
	for scheme := range urns.ValidSchemes {
		if _, seen := byScheme[scheme]; !seen {
			byScheme[scheme] = nil
		}
	}

	return byScheme
}
