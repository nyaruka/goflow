package contactql

import (
	"github.com/nyaruka/goflow/utils"
)

// error codes with values included in extra
const (
	ErrSyntax                = "syntax"
	ErrUnexpectedToken       = "unexpected_token"       // `token` the unexpected token
	ErrInvalidNumber         = "invalid_number"         // `value` the value we tried to parse as a number
	ErrInvalidDate           = "invalid_date"           // `value` the value we tried to parse as a date
	ErrInvalidStatus         = "invalid_status"         // `value` the value we tried to parse as a contact status
	ErrInvalidLanguage       = "invalid_language"       // `value` the value we tried to parse as a language code
	ErrInvalidGroup          = "invalid_group"          // `value` the value we tried to parse as a group name
	ErrInvalidFlow           = "invalid_flow"           // `value` the value we tried to parse as a flow name
	ErrInvalidPartialName    = "invalid_partial_name"   // `min_token_length` the minimum length of token required for name contains condition
	ErrInvalidPartialURN     = "invalid_partial_urn"    // `min_value_length` the minimum length of value required for URN contains condition
	ErrUnsupportedContains   = "unsupported_contains"   // `property` the property key
	ErrUnsupportedComparison = "unsupported_comparison" // `property` the property key, `operator` one of =>, <, >=, <=
	ErrUnsupportedSetCheck   = "unsupported_setcheck"   // `property` the property key, `operator` one of =, !=
	ErrUnknownPropertyType   = "unknown_property_type"  // `type` the property type
	ErrUnknownProperty       = "unknown_property"       // `property` the property key
	ErrRedactedURNs          = "redacted_urns"
)

// creates a new query error
func newQueryError(code, err string, args ...any) *utils.RichError {
	return utils.NewRichError("query", code, err, args...)
}
