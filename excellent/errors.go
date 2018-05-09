package excellent

import (
	"fmt"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// TemplateError is an error which occurs during evaluation of an expression
type TemplateError struct {
	expression string
	message    string
}

func (e TemplateError) Error() string {
	return fmt.Sprintf("error evaluating '%s': %s", e.expression, e.message)
}

// TemplateErrors represents the list of all errors encountered during evaluation of a template
type TemplateErrors struct {
	errors []*TemplateError
}

func NewTemplateErrors() *TemplateErrors {
	return &TemplateErrors{}
}

func (e *TemplateErrors) Add(expression, message string) {
	e.errors = append(e.errors, &TemplateError{expression: expression, message: message})
}

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

type errorListener struct {
	*antlr.DefaultErrorListener

	expression string
	errors     []error
}

func NewErrorListener(expression string) *errorListener {
	return &errorListener{expression: expression}
}

func (l *errorListener) HasErrors() bool {
	return len(l.errors) > 0
}

func (l *errorListener) FirstError() error {
	return l.errors[0]
}

func (l *errorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	// extract the part of the original expression where this error has occured
	lines := strings.Split(l.expression, "\n")
	lineOfError := lines[line-1]
	contextOfError := lineOfError[column:min(column+10, len(lineOfError))]

	l.errors = append(l.errors, fmt.Errorf("syntax error at %s", contextOfError))
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
