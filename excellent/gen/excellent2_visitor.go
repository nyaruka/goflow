// Generated from Excellent2.g4 by ANTLR 4.7.

package gen // Excellent2
import "github.com/antlr/antlr4/runtime/Go/antlr"

// A complete Visitor for a parse tree produced by Excellent2Parser.
type Excellent2Visitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by Excellent2Parser#parse.
	VisitParse(ctx *ParseContext) interface{}

	// Visit a parse tree produced by Excellent2Parser#decimalLiteral.
	VisitDecimalLiteral(ctx *DecimalLiteralContext) interface{}

	// Visit a parse tree produced by Excellent2Parser#dotLookup.
	VisitDotLookup(ctx *DotLookupContext) interface{}

	// Visit a parse tree produced by Excellent2Parser#stringLiteral.
	VisitStringLiteral(ctx *StringLiteralContext) interface{}

	// Visit a parse tree produced by Excellent2Parser#functionCall.
	VisitFunctionCall(ctx *FunctionCallContext) interface{}

	// Visit a parse tree produced by Excellent2Parser#true.
	VisitTrue(ctx *TrueContext) interface{}

	// Visit a parse tree produced by Excellent2Parser#false.
	VisitFalse(ctx *FalseContext) interface{}

	// Visit a parse tree produced by Excellent2Parser#arrayLookup.
	VisitArrayLookup(ctx *ArrayLookupContext) interface{}

	// Visit a parse tree produced by Excellent2Parser#contextReference.
	VisitContextReference(ctx *ContextReferenceContext) interface{}

	// Visit a parse tree produced by Excellent2Parser#parentheses.
	VisitParentheses(ctx *ParenthesesContext) interface{}

	// Visit a parse tree produced by Excellent2Parser#negation.
	VisitNegation(ctx *NegationContext) interface{}

	// Visit a parse tree produced by Excellent2Parser#comparison.
	VisitComparison(ctx *ComparisonContext) interface{}

	// Visit a parse tree produced by Excellent2Parser#concatenation.
	VisitConcatenation(ctx *ConcatenationContext) interface{}

	// Visit a parse tree produced by Excellent2Parser#multiplicationOrDivision.
	VisitMultiplicationOrDivision(ctx *MultiplicationOrDivisionContext) interface{}

	// Visit a parse tree produced by Excellent2Parser#atomReference.
	VisitAtomReference(ctx *AtomReferenceContext) interface{}

	// Visit a parse tree produced by Excellent2Parser#additionOrSubtraction.
	VisitAdditionOrSubtraction(ctx *AdditionOrSubtractionContext) interface{}

	// Visit a parse tree produced by Excellent2Parser#equality.
	VisitEquality(ctx *EqualityContext) interface{}

	// Visit a parse tree produced by Excellent2Parser#exponent.
	VisitExponent(ctx *ExponentContext) interface{}

	// Visit a parse tree produced by Excellent2Parser#fnname.
	VisitFnname(ctx *FnnameContext) interface{}

	// Visit a parse tree produced by Excellent2Parser#functionParameters.
	VisitFunctionParameters(ctx *FunctionParametersContext) interface{}
}
