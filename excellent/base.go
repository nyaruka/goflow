package excellent

import (
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	gen "github.com/nyaruka/goflow/antlr/gen/excellent3"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
)

// Evaluator evaluates templates and expressions.
type Evaluator struct {
	onDeprecated func(types.XValue)
}

// NewEvaluator creates a new evaluator
func NewEvaluator(opts ...func(*Evaluator)) *Evaluator {
	e := &Evaluator{}
	for _, o := range opts {
		o(e)
	}
	return e
}

// WithDeprecatedCallback sets a callback to be called when a deprecated value in the context is accessed.
func WithDeprecatedCallback(onDeprecated func(types.XValue)) func(*Evaluator) {
	return func(e *Evaluator) {
		e.onDeprecated = onDeprecated
	}
}

// Escaping is a function applied to expressions in a template after they've been evaluated
type Escaping func(string) string

// Template evaluates the passed in template
func (e *Evaluator) Template(env envs.Environment, ctx *types.XObject, template string, escaping Escaping) (string, error) {
	var buf strings.Builder

	err := VisitTemplate(template, ctx.Properties(), func(tokenType XTokenType, token string) error {
		switch tokenType {
		case BODY:
			buf.WriteString(token)
		case IDENTIFIER, EXPRESSION:
			value := e.Expression(env, ctx, token)

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

	return buf.String(), err
}

// TemplateValue is equivalent to Template except in the case where the template contains
// a single identifier or expression, ie: "@contact" or "@(first(contact.urns))". In these cases we return
// the typed value from EvaluateExpression instead of stringifying the result.
func (e *Evaluator) TemplateValue(env envs.Environment, ctx *types.XObject, template string) (types.XValue, error) {
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
			return e.Expression(env, ctx, token), nil
		}
	}

	// otherwise fallback to full template evaluation
	asStr, err := e.Template(env, ctx, template, nil)
	return types.NewXText(asStr), err
}

// Expression evalutes the passed in Excellent expression, returning the typed value it evaluates to,
// which might be an error, e.g. "2 / 3" or "contact.fields.age"
func (e *Evaluator) Expression(env envs.Environment, ctx *types.XObject, expression string) types.XValue {
	parsed, err := Parse(expression, nil)
	if err != nil {
		return types.NewXError(err)
	}

	scope := NewScope(ctx, nil)

	return parsed.Evaluate(e, env, scope)
}

type lookupNotation string

const (
	lookupNotationDot   lookupNotation = "dot"
	lookupNotationArray lookupNotation = "array"
)

func (e *Evaluator) onContextRead(val types.XValue) {
	if !utils.IsNil(val) && val.Deprecated() != "" && e.onDeprecated != nil {
		e.onDeprecated(val)
	}
}

func (e *Evaluator) resolveLookup(env envs.Environment, container types.XValue, lookup types.XValue, notation lookupNotation) types.XValue {
	array, isArray := container.(*types.XArray)
	object, isObject := container.(*types.XObject)
	var resolved types.XValue

	if isArray && array != nil {
		// if left-hand side is an array, then this is an index
		index, xerr := types.ToInteger(env, lookup)
		if xerr != nil {
			return xerr
		}

		if index >= array.Count() || index < -array.Count() {
			return types.NewXErrorf("index %d out of range for %d items", index, array.Count())
		}
		if index < 0 {
			index += array.Count()
		}

		resolved = array.Get(index)

	} else if isObject && object != nil {
		// if left-hand side is an object, then this is a property lookup
		property, xerr := types.ToXText(env, lookup)
		if xerr != nil {
			return xerr
		}

		value, exists := object.Get(property.Native())

		// [] notation doesn't error for non-existent properties, . does
		if !exists && notation == lookupNotationDot {
			return types.NewXErrorf("%s has no property '%s'", types.Describe(container), property.Native())
		}

		resolved = value

	} else {
		return types.NewXErrorf("%s doesn't support lookups", types.Describe(container))
	}

	e.onContextRead(resolved)

	return resolved
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
func VisitTemplate(template string, allowedTopLevels []string, callback func(XTokenType, string) error) error {
	// nothing todo for an empty template
	if template == "" {
		return nil
	}

	scanner := NewXScanner(strings.NewReader(template), allowedTopLevels)
	errors := NewTemplateErrors()

	for tokenType, token := scanner.Scan(); tokenType != EOF; tokenType, token = scanner.Scan() {
		err := callback(tokenType, token)
		if err != nil {
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
	VisitTemplate(template, allowedTopLevels, func(tokenType XTokenType, token string) error {
		switch tokenType {
		case IDENTIFIER, EXPRESSION:
			found = true
			return nil
		}
		return nil
	})
	return found
}
