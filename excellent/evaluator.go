package excellent

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/excellent/gen"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// EvaluateExpression evalutes the passed in template, returning the typed value it evaluates to, which might be an error
func EvaluateExpression(env utils.Environment, context types.XValue, expression string) types.XValue {
	errListener := NewErrorListener(expression)

	input := antlr.NewInputStream(expression)
	lexer := gen.NewExcellent2Lexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := gen.NewExcellent2Parser(stream)
	p.RemoveErrorListeners()
	p.AddErrorListener(errListener)
	tree := p.Parse()

	// if we ran into errors parsing, return the first one
	if len(errListener.Errors()) > 0 {
		return errListener.Errors()[0]
	}

	visitor := NewVisitor(env, context)
	return toXValue(visitor.Visit(tree))
}

// EvaluateTemplate tries to evaluate the passed in template into an object, this only works if the template
// is a single identifier or expression, ie: "@contact" or "@(first(contact.urns))". In cases
// which are not a single identifier or expression, we return the stringified value
func EvaluateTemplate(env utils.Environment, context types.XValue, template string, allowedTopLevels []string) (types.XValue, error) {
	template = strings.TrimSpace(template)
	scanner := NewXScanner(strings.NewReader(template), allowedTopLevels)

	// parse our first token
	tokenType, token := scanner.Scan()

	// try to scan to our next token
	nextTT, _ := scanner.Scan()

	// if we only have an identifier or an expression, evaluate it on its own
	if nextTT == EOF {
		switch tokenType {
		case IDENTIFIER:
			return evaluateIdentifier(env, context, token), nil
		case EXPRESSION:
			return EvaluateExpression(env, context, token), nil
		}
	}

	// otherwise fallback to full template evaluation
	asStr, err := EvaluateTemplateAsString(env, context, template, allowedTopLevels)
	return types.NewXText(asStr), err
}

// EvaluateTemplateAsString evaluates the passed in template returning the string value of its execution
func EvaluateTemplateAsString(env utils.Environment, context types.XValue, template string, allowedTopLevels []string) (string, error) {
	var buf bytes.Buffer
	scanner := NewXScanner(strings.NewReader(template), allowedTopLevels)
	errors := NewTemplateErrors()

	for tokenType, token := scanner.Scan(); tokenType != EOF; tokenType, token = scanner.Scan() {
		switch tokenType {
		case BODY:
			buf.WriteString(token)
		case IDENTIFIER:
			value := evaluateIdentifier(env, context, token)

			if types.IsXError(value) {
				errors.Add(fmt.Sprintf("@%s", token), value.(error).Error())
			} else {
				strValue, _ := types.ToXText(env, value)

				buf.WriteString(strValue.Native())
			}
		case EXPRESSION:
			value := EvaluateExpression(env, context, token)

			if types.IsXError(value) {
				errors.Add(fmt.Sprintf("@(%s)", token), value.(error).Error())
			} else {
				strValue, _ := types.ToXText(env, value)

				buf.WriteString(strValue.Native())
			}
		}
	}

	if errors.HasErrors() {
		return buf.String(), errors
	}
	return buf.String(), nil
}

// Evaluates an identifier like "foo.bar.zed".. these could be passed through the full excellent parser
// but as an optimization we handle them separately.
func evaluateIdentifier(env utils.Environment, context types.XValue, identifier string) types.XValue {
	parts := strings.Split(identifier, ".")
	value := context
	for _, part := range parts {
		value = lookupProperty(env, value, part)
		if types.IsXError(value) {
			break
		}
	}
	return value
}
