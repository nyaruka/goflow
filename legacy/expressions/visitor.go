package expressions

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/nyaruka/goflow/legacy/gen"
	"github.com/nyaruka/goflow/utils"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type legacyVisitor struct {
	gen.BaseExcellent1Visitor
	env      utils.Environment
	resolver interface{}
}

func newLegacyVisitor(env utils.Environment, resolver interface{}) *legacyVisitor {
	return &legacyVisitor{env: env, resolver: resolver}
}

// ---------------------------------------------------------------

// Visit the top level parse tree
func (v *legacyVisitor) Visit(tree antlr.ParseTree) interface{} {
	return tree.Accept(v)
}

// VisitParse handles our top level parser
func (v *legacyVisitor) VisitParse(ctx *gen.ParseContext) interface{} {
	return v.Visit(ctx.Expression())
}

// VisitDecimalLiteral deals with decimals like 1.5
func (v *legacyVisitor) VisitDecimalLiteral(ctx *gen.DecimalLiteralContext) interface{} {
	decStr, _ := toString(ctx.GetText())
	return decStr
}

// VisitDotLookup deals with lookups like foo.0 or foo.bar
func (v *legacyVisitor) VisitDotLookup(ctx *gen.DotLookupContext) interface{} {
	value := v.Visit(ctx.Atom(0))
	expression := v.Visit(ctx.Atom(1))
	lookup, err := toString(expression)
	if err != nil {
		return err
	}
	return resolveLookup(v.env, value, lookup)
}

// VisitStringLiteral deals with string literals such as "asdf"
func (v *legacyVisitor) VisitStringLiteral(ctx *gen.StringLiteralContext) interface{} {
	return MigrateStringLiteral(ctx.GetText())
}

// VisitFunctionCall deals with function calls like TITLE(foo.bar)
func (v *legacyVisitor) VisitFunctionCall(ctx *gen.FunctionCallContext) interface{} {
	functionName := strings.ToLower(ctx.Fnname().GetText())

	var params []interface{}
	if ctx.Parameters() != nil {
		funcParams := v.Visit(ctx.Parameters())
		switch funcParams.(type) {
		case error:
			return funcParams
		default:
			params = funcParams.([]interface{})
		}
	}

	paramsAsStrs := make([]string, len(params))
	var err error
	for p := range params {
		paramsAsStrs[p], err = toString(params[p])
		if err != nil {
			return err
		}
	}

	rewrittenFuncCall, err := migrateFunctionCall(functionName, paramsAsStrs)
	if err != nil {
		return err
	}
	return rewrittenFuncCall
}

// VisitTrue deals with the "true" literal
func (v *legacyVisitor) VisitTrue(ctx *gen.TrueContext) interface{} {
	return "true"
}

// VisitFalse deals with the "false" literal
func (v *legacyVisitor) VisitFalse(ctx *gen.FalseContext) interface{} {
	return "false"
}

// VisitContextReference deals with references to variables in the context such as "foo"
func (v *legacyVisitor) VisitContextReference(ctx *gen.ContextReferenceContext) interface{} {
	key := strings.ToLower(ctx.GetText())
	val := resolveLookup(v.env, v.resolver, key)
	if val == nil {
		return fmt.Errorf("Invalid key: '%s'", key)
	}

	err, isErr := val.(error)
	if isErr {
		return err
	}

	return val
}

// VisitParentheses deals with expressions in parentheses such as (1+2)
func (v *legacyVisitor) VisitParentheses(ctx *gen.ParenthesesContext) interface{} {
	return fmt.Sprintf("(%s)", v.Visit(ctx.Expression()))
}

// VisitNegation deals with negations such as -5
func (v *legacyVisitor) VisitNegation(ctx *gen.NegationContext) interface{} {
	dec, err := toString(v.Visit(ctx.Expression()))
	if err != nil {
		return err
	}
	return "-" + dec
}

// VisitExponent deals with exponenets such as 5^5
func (v *legacyVisitor) VisitExponent(ctx *gen.ExponentContext) interface{} {
	arg1, err := toString(v.Visit(ctx.Expression(0)))
	if err != nil {
		return err
	}

	arg2, err := toString(v.Visit(ctx.Expression(1)))
	if err != nil {
		return err
	}

	return fmt.Sprintf("%s ^ %s", arg1, arg2)
}

// VisitConcatenation deals with string concatenations like "foo" & "bar"
func (v *legacyVisitor) VisitConcatenation(ctx *gen.ConcatenationContext) interface{} {
	arg1, err := toString(v.Visit(ctx.Expression(0)))
	if err != nil {
		return err
	}

	arg2, err := toString(v.Visit(ctx.Expression(1)))
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	buffer.WriteString(arg1)
	buffer.WriteString(" & ")
	buffer.WriteString(arg2)

	return buffer.String()
}

// VisitAdditionOrSubtraction deals with addition and subtraction like 5+5 and 5-3
func (v *legacyVisitor) VisitAdditionOrSubtraction(ctx *gen.AdditionOrSubtractionContext) interface{} {
	value, err := toString(v.Visit(ctx.Expression(0)))
	if err != nil {
		return err
	}

	dateUnit := "D"
	firstIsDate := isDate(value)
	if firstIsDate {
		firstSeconds, ok := convertTimeToSeconds(value)
		if ok {
			value = firstSeconds
			dateUnit = "s"
		}
	}

	// see if our first param is an int
	_, firstNumberErr := strconv.Atoi(value)

	next, err := toString(v.Visit(ctx.Expression(1)))
	if err != nil {
		return err
	}

	op := "+"
	if ctx.MINUS() != nil {
		op = "-"
	}

	secondIsDate := isDate(next)
	if secondIsDate {
		secondSeconds, ok := convertTimeToSeconds(next)
		if ok {
			next = secondSeconds
			dateUnit = "s"
		}
	}

	// see if our second param is an int
	_, secondNumberErr := strconv.Atoi(next)
	if (firstIsDate || secondIsDate) && (firstNumberErr != nil || secondNumberErr != nil) {

		// we are adding two values where we know at least one side is a date
		template := "datetime_add(%s, %s, \"%s\")"
		if op == "-" {
			template = "datetime_add(%s, -%s, \"%s\")"
		}

		// determine the order of our parameters
		replacements := []interface{}{value, next, dateUnit}
		if firstNumberErr == nil {
			replacements = []interface{}{next, value, dateUnit}
		}

		value = fmt.Sprintf(template, replacements...)

	} else if firstNumberErr == nil && secondNumberErr == nil {
		// we are adding two numbers
		if op == "+" {
			value = fmt.Sprintf("%s + %s", value, next)
		} else {
			value = fmt.Sprintf("%s - %s", value, next)
		}
	} else {
		// we are adding a field of unknown type with an integer
		if op == "+" {
			value = fmt.Sprintf("legacy_add(%s, %s)", value, next)
		} else {
			value = fmt.Sprintf("legacy_add(%s, -%s)", value, next)
		}
	}

	return value
}

// VisitEquality deals with equality or inequality tests 5 = 5 and 5 != 5
func (v *legacyVisitor) VisitEquality(ctx *gen.EqualityContext) interface{} {
	arg1 := v.Visit(ctx.Expression(0))
	err, isErr := arg1.(error)
	if isErr {
		return err
	}

	arg2 := v.Visit(ctx.Expression(1))
	err, isErr = arg2.(error)
	if isErr {
		return err
	}

	if ctx.EQ() != nil {
		return fmt.Sprintf("%s = %s", arg1, arg2)
	}

	return fmt.Sprintf("%s != %s", arg1, arg2)
}

// VisitAtomReference deals with visiting a single atom in our expression
func (v *legacyVisitor) VisitAtomReference(ctx *gen.AtomReferenceContext) interface{} {
	return v.Visit(ctx.Atom())
}

// VisitMultiplicationOrDivision deals with division and multiplication such as 5*5 or 5/2
func (v *legacyVisitor) VisitMultiplicationOrDivision(ctx *gen.MultiplicationOrDivisionContext) interface{} {
	arg1 := v.Visit(ctx.Expression(0))
	str1, err := toString(arg1)
	if err != nil {
		return err
	}

	arg2 := v.Visit(ctx.Expression(1))
	str2, err := toString(arg2)
	if err != nil {
		return err
	}

	if ctx.TIMES() != nil {
		return fmt.Sprintf("%s * %s", str1, str2)
	}

	return fmt.Sprintf("%s / %s", str1, str2)
}

// VisitComparison deals with visiting a comparison between two values, such as 5<3 or 3>5
func (v *legacyVisitor) VisitComparison(ctx *gen.ComparisonContext) interface{} {
	arg1 := v.Visit(ctx.Expression(0))
	arg2 := v.Visit(ctx.Expression(1))

	err, isErr := arg1.(error)
	if isErr {
		return err
	}

	err, isErr = arg2.(error)
	if isErr {
		return err
	}

	return fmt.Sprintf("%s %s %s", arg1, ctx.GetOp().GetText(), arg2)
}

// VisitFunctionParameters deals with the parameters to a function call
func (v *legacyVisitor) VisitFunctionParameters(ctx *gen.FunctionParametersContext) interface{} {
	expressions := ctx.AllExpression()
	params := make([]interface{}, len(expressions))

	for i := range expressions {
		params[i] = v.Visit(expressions[i])
		error, isError := params[i].(error)
		if isError {
			return error
		}
	}
	return params
}
