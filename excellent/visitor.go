package excellent

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/nyaruka/goflow/excellent/functions"
	"github.com/nyaruka/goflow/excellent/gen"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type Visitor struct {
	gen.BaseExcellent2Visitor
	env      utils.Environment
	resolver types.XValue
}

// NewVisitor creates a new Excellent visitor
func NewVisitor(env utils.Environment, resolver types.XValue) *Visitor {
	visitor := Visitor{env: env, resolver: resolver}
	return &visitor
}

// Visit the top level parse tree
func (v *Visitor) Visit(tree antlr.ParseTree) interface{} {
	return tree.Accept(v)
}

// VisitParse handles our top level parser
func (v *Visitor) VisitParse(ctx *gen.ParseContext) interface{} {
	return v.Visit(ctx.Expression())
}

// VisitTextLiteral deals with string literals such as "asdf"
func (v *Visitor) VisitTextLiteral(ctx *gen.TextLiteralContext) interface{} {
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
func (v *Visitor) VisitNumberLiteral(ctx *gen.NumberLiteralContext) interface{} {
	return types.RequireXNumberFromString(ctx.GetText())
}

// VisitDotLookup deals with lookups like foo.0 or foo.bar
func (v *Visitor) VisitDotLookup(ctx *gen.DotLookupContext) interface{} {
	context := toXValue(v.Visit(ctx.Atom(0)))
	if types.IsXError(context) {
		return context
	}

	lookup := ctx.Atom(1).GetText()
	return ResolveValue(v.env, context, lookup)
}

// VisitFunctionCall deals with function calls like TITLE(foo.bar)
func (v *Visitor) VisitFunctionCall(ctx *gen.FunctionCallContext) interface{} {
	functionName := strings.ToLower(ctx.Fnname().GetText())

	var function functions.XFunction
	var found bool

	function, found = functions.XFUNCTIONS[functionName]
	if !found {
		return types.NewXErrorf("no function with name '%s'", functionName)
	}

	var params []types.XValue
	if ctx.Parameters() != nil {
		params, _ = v.Visit(ctx.Parameters()).([]types.XValue)
	}

	val := function(v.env, params...)

	// if function returned an error, wrap the error with the function name
	if types.IsXError(val) {
		return types.NewXErrorf("error calling %s: %s", strings.ToUpper(functionName), val.(types.XError).Error())
	}

	return val
}

// VisitTrue deals with the `true` reserved word
func (v *Visitor) VisitTrue(ctx *gen.TrueContext) interface{} {
	return types.XBooleanTrue
}

// VisitFalse deals with the `false` reserved word
func (v *Visitor) VisitFalse(ctx *gen.FalseContext) interface{} {
	return types.XBooleanFalse
}

// VisitNull deals with the `null` reserved word
func (v *Visitor) VisitNull(ctx *gen.NullContext) interface{} {
	return nil
}

// VisitArrayLookup deals with lookups such as foo[5] or foo["key with spaces"]
func (v *Visitor) VisitArrayLookup(ctx *gen.ArrayLookupContext) interface{} {
	context := toXValue(v.Visit(ctx.Atom()))
	if types.IsXError(context) {
		return context
	}

	expression := toXValue(v.Visit(ctx.Expression()))

	lookup, xerr := types.ToXText(v.env, expression)
	if xerr != nil {
		return xerr
	}

	return ResolveValue(v.env, context, lookup.Native())
}

// VisitContextReference deals with references to variables in the context such as "foo"
func (v *Visitor) VisitContextReference(ctx *gen.ContextReferenceContext) interface{} {
	key := strings.ToLower(ctx.GetText())

	return ResolveValue(v.env, v.resolver, key)
}

// VisitParentheses deals with expressions in parentheses such as (1+2)
func (v *Visitor) VisitParentheses(ctx *gen.ParenthesesContext) interface{} {
	return v.Visit(ctx.Expression())
}

// VisitNegation deals with negations such as -5
func (v *Visitor) VisitNegation(ctx *gen.NegationContext) interface{} {
	arg := toXValue(v.Visit(ctx.Expression()))

	number, xerr := types.ToXNumber(v.env, arg)
	if xerr != nil {
		return xerr
	}

	return types.NewXNumber(number.Native().Neg())
}

// VisitExponent deals with exponenets such as 5^5
func (v *Visitor) VisitExponent(ctx *gen.ExponentContext) interface{} {
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
func (v *Visitor) VisitConcatenation(ctx *gen.ConcatenationContext) interface{} {
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

	var buffer bytes.Buffer
	buffer.WriteString(str1.Native())
	buffer.WriteString(str2.Native())

	return types.NewXText(buffer.String())
}

// VisitAdditionOrSubtraction deals with addition and subtraction like 5+5 and 5-3
func (v *Visitor) VisitAdditionOrSubtraction(ctx *gen.AdditionOrSubtractionContext) interface{} {
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
func (v *Visitor) VisitEquality(ctx *gen.EqualityContext) interface{} {
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
func (v *Visitor) VisitAtomReference(ctx *gen.AtomReferenceContext) interface{} {
	return v.Visit(ctx.Atom())
}

// VisitMultiplicationOrDivision deals with division and multiplication such as 5*5 or 5/2
func (v *Visitor) VisitMultiplicationOrDivision(ctx *gen.MultiplicationOrDivisionContext) interface{} {
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
func (v *Visitor) VisitComparison(ctx *gen.ComparisonContext) interface{} {
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
func (v *Visitor) VisitFunctionParameters(ctx *gen.FunctionParametersContext) interface{} {
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
