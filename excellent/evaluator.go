package excellent

import (
	"strconv"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/gen"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
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
func EvaluateExpression(env envs.Environment, ctx *types.XObject, expression string) types.XValue {
	parsed, err := Parse(expression)
	if err != nil {
		return types.NewXError(err)
	}
	return parsed.Evaluate(env, ctx)
}

// visitor which evaluates each part of an expression as a value
type visitor struct {
	gen.BaseExcellent2Visitor
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

	return &TextLiteral{val: types.NewXText(unquoted)}
}

// VisitNumberLiteral deals with numbers like 123 or 1.5
func (v *visitor) VisitNumberLiteral(ctx *gen.NumberLiteralContext) interface{} {
	return &NumberLiteral{val: types.RequireXNumberFromString(ctx.GetText())}
}

// VisitContextReference deals with identifiers which are function names or root variables in the context
func (v *visitor) VisitContextReference(ctx *gen.ContextReferenceContext) interface{} {
	return &ContextReference{name: ctx.GetText()}
}

// VisitDotLookup deals with property lookups like foo.bar
func (v *visitor) VisitDotLookup(ctx *gen.DotLookupContext) interface{} {
	container := toExpression(v.Visit(ctx.Atom()))
	var lookup string

	if ctx.NAME() != nil {
		lookup = ctx.NAME().GetText()
	} else {
		lookup = ctx.INTEGER().GetText()
	}

	return &DotLookup{container: container, lookup: lookup}
}

// VisitArrayLookup deals with lookups such as foo[5] or foo["key with spaces"]
func (v *visitor) VisitArrayLookup(ctx *gen.ArrayLookupContext) interface{} {
	container := toExpression(v.Visit(ctx.Atom()))
	lookup := toExpression(v.Visit(ctx.Expression()))

	return &ArrayLookup{container: container, lookup: lookup}
}

// VisitFunctionCall deals with function calls like TITLE(foo.bar)
func (v *visitor) VisitFunctionCall(ctx *gen.FunctionCallContext) interface{} {
	function := toExpression(v.Visit(ctx.Atom()))

	var params []Expression
	if ctx.Parameters() != nil {
		params, _ = v.Visit(ctx.Parameters()).([]Expression)
	}

	return &FunctionCall{function: function, params: params}
}

// VisitTrue deals with the `true` reserved word
func (v *visitor) VisitTrue(ctx *gen.TrueContext) interface{} {
	return &BooleanLiteral{val: types.XBooleanTrue}
}

// VisitFalse deals with the `false` reserved word
func (v *visitor) VisitFalse(ctx *gen.FalseContext) interface{} {
	return &BooleanLiteral{val: types.XBooleanFalse}
}

// VisitNull deals with the `null` reserved word
func (v *visitor) VisitNull(ctx *gen.NullContext) interface{} {
	return &NullLiteral{}
}

// VisitParentheses deals with expressions in parentheses such as (1+2)
func (v *visitor) VisitParentheses(ctx *gen.ParenthesesContext) interface{} {
	return &Parentheses{exp: toExpression(v.Visit(ctx.Expression()))}
}

// VisitNegation deals with negations such as -5
func (v *visitor) VisitNegation(ctx *gen.NegationContext) interface{} {
	return &Negation{exp: toExpression(v.Visit(ctx.Expression()))}
}

// VisitExponent deals with exponenets such as 5^5
func (v *visitor) VisitExponent(ctx *gen.ExponentContext) interface{} {
	return &Exponent{
		expression: toExpression(v.Visit(ctx.Expression(0))),
		exponent:   toExpression(v.Visit(ctx.Expression(1))),
	}
}

// VisitConcatenation deals with string concatenations like "foo" & "bar"
func (v *visitor) VisitConcatenation(ctx *gen.ConcatenationContext) interface{} {
	return &Concatenation{
		exp1: toExpression(v.Visit(ctx.Expression(0))),
		exp2: toExpression(v.Visit(ctx.Expression(1))),
	}
}

// VisitAdditionOrSubtraction deals with addition and subtraction like 5+5 and 5-3
func (v *visitor) VisitAdditionOrSubtraction(ctx *gen.AdditionOrSubtractionContext) interface{} {
	exp1 := toExpression(v.Visit(ctx.Expression(0)))
	exp2 := toExpression(v.Visit(ctx.Expression(1)))

	if ctx.PLUS() != nil {
		return &Addition{exp1: exp1, exp2: exp2}
	}
	return &Subtraction{exp1: exp1, exp2: exp2}
}

// VisitMultiplicationOrDivision deals with division and multiplication such as 5*5 or 5/2
func (v *visitor) VisitMultiplicationOrDivision(ctx *gen.MultiplicationOrDivisionContext) interface{} {
	exp1 := toExpression(v.Visit(ctx.Expression(0)))
	exp2 := toExpression(v.Visit(ctx.Expression(1)))

	if ctx.TIMES() != nil {
		return &Multiplication{exp1: exp1, exp2: exp2}
	}
	return &Division{exp1: exp1, exp2: exp2}
}

// VisitEquality deals with equality or inequality tests 5 = 5 and 5 != 5
func (v *visitor) VisitEquality(ctx *gen.EqualityContext) interface{} {
	exp1 := toExpression(v.Visit(ctx.Expression(0)))
	exp2 := toExpression(v.Visit(ctx.Expression(1)))

	if ctx.EQ() != nil {
		return &Equality{exp1: exp1, exp2: exp2}
	}
	return &InEquality{exp1: exp1, exp2: exp2}
}

// VisitComparison deals with visiting a comparison between two values, such as 5<3 or 3>5
func (v *visitor) VisitComparison(ctx *gen.ComparisonContext) interface{} {
	exp1 := toExpression(v.Visit(ctx.Expression(0)))
	exp2 := toExpression(v.Visit(ctx.Expression(1)))

	switch {
	case ctx.LT() != nil:
		return &LessComparison{exp1: exp1, exp2: exp2}
	case ctx.LTE() != nil:
		return &LessOrEqualComparison{exp1: exp1, exp2: exp2}
	case ctx.GTE() != nil:
		return &GreaterOrEqualComparison{exp1: exp1, exp2: exp2}
	default:
		return &GreaterComparison{exp1: exp1, exp2: exp2}
	}
}

// VisitAtomReference deals with visiting a single atom in our expression
func (v *visitor) VisitAtomReference(ctx *gen.AtomReferenceContext) interface{} {
	return v.Visit(ctx.Atom())
}

// VisitFunctionParameters deals with the parameters to a function call
func (v *visitor) VisitFunctionParameters(ctx *gen.FunctionParametersContext) interface{} {
	expressions := ctx.AllExpression()
	params := make([]Expression, len(expressions))

	for i := range expressions {
		params[i] = toExpression(v.Visit(expressions[i]))
	}
	return params
}

// convenience utility to convert the given value to an Expression
func toExpression(val interface{}) Expression {
	asExp, isExp := val.(Expression)
	if !isExp && val != nil {
		panic("attempt to convert a non-expression to an Expression")
	}
	return asExp
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
