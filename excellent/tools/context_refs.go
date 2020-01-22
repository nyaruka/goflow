package tools

import (
	"strconv"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/functions"
	"github.com/nyaruka/goflow/excellent/gen"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// FindContextRefsInTemplate audits context references in the given template. Note that the case of
// the found references is preserved as these may be significant, e.g. ["X"] vs ["x"] in JSON
func FindContextRefsInTemplate(template string, allowedTopLevels []string, callback func([]string)) error {
	return excellent.VisitTemplate(template, allowedTopLevels, func(tokenType excellent.XTokenType, token string) error {
		switch tokenType {
		case excellent.IDENTIFIER, excellent.EXPRESSION:
			return findContextRefsInTemplate(token, callback)
		}
		return nil
	})
}

func findContextRefsInTemplate(expression string, callback func([]string)) error {
	visitor := &auditContextVisitor{callback: callback}

	_, err := excellent.VisitExpression(expression, visitor)

	return err
}

// visitor which audits access to the context
type auditContextVisitor struct {
	gen.BaseExcellent2Visitor

	callback func([]string)
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

// VisitContextReference deals with root variables in the context
func (v *auditContextVisitor) VisitContextReference(ctx *gen.ContextReferenceContext) interface{} {
	name := ctx.NAME().GetText()

	function := functions.Lookup(name)
	if function == nil {
		path := []string{name}
		v.callback(path)
		return path
	}

	return nil
}

// VisitDotLookup deals with lookups like foo.bar
func (v *auditContextVisitor) VisitDotLookup(ctx *gen.DotLookupContext) interface{} {
	path, isPath := v.Visit(ctx.Atom()).([]string)

	var lookup string

	if ctx.NAME() != nil {
		lookup = ctx.NAME().GetText()
	} else {
		lookup = ctx.INTEGER().GetText()
	}

	if isPath {
		path = append(path, lookup)
		v.callback(path)
		return path
	}
	return nil
}

// VisitArrayLookup deals with lookups such as foo[5] or foo["key with spaces"]
func (v *auditContextVisitor) VisitArrayLookup(ctx *gen.ArrayLookupContext) interface{} {
	path, isPath := v.Visit(ctx.Atom()).([]string)
	key, isString := v.Visit(ctx.Expression()).(string)
	if isPath && isString {
		path = append(path, key)
		v.callback(path)
		return path
	}
	return nil
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
