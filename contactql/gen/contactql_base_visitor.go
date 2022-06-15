// Code generated from ContactQL.g4 by ANTLR 4.10.1. DO NOT EDIT.

package gen // ContactQL
import "github.com/antlr/antlr4/runtime/Go/antlr"

type BaseContactQLVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseContactQLVisitor) VisitParse(ctx *ParseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseContactQLVisitor) VisitImplicitCondition(ctx *ImplicitConditionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseContactQLVisitor) VisitCondition(ctx *ConditionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseContactQLVisitor) VisitCombinationAnd(ctx *CombinationAndContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseContactQLVisitor) VisitCombinationImpicitAnd(ctx *CombinationImpicitAndContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseContactQLVisitor) VisitCombinationOr(ctx *CombinationOrContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseContactQLVisitor) VisitExpressionGrouping(ctx *ExpressionGroupingContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseContactQLVisitor) VisitTextLiteral(ctx *TextLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseContactQLVisitor) VisitStringLiteral(ctx *StringLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}
