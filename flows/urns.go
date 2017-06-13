package flows

import (
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/utils"
)

type URN string
type URNScheme string
type URNPath string

type URNList []URN

// List of schemes we support for URNs
const (
	TelScheme      = "tel"
	FacebookScheme = "facebook"
	TwitterScheme  = "twitter"
	ViberScheme    = "viber"
	TelegramScheme = "telegram"
	EmailScheme    = "email"
)

// Used as a lookup for faster checks whether a Scheme is supported
var schemes = map[string]bool{
	TelScheme:      true,
	FacebookScheme: true,
	TwitterScheme:  true,
	ViberScheme:    true,
	TelegramScheme: true,
	EmailScheme:    true,
}

func (u URNScheme) String() string { return string(u) }
func (u URNPath) String() string   { return string(u) }

func GetScheme(scheme string) URNScheme {
	lowered := strings.ToLower(scheme)
	if schemes[lowered] {
		return URNScheme(lowered)
	}
	return ""
}

func (u URN) Path() URNPath {
	offset := strings.Index(string(u), ":")
	if offset >= 0 {
		return URNPath(u[offset+1:])
	}
	return URNPath(offset)
}

func (u URN) Scheme() URNScheme {
	offset := strings.Index(string(u), ":")
	if offset >= 0 {
		return URNScheme(strings.ToLower(string(u[:offset])))
	}
	return URNScheme("")
}

func (u URN) Resolve(key string) interface{} {
	switch key {

	case "path":
		return u.Path()

	case "scheme":
		return u.Scheme()

	case "urn":
		return string(u)
	}

	return fmt.Errorf("No field '%s' on URN", key)
}

func (u URN) Default() interface{} { return u }
func (u URN) String() string       { return string(u.Path()) }

var _ utils.VariableResolver = (URN)("")

func (l URNList) Resolve(key string) interface{} {
	// If this isn't a valid scheme, bail
	scheme := GetScheme(key)
	if scheme == "" {
		return fmt.Errorf("Unknown URN scheme: %s", key)
	}

	// This is a specific scheme, look up all matches
	var found URNList
	for _, u := range l {
		if u.Scheme() == scheme {
			found = append(found, u)
		}
	}

	return found
}

func (l URNList) Default() interface{} {
	if len(l) > 0 {
		return l[0]
	}
	return nil
}

func (l URNList) String() string {
	if len(l) > 0 {
		return l[0].String()
	}
	return ""
}

var _ utils.VariableResolver = (URNList)(nil)
