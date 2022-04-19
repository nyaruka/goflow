package contactql

import (
	"fmt"

	"github.com/nyaruka/goflow/utils"
	"github.com/pkg/errors"
)

// error codes with values included in extra
const (
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
	ErrUnknownProperty       = "unknown_property"       // `property` the property key
	ErrRedactedURNs          = "redacted_urns"
)

// QueryError is used when an error is a result of an invalid query
type QueryError struct {
	msg   string
	code  string
	extra map[string]string
}

// NewQueryError creates a new query error
func NewQueryError(code, err string, args ...interface{}) *QueryError {
	return &QueryError{code: code, msg: fmt.Sprintf(err, args...)}
}

func (e *QueryError) withExtra(k, v string) *QueryError {
	if e.extra == nil {
		e.extra = make(map[string]string)
	}
	e.extra[k] = v
	return e
}

// Error returns the error message
func (e *QueryError) Error() string {
	return e.msg
}

// Code returns a code representing this error condition
func (e *QueryError) Code() string {
	return e.code
}

// Extra returns additional data about the error
func (e *QueryError) Extra() map[string]string {
	return e.extra
}

// IsQueryError is a utility to determine if the cause of an error was a query error
func IsQueryError(err error) (bool, error) {
	switch cause := errors.Cause(err).(type) {
	case *QueryError:
		return true, cause
	default:
		return false, nil
	}
}

var _ utils.RichError = (*QueryError)(nil)
