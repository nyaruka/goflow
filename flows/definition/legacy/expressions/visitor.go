package expressions

import (
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows/definition/legacy/gen"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type legacyVisitor struct {
	gen.BaseExcellent1Visitor
	env     envs.Environment
	options *MigrateOptions
}

func newLegacyVisitor(env envs.Environment, options *MigrateOptions) *legacyVisitor {
	return &legacyVisitor{env: env, options: options}
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
	return ctx.GetText()
}

// VisitStringLiteral deals with string literals such as "asdf"
func (v *legacyVisitor) VisitStringLiteral(ctx *gen.StringLiteralContext) interface{} {
	return MigrateStringLiteral(ctx.GetText())
}

// VisitFunctionCall deals with function calls like TITLE(foo.bar)
func (v *legacyVisitor) VisitFunctionCall(ctx *gen.FunctionCallContext) interface{} {
	functionName := strings.ToLower(ctx.Fnname().GetText())

	var params []string
	if ctx.Parameters() != nil {
		params = v.Visit(ctx.Parameters()).([]string)
	}

	rewrittenFuncCall, _ := migrateFunctionCall(functionName, params)
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
	return MigrateContextReference(ctx.GetText(), v.options.RawDates)
}

// VisitParentheses deals with expressions in parentheses such as (1+2)
func (v *legacyVisitor) VisitParentheses(ctx *gen.ParenthesesContext) interface{} {
	return fmt.Sprintf("(%s)", v.Visit(ctx.Expression()))
}

// VisitNegation deals with negations such as -5
func (v *legacyVisitor) VisitNegation(ctx *gen.NegationContext) interface{} {
	return fmt.Sprintf("-%s", v.Visit(ctx.Expression()))
}

// VisitExponentExpression deals with exponenets such as 5^5
func (v *legacyVisitor) VisitExponentExpression(ctx *gen.ExponentExpressionContext) interface{} {
	arg1 := v.Visit(ctx.Expression(0))
	arg2 := v.Visit(ctx.Expression(1))

	return fmt.Sprintf("%s ^ %s", arg1, arg2)
}

// VisitConcatenation deals with string concatenations like "foo" & "bar"
func (v *legacyVisitor) VisitConcatenation(ctx *gen.ConcatenationContext) interface{} {
	arg1 := v.Visit(ctx.Expression(0))
	arg2 := v.Visit(ctx.Expression(1))

	return fmt.Sprintf("%s & %s", arg1, arg2)
}

// VisitAdditionOrSubtractionExpression deals with addition and subtraction like 5+5 and 5-3
func (v *legacyVisitor) VisitAdditionOrSubtractionExpression(ctx *gen.AdditionOrSubtractionExpressionContext) interface{} {
	arg1 := v.Visit(ctx.Expression(0)).(string)
	arg2 := v.Visit(ctx.Expression(1)).(string)

	op := "+"
	if ctx.MINUS() != nil {
		op = "-"
	}

	// see if either of our arguments is a date value
	arg1Type := inferType(arg1)
	arg2Type := inferType(arg2)

	//fmt.Printf("Migrating add/sub with types: %s => %s, %s =>%s\n", arg1, arg1Type, arg2, arg2Type)

	if arg1Type == "number" && arg2Type == "number" {
		// we are adding two numbers
		return fmt.Sprintf("%s %s %s", arg1, op, arg2)

	} else if arg1Type == "datetime" && arg2Type == "number" {
		// we are adding a datetime and a number (of days)
		template := `datetime_add(%s, %s, "D")`
		if op == "-" {
			template = `datetime_add(%s, -%s, "D")`
		}

		return fmt.Sprintf(template, arg1, arg2)

	} else if arg1Type == "date" && arg2Type == "number" {
		// we are adding a date and a number (of days)
		template := `datetime_add(%s, %s, "D")`
		if op == "-" {
			template = `datetime_add(%s, -%s, "D")`
		}

		if !v.options.RawDates {
			template = wrap(template, "format_date")
		}

		return fmt.Sprintf(template, arg1, arg2)

	} else if arg1Type == "datetime" && arg2Type == "time" {
		// we are adding a datetime and a time

		// create expression which converts arg2 to minutes
		asMinutes := fmt.Sprintf(`format_time(%[1]s, "tt") * 60 + format_time(%[1]s, "m")`, arg2)

		template := `datetime_add(%s, %s, "m")`
		if op == "-" {
			template = `datetime_add(%s, -(%s), "m")`
		}
		return fmt.Sprintf(template, arg1, asMinutes)

	} else if arg2Type == "time" && op == "+" {
		// we are adding a date and a time
		return fmt.Sprintf(`replace_time(%s, %s)`, arg1, arg2)
	}

	// we don't know what we are adding so fallback to legacy_add
	if op == "+" {
		return fmt.Sprintf("legacy_add(%s, %s)", arg1, arg2)
	}
	return fmt.Sprintf("legacy_add(%s, -%s)", arg1, arg2)
}

// VisitEquality deals with equality or inequality tests 5 = 5 and 5 != 5
func (v *legacyVisitor) VisitEqualityExpression(ctx *gen.EqualityExpressionContext) interface{} {
	arg1 := v.Visit(ctx.Expression(0))
	arg2 := v.Visit(ctx.Expression(1))

	if ctx.EQ() != nil {
		return fmt.Sprintf("%s = %s", arg1, arg2)
	}

	return fmt.Sprintf("%s != %s", arg1, arg2)
}

// VisitMultiplicationOrDivision deals with division and multiplication such as 5*5 or 5/2
func (v *legacyVisitor) VisitMultiplicationOrDivisionExpression(ctx *gen.MultiplicationOrDivisionExpressionContext) interface{} {
	arg1 := v.Visit(ctx.Expression(0))
	arg2 := v.Visit(ctx.Expression(1))

	if ctx.TIMES() != nil {
		return fmt.Sprintf("%s * %s", arg1, arg2)
	}

	return fmt.Sprintf("%s / %s", arg1, arg2)
}

// VisitComparison deals with visiting a comparison between two values, such as 5<3 or 3>5
func (v *legacyVisitor) VisitComparisonExpression(ctx *gen.ComparisonExpressionContext) interface{} {
	arg1 := v.Visit(ctx.Expression(0))
	arg2 := v.Visit(ctx.Expression(1))

	return fmt.Sprintf("%s %s %s", arg1, ctx.GetOp().GetText(), arg2)
}

// VisitFunctionParameters deals with the parameters to a function call
func (v *legacyVisitor) VisitFunctionParameters(ctx *gen.FunctionParametersContext) interface{} {
	expressions := ctx.AllExpression()
	params := make([]string, len(expressions))

	for i := range expressions {
		params[i] = v.Visit(expressions[i]).(string)
	}
	return params
}
