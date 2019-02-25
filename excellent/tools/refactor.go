package tools

import (
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/excellent/gen"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// RefactorVisitor which rewrites each part of an expression
type RefactorVisitor struct {
	gen.BaseExcellent2Visitor
}

// NewRefactorVisitor creates a new refactor visitor
func NewRefactorVisitor() *RefactorVisitor {
	return &RefactorVisitor{}
}

// Visit the top level parse tree
func (v *RefactorVisitor) Visit(tree antlr.ParseTree) interface{} {
	return tree.Accept(v)
}

// VisitParse handles our top level parser
func (v *RefactorVisitor) VisitParse(ctx *gen.ParseContext) interface{} {
	return v.Visit(ctx.Expression())
}

// VisitTextLiteral deals with string literals such as "asdf"
func (v *RefactorVisitor) VisitTextLiteral(ctx *gen.TextLiteralContext) interface{} {
	return ctx.GetText()
}

// VisitNumberLiteral deals with numbers like 123 or 1.5
func (v *RefactorVisitor) VisitNumberLiteral(ctx *gen.NumberLiteralContext) interface{} {
	return ctx.GetText()
}

// VisitDotLookup deals with lookups like foo.0 or foo.bar
func (v *RefactorVisitor) VisitDotLookup(ctx *gen.DotLookupContext) interface{} {
	return fmt.Sprintf("%s.%s", v.Visit(ctx.Atom(0)), v.Visit(ctx.Atom(1)))
}

// VisitFunctionCall deals with function calls like TITLE(foo.bar)
func (v *RefactorVisitor) VisitFunctionCall(ctx *gen.FunctionCallContext) interface{} {
	functionName := strings.ToLower(ctx.Fnname().GetText())

	var params []string
	if ctx.Parameters() != nil {
		params, _ = v.Visit(ctx.Parameters()).([]string)
	}

	return fmt.Sprintf("%s(%s)", functionName, strings.Join(params, ", "))
}

// VisitTrue deals with the `true` reserved word
func (v *RefactorVisitor) VisitTrue(ctx *gen.TrueContext) interface{} {
	return "true"
}

// VisitFalse deals with the `false` reserved word
func (v *RefactorVisitor) VisitFalse(ctx *gen.FalseContext) interface{} {
	return "false"
}

// VisitNull deals with the `null` reserved word
func (v *RefactorVisitor) VisitNull(ctx *gen.NullContext) interface{} {
	return "null"
}

// VisitArrayLookup deals with lookups such as foo[5] or foo["key with spaces"]
func (v *RefactorVisitor) VisitArrayLookup(ctx *gen.ArrayLookupContext) interface{} {
	return fmt.Sprintf("%s[%s]", v.Visit(ctx.Atom()), v.Visit(ctx.Expression()))
}

// VisitContextReference deals with references to variables in the context such as "foo"
func (v *RefactorVisitor) VisitContextReference(ctx *gen.ContextReferenceContext) interface{} {
	return ctx.GetText()
}

// VisitParentheses deals with expressions in parentheses such as (1+2)
func (v *RefactorVisitor) VisitParentheses(ctx *gen.ParenthesesContext) interface{} {
	return fmt.Sprintf("(%s)", v.Visit(ctx.Expression()))
}

// VisitNegation deals with negations such as -5
func (v *RefactorVisitor) VisitNegation(ctx *gen.NegationContext) interface{} {
	return fmt.Sprintf("-%s", v.Visit(ctx.Expression()))
}

// VisitExponent deals with exponenets such as 5^5
func (v *RefactorVisitor) VisitExponent(ctx *gen.ExponentContext) interface{} {
	return fmt.Sprintf("%s ^ %s", v.Visit(ctx.Expression(0)), v.Visit(ctx.Expression(1)))
}

// VisitConcatenation deals with string concatenations like "foo" & "bar"
func (v *RefactorVisitor) VisitConcatenation(ctx *gen.ConcatenationContext) interface{} {
	return fmt.Sprintf("%s & %s", v.Visit(ctx.Expression(0)), v.Visit(ctx.Expression(1)))
}

// VisitAdditionOrSubtraction deals with addition and subtraction like 5+5 and 5-3
func (v *RefactorVisitor) VisitAdditionOrSubtraction(ctx *gen.AdditionOrSubtractionContext) interface{} {
	return fmt.Sprintf("%s %s %s", v.Visit(ctx.Expression(0)), ctx.GetOp().GetText(), v.Visit(ctx.Expression(1)))
}

// VisitEquality deals with equality or inequality tests 5 = 5 and 5 != 5
func (v *RefactorVisitor) VisitEquality(ctx *gen.EqualityContext) interface{} {
	return fmt.Sprintf("%s %s %s", v.Visit(ctx.Expression(0)), ctx.GetOp().GetText(), v.Visit(ctx.Expression(1)))
}

// VisitAtomReference deals with visiting a single atom in our expression
func (v *RefactorVisitor) VisitAtomReference(ctx *gen.AtomReferenceContext) interface{} {
	return v.Visit(ctx.Atom())
}

// VisitMultiplicationOrDivision deals with division and multiplication such as 5*5 or 5/2
func (v *RefactorVisitor) VisitMultiplicationOrDivision(ctx *gen.MultiplicationOrDivisionContext) interface{} {
	return fmt.Sprintf("%s %s %s", v.Visit(ctx.Expression(0)), ctx.GetOp().GetText(), v.Visit(ctx.Expression(1)))
}

// VisitComparison deals with visiting a comparison between two values, such as 5<3 or 3>5
func (v *RefactorVisitor) VisitComparison(ctx *gen.ComparisonContext) interface{} {
	return fmt.Sprintf("%s %s %s", v.Visit(ctx.Expression(0)), ctx.GetOp().GetText(), v.Visit(ctx.Expression(1)))
}

// VisitFunctionParameters deals with the parameters to a function call
func (v *RefactorVisitor) VisitFunctionParameters(ctx *gen.FunctionParametersContext) interface{} {
	params := make([]string, len(ctx.AllExpression()))
	for i, exp := range ctx.AllExpression() {
		params[i] = fmt.Sprintf("%s", v.Visit(exp))
	}

	// return as slice of strings so that function name and params can be refactored together
	// in VisitFunctionCall
	return params
}
