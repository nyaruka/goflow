// Code generated from Excellent3.g4 by ANTLR 4.10.1. DO NOT EDIT.

package gen // Excellent3
import "github.com/antlr/antlr4/runtime/Go/antlr"

// A complete Visitor for a parse tree produced by Excellent3Parser.
type Excellent3Visitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by Excellent3Parser#parse.
	VisitParse(ctx *ParseContext) interface{}

	// Visit a parse tree produced by Excellent3Parser#negation.
	VisitNegation(ctx *NegationContext) interface{}

	// Visit a parse tree produced by Excellent3Parser#comparison.
	VisitComparison(ctx *ComparisonContext) interface{}

	// Visit a parse tree produced by Excellent3Parser#false.
	VisitFalse(ctx *FalseContext) interface{}

	// Visit a parse tree produced by Excellent3Parser#additionOrSubtraction.
	VisitAdditionOrSubtraction(ctx *AdditionOrSubtractionContext) interface{}

	// Visit a parse tree produced by Excellent3Parser#textLiteral.
	VisitTextLiteral(ctx *TextLiteralContext) interface{}

	// Visit a parse tree produced by Excellent3Parser#concatenation.
	VisitConcatenation(ctx *ConcatenationContext) interface{}

	// Visit a parse tree produced by Excellent3Parser#null.
	VisitNull(ctx *NullContext) interface{}

	// Visit a parse tree produced by Excellent3Parser#multiplicationOrDivision.
	VisitMultiplicationOrDivision(ctx *MultiplicationOrDivisionContext) interface{}

	// Visit a parse tree produced by Excellent3Parser#true.
	VisitTrue(ctx *TrueContext) interface{}

	// Visit a parse tree produced by Excellent3Parser#atomReference.
	VisitAtomReference(ctx *AtomReferenceContext) interface{}

	// Visit a parse tree produced by Excellent3Parser#anonFunction.
	VisitAnonFunction(ctx *AnonFunctionContext) interface{}

	// Visit a parse tree produced by Excellent3Parser#equality.
	VisitEquality(ctx *EqualityContext) interface{}

	// Visit a parse tree produced by Excellent3Parser#numberLiteral.
	VisitNumberLiteral(ctx *NumberLiteralContext) interface{}

	// Visit a parse tree produced by Excellent3Parser#exponent.
	VisitExponent(ctx *ExponentContext) interface{}

	// Visit a parse tree produced by Excellent3Parser#parentheses.
	VisitParentheses(ctx *ParenthesesContext) interface{}

	// Visit a parse tree produced by Excellent3Parser#dotLookup.
	VisitDotLookup(ctx *DotLookupContext) interface{}

	// Visit a parse tree produced by Excellent3Parser#functionCall.
	VisitFunctionCall(ctx *FunctionCallContext) interface{}

	// Visit a parse tree produced by Excellent3Parser#arrayLookup.
	VisitArrayLookup(ctx *ArrayLookupContext) interface{}

	// Visit a parse tree produced by Excellent3Parser#contextReference.
	VisitContextReference(ctx *ContextReferenceContext) interface{}

	// Visit a parse tree produced by Excellent3Parser#functionParameters.
	VisitFunctionParameters(ctx *FunctionParametersContext) interface{}

	// Visit a parse tree produced by Excellent3Parser#nameList.
	VisitNameList(ctx *NameListContext) interface{}
}
