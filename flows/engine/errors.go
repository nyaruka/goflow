package engine

import "fmt"

const (
	ErrorResumeNonWaitingSession int = 101
	ErrorResumeNoWaitingRun      int = 102
	ErrorResumeRejectedByWait    int = 103
)

type Error struct {
	code int
	msg  string
}

func newError(code int, msg string, args ...interface{}) error {
	return &Error{code, fmt.Sprintf(msg, args...)}
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Error() string {
	return e.msg
}
