package excellent

import (
	"strings"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
)

// Escaping is a function applied to expressions in a template after they've been evaluated
type Escaping func(string) string

type Context struct {
	root   *types.XObject
	parent *Context
}

func NewContext(root *types.XObject, parent *Context) *Context {
	return &Context{root: root, parent: parent}
}

func (c *Context) TopLevels() []string {
	return c.root.Properties()
}

func (c *Context) Get(name string) (types.XValue, bool) {
	v, exists := c.root.Get(name)
	if exists {
		return v, true
	}

	if c.parent != nil {
		return c.parent.Get(name)
	}

	return nil, false
}

// EvaluateTemplate evaluates the passed in template
func EvaluateTemplate(env envs.Environment, ctx *Context, template string, escaping Escaping) (string, error) {
	var buf strings.Builder

	err := VisitTemplate(template, ctx.TopLevels(), func(tokenType XTokenType, token string) error {
		switch tokenType {
		case BODY:
			buf.WriteString(token)
		case IDENTIFIER, EXPRESSION:
			value := EvaluateExpression(env, ctx, token)

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

// EvaluateTemplateValue is equivalent to EvaluateTemplate except in the case where the template contains
// a single identifier or expression, ie: "@contact" or "@(first(contact.urns))". In these cases we return
// the typed value from EvaluateExpression instead of stringifying the result.
func EvaluateTemplateValue(env envs.Environment, ctx *Context, template string) (types.XValue, error) {
	template = strings.TrimSpace(template)
	scanner := NewXScanner(strings.NewReader(template), ctx.TopLevels())

	// parse our first token
	tokenType, token := scanner.Scan()

	// try to scan to our next token
	nextTT, _ := scanner.Scan()

	// if we only have an identifier or an expression, evaluate it on its own
	if nextTT == EOF {
		switch tokenType {
		case IDENTIFIER, EXPRESSION:
			return EvaluateExpression(env, ctx, token), nil
		}
	}

	// otherwise fallback to full template evaluation
	asStr, err := EvaluateTemplate(env, ctx, template, nil)
	return types.NewXText(asStr), err
}

// EvaluateExpression evalutes the passed in Excellent expression, returning the typed value it evaluates to,
// which might be an error, e.g. "2 / 3" or "contact.fields.age"
func EvaluateExpression(env envs.Environment, ctx *Context, expression string) types.XValue {
	parsed, err := Parse(expression, nil)
	if err != nil {
		return types.NewXError(err)
	}
	return parsed.Evaluate(env, ctx)
}

type lookupNotation string

const (
	lookupNotationDot   lookupNotation = "dot"
	lookupNotationArray lookupNotation = "array"
)

func resolveLookup(env envs.Environment, container types.XValue, lookup types.XValue, notation lookupNotation) types.XValue {
	// if left-hand side is an array, then this is an index
	array, isArray := container.(*types.XArray)
	if isArray && array != nil {
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
		return array.Get(index)
	}

	// if left-hand side is an object, then this is a property lookup
	object, isObject := container.(*types.XObject)
	if isObject && object != nil {
		property, xerr := types.ToXText(env, lookup)
		if xerr != nil {
			return xerr
		}

		value, exists := object.Get(property.Native())

		// [] notation doesn't error for non-existent properties, . does
		if !exists && notation == lookupNotationDot {
			return types.NewXErrorf("%s has no property '%s'", types.Describe(container), property.Native())
		}

		return value
	}

	return types.NewXErrorf("%s doesn't support lookups", types.Describe(container))
}
