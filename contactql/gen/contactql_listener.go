// Code generated from ContactQL.g4 by ANTLR 4.10.1. DO NOT EDIT.

package gen // ContactQL
import "github.com/antlr/antlr4/runtime/Go/antlr"

// ContactQLListener is a complete listener for a parse tree produced by ContactQLParser.
type ContactQLListener interface {
	antlr.ParseTreeListener

	// EnterParse is called when entering the parse production.
	EnterParse(c *ParseContext)

	// EnterImplicitCondition is called when entering the implicitCondition production.
	EnterImplicitCondition(c *ImplicitConditionContext)

	// EnterCondition is called when entering the condition production.
	EnterCondition(c *ConditionContext)

	// EnterCombinationAnd is called when entering the combinationAnd production.
	EnterCombinationAnd(c *CombinationAndContext)

	// EnterCombinationImpicitAnd is called when entering the combinationImpicitAnd production.
	EnterCombinationImpicitAnd(c *CombinationImpicitAndContext)

	// EnterCombinationOr is called when entering the combinationOr production.
	EnterCombinationOr(c *CombinationOrContext)

	// EnterExpressionGrouping is called when entering the expressionGrouping production.
	EnterExpressionGrouping(c *ExpressionGroupingContext)

	// EnterTextLiteral is called when entering the textLiteral production.
	EnterTextLiteral(c *TextLiteralContext)

	// EnterStringLiteral is called when entering the stringLiteral production.
	EnterStringLiteral(c *StringLiteralContext)

	// ExitParse is called when exiting the parse production.
	ExitParse(c *ParseContext)

	// ExitImplicitCondition is called when exiting the implicitCondition production.
	ExitImplicitCondition(c *ImplicitConditionContext)

	// ExitCondition is called when exiting the condition production.
	ExitCondition(c *ConditionContext)

	// ExitCombinationAnd is called when exiting the combinationAnd production.
	ExitCombinationAnd(c *CombinationAndContext)

	// ExitCombinationImpicitAnd is called when exiting the combinationImpicitAnd production.
	ExitCombinationImpicitAnd(c *CombinationImpicitAndContext)

	// ExitCombinationOr is called when exiting the combinationOr production.
	ExitCombinationOr(c *CombinationOrContext)

	// ExitExpressionGrouping is called when exiting the expressionGrouping production.
	ExitExpressionGrouping(c *ExpressionGroupingContext)

	// ExitTextLiteral is called when exiting the textLiteral production.
	ExitTextLiteral(c *TextLiteralContext)

	// ExitStringLiteral is called when exiting the stringLiteral production.
	ExitStringLiteral(c *StringLiteralContext)
}
