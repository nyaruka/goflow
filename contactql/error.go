package contactql

import (
	"fmt"

	"github.com/nyaruka/goflow/utils"
	"github.com/pkg/errors"
)

// QueryError is used when an error is a result of an invalid query
type QueryError struct {
	msg   string
	code  string
	extra map[string]string
}

// NewQueryErrorf creates a new query error
func NewQueryErrorf(err string, args ...interface{}) *QueryError {
	return &QueryError{msg: fmt.Sprintf(err, args...)}
}

// NewQueryError creates a new query error
func NewQueryError(msg string, code string, extra map[string]string) *QueryError {
	return &QueryError{msg: msg, code: code, extra: extra}
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
