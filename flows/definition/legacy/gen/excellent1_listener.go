// Code generated from Excellent1.g4 by ANTLR 4.10.1. DO NOT EDIT.

package gen // Excellent1
import "github.com/antlr/antlr4/runtime/Go/antlr"

// Excellent1Listener is a complete listener for a parse tree produced by Excellent1Parser.
type Excellent1Listener interface {
	antlr.ParseTreeListener

	// EnterParse is called when entering the parse production.
	EnterParse(c *ParseContext)

	// EnterDecimalLiteral is called when entering the decimalLiteral production.
	EnterDecimalLiteral(c *DecimalLiteralContext)

	// EnterParentheses is called when entering the parentheses production.
	EnterParentheses(c *ParenthesesContext)

	// EnterNegation is called when entering the negation production.
	EnterNegation(c *NegationContext)

	// EnterExponentExpression is called when entering the exponentExpression production.
	EnterExponentExpression(c *ExponentExpressionContext)

	// EnterAdditionOrSubtractionExpression is called when entering the additionOrSubtractionExpression production.
	EnterAdditionOrSubtractionExpression(c *AdditionOrSubtractionExpressionContext)

	// EnterFalse is called when entering the false production.
	EnterFalse(c *FalseContext)

	// EnterContextReference is called when entering the contextReference production.
	EnterContextReference(c *ContextReferenceContext)

	// EnterComparisonExpression is called when entering the comparisonExpression production.
	EnterComparisonExpression(c *ComparisonExpressionContext)

	// EnterConcatenation is called when entering the concatenation production.
	EnterConcatenation(c *ConcatenationContext)

	// EnterStringLiteral is called when entering the stringLiteral production.
	EnterStringLiteral(c *StringLiteralContext)

	// EnterFunctionCall is called when entering the functionCall production.
	EnterFunctionCall(c *FunctionCallContext)

	// EnterTrue is called when entering the true production.
	EnterTrue(c *TrueContext)

	// EnterEqualityExpression is called when entering the equalityExpression production.
	EnterEqualityExpression(c *EqualityExpressionContext)

	// EnterMultiplicationOrDivisionExpression is called when entering the multiplicationOrDivisionExpression production.
	EnterMultiplicationOrDivisionExpression(c *MultiplicationOrDivisionExpressionContext)

	// EnterFnname is called when entering the fnname production.
	EnterFnname(c *FnnameContext)

	// EnterFunctionParameters is called when entering the functionParameters production.
	EnterFunctionParameters(c *FunctionParametersContext)

	// ExitParse is called when exiting the parse production.
	ExitParse(c *ParseContext)

	// ExitDecimalLiteral is called when exiting the decimalLiteral production.
	ExitDecimalLiteral(c *DecimalLiteralContext)

	// ExitParentheses is called when exiting the parentheses production.
	ExitParentheses(c *ParenthesesContext)

	// ExitNegation is called when exiting the negation production.
	ExitNegation(c *NegationContext)

	// ExitExponentExpression is called when exiting the exponentExpression production.
	ExitExponentExpression(c *ExponentExpressionContext)

	// ExitAdditionOrSubtractionExpression is called when exiting the additionOrSubtractionExpression production.
	ExitAdditionOrSubtractionExpression(c *AdditionOrSubtractionExpressionContext)

	// ExitFalse is called when exiting the false production.
	ExitFalse(c *FalseContext)

	// ExitContextReference is called when exiting the contextReference production.
	ExitContextReference(c *ContextReferenceContext)

	// ExitComparisonExpression is called when exiting the comparisonExpression production.
	ExitComparisonExpression(c *ComparisonExpressionContext)

	// ExitConcatenation is called when exiting the concatenation production.
	ExitConcatenation(c *ConcatenationContext)

	// ExitStringLiteral is called when exiting the stringLiteral production.
	ExitStringLiteral(c *StringLiteralContext)

	// ExitFunctionCall is called when exiting the functionCall production.
	ExitFunctionCall(c *FunctionCallContext)

	// ExitTrue is called when exiting the true production.
	ExitTrue(c *TrueContext)

	// ExitEqualityExpression is called when exiting the equalityExpression production.
	ExitEqualityExpression(c *EqualityExpressionContext)

	// ExitMultiplicationOrDivisionExpression is called when exiting the multiplicationOrDivisionExpression production.
	ExitMultiplicationOrDivisionExpression(c *MultiplicationOrDivisionExpressionContext)

	// ExitFnname is called when exiting the fnname production.
	ExitFnname(c *FnnameContext)

	// ExitFunctionParameters is called when exiting the functionParameters production.
	ExitFunctionParameters(c *FunctionParametersContext)
}
