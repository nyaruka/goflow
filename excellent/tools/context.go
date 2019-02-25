package tools

import (
	"fmt"
	"strconv"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/gen"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// AuditContextInTemplate audits context references in the given template
func AuditContextInTemplate(template string, allowedTopLevels []string, callback func(string)) error {
	return excellent.VisitTemplate(template, allowedTopLevels, func(tokenType excellent.XTokenType, token string) error {
		switch tokenType {
		case excellent.IDENTIFIER, excellent.EXPRESSION:
			return auditContextInExpression(token, callback)
		}
		return nil
	})
}

func auditContextInExpression(expression string, callback func(string)) error {
	visitor := &auditContextVisitor{callback: callback}

	_, err := excellent.VisitExpression(expression, visitor)

	return err
}

// visitor which audits access to the context
type auditContextVisitor struct {
	excellent.BaseVisitor

	callback func(string)
}

// Visit the top level parse tree
func (v *auditContextVisitor) Visit(tree antlr.ParseTree) interface{} {
	return tree.Accept(v)
}

func (v *auditContextVisitor) VisitChildren(node antlr.RuleNode) interface{} {
	for _, c := range node.GetChildren() {
		c.(antlr.ParseTree).Accept(v)
	}
	return nil
}

// VisitParse handles our top level parser
func (v *auditContextVisitor) VisitParse(ctx *gen.ParseContext) interface{} {
	return v.Visit(ctx.Expression())
}

// VisitDotLookup deals with lookups like foo.0 or foo.bar
func (v *auditContextVisitor) VisitDotLookup(ctx *gen.DotLookupContext) interface{} {
	base := v.Visit(ctx.Atom(0)).(string)
	key := ctx.Atom(1).GetText()
	path := fmt.Sprintf("%s.%s", base, key)
	v.callback(path)
	return path
}

// VisitArrayLookup deals with lookups such as foo[5] or foo["key with spaces"]
func (v *auditContextVisitor) VisitArrayLookup(ctx *gen.ArrayLookupContext) interface{} {
	base := v.Visit(ctx.Atom()).(string)
	key := v.Visit(ctx.Expression())

	if key != nil {
		path := fmt.Sprintf("%s.%s", base, key)
		v.callback(path)
		return path
	}

	return nil
}

// VisitContextReference deals with references to variables in the context such as "foo"
func (v *auditContextVisitor) VisitContextReference(ctx *gen.ContextReferenceContext) interface{} {
	path := ctx.GetText()
	v.callback(path)
	return path
}

// VisitTextLiteral deals with string literals such as "asdf"
func (v *auditContextVisitor) VisitTextLiteral(ctx *gen.TextLiteralContext) interface{} {
	// unquote, this takes care of escape sequences as well
	unquoted, _ := strconv.Unquote(ctx.GetText())
	return unquoted
}

// VisitAtomReference deals with visiting a single atom in our expression
func (v *auditContextVisitor) VisitAtomReference(ctx *gen.AtomReferenceContext) interface{} {
	return v.Visit(ctx.Atom())
}

// VisitFunctionCall deals with function calls like TITLE(foo.bar)
func (v *auditContextVisitor) VisitFunctionCall(ctx *gen.FunctionCallContext) interface{} {
	return v.VisitChildren(ctx)
}

// VisitFunctionParameters deals with the parameters to a function call
func (v *auditContextVisitor) VisitFunctionParameters(ctx *gen.FunctionParametersContext) interface{} {
	return v.VisitChildren(ctx)
}

// VisitParentheses deals with expressions in parentheses such as (1+2)
func (v *auditContextVisitor) VisitParentheses(ctx *gen.ParenthesesContext) interface{} {
	return v.VisitChildren(ctx)
}

// VisitAdditionOrSubtraction deals with addition and subtraction like 5+5 and 5-3
func (v *auditContextVisitor) VisitAdditionOrSubtraction(ctx *gen.AdditionOrSubtractionContext) interface{} {
	return v.VisitChildren(ctx)
}

// VisitMultiplicationOrDivision deals with division and multiplication such as 5*5 or 5/2
func (v *auditContextVisitor) VisitMultiplicationOrDivision(ctx *gen.MultiplicationOrDivisionContext) interface{} {
	return v.VisitChildren(ctx)
}

// VisitNegation deals with negations such as -5
func (v *auditContextVisitor) VisitNegation(ctx *gen.NegationContext) interface{} {
	return v.VisitChildren(ctx)
}

// VisitExponent deals with exponenets such as 5^5
func (v *auditContextVisitor) VisitExponent(ctx *gen.ExponentContext) interface{} {
	return v.VisitChildren(ctx)
}

// VisitConcatenation deals with string concatenations like "foo" & "bar"
func (v *auditContextVisitor) VisitConcatenation(ctx *gen.ConcatenationContext) interface{} {
	return v.VisitChildren(ctx)
}

// VisitEquality deals with equality or inequality tests 5 = 5 and 5 != 5
func (v *auditContextVisitor) VisitEquality(ctx *gen.EqualityContext) interface{} {
	return v.VisitChildren(ctx)
}

// VisitComparison deals with visiting a comparison between two values, such as 5<3 or 3>5
func (v *auditContextVisitor) VisitComparison(ctx *gen.ComparisonContext) interface{} {
	return v.VisitChildren(ctx)
}
