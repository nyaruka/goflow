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
	"github.com/shopspring/decimal"
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

// VisitDecimalLiteral deals with decimals like 1.5
func (v *Visitor) VisitDecimalLiteral(ctx *gen.DecimalLiteralContext) interface{} {
	return types.RequireXNumberFromString(ctx.GetText())
}

// VisitDotLookup deals with lookups like foo.0 or foo.bar
func (v *Visitor) VisitDotLookup(ctx *gen.DotLookupContext) interface{} {
	context, _ := v.Visit(ctx.Atom(0)).(types.XValue)
	if types.IsError(context) {
		return context
	}

	lookup := ctx.Atom(1).GetText()
	return ResolveXValue(v.env, context, lookup)
}

// VisitStringLiteral deals with string literals such as "asdf"
func (v *Visitor) VisitStringLiteral(ctx *gen.StringLiteralContext) interface{} {
	value := ctx.GetText()

	// unquote, this takes care of escape sequences as well
	unquoted, err := strconv.Unquote(value)

	// if we had an error, just strip surrounding quotes
	if err != nil {
		unquoted = value[1 : len(value)-1]
	}

	// replace "" with "
	unquoted = strings.Replace(unquoted, "\"\"", "\"", -1)

	return types.NewXString(unquoted)
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
	return val
}

// VisitTrue deals with the "true" literal
func (v *Visitor) VisitTrue(ctx *gen.TrueContext) interface{} {
	return types.XBoolTrue
}

// VisitFalse deals with the "false" literal
func (v *Visitor) VisitFalse(ctx *gen.FalseContext) interface{} {
	return types.XBoolFalse
}

// VisitArrayLookup deals with lookups such as foo[5]
func (v *Visitor) VisitArrayLookup(ctx *gen.ArrayLookupContext) interface{} {
	context, _ := v.Visit(ctx.Atom()).(types.XValue)
	if types.IsError(context) {
		return context
	}

	expression, _ := v.Visit(ctx.Expression()).(types.XValue)
	if types.IsError(expression) {
		return expression
	}

	lookup := types.ToXString(expression).Native()

	return ResolveXValue(v.env, context, lookup)
}

// VisitContextReference deals with references to variables in the context such as "foo"
func (v *Visitor) VisitContextReference(ctx *gen.ContextReferenceContext) interface{} {
	key := strings.ToLower(ctx.GetText())

	return ResolveXValue(v.env, v.resolver, key)
}

// VisitParentheses deals with expressions in parentheses such as (1+2)
func (v *Visitor) VisitParentheses(ctx *gen.ParenthesesContext) interface{} {
	return v.Visit(ctx.Expression())
}

// VisitNegation deals with negations such as -5
func (v *Visitor) VisitNegation(ctx *gen.NegationContext) interface{} {
	arg, _ := v.Visit(ctx.Expression()).(types.XValue)
	if types.IsError(arg) {
		return arg
	}

	number, err := types.ToXNumber(arg)
	if err != nil {
		return types.NewXError(err)
	}

	if ctx.MINUS() != nil {
		return types.NewXNumber(number.Native().Neg())
	}
	return number
}

// VisitExponent deals with exponenets such as 5^5
func (v *Visitor) VisitExponent(ctx *gen.ExponentContext) interface{} {
	arg1, _ := v.Visit(ctx.Expression(0)).(types.XValue)
	if types.IsError(arg1) {
		return arg1
	}

	arg2, _ := v.Visit(ctx.Expression(1)).(types.XValue)
	if types.IsError(arg2) {
		return arg2
	}

	num1, err := types.ToXNumber(arg1)
	if err != nil {
		return types.NewXError(err)
	}

	num2, err := types.ToXNumber(arg2)
	if err != nil {
		return types.NewXError(err)
	}

	return types.NewXNumber(num1.Native().Pow(num2.Native()))
}

// VisitConcatenation deals with string concatenations like "foo" & "bar"
func (v *Visitor) VisitConcatenation(ctx *gen.ConcatenationContext) interface{} {
	arg1, _ := v.Visit(ctx.Expression(0)).(types.XValue)
	if types.IsError(arg1) {
		return arg1
	}

	arg2, _ := v.Visit(ctx.Expression(1)).(types.XValue)
	if types.IsError(arg2) {
		return arg2
	}

	str1 := types.ToXString(arg1)
	str2 := types.ToXString(arg2)

	var buffer bytes.Buffer
	buffer.WriteString(str1.Native())
	buffer.WriteString(str2.Native())

	return types.NewXString(buffer.String())
}

// VisitAdditionOrSubtraction deals with addition and subtraction like 5+5 and 5-3
func (v *Visitor) VisitAdditionOrSubtraction(ctx *gen.AdditionOrSubtractionContext) interface{} {
	arg1, _ := v.Visit(ctx.Expression(0)).(types.XValue)
	if types.IsError(arg1) {
		return arg1
	}

	arg2, _ := v.Visit(ctx.Expression(1)).(types.XValue)
	if types.IsError(arg2) {
		return arg2
	}

	num1, err := types.ToXNumber(arg1)
	if err != nil {
		return types.NewXError(err)
	}

	num2, err := types.ToXNumber(arg2)
	if err != nil {
		return types.NewXError(err)
	}

	if ctx.PLUS() != nil {
		return types.NewXNumber(num1.Native().Add(num2.Native()))
	}
	return types.NewXNumber(num1.Native().Sub(num2.Native()))
}

// VisitEquality deals with equality or inequality tests 5 = 5 and 5 != 5
func (v *Visitor) VisitEquality(ctx *gen.EqualityContext) interface{} {
	arg1, _ := v.Visit(ctx.Expression(0)).(types.XValue)
	if types.IsError(arg1) {
		return arg1
	}

	arg2, _ := v.Visit(ctx.Expression(1)).(types.XValue)
	if types.IsError(arg2) {
		return arg2
	}

	num1, err := types.ToXNumber(arg1)
	if err != nil {
		return types.NewXError(err)
	}

	num2, err := types.ToXNumber(arg2)
	if err != nil {
		return types.NewXError(err)
	}

	if ctx.EQ() != nil {
		return types.NewXBool(num1.Native().Equal(num2.Native()))
	}

	return types.NewXBool(!num1.Native().Equal(num2.Native()))
}

// VisitAtomReference deals with visiting a single atom in our expression
func (v *Visitor) VisitAtomReference(ctx *gen.AtomReferenceContext) interface{} {
	return v.Visit(ctx.Atom())
}

// VisitMultiplicationOrDivision deals with division and multiplication such as 5*5 or 5/2
func (v *Visitor) VisitMultiplicationOrDivision(ctx *gen.MultiplicationOrDivisionContext) interface{} {
	arg1, _ := v.Visit(ctx.Expression(0)).(types.XValue)
	if types.IsError(arg1) {
		return arg1
	}

	arg2, _ := v.Visit(ctx.Expression(1)).(types.XValue)
	if types.IsError(arg2) {
		return arg2
	}

	num1, err := types.ToXNumber(arg1)
	if err != nil {
		return types.NewXError(err)
	}

	num2, err := types.ToXNumber(arg2)
	if err != nil {
		return types.NewXError(err)
	}

	if ctx.TIMES() != nil {
		return types.NewXNumber(num1.Native().Mul(num2.Native()))
	}

	// division!
	if num2.Native().Equals(decimal.Zero) {
		return types.NewXErrorf("division by zero")
	}

	return types.NewXNumber(num1.Native().Div(num2.Native()))
}

// VisitComparison deals with visiting a comparison between two values, such as 5<3 or 3>5
func (v *Visitor) VisitComparison(ctx *gen.ComparisonContext) interface{} {
	arg1, _ := v.Visit(ctx.Expression(0)).(types.XValue)
	if types.IsError(arg1) {
		return arg1
	}

	arg2, _ := v.Visit(ctx.Expression(1)).(types.XValue)
	if types.IsError(arg2) {
		return arg2
	}

	num1, err := types.ToXNumber(arg1)
	if err != nil {
		return types.NewXError(err)
	}

	num2, err := types.ToXNumber(arg2)
	if err != nil {
		return types.NewXError(err)
	}

	switch {
	case ctx.LT() != nil:
		return types.NewXBool(num1.Native().LessThan(num2.Native()))
	case ctx.LTE() != nil:
		return types.NewXBool(num1.Native().LessThanOrEqual(num2.Native()))
	case ctx.GTE() != nil:
		return types.NewXBool(num1.Native().GreaterThanOrEqual(num2.Native()))
	case ctx.GT() != nil:
		return types.NewXBool(num1.Native().GreaterThan(num2.Native()))
	}

	return types.NewXErrorf("unknown comparison operator: %s", ctx.GetText())
}

// VisitFunctionParameters deals with the parameters to a function call
func (v *Visitor) VisitFunctionParameters(ctx *gen.FunctionParametersContext) interface{} {
	expressions := ctx.AllExpression()
	params := make([]types.XValue, len(expressions))

	for i := range expressions {
		params[i], _ = v.Visit(expressions[i]).(types.XValue)
	}
	return params
}
