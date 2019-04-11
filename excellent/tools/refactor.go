package tools

import (
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/gen"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// RefactorTemplate refactors the passed in template
func RefactorTemplate(template string, allowedTopLevels []string) (string, error) {
	buf := &strings.Builder{}

	err := excellent.VisitTemplate(template, allowedTopLevels, func(tokenType excellent.XTokenType, token string) error {
		switch tokenType {
		case excellent.BODY:
			buf.WriteString(token)
		case excellent.IDENTIFIER, excellent.EXPRESSION:
			refactored, err := refactorExpression(token)

			// if we got an error, return that, and rewrite original expression
			if err != nil {
				buf.WriteString(wrapExpression(tokenType, token))
				return err
			}

			// if not, append refactored expresion to the output
			buf.WriteString(wrapExpression(tokenType, refactored))
		}
		return nil
	})

	return buf.String(), err
}

// RefactorTemplate refactors the passed in template
func refactorExpression(expression string) (string, error) {
	visitor := &refactorVisitor{}
	output, err := excellent.VisitExpression(expression, visitor)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", output), nil
}

func wrapExpression(tokenType excellent.XTokenType, token string) string {
	if tokenType == excellent.IDENTIFIER {
		return "@" + token
	}
	return "@(" + token + ")"
}

// visitor which rewrites each part of an expression
type refactorVisitor struct {
	gen.BaseExcellent2Visitor
}

// Visit the top level parse tree
func (v *refactorVisitor) Visit(tree antlr.ParseTree) interface{} {
	return tree.Accept(v)
}

// VisitParse handles our top level parser
func (v *refactorVisitor) VisitParse(ctx *gen.ParseContext) interface{} {
	return v.Visit(ctx.Expression())
}

// VisitTextLiteral deals with string literals such as "asdf"
func (v *refactorVisitor) VisitTextLiteral(ctx *gen.TextLiteralContext) interface{} {
	return ctx.GetText()
}

// VisitNumberLiteral deals with numbers like 123 or 1.5
func (v *refactorVisitor) VisitNumberLiteral(ctx *gen.NumberLiteralContext) interface{} {
	return ctx.GetText()
}

// VisitDotLookup deals with lookups like foo.0 or foo.bar
func (v *refactorVisitor) VisitContextReference(ctx *gen.ContextReferenceContext) interface{} {
	return strings.ToLower(ctx.NAME().GetText())
}

// VisitDotLookup deals with lookups like foo.bar
func (v *refactorVisitor) VisitDotLookup(ctx *gen.DotLookupContext) interface{} {
	property := ctx.NAME().GetText()

	return fmt.Sprintf("%s.%s", v.Visit(ctx.Atom()), strings.ToLower(property))
}

// VisitFunctionCall deals with function calls like TITLE(foo.bar)
func (v *refactorVisitor) VisitFunctionCall(ctx *gen.FunctionCallContext) interface{} {
	functionName := v.Visit(ctx.Atom())

	var params []string
	if ctx.Parameters() != nil {
		params, _ = v.Visit(ctx.Parameters()).([]string)
	}

	return fmt.Sprintf("%s(%s)", functionName, strings.Join(params, ", "))
}

// VisitFunctionParameters deals with the parameters to a function call
func (v *refactorVisitor) VisitFunctionParameters(ctx *gen.FunctionParametersContext) interface{} {
	params := make([]string, len(ctx.AllExpression()))
	for i, exp := range ctx.AllExpression() {
		params[i] = fmt.Sprintf("%s", v.Visit(exp))
	}

	// return as slice of strings so that function name and params can be refactored together
	// in VisitFunctionCall
	return params
}

// VisitTrue deals with the `true` reserved word
func (v *refactorVisitor) VisitTrue(ctx *gen.TrueContext) interface{} {
	return "true"
}

// VisitFalse deals with the `false` reserved word
func (v *refactorVisitor) VisitFalse(ctx *gen.FalseContext) interface{} {
	return "false"
}

// VisitNull deals with the `null` reserved word
func (v *refactorVisitor) VisitNull(ctx *gen.NullContext) interface{} {
	return "null"
}

// VisitArrayLookup deals with lookups such as foo[5] or foo["key with spaces"]
func (v *refactorVisitor) VisitArrayLookup(ctx *gen.ArrayLookupContext) interface{} {
	return fmt.Sprintf("%s[%s]", v.Visit(ctx.Atom()), v.Visit(ctx.Expression()))
}

// VisitAtomReference deals with visiting a single atom in our expression
func (v *refactorVisitor) VisitAtomReference(ctx *gen.AtomReferenceContext) interface{} {
	return v.Visit(ctx.Atom())
}

// VisitParentheses deals with expressions in parentheses such as (1+2)
func (v *refactorVisitor) VisitParentheses(ctx *gen.ParenthesesContext) interface{} {
	return fmt.Sprintf("(%s)", v.Visit(ctx.Expression()))
}

// VisitAdditionOrSubtraction deals with addition and subtraction like 5+5 and 5-3
func (v *refactorVisitor) VisitAdditionOrSubtraction(ctx *gen.AdditionOrSubtractionContext) interface{} {
	return fmt.Sprintf("%s %s %s", v.Visit(ctx.Expression(0)), ctx.GetOp().GetText(), v.Visit(ctx.Expression(1)))
}

// VisitMultiplicationOrDivision deals with division and multiplication such as 5*5 or 5/2
func (v *refactorVisitor) VisitMultiplicationOrDivision(ctx *gen.MultiplicationOrDivisionContext) interface{} {
	return fmt.Sprintf("%s %s %s", v.Visit(ctx.Expression(0)), ctx.GetOp().GetText(), v.Visit(ctx.Expression(1)))
}

// VisitNegation deals with negations such as -5
func (v *refactorVisitor) VisitNegation(ctx *gen.NegationContext) interface{} {
	return fmt.Sprintf("-%s", v.Visit(ctx.Expression()))
}

// VisitExponent deals with exponenets such as 5^5
func (v *refactorVisitor) VisitExponent(ctx *gen.ExponentContext) interface{} {
	return fmt.Sprintf("%s ^ %s", v.Visit(ctx.Expression(0)), v.Visit(ctx.Expression(1)))
}

// VisitConcatenation deals with string concatenations like "foo" & "bar"
func (v *refactorVisitor) VisitConcatenation(ctx *gen.ConcatenationContext) interface{} {
	return fmt.Sprintf("%s & %s", v.Visit(ctx.Expression(0)), v.Visit(ctx.Expression(1)))
}

// VisitEquality deals with equality or inequality tests 5 = 5 and 5 != 5
func (v *refactorVisitor) VisitEquality(ctx *gen.EqualityContext) interface{} {
	return fmt.Sprintf("%s %s %s", v.Visit(ctx.Expression(0)), ctx.GetOp().GetText(), v.Visit(ctx.Expression(1)))
}

// VisitComparison deals with visiting a comparison between two values, such as 5<3 or 3>5
func (v *refactorVisitor) VisitComparison(ctx *gen.ComparisonContext) interface{} {
	return fmt.Sprintf("%s %s %s", v.Visit(ctx.Expression(0)), ctx.GetOp().GetText(), v.Visit(ctx.Expression(1)))
}
