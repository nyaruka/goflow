// Code generated from ContactQL.g4 by ANTLR 4.10.1. DO NOT EDIT.

package gen // ContactQL
import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseContactQLListener is a complete listener for a parse tree produced by ContactQLParser.
type BaseContactQLListener struct{}

var _ ContactQLListener = &BaseContactQLListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseContactQLListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseContactQLListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseContactQLListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseContactQLListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterParse is called when production parse is entered.
func (s *BaseContactQLListener) EnterParse(ctx *ParseContext) {}

// ExitParse is called when production parse is exited.
func (s *BaseContactQLListener) ExitParse(ctx *ParseContext) {}

// EnterImplicitCondition is called when production implicitCondition is entered.
func (s *BaseContactQLListener) EnterImplicitCondition(ctx *ImplicitConditionContext) {}

// ExitImplicitCondition is called when production implicitCondition is exited.
func (s *BaseContactQLListener) ExitImplicitCondition(ctx *ImplicitConditionContext) {}

// EnterCondition is called when production condition is entered.
func (s *BaseContactQLListener) EnterCondition(ctx *ConditionContext) {}

// ExitCondition is called when production condition is exited.
func (s *BaseContactQLListener) ExitCondition(ctx *ConditionContext) {}

// EnterCombinationAnd is called when production combinationAnd is entered.
func (s *BaseContactQLListener) EnterCombinationAnd(ctx *CombinationAndContext) {}

// ExitCombinationAnd is called when production combinationAnd is exited.
func (s *BaseContactQLListener) ExitCombinationAnd(ctx *CombinationAndContext) {}

// EnterCombinationImpicitAnd is called when production combinationImpicitAnd is entered.
func (s *BaseContactQLListener) EnterCombinationImpicitAnd(ctx *CombinationImpicitAndContext) {}

// ExitCombinationImpicitAnd is called when production combinationImpicitAnd is exited.
func (s *BaseContactQLListener) ExitCombinationImpicitAnd(ctx *CombinationImpicitAndContext) {}

// EnterCombinationOr is called when production combinationOr is entered.
func (s *BaseContactQLListener) EnterCombinationOr(ctx *CombinationOrContext) {}

// ExitCombinationOr is called when production combinationOr is exited.
func (s *BaseContactQLListener) ExitCombinationOr(ctx *CombinationOrContext) {}

// EnterExpressionGrouping is called when production expressionGrouping is entered.
func (s *BaseContactQLListener) EnterExpressionGrouping(ctx *ExpressionGroupingContext) {}

// ExitExpressionGrouping is called when production expressionGrouping is exited.
func (s *BaseContactQLListener) ExitExpressionGrouping(ctx *ExpressionGroupingContext) {}

// EnterTextLiteral is called when production textLiteral is entered.
func (s *BaseContactQLListener) EnterTextLiteral(ctx *TextLiteralContext) {}

// ExitTextLiteral is called when production textLiteral is exited.
func (s *BaseContactQLListener) ExitTextLiteral(ctx *TextLiteralContext) {}

// EnterStringLiteral is called when production stringLiteral is entered.
func (s *BaseContactQLListener) EnterStringLiteral(ctx *StringLiteralContext) {}

// ExitStringLiteral is called when production stringLiteral is exited.
func (s *BaseContactQLListener) ExitStringLiteral(ctx *StringLiteralContext) {}
