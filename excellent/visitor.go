package excellent

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	gen "github.com/nyaruka/goflow/antlr/gen/excellent3"
	"github.com/nyaruka/goflow/excellent/types"
)

// visitor which evaluates each part of an expression as a value
type visitor struct {
	gen.BaseExcellent3Visitor

	// tracks where we are in the context
	currContext     []string
	contextCallback func([]string)
}

func (v *visitor) context(part string, reset bool) {
	part = strings.ToLower(part)
	if reset {
		v.currContext = []string{part}
	} else {
		v.currContext = append(v.currContext, part)
	}
	if v.contextCallback != nil {
		v.contextCallback(v.currContext)
	}
}

// Visit the top level parse tree
func (v *visitor) Visit(tree antlr.ParseTree) any {
	return tree.Accept(v)
}

// VisitParse handles our top level parser
func (v *visitor) VisitParse(ctx *gen.ParseContext) any {
	return v.Visit(ctx.Expression())
}

// VisitAtomReference deals with visiting a single atom in our expression
func (v *visitor) VisitAtomReference(ctx *gen.AtomReferenceContext) any {
	return v.Visit(ctx.Atom())
}

// VisitContextReference deals with identifiers which are function names or root variables in the context
func (v *visitor) VisitContextReference(ctx *gen.ContextReferenceContext) any {
	name := ctx.GetText()
	v.context(name, true)

	return &ContextReference{Name: name}
}

// VisitDotLookup deals with property lookups like foo.bar
func (v *visitor) VisitDotLookup(ctx *gen.DotLookupContext) any {
	container := toExpression(v.Visit(ctx.Atom()))
	var lookup string

	if ctx.NAME() != nil {
		lookup = ctx.NAME().GetText()
	} else {
		lookup = ctx.INTEGER().GetText()
	}

	v.context(lookup, false)

	return &DotLookup{Container: container, Lookup: lookup}
}

// VisitArrayLookup deals with lookups such as foo[5] or foo["key with spaces"]
func (v *visitor) VisitArrayLookup(ctx *gen.ArrayLookupContext) any {
	container := toExpression(v.Visit(ctx.Atom()))
	lookup := toExpression(v.Visit(ctx.Expression()))

	asText, isText := lookup.(*TextLiteral)
	if isText {
		v.context(asText.Value.Native(), false)
	}

	return &ArrayLookup{Container: container, Lookup: lookup}
}

// VisitFunctionCall deals with function calls like TITLE(foo.bar)
func (v *visitor) VisitFunctionCall(ctx *gen.FunctionCallContext) any {
	function := toExpression(v.Visit(ctx.Atom()))

	var params []Expression
	if ctx.Parameters() != nil {
		params, _ = v.Visit(ctx.Parameters()).([]Expression)
	}

	return &FunctionCall{Func: function, Params: params}
}

// VisitFunctionParameters deals with the parameters to a function call
func (v *visitor) VisitFunctionParameters(ctx *gen.FunctionParametersContext) any {
	expressions := ctx.AllExpression()
	params := make([]Expression, len(expressions))

	for i := range expressions {
		params[i] = toExpression(v.Visit(expressions[i]))
	}
	return params
}

// VisitAnonFunction deals with anonymous functions, e.g. (x) => 2 * x
func (v *visitor) VisitAnonFunction(ctx *gen.AnonFunctionContext) any {
	return &AnonFunction{
		Args: v.Visit(ctx.NameList()).([]string),
		Body: toExpression(v.Visit(ctx.Expression())),
	}
}

func (v *visitor) VisitNameList(ctx *gen.NameListContext) any {
	names := ctx.AllNAME()
	args := make([]string, len(names))

	for i := range names {
		args[i] = names[i].GetText()
	}
	return args
}

// VisitConcatenation deals with string concatenations like "foo" & "bar"
func (v *visitor) VisitConcatenation(ctx *gen.ConcatenationContext) any {
	return &Concatenation{
		Exp1: toExpression(v.Visit(ctx.Expression(0))),
		Exp2: toExpression(v.Visit(ctx.Expression(1))),
	}
}

// VisitAdditionOrSubtraction deals with addition and subtraction like 5+5 and 5-3
func (v *visitor) VisitAdditionOrSubtraction(ctx *gen.AdditionOrSubtractionContext) any {
	exp1 := toExpression(v.Visit(ctx.Expression(0)))
	exp2 := toExpression(v.Visit(ctx.Expression(1)))

	if ctx.PLUS() != nil {
		return &Addition{Exp1: exp1, Exp2: exp2}
	}
	return &Subtraction{Exp1: exp1, Exp2: exp2}
}

// VisitMultiplicationOrDivision deals with division and multiplication such as 5*5 or 5/2
func (v *visitor) VisitMultiplicationOrDivision(ctx *gen.MultiplicationOrDivisionContext) any {
	exp1 := toExpression(v.Visit(ctx.Expression(0)))
	exp2 := toExpression(v.Visit(ctx.Expression(1)))

	if ctx.TIMES() != nil {
		return &Multiplication{Exp1: exp1, Exp2: exp2}
	}
	return &Division{Exp1: exp1, Exp2: exp2}
}

// VisitExponent deals with exponenets such as 5^5
func (v *visitor) VisitExponent(ctx *gen.ExponentContext) any {
	return &Exponent{
		Expression: toExpression(v.Visit(ctx.Expression(0))),
		Exponent:   toExpression(v.Visit(ctx.Expression(1))),
	}
}

// VisitNegation deals with negations such as -5
func (v *visitor) VisitNegation(ctx *gen.NegationContext) any {
	return &Negation{Exp: toExpression(v.Visit(ctx.Expression()))}
}

// VisitEquality deals with equality or inequality tests 5 = 5 and 5 != 5
func (v *visitor) VisitEquality(ctx *gen.EqualityContext) any {
	exp1 := toExpression(v.Visit(ctx.Expression(0)))
	exp2 := toExpression(v.Visit(ctx.Expression(1)))

	if ctx.EQ() != nil {
		return &Equality{Exp1: exp1, Exp2: exp2}
	}
	return &InEquality{Exp1: exp1, Exp2: exp2}
}

// VisitComparison deals with visiting a comparison between two values, such as 5<3 or 3>5
func (v *visitor) VisitComparison(ctx *gen.ComparisonContext) any {
	exp1 := toExpression(v.Visit(ctx.Expression(0)))
	exp2 := toExpression(v.Visit(ctx.Expression(1)))

	switch {
	case ctx.LT() != nil:
		return &LessThan{Exp1: exp1, Exp2: exp2}
	case ctx.LTE() != nil:
		return &LessThanOrEqual{Exp1: exp1, Exp2: exp2}
	case ctx.GTE() != nil:
		return &GreaterThanOrEqual{Exp1: exp1, Exp2: exp2}
	default:
		return &GreaterThan{Exp1: exp1, Exp2: exp2}
	}
}

// VisitParentheses deals with expressions in parentheses such as (1+2)
func (v *visitor) VisitParentheses(ctx *gen.ParenthesesContext) any {
	return &Parentheses{Exp: toExpression(v.Visit(ctx.Expression()))}
}

// VisitTextLiteral deals with string literals such as "asdf"
func (v *visitor) VisitTextLiteral(ctx *gen.TextLiteralContext) any {
	value := ctx.GetText()

	// unquote, this takes care of escape sequences as well
	unquoted, err := strconv.Unquote(value)

	// if we had an error, just strip surrounding quotes. It's fairly common for text literals
	// to contain escape sequences which aren't legal in go, e.g. a regex \w+
	if err != nil {
		unquoted = value[1 : len(value)-1]
	}

	return &TextLiteral{Value: types.NewXText(unquoted)}
}

// VisitNumberLiteral deals with numbers like 123 or 1.5
func (v *visitor) VisitNumberLiteral(ctx *gen.NumberLiteralContext) any {
	num, err := types.NewXNumberFromString(ctx.GetText())
	if err != nil {
		return &ErrorLiteral{Err: types.NewXErrorf("number %s is out of range", ctx.GetText())}
	}
	return &NumberLiteral{Value: num}
}

// VisitTrue deals with the `true` reserved word
func (v *visitor) VisitTrue(ctx *gen.TrueContext) any {
	return &BooleanLiteral{Value: types.XBooleanTrue}
}

// VisitFalse deals with the `false` reserved word
func (v *visitor) VisitFalse(ctx *gen.FalseContext) any {
	return &BooleanLiteral{Value: types.XBooleanFalse}
}

// VisitNull deals with the `null` reserved word
func (v *visitor) VisitNull(ctx *gen.NullContext) any {
	return &NullLiteral{}
}

// convenience utility to convert the given value to an Expression
func toExpression(val any) Expression {
	asExp, isExp := val.(Expression)
	if !isExp && val != nil {
		panic(fmt.Sprintf("attempt to convert a %T to an Expression", val))
	}
	return asExp
}
