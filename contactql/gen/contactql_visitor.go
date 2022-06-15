// Code generated from ContactQL.g4 by ANTLR 4.10.1. DO NOT EDIT.

package gen // ContactQL
import "github.com/antlr/antlr4/runtime/Go/antlr"

// A complete Visitor for a parse tree produced by ContactQLParser.
type ContactQLVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by ContactQLParser#parse.
	VisitParse(ctx *ParseContext) interface{}

	// Visit a parse tree produced by ContactQLParser#implicitCondition.
	VisitImplicitCondition(ctx *ImplicitConditionContext) interface{}

	// Visit a parse tree produced by ContactQLParser#condition.
	VisitCondition(ctx *ConditionContext) interface{}

	// Visit a parse tree produced by ContactQLParser#combinationAnd.
	VisitCombinationAnd(ctx *CombinationAndContext) interface{}

	// Visit a parse tree produced by ContactQLParser#combinationImpicitAnd.
	VisitCombinationImpicitAnd(ctx *CombinationImpicitAndContext) interface{}

	// Visit a parse tree produced by ContactQLParser#combinationOr.
	VisitCombinationOr(ctx *CombinationOrContext) interface{}

	// Visit a parse tree produced by ContactQLParser#expressionGrouping.
	VisitExpressionGrouping(ctx *ExpressionGroupingContext) interface{}

	// Visit a parse tree produced by ContactQLParser#textLiteral.
	VisitTextLiteral(ctx *TextLiteralContext) interface{}

	// Visit a parse tree produced by ContactQLParser#stringLiteral.
	VisitStringLiteral(ctx *StringLiteralContext) interface{}
}
