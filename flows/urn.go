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

// URN adds additional functionality around urns.URN such as channel affinity.
type URN struct {
	Scheme  string
	Path    string
	Display string
	Channel *Channel
}

// NewURN creates a new contact URN with associated channel
func NewURN(scheme, path, display string, channel *Channel) *URN {
	return &URN{Scheme: scheme, Path: path, Display: display, Channel: channel}
}

// ParseURN converts an encoded urns.URN to a URN
func ParseURN(ca *ChannelAssets, encoded urns.URN, missing assets.MissingCallback) (*URN, error) {
	scheme, path, query, display := encoded.ToParts()

	parsedQuery, err := url.ParseQuery(query)
	if err != nil {
		return nil, err
	}

	var channel *Channel
	if chUUID := parsedQuery.Get("channel"); chUUID != "" {
		if channel = ca.Get(assets.ChannelUUID(chUUID)); channel == nil {
			missing(assets.NewChannelReference(assets.ChannelUUID(chUUID), ""), nil)
		}
	}

	return NewURN(scheme, path, display, channel), nil
}

// Encode gets the URN as an encoded string
func (u *URN) Encode() urns.URN {
	return u.encode(true, true)
}

// Identity gets the URN as an identity URN (just scheme and path)
func (u *URN) Identity() urns.URN {
	return u.encode(false, false)
}

func (u *URN) encode(display, channel bool) urns.URN {
	var fragment string
	if display {
		fragment = u.Display
	}

	query := url.Values{}
	if channel && u.Channel != nil {
		query.Set("channel", string(u.Channel.UUID()))
	}

	urn, _ := urns.NewFromParts(u.Scheme, u.Path, query, fragment)
	return urn
}

// Equal determines if this URN is equal to another
func (u *URN) Equal(other *URN) bool {
	return other != nil && u.Encode() == other.Encode()
}

func (u *URN) clone() *URN {
	return NewURN(u.Scheme, u.Path, u.Display, u.Channel)
}

// ToXValue returns a representation of this object for use in expressions
func (u *URN) ToXValue(env envs.Environment) types.XValue {
	if env.RedactionPolicy() == envs.RedactionPolicyURNs {
		return types.NewXText(fmt.Sprintf("%s:%s", u.Scheme, redacted))
	}

	return types.NewXText(string(u.encode(true, false)))
}

// URNList is the list of a contact's URNs
type URNList []*URN

func NewURNList(sa SessionAssets, encoded []urns.URN, missing assets.MissingCallback) (URNList, error) {
	urns := make(URNList, len(encoded))
	for i, e := range encoded {
		parsed, err := ParseURN(sa.Channels(), e, missing)
		if err != nil {
			return nil, fmt.Errorf("unable to parse URN %s: %w", e, err)
		}
		urns[i] = parsed
	}
	return urns, nil
}

// Encode returns encoded URNs
func (l URNList) Encode() []urns.URN {
	encoded := make([]urns.URN, len(l))
	for i := range l {
		encoded[i] = l[i].Encode()
	}
	return encoded
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
	for i, u := range l {
		urns[i] = u.clone()
	}
	return urns
}

// WithScheme returns a new URN list containing of only URNs of the given schemes
func (l URNList) WithScheme(schemes ...string) URNList {
	var matching URNList
	for _, u := range l {
		if slices.Contains(schemes, u.Scheme) {
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
		scheme := u.Scheme
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
