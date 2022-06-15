// Code generated from Excellent1.g4 by ANTLR 4.10.1. DO NOT EDIT.

package gen // Excellent1
import "github.com/antlr/antlr4/runtime/Go/antlr"

// A complete Visitor for a parse tree produced by Excellent1Parser.
type Excellent1Visitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by Excellent1Parser#parse.
	VisitParse(ctx *ParseContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#decimalLiteral.
	VisitDecimalLiteral(ctx *DecimalLiteralContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#parentheses.
	VisitParentheses(ctx *ParenthesesContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#negation.
	VisitNegation(ctx *NegationContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#exponentExpression.
	VisitExponentExpression(ctx *ExponentExpressionContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#additionOrSubtractionExpression.
	VisitAdditionOrSubtractionExpression(ctx *AdditionOrSubtractionExpressionContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#false.
	VisitFalse(ctx *FalseContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#contextReference.
	VisitContextReference(ctx *ContextReferenceContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#comparisonExpression.
	VisitComparisonExpression(ctx *ComparisonExpressionContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#concatenation.
	VisitConcatenation(ctx *ConcatenationContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#stringLiteral.
	VisitStringLiteral(ctx *StringLiteralContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#functionCall.
	VisitFunctionCall(ctx *FunctionCallContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#true.
	VisitTrue(ctx *TrueContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#equalityExpression.
	VisitEqualityExpression(ctx *EqualityExpressionContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#multiplicationOrDivisionExpression.
	VisitMultiplicationOrDivisionExpression(ctx *MultiplicationOrDivisionExpressionContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#fnname.
	VisitFnname(ctx *FnnameContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#functionParameters.
	VisitFunctionParameters(ctx *FunctionParametersContext) interface{}
}
