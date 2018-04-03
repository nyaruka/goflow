package contactql

import (
	"strings"

	"github.com/nyaruka/goflow/contactql/gen"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

var comparatorAliases = map[string]string{
	"has": "~",
	"is":  "=",
}

type Visitor struct {
	gen.BaseContactQLVisitor
}

// NewVisitor creates a new ContactQL visitor
func NewVisitor() *Visitor {
	return &Visitor{}
}

// Visit the top level parse tree
func (v *Visitor) Visit(tree antlr.ParseTree) interface{} {
	return tree.Accept(v)
}

// parse: expression
func (v *Visitor) VisitParse(ctx *gen.ParseContext) interface{} {
	return v.Visit(ctx.Expression())
}

// expression : TEXT
func (v *Visitor) VisitImplicitCondition(ctx *gen.ImplicitConditionContext) interface{} {
	return &Condition{key: ImplicitKey, comparator: "=", value: ctx.TEXT().GetText()}
}

// expression : TEXT COMPARATOR literal
func (v *Visitor) VisitCondition(ctx *gen.ConditionContext) interface{} {
	key := strings.ToLower(ctx.TEXT().GetText())
	comparator := strings.ToLower(ctx.COMPARATOR().GetText())
	value := v.Visit(ctx.Literal()).(string)

	resolvedAlias, isAlias := comparatorAliases[comparator]
	if isAlias {
		comparator = resolvedAlias
	}

	return &Condition{key: key, comparator: comparator, value: value}
}

// expression : expression AND expression
func (v *Visitor) VisitCombinationAnd(ctx *gen.CombinationAndContext) interface{} {
	child1 := v.Visit(ctx.Expression(0)).(QueryNode)
	child2 := v.Visit(ctx.Expression(1)).(QueryNode)
	return NewBoolCombination(boolOpAnd, child1, child2)
}

// expression : expression expression
func (v *Visitor) VisitCombinationImpicitAnd(ctx *gen.CombinationImpicitAndContext) interface{} {
	child1 := v.Visit(ctx.Expression(0)).(QueryNode)
	child2 := v.Visit(ctx.Expression(1)).(QueryNode)
	return NewBoolCombination(boolOpAnd, child1, child2)
}

// expression : expression OR expression
func (v *Visitor) VisitCombinationOr(ctx *gen.CombinationOrContext) interface{} {
	child1 := v.Visit(ctx.Expression(0)).(QueryNode)
	child2 := v.Visit(ctx.Expression(1)).(QueryNode)
	return NewBoolCombination(boolOpOr, child1, child2)
}

// expression : LPAREN expression RPAREN
func (v *Visitor) VisitExpressionGrouping(ctx *gen.ExpressionGroupingContext) interface{} {
	return v.Visit(ctx.Expression())
}

// literal : TEXT
func (v *Visitor) VisitTextLiteral(ctx *gen.TextLiteralContext) interface{} {
	return ctx.GetText()
}

// literal : STRING
func (v *Visitor) VisitStringLiteral(ctx *gen.StringLiteralContext) interface{} {
	value := ctx.GetText()
	value = value[1 : len(value)-1]
	return strings.Replace(value, `""`, `"`, -1) // unescape embedded quotes
}
