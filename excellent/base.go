package excellent

import (
	"strings"

	"github.com/antlr4-go/antlr/v4"
	gen "github.com/nyaruka/goflow/antlr/gen/excellent3"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
)

// Evaluator evaluates templates and expressions.
type Evaluator struct{}

// NewEvaluator creates a new evaluator
func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

// Escaping is a function applied to expressions in a template after they've been evaluated
type Escaping func(string) string

// Template evaluates the passed in template
func (e *Evaluator) Template(env envs.Environment, ctx *types.XObject, template string, escaping Escaping) (string, []string, error) {
	var buf strings.Builder
	var allWarnings []string

	err := VisitTemplate(template, ctx.Properties(), true, func(tokenType XTokenType, token string) error {
		switch tokenType {
		case BODY:
			buf.WriteString(token)
		case IDENTIFIER, EXPRESSION:
			value, warnings := e.Expression(env, ctx, token)

			allWarnings = append(allWarnings, warnings...)

			// if we got an error, return that
			if types.IsXError(value) {
				return value.(error)
			}

			// if not, stringify value and append to the output
			asText, _ := types.ToXText(env, value)
			asString := asText.Native()

			if escaping != nil {
				asString = escaping(asString)
			}

			buf.WriteString(asString)
		}
		return nil
	})

	return buf.String(), allWarnings, err
}

// TemplateValue is equivalent to Template except in the case where the template contains
// a single identifier or expression, ie: "@contact" or "@(first(contact.urns))". In these cases we return
// the typed value from EvaluateExpression instead of stringifying the result.
func (e *Evaluator) TemplateValue(env envs.Environment, ctx *types.XObject, template string) (types.XValue, []string, error) {
	template = strings.TrimSpace(template)
	scanner := NewXScanner(strings.NewReader(template), ctx.Properties())

	// parse our first token
	tokenType, token := scanner.Scan()

	// try to scan to our next token
	nextTT, _ := scanner.Scan()

	// if we only have an identifier or an expression, evaluate it on its own
	if nextTT == EOF {
		switch tokenType {
		case IDENTIFIER, EXPRESSION:
			val, warnings := e.Expression(env, ctx, token)
			return val, warnings, nil
		}
	}

	// otherwise fallback to full template evaluation
	asStr, warnings, err := e.Template(env, ctx, template, nil)
	return types.NewXText(asStr), warnings, err
}

// Expression evalutes the passed in Excellent expression, returning the typed value it evaluates to,
// which might be an error, e.g. "2 / 3" or "contact.fields.age"
func (e *Evaluator) Expression(env envs.Environment, ctx *types.XObject, expression string) (types.XValue, []string) {
	parsed, err := Parse(expression, nil)
	if err != nil {
		return types.NewXError(err), nil
	}

	scope := NewScope(ctx, nil)

	warnings := &Warnings{}

	return parsed.Evaluate(env, scope, warnings), warnings.all
}

// Parse parses an expression
func Parse(expression string, contextCallback func([]string)) (Expression, error) {
	errListener := NewErrorListener(expression)

	input := antlr.NewInputStream(expression)
	lexer := gen.NewExcellent3Lexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := gen.NewExcellent3Parser(stream)
	p.RemoveErrorListeners()
	p.AddErrorListener(errListener)
	tree := p.Parse()

	// if we ran into errors parsing, return the first one
	if len(errListener.Errors()) > 0 {
		return nil, errListener.Errors()[0]
	}

	visitor := &visitor{contextCallback: contextCallback}
	output := visitor.Visit(tree)
	return toExpression(output), nil
}

// VisitTemplate scans the given template and calls the callback for each token encountered
func VisitTemplate(template string, allowedTopLevels []string, unescapeBody bool, callback func(XTokenType, string) error) error {
	// nothing todo for an empty template
	if template == "" {
		return nil
	}

	scanner := NewXScanner(strings.NewReader(template), allowedTopLevels)
	scanner.SetUnescapeBody(unescapeBody)
	errors := NewTemplateErrors()

	for tokenType, token := scanner.Scan(); tokenType != EOF; tokenType, token = scanner.Scan() {
		if err := callback(tokenType, token); err != nil {
			var repr string
			if tokenType == IDENTIFIER {
				repr = "@" + token
			} else {
				repr = "@(" + token + ")"
			}

			errors.Add(repr, err.Error())
		}
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}

// HasExpressions returns whether the given template contains any expressions or identifiers
func HasExpressions(template string, allowedTopLevels []string) bool {
	found := false
	VisitTemplate(template, allowedTopLevels, false, func(tokenType XTokenType, token string) error {
		switch tokenType {
		case IDENTIFIER, EXPRESSION:
			found = true
			return nil
		}
		return nil
	})
	return found
}
