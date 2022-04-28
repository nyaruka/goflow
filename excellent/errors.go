package excellent

import (
	"fmt"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/nyaruka/goflow/utils"
	"github.com/pkg/errors"
)

// TemplateError is an error which occurs during evaluation of an expression
type TemplateError struct {
	expression string
	message    string
}

func (e TemplateError) Error() string {
	return fmt.Sprintf("error evaluating %s: %s", e.expression, e.message)
}

// TemplateErrors represents the list of all errors encountered during evaluation of a template
type TemplateErrors struct {
	errors []*TemplateError
}

// NewTemplateErrors creates a new empty lsit of template errors
func NewTemplateErrors() *TemplateErrors {
	return &TemplateErrors{}
}

// Add adds an error for the given expression
func (e *TemplateErrors) Add(expression, message string) {
	e.errors = append(e.errors, &TemplateError{expression: expression, message: message})
}

// HasErrors returns whether there are errors
func (e *TemplateErrors) HasErrors() bool {
	return len(e.errors) > 0
}

// Error returns a single string describing all the errors encountered
func (e *TemplateErrors) Error() string {
	messages := make([]string, len(e.errors))
	for i, err := range e.errors {
		messages[i] = err.Error()
	}
	return strings.Join(messages, ", ")
}

// ErrorListener records syntax errors
type ErrorListener struct {
	*antlr.DefaultErrorListener

	expression string
	errors     []error
}

// NewErrorListener creates a new error listener
func NewErrorListener(expression string) *ErrorListener {
	return &ErrorListener{expression: expression}
}

// Errors returns the errors encountered so far
func (l *ErrorListener) Errors() []error {
	return l.errors
}

// SyntaxError handles a new syntax error encountered by the recognizer
func (l *ErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	// extract the part of the original expression where this error has occurred
	lines := strings.Split(l.expression, "\n")
	lineOfError := lines[line-1]
	contextOfError := lineOfError[column:utils.Min(column+10, len(lineOfError))]

	l.errors = append(l.errors, errors.Errorf("syntax error at %s", contextOfError))
}
