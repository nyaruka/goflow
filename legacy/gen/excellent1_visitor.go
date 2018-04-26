// Code generated from Excellent1.g4 by ANTLR 4.7.1. DO NOT EDIT.

package gen // Excellent1
import "github.com/antlr/antlr4/runtime/Go/antlr"

// A complete Visitor for a parse tree produced by Excellent1Parser.
type Excellent1Visitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by Excellent1Parser#parse.
	VisitParse(ctx *ParseContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#decimalLiteral.
	VisitDecimalLiteral(ctx *DecimalLiteralContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#dotLookup.
	VisitDotLookup(ctx *DotLookupContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#null.
	VisitNull(ctx *NullContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#stringLiteral.
	VisitStringLiteral(ctx *StringLiteralContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#functionCall.
	VisitFunctionCall(ctx *FunctionCallContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#true.
	VisitTrue(ctx *TrueContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#false.
	VisitFalse(ctx *FalseContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#arrayLookup.
	VisitArrayLookup(ctx *ArrayLookupContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#contextReference.
	VisitContextReference(ctx *ContextReferenceContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#parentheses.
	VisitParentheses(ctx *ParenthesesContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#negation.
	VisitNegation(ctx *NegationContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#comparison.
	VisitComparison(ctx *ComparisonContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#concatenation.
	VisitConcatenation(ctx *ConcatenationContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#multiplicationOrDivision.
	VisitMultiplicationOrDivision(ctx *MultiplicationOrDivisionContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#atomReference.
	VisitAtomReference(ctx *AtomReferenceContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#additionOrSubtraction.
	VisitAdditionOrSubtraction(ctx *AdditionOrSubtractionContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#equality.
	VisitEquality(ctx *EqualityContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#exponent.
	VisitExponent(ctx *ExponentContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#fnname.
	VisitFnname(ctx *FnnameContext) interface{}

	// Visit a parse tree produced by Excellent1Parser#functionParameters.
	VisitFunctionParameters(ctx *FunctionParametersContext) interface{}
}
