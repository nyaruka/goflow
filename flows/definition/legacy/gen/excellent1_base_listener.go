// Code generated from Excellent1.g4 by ANTLR 4.10.1. DO NOT EDIT.

package gen // Excellent1
import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseExcellent1Listener is a complete listener for a parse tree produced by Excellent1Parser.
type BaseExcellent1Listener struct{}

var _ Excellent1Listener = &BaseExcellent1Listener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseExcellent1Listener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseExcellent1Listener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseExcellent1Listener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseExcellent1Listener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterParse is called when production parse is entered.
func (s *BaseExcellent1Listener) EnterParse(ctx *ParseContext) {}

// ExitParse is called when production parse is exited.
func (s *BaseExcellent1Listener) ExitParse(ctx *ParseContext) {}

// EnterDecimalLiteral is called when production decimalLiteral is entered.
func (s *BaseExcellent1Listener) EnterDecimalLiteral(ctx *DecimalLiteralContext) {}

// ExitDecimalLiteral is called when production decimalLiteral is exited.
func (s *BaseExcellent1Listener) ExitDecimalLiteral(ctx *DecimalLiteralContext) {}

// EnterParentheses is called when production parentheses is entered.
func (s *BaseExcellent1Listener) EnterParentheses(ctx *ParenthesesContext) {}

// ExitParentheses is called when production parentheses is exited.
func (s *BaseExcellent1Listener) ExitParentheses(ctx *ParenthesesContext) {}

// EnterNegation is called when production negation is entered.
func (s *BaseExcellent1Listener) EnterNegation(ctx *NegationContext) {}

// ExitNegation is called when production negation is exited.
func (s *BaseExcellent1Listener) ExitNegation(ctx *NegationContext) {}

// EnterExponentExpression is called when production exponentExpression is entered.
func (s *BaseExcellent1Listener) EnterExponentExpression(ctx *ExponentExpressionContext) {}

// ExitExponentExpression is called when production exponentExpression is exited.
func (s *BaseExcellent1Listener) ExitExponentExpression(ctx *ExponentExpressionContext) {}

// EnterAdditionOrSubtractionExpression is called when production additionOrSubtractionExpression is entered.
func (s *BaseExcellent1Listener) EnterAdditionOrSubtractionExpression(ctx *AdditionOrSubtractionExpressionContext) {
}

// ExitAdditionOrSubtractionExpression is called when production additionOrSubtractionExpression is exited.
func (s *BaseExcellent1Listener) ExitAdditionOrSubtractionExpression(ctx *AdditionOrSubtractionExpressionContext) {
}

// EnterFalse is called when production false is entered.
func (s *BaseExcellent1Listener) EnterFalse(ctx *FalseContext) {}

// ExitFalse is called when production false is exited.
func (s *BaseExcellent1Listener) ExitFalse(ctx *FalseContext) {}

// EnterContextReference is called when production contextReference is entered.
func (s *BaseExcellent1Listener) EnterContextReference(ctx *ContextReferenceContext) {}

// ExitContextReference is called when production contextReference is exited.
func (s *BaseExcellent1Listener) ExitContextReference(ctx *ContextReferenceContext) {}

// EnterComparisonExpression is called when production comparisonExpression is entered.
func (s *BaseExcellent1Listener) EnterComparisonExpression(ctx *ComparisonExpressionContext) {}

// ExitComparisonExpression is called when production comparisonExpression is exited.
func (s *BaseExcellent1Listener) ExitComparisonExpression(ctx *ComparisonExpressionContext) {}

// EnterConcatenation is called when production concatenation is entered.
func (s *BaseExcellent1Listener) EnterConcatenation(ctx *ConcatenationContext) {}

// ExitConcatenation is called when production concatenation is exited.
func (s *BaseExcellent1Listener) ExitConcatenation(ctx *ConcatenationContext) {}

// EnterStringLiteral is called when production stringLiteral is entered.
func (s *BaseExcellent1Listener) EnterStringLiteral(ctx *StringLiteralContext) {}

// ExitStringLiteral is called when production stringLiteral is exited.
func (s *BaseExcellent1Listener) ExitStringLiteral(ctx *StringLiteralContext) {}

// EnterFunctionCall is called when production functionCall is entered.
func (s *BaseExcellent1Listener) EnterFunctionCall(ctx *FunctionCallContext) {}

// ExitFunctionCall is called when production functionCall is exited.
func (s *BaseExcellent1Listener) ExitFunctionCall(ctx *FunctionCallContext) {}

// EnterTrue is called when production true is entered.
func (s *BaseExcellent1Listener) EnterTrue(ctx *TrueContext) {}

// ExitTrue is called when production true is exited.
func (s *BaseExcellent1Listener) ExitTrue(ctx *TrueContext) {}

// EnterEqualityExpression is called when production equalityExpression is entered.
func (s *BaseExcellent1Listener) EnterEqualityExpression(ctx *EqualityExpressionContext) {}

// ExitEqualityExpression is called when production equalityExpression is exited.
func (s *BaseExcellent1Listener) ExitEqualityExpression(ctx *EqualityExpressionContext) {}

// EnterMultiplicationOrDivisionExpression is called when production multiplicationOrDivisionExpression is entered.
func (s *BaseExcellent1Listener) EnterMultiplicationOrDivisionExpression(ctx *MultiplicationOrDivisionExpressionContext) {
}

// ExitMultiplicationOrDivisionExpression is called when production multiplicationOrDivisionExpression is exited.
func (s *BaseExcellent1Listener) ExitMultiplicationOrDivisionExpression(ctx *MultiplicationOrDivisionExpressionContext) {
}

// EnterFnname is called when production fnname is entered.
func (s *BaseExcellent1Listener) EnterFnname(ctx *FnnameContext) {}

// ExitFnname is called when production fnname is exited.
func (s *BaseExcellent1Listener) ExitFnname(ctx *FnnameContext) {}

// EnterFunctionParameters is called when production functionParameters is entered.
func (s *BaseExcellent1Listener) EnterFunctionParameters(ctx *FunctionParametersContext) {}

// ExitFunctionParameters is called when production functionParameters is exited.
func (s *BaseExcellent1Listener) ExitFunctionParameters(ctx *FunctionParametersContext) {}
