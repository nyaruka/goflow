package contactql

import (
	"fmt"

	"github.com/pkg/errors"
)

// QueryError is used when an error is a result of an invalid query
type QueryError struct {
	msg string
}

func (e *QueryError) Error() string {
	return e.msg
}

// NewQueryErrorf creates a new query error
func NewQueryErrorf(err string, args ...interface{}) *QueryError {
	return &QueryError{fmt.Sprintf(err, args...)}
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
