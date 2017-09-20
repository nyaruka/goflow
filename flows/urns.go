package flows

import (
	"fmt"
	"strconv"
	"strings"

	validator "gopkg.in/go-playground/validator.v9"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	utils.Validator.RegisterValidation("urn", ValidateURN)
	utils.Validator.RegisterValidation("urnscheme", ValidateURNScheme)
}

// ValidateURN validates whether the field value is a valid URN
func ValidateURN(fl validator.FieldLevel) bool {
	fmt.Printf("ValidateURN(%s) -> %s\n", fl.Field().String(), strconv.FormatBool(urns.URN(fl.Field().String()).Validate("")))
	return urns.URN(fl.Field().String()).Validate("")
}

// ValidateURNScheme validates whether the field value is a valid URN scheme
func ValidateURNScheme(fl validator.FieldLevel) bool {
	return urns.IsValidScheme(fl.Field().String())
}

// URNList is the list of a contact's URNs
type URNList []urns.URN

func (l URNList) Resolve(key string) interface{} {
	scheme := strings.ToLower(key)

	// if this isn't a valid scheme, bail
	if !urns.IsValidScheme(scheme) {
		return fmt.Errorf("unknown URN scheme: %s", key)
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

var _ utils.VariableResolver = (urns.URN)("")
var _ utils.VariableResolver = (URNList)(nil)
