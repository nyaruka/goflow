package flows

import (
	"fmt"
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
	return urns.URN(fl.Field().String()).Validate()
}

// ValidateURNScheme validates whether the field value is a valid URN scheme
func ValidateURNScheme(fl validator.FieldLevel) bool {
	return urns.IsValidScheme(fl.Field().String())
}

// URNList is the list of a contact's URNs
type URNList []urns.URN

func (l URNList) Clone() URNList {
	urns := make(URNList, len(l))
	copy(urns, l)
	return urns
}

func (l URNList) WithScheme(scheme string) URNList {
	var matching URNList
	for _, u := range l {
		if u.Scheme() == scheme {
			matching = append(matching, u)
		}
	}
	return matching
}

func (l URNList) Resolve(key string) interface{} {
	scheme := strings.ToLower(key)

	// if this isn't a valid scheme, bail
	if !urns.IsValidScheme(scheme) {
		return fmt.Errorf("unknown URN scheme: %s", key)
	}

	// this is a specific scheme, return all matches
	return l.WithScheme(scheme)
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
