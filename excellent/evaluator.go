package excellent

import (
	"strconv"
	"strings"

	"github.com/nyaruka/goflow/excellent/functions"
	"github.com/nyaruka/goflow/excellent/gen"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// EvaluateTemplate evaluates the passed in template
func EvaluateTemplate(env utils.Environment, context types.XValue, template string, allowedTopLevels []string) (string, error) {
	var buf strings.Builder

	err := VisitTemplate(template, allowedTopLevels, func(tokenType XTokenType, token string) error {
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
			strValue, _ := types.ToXText(env, value)
			buf.WriteString(strValue.Native())
		}
		return nil
	})

	return buf.String(), err
}

// EvaluateTemplateValue is equivalent to EvaluateTemplate except in the case where the template contains
// a single identifier or expression, ie: "@contact" or "@(first(contact.urns))". In these cases we return
// the typed value from EvaluateExpression instead of stringifying the result.
func EvaluateTemplateValue(env utils.Environment, context types.XValue, template string, allowedTopLevels []string) (types.XValue, error) {
	template = strings.TrimSpace(template)
	scanner := NewXScanner(strings.NewReader(template), allowedTopLevels)

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
	asStr, err := EvaluateTemplate(env, context, template, allowedTopLevels)
	return types.NewXText(asStr), err
}

// EvaluateExpression evalutes the passed in Excellent expression, returning the typed value it evaluates to,
// which might be an error, e.g. "2 / 3" or "contact.fields.age"
func EvaluateExpression(env utils.Environment, context types.XValue, expression string) types.XValue {
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

	env      utils.Environment
	resolver types.XValue
}

// creates a new visitor for evaluation
func newEvaluationVisitor(env utils.Environment, resolver types.XValue) *visitor {
	return &visitor{env: env, resolver: resolver}
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

	return types.Resolve(v.env, v.resolver, name)
}

// VisitDotLookup deals with lookups like foo.0 or foo.bar
func (v *visitor) VisitDotLookup(ctx *gen.DotLookupContext) interface{} {
	context := toXValue(v.Visit(ctx.Atom()))
	if types.IsXError(context) {
		return context
	}

	var lookup string
	if ctx.NAME() != nil {
		lookup = ctx.NAME().GetText()
	} else {
		lookup = ctx.NUMBER().GetText()
	}

	return types.Resolve(v.env, context, lookup)
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

// VisitArrayLookup deals with lookups such as foo[5] or foo["key with spaces"]
func (v *visitor) VisitArrayLookup(ctx *gen.ArrayLookupContext) interface{} {
	context := toXValue(v.Visit(ctx.Atom()))
	if types.IsXError(context) {
		return context
	}

	expression := toXValue(v.Visit(ctx.Expression()))

	// if left-hand side is an array, then this is an index
	asArray, isArray := context.(types.XArray)
	if isArray {
		index, xerr := types.ToInteger(v.env, expression)
		if xerr != nil {
			return xerr
		}

		return lookupIndex(v.env, asArray, index)
	}

	// if not it is a property lookup so stringify the key
	lookup, xerr := types.ToXText(v.env, expression)
	if xerr != nil {
		return xerr
	}

	return types.Resolve(v.env, context, lookup.Native())
}

// VisitParentheses deals with expressions in parentheses such as (1+2)
func (v *visitor) VisitParentheses(ctx *gen.ParenthesesContext) interface{} {
	return v.Visit(ctx.Expression())
}

// VisitNegation deals with negations such as -5
func (v *visitor) VisitNegation(ctx *gen.NegationContext) interface{} {
	arg := toXValue(v.Visit(ctx.Expression()))

	number, xerr := types.ToXNumber(v.env, arg)
	if xerr != nil {
		return xerr
	}

	return types.NewXNumber(number.Native().Neg())
}

// VisitExponent deals with exponenets such as 5^5
func (v *visitor) VisitExponent(ctx *gen.ExponentContext) interface{} {
	arg1 := toXValue(v.Visit(ctx.Expression(0)))
	arg2 := toXValue(v.Visit(ctx.Expression(1)))

	num1, xerr := types.ToXNumber(v.env, arg1)
	if xerr != nil {
		return xerr
	}
	num2, xerr := types.ToXNumber(v.env, arg2)
	if xerr != nil {
		return xerr
	}

	return types.NewXNumber(num1.Native().Pow(num2.Native()))
}

// VisitConcatenation deals with string concatenations like "foo" & "bar"
func (v *visitor) VisitConcatenation(ctx *gen.ConcatenationContext) interface{} {
	arg1 := toXValue(v.Visit(ctx.Expression(0)))
	arg2 := toXValue(v.Visit(ctx.Expression(1)))

	str1, xerr := types.ToXText(v.env, arg1)
	if xerr != nil {
		return xerr
	}
	str2, xerr := types.ToXText(v.env, arg2)
	if xerr != nil {
		return xerr
	}

	var buffer strings.Builder
	buffer.WriteString(str1.Native())
	buffer.WriteString(str2.Native())

	return types.NewXText(buffer.String())
}

// VisitAdditionOrSubtraction deals with addition and subtraction like 5+5 and 5-3
func (v *visitor) VisitAdditionOrSubtraction(ctx *gen.AdditionOrSubtractionContext) interface{} {
	arg1 := toXValue(v.Visit(ctx.Expression(0)))
	arg2 := toXValue(v.Visit(ctx.Expression(1)))

	num1, xerr := types.ToXNumber(v.env, arg1)
	if xerr != nil {
		return xerr
	}
	num2, xerr := types.ToXNumber(v.env, arg2)
	if xerr != nil {
		return xerr
	}

	if ctx.PLUS() != nil {
		return types.NewXNumber(num1.Native().Add(num2.Native()))
	}
	return types.NewXNumber(num1.Native().Sub(num2.Native()))
}

// VisitEquality deals with equality or inequality tests 5 = 5 and 5 != 5
func (v *visitor) VisitEquality(ctx *gen.EqualityContext) interface{} {
	arg1 := toXValue(v.Visit(ctx.Expression(0)))
	arg2 := toXValue(v.Visit(ctx.Expression(1)))

	str1, xerr := types.ToXText(v.env, arg1)
	if xerr != nil {
		return xerr
	}
	str2, xerr := types.ToXText(v.env, arg2)
	if xerr != nil {
		return xerr
	}

	isEqual := str1.Equals(str2)

	if ctx.EQ() != nil {
		return types.NewXBoolean(isEqual)
	}

	return types.NewXBoolean(!isEqual)
}

// VisitAtomReference deals with visiting a single atom in our expression
func (v *visitor) VisitAtomReference(ctx *gen.AtomReferenceContext) interface{} {
	return v.Visit(ctx.Atom())
}

// VisitMultiplicationOrDivision deals with division and multiplication such as 5*5 or 5/2
func (v *visitor) VisitMultiplicationOrDivision(ctx *gen.MultiplicationOrDivisionContext) interface{} {
	arg1 := toXValue(v.Visit(ctx.Expression(0)))
	arg2 := toXValue(v.Visit(ctx.Expression(1)))

	num1, xerr := types.ToXNumber(v.env, arg1)
	if xerr != nil {
		return xerr
	}
	num2, xerr := types.ToXNumber(v.env, arg2)
	if xerr != nil {
		return xerr
	}

	if ctx.TIMES() != nil {
		return types.NewXNumber(num1.Native().Mul(num2.Native()))
	}

	// division!
	if num2.Equals(types.XNumberZero) {
		return types.NewXErrorf("division by zero")
	}

	return types.NewXNumber(num1.Native().Div(num2.Native()))
}

// VisitComparison deals with visiting a comparison between two values, such as 5<3 or 3>5
func (v *visitor) VisitComparison(ctx *gen.ComparisonContext) interface{} {
	arg1 := toXValue(v.Visit(ctx.Expression(0)))
	arg2 := toXValue(v.Visit(ctx.Expression(1)))

	num1, xerr := types.ToXNumber(v.env, arg1)
	if xerr != nil {
		return xerr
	}
	num2, xerr := types.ToXNumber(v.env, arg2)
	if xerr != nil {
		return xerr
	}

	cmp := num1.Compare(num2)

	switch {
	case ctx.LT() != nil:
		return types.NewXBoolean(cmp < 0)
	case ctx.LTE() != nil:
		return types.NewXBoolean(cmp <= 0)
	case ctx.GTE() != nil:
		return types.NewXBoolean(cmp >= 0)
	default: // ctx.GT() != nil
		return types.NewXBoolean(cmp > 0)
	}
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

// lookup an index on the given value
func lookupIndex(env utils.Environment, value types.XValue, index int) types.XValue {
	indexable, isIndexable := value.(types.XIndexable)

	if !isIndexable || utils.IsNil(indexable) {
		return types.NewXErrorf("%s is not indexable", value.Describe())
	}

	if index >= indexable.Length() || index < -indexable.Length() {
		return types.NewXErrorf("index %d out of range for %d items", index, indexable.Length())
	}
	if index < 0 {
		index += indexable.Length()
	}
	return indexable.Index(index)
}
