package flows

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/nyaruka/gocommon/urns"
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

// ContactURN holds a URN for a contact with the channel parsed out
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
func (u *ContactURN) Resolve(key string) interface{} {
	switch key {
	case "scheme":
		return u.URN.Scheme()
	case "path":
		return u.URN.Path()
	case "display":
		return u.URN.Display()
	case "channel":
		return u.Channel()
	}
	return fmt.Errorf("no field '%s' on URN", key)
}

// Default returns the value of this URN when it is the result of an expression
func (u *ContactURN) Default() interface{} { return u }

func (u *ContactURN) String() string { return string(u.URN) }

// URNList is the list of a contact's URNs
type URNList []*ContactURN

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

func (l URNList) Clone() URNList {
	urns := make(URNList, len(l))
	copy(urns, l)
	return urns
}

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
func (l URNList) Resolve(key string) interface{} {
	// first try as numeric index to a single URN
	index, err := strconv.Atoi(key)
	if err == nil {
		if index < len(l) {
			return l[index]
		}
		return fmt.Errorf("index out of range: %d", index)
	}

	// next try as a URN scheme
	scheme := strings.ToLower(key)

	// if this isn't a valid scheme, bail
	if !urns.IsValidScheme(scheme) {
		return fmt.Errorf("unknown URN scheme: %s", key)
	}

	return l.WithScheme(scheme)
}

// Default returns the value of this URN list when it is the result of an expression
func (l URNList) Default() interface{} {
	return l
}

func (l URNList) String() string {
	if len(l) > 0 {
		return l[0].String()
	}
	return ""
}

var _ utils.VariableResolver = &ContactURN{}
var _ utils.VariableResolver = (URNList)(nil)
