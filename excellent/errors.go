package excellent

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// TemplateErrors represents the list of errors we may have received during execution
type TemplateErrors []error

// Error returns a single string describing all the errors encountered
func (e TemplateErrors) Error() string {
	if len(e) == 1 {
		return e[0].Error()
	}

	msg := "multiple errors:"
	for _, err := range e {
		msg += "\n" + err.Error()
	}
	return msg
}

type errorListener struct {
	errors bytes.Buffer
	*antlr.DefaultErrorListener
}

func NewErrorListener() *errorListener {
	return &errorListener{}
}

func (l *errorListener) HasErrors() bool {
	return l.errors.Len() > 0
}

func (l *errorListener) Errors() string {
	return l.errors.String()
}

func (l *errorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	l.errors.WriteString(fmt.Sprintln("line " + strconv.Itoa(line) + ":" + strconv.Itoa(column) + " " + msg))
}
