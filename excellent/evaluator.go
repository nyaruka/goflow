package excellent

import (
	"strconv"
	"strings"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/functions"
	"github.com/nyaruka/goflow/excellent/gen"
	"github.com/nyaruka/goflow/excellent/operators"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Escaping is a function applied to expressions in a template after they've been evaluated
type Escaping func(string) string

// EvaluateTemplate evaluates the passed in template
func EvaluateTemplate(env envs.Environment, context *types.XObject, template string, escaping Escaping) (string, error) {
	var buf strings.Builder

	err := VisitTemplate(template, context.Properties(), func(tokenType XTokenType, token string) error {
		switch tokenType {
		case BODY:
			buf.WriteString(token)
		case IDENTIFIER, EXPRESSION:
			value := EvaluateExpression(env, context, token)

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
func EvaluateTemplateValue(env envs.Environment, context *types.XObject, template string) (types.XValue, error) {
	template = strings.TrimSpace(template)
	scanner := NewXScanner(strings.NewReader(template), context.Properties())

	// parse our first token
	tokenType, token := scanner.Scan()

	// try to scan to our next token
	nextTT, _ := scanner.Scan()

	// if we only have an identifier or an expression, evaluate it on its own
	if nextTT == EOF {
		switch tokenType {
		case IDENTIFIER, EXPRESSION:
			return EvaluateExpression(env, context, token), nil
		}
	}

	// otherwise fallback to full template evaluation
	asStr, err := EvaluateTemplate(env, context, template, nil)
	return types.NewXText(asStr), err
}

// EvaluateExpression evalutes the passed in Excellent expression, returning the typed value it evaluates to,
// which might be an error, e.g. "2 / 3" or "contact.fields.age"
func EvaluateExpression(env envs.Environment, context *types.XObject, expression string) types.XValue {
	visitor := newEvaluationVisitor(env, context)
	output, err := VisitExpression(expression, visitor)
	if err != nil {
		return types.NewXError(err)
	}

	return toXValue(output)
}

// visitor which evaluates each part of an expression as a value
type visitor struct {
	gen.BaseExcellent2Visitor

	env     envs.Environment
	context *types.XObject
}

// creates a new visitor for evaluation
func newEvaluationVisitor(env envs.Environment, context *types.XObject) *visitor {
	return &visitor{env: env, context: context}
}

// Visit the top level parse tree
func (v *visitor) Visit(tree antlr.ParseTree) interface{} {
	return tree.Accept(v)
}

// VisitParse handles our top level parser
func (v *visitor) VisitParse(ctx *gen.ParseContext) interface{} {
	return v.Visit(ctx.Expression())
}

// VisitTextLiteral deals with string literals such as "asdf"
func (v *visitor) VisitTextLiteral(ctx *gen.TextLiteralContext) interface{} {
	value := ctx.GetText()

	// unquote, this takes care of escape sequences as well
	unquoted, err := strconv.Unquote(value)

	// if we had an error, just strip surrounding quotes
	if err != nil {
		unquoted = value[1 : len(value)-1]
	}

	return types.NewXText(unquoted)
}

// VisitNumberLiteral deals with numbers like 123 or 1.5
func (v *visitor) VisitNumberLiteral(ctx *gen.NumberLiteralContext) interface{} {
	return types.RequireXNumberFromString(ctx.GetText())
}

// VisitContextReference deals with identifiers which are function names or root variables in the context
func (v *visitor) VisitContextReference(ctx *gen.ContextReferenceContext) interface{} {
	name := strings.ToLower(ctx.GetText())

	// first of all try to look this up as a function
	function := functions.Lookup(name)
	if function != nil {
		return toXValue(function)
	}

	value, exists := v.context.Get(name)
	if !exists {
		return types.NewXErrorf("context has no property '%s'", name)
	}

	return value
}

// VisitDotLookup deals with property lookups like foo.bar
func (v *visitor) VisitDotLookup(ctx *gen.DotLookupContext) interface{} {
	container := toXValue(v.Visit(ctx.Atom()))
	if types.IsXError(container) {
		return container
	}

	var lookup types.XText

	if ctx.NAME() != nil {
		lookup = types.NewXText(ctx.NAME().GetText())
	} else {
		lookup = types.NewXText(ctx.INTEGER().GetText())
	}

	return resolveLookup(v.env, container, lookup, lookupNotationDot)
}

// VisitArrayLookup deals with lookups such as foo[5] or foo["key with spaces"]
func (v *visitor) VisitArrayLookup(ctx *gen.ArrayLookupContext) interface{} {
	container := toXValue(v.Visit(ctx.Atom()))
	if types.IsXError(container) {
		return container
	}

	lookup := toXValue(v.Visit(ctx.Expression()))

	return resolveLookup(v.env, container, lookup, lookupNotationArray)
}

// VisitFunctionCall deals with function calls like TITLE(foo.bar)
func (v *visitor) VisitFunctionCall(ctx *gen.FunctionCallContext) interface{} {
	function := toXValue(v.Visit(ctx.Atom()))
	if types.IsXError(function) {
		return function
	}

	asFunction, isFunction := function.(types.XFunction)
	if !isFunction {
		return types.NewXErrorf("%s is not a function", ctx.Atom().GetText())
	}

	name := strings.ToLower(ctx.Atom().GetText())

	var params []types.XValue
	if ctx.Parameters() != nil {
		params, _ = v.Visit(ctx.Parameters()).([]types.XValue)
	}

	return functions.Call(v.env, name, asFunction, params)
}

// VisitTrue deals with the `true` reserved word
func (v *visitor) VisitTrue(ctx *gen.TrueContext) interface{} {
	return types.XBooleanTrue
}

// VisitFalse deals with the `false` reserved word
func (v *visitor) VisitFalse(ctx *gen.FalseContext) interface{} {
	return types.XBooleanFalse
}

// VisitNull deals with the `null` reserved word
func (v *visitor) VisitNull(ctx *gen.NullContext) interface{} {
	return nil
}

// VisitParentheses deals with expressions in parentheses such as (1+2)
func (v *visitor) VisitParentheses(ctx *gen.ParenthesesContext) interface{} {
	return v.Visit(ctx.Expression())
}

// VisitNegation deals with negations such as -5
func (v *visitor) VisitNegation(ctx *gen.NegationContext) interface{} {
	arg := toXValue(v.Visit(ctx.Expression()))

	return operators.Negate(v.env, arg)
}

// VisitExponent deals with exponenets such as 5^5
func (v *visitor) VisitExponent(ctx *gen.ExponentContext) interface{} {
	arg1 := toXValue(v.Visit(ctx.Expression(0)))
	arg2 := toXValue(v.Visit(ctx.Expression(1)))

	return operators.Exponent(v.env, arg1, arg2)
}

// VisitConcatenation deals with string concatenations like "foo" & "bar"
func (v *visitor) VisitConcatenation(ctx *gen.ConcatenationContext) interface{} {
	arg1 := toXValue(v.Visit(ctx.Expression(0)))
	arg2 := toXValue(v.Visit(ctx.Expression(1)))

	return operators.Concatenate(v.env, arg1, arg2)
}

// VisitAdditionOrSubtraction deals with addition and subtraction like 5+5 and 5-3
func (v *visitor) VisitAdditionOrSubtraction(ctx *gen.AdditionOrSubtractionContext) interface{} {
	arg1 := toXValue(v.Visit(ctx.Expression(0)))
	arg2 := toXValue(v.Visit(ctx.Expression(1)))

	if ctx.PLUS() != nil {
		return operators.Add(v.env, arg1, arg2)
	}
	return operators.Subtract(v.env, arg1, arg2)
}

// VisitMultiplicationOrDivision deals with division and multiplication such as 5*5 or 5/2
func (v *visitor) VisitMultiplicationOrDivision(ctx *gen.MultiplicationOrDivisionContext) interface{} {
	arg1 := toXValue(v.Visit(ctx.Expression(0)))
	arg2 := toXValue(v.Visit(ctx.Expression(1)))

	if ctx.TIMES() != nil {
		return operators.Multiply(v.env, arg1, arg2)
	}
	return operators.Divide(v.env, arg1, arg2)
}

// VisitEquality deals with equality or inequality tests 5 = 5 and 5 != 5
func (v *visitor) VisitEquality(ctx *gen.EqualityContext) interface{} {
	arg1 := toXValue(v.Visit(ctx.Expression(0)))
	arg2 := toXValue(v.Visit(ctx.Expression(1)))

	if ctx.EQ() != nil {
		return operators.Equal(v.env, arg1, arg2)
	}
	return operators.NotEqual(v.env, arg1, arg2)
}

// VisitComparison deals with visiting a comparison between two values, such as 5<3 or 3>5
func (v *visitor) VisitComparison(ctx *gen.ComparisonContext) interface{} {
	arg1 := toXValue(v.Visit(ctx.Expression(0)))
	arg2 := toXValue(v.Visit(ctx.Expression(1)))

	switch {
	case ctx.LT() != nil:
		return operators.LessThan(v.env, arg1, arg2)
	case ctx.LTE() != nil:
		return operators.LessThanOrEqual(v.env, arg1, arg2)
	case ctx.GTE() != nil:
		return operators.GreaterThanOrEqual(v.env, arg1, arg2)
	default:
		return operators.GreaterThan(v.env, arg1, arg2)
	}
}

// VisitAtomReference deals with visiting a single atom in our expression
func (v *visitor) VisitAtomReference(ctx *gen.AtomReferenceContext) interface{} {
	return v.Visit(ctx.Atom())
}

// VisitFunctionParameters deals with the parameters to a function call
func (v *visitor) VisitFunctionParameters(ctx *gen.FunctionParametersContext) interface{} {
	expressions := ctx.AllExpression()
	params := make([]types.XValue, len(expressions))

	for i := range expressions {
		params[i] = toXValue(v.Visit(expressions[i]))
	}
	return params
}

// convenience utility to convert the given value to an XValue. Might be able to rewrite the visitor in future
// to only pass around XValues and then wouldn't need this
func toXValue(val interface{}) types.XValue {
	asX, isXValue := val.(types.XValue)
	if !isXValue && !utils.IsNil(val) {
		panic("Attempt to convert a non XValue to an XValue")
	}
	return asX
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
