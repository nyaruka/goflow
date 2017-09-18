package flows

import (
	"fmt"
	"strings"

	"gopkg.in/go-playground/validator.v9"

	"github.com/nyaruka/goflow/utils"
)

func init() {
	utils.Validator.RegisterValidation("urnscheme", ValidateURNScheme)
}

type URN string
type URNScheme string
type URNPath string

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

func ValidateURNScheme(fl validator.FieldLevel) bool {
	_, valid := schemes[fl.Field().String()]
	return valid
}

func GetScheme(scheme string) URNScheme {
	lowered := strings.ToLower(scheme)
	if schemes[lowered] {
		return URNScheme(lowered)
	}
	return ""
}

func NewURNFromParts(scheme URNScheme, path URNPath) URN {
	return URN(fmt.Sprintf("%s:%s", scheme, path))
}

func (u URN) Path() URNPath {
	offset := strings.Index(string(u), ":")
	if offset >= 0 {
		return URNPath(u[offset+1:])
	}
	return URNPath(u)
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

// URNList is a list of URNs on a contact
type URNList []URN

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
