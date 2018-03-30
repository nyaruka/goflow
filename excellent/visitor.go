package excellent

import (
	"bytes"
	"fmt"
	"strings"

	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/nyaruka/goflow/excellent/gen"
	"github.com/nyaruka/goflow/utils"
	"github.com/shopspring/decimal"
)

type Visitor struct {
	gen.BaseExcellent2Visitor
	env      utils.Environment
	resolver utils.VariableResolver
}

// NewVisitor creates a new Excellent visitor
func NewVisitor(env utils.Environment, resolver utils.VariableResolver) *Visitor {
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
	dec, _ := utils.ToDecimal(v.env, ctx.GetText())
	return dec
}

// VisitDotLookup deals with lookups like foo.0 or foo.bar
func (v *Visitor) VisitDotLookup(ctx *gen.DotLookupContext) interface{} {
	context := v.Visit(ctx.Atom(0))
	lookup := ctx.Atom(1).GetText()
	return utils.ResolveVariable(v.env, context, lookup)
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
	return strings.Replace(unquoted, "\"\"", "\"", -1)
}

// VisitFunctionCall deals with function calls like TITLE(foo.bar)
func (v *Visitor) VisitFunctionCall(ctx *gen.FunctionCallContext) interface{} {
	functionName := strings.ToLower(ctx.Fnname().GetText())

	var function XFunction
	var found bool

	// this is a test, look it up from those
	if strings.HasPrefix(functionName, "has_") {
		function, found = XTESTS[functionName]
	} else {
		function, found = XFUNCTIONS[functionName]
	}

	if !found {
		return fmt.Errorf("No function with name '%s'", functionName)
	}

	var params []interface{}
	if ctx.Parameters() != nil {
		params = v.Visit(ctx.Parameters()).([]interface{})
	}
	val := function(v.env, params...)
	return val
}

// VisitTrue deals with the "true" literal
func (v *Visitor) VisitTrue(ctx *gen.TrueContext) interface{} {
	return true
}

// VisitFalse deals with the "false" literal
func (v *Visitor) VisitFalse(ctx *gen.FalseContext) interface{} {
	return false
}

// VisitArrayLookup deals with lookups such as foo[5]
func (v *Visitor) VisitArrayLookup(ctx *gen.ArrayLookupContext) interface{} {
	context := v.Visit(ctx.Atom())
	lookup, err := utils.ToString(v.env, v.Visit(ctx.Expression()))
	if err != nil {
		return err
	}

	return utils.ResolveVariable(v.env, context, lookup)
}

// VisitContextReference deals with references to variables in the context such as "foo"
func (v *Visitor) VisitContextReference(ctx *gen.ContextReferenceContext) interface{} {
	key := strings.ToLower(ctx.GetText())
	val := utils.ResolveVariable(v.env, v.resolver, key)
	return val
}

// VisitParentheses deals with expressions in parentheses such as (1+2)
func (v *Visitor) VisitParentheses(ctx *gen.ParenthesesContext) interface{} {
	return v.Visit(ctx.Expression())
}

// VisitNegation deals with negations such as -5
func (v *Visitor) VisitNegation(ctx *gen.NegationContext) interface{} {
	dec, err := utils.ToDecimal(v.env, v.Visit(ctx.Expression()))
	if err != nil {
		return err
	}

	if ctx.MINUS() != nil {
		return dec.Neg()
	}
	return dec
}

// VisitExponent deals with exponenets such as 5^5
func (v *Visitor) VisitExponent(ctx *gen.ExponentContext) interface{} {
	arg1, err := utils.ToDecimal(v.env, v.Visit(ctx.Expression(0)))
	if err != nil {
		return err
	}

	arg2, err := utils.ToDecimal(v.env, v.Visit(ctx.Expression(1)))
	if err != nil {
		return err
	}

	return arg1.Pow(arg2)
}

// VisitConcatenation deals with string concatenations like "foo" & "bar"
func (v *Visitor) VisitConcatenation(ctx *gen.ConcatenationContext) interface{} {
	arg1, err := utils.ToString(v.env, v.Visit(ctx.Expression(0)))
	if err != nil {
		return err
	}

	arg2, err := utils.ToString(v.env, v.Visit(ctx.Expression(1)))
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	buffer.WriteString(arg1)
	buffer.WriteString(arg2)

	return buffer.String()
}

// VisitAdditionOrSubtraction deals with addition and subtraction like 5+5 and 5-3
func (v *Visitor) VisitAdditionOrSubtraction(ctx *gen.AdditionOrSubtractionContext) interface{} {
	arg1 := v.Visit(ctx.Expression(0))
	arg2 := v.Visit(ctx.Expression(1))

	arg1Dec, err := utils.ToDecimal(v.env, arg1)
	if err != nil {
		return err
	}

	arg2Dec, err := utils.ToDecimal(v.env, arg2)
	if err != nil {
		return err
	}

	if ctx.PLUS() != nil {
		return arg1Dec.Add(arg2Dec)
	}
	return arg1Dec.Sub(arg2Dec)
}

// VisitEquality deals with equality or inequality tests 5 = 5 and 5 != 5
func (v *Visitor) VisitEquality(ctx *gen.EqualityContext) interface{} {
	arg1 := v.Visit(ctx.Expression(0))
	arg2 := v.Visit(ctx.Expression(1))

	cmp, err := utils.Compare(v.env, arg1, arg2)
	if err != nil {
		return err
	}

	if ctx.EQ() != nil {
		return cmp == 0
	}

	return cmp != 0
}

// VisitAtomReference deals with visiting a single atom in our expression
func (v *Visitor) VisitAtomReference(ctx *gen.AtomReferenceContext) interface{} {
	return v.Visit(ctx.Atom())
}

// VisitMultiplicationOrDivision deals with division and multiplication such as 5*5 or 5/2
func (v *Visitor) VisitMultiplicationOrDivision(ctx *gen.MultiplicationOrDivisionContext) interface{} {
	arg1, err := utils.ToDecimal(v.env, v.Visit(ctx.Expression(0)))
	if err != nil {
		return err
	}

	arg2, err := utils.ToDecimal(v.env, v.Visit(ctx.Expression(1)))
	if err != nil {
		return err
	}

	if ctx.TIMES() != nil {
		return arg1.Mul(arg2)
	}

	// division!
	if arg2.Equals(decimal.Zero) {
		return fmt.Errorf("Division by zero")
	}

	return arg1.Div(arg2)
}

// VisitComparison deals with visiting a comparison between two values, such as 5<3 or 3>5
func (v *Visitor) VisitComparison(ctx *gen.ComparisonContext) interface{} {
	arg1 := v.Visit(ctx.Expression(0))
	arg2 := v.Visit(ctx.Expression(1))

	cmp, err := utils.Compare(v.env, arg1, arg2)
	if err != nil {
		return err
	}

	switch {
	case ctx.LT() != nil:
		return cmp < 0
	case ctx.LTE() != nil:
		return cmp <= 0
	case ctx.GTE() != nil:
		return cmp >= 0
	case ctx.GT() != nil:
		return cmp > 0
	}

	return fmt.Errorf("Unknown comparison operator: %s", ctx.GetText())
}

// VisitFunctionParameters deals with the parameters to a function call
func (v *Visitor) VisitFunctionParameters(ctx *gen.FunctionParametersContext) interface{} {
	expressions := ctx.AllExpression()
	params := make([]interface{}, len(expressions))

	for i := range expressions {
		params[i] = v.Visit(expressions[i])
	}
	return params
}
