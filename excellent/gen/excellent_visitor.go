// Generated from excellent/gen/Excellent.g4 by ANTLR 4.7.

package gen // Excellent
import "github.com/antlr/antlr4/runtime/Go/antlr"

// A complete Visitor for a parse tree produced by ExcellentParser.
type ExcellentVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by ExcellentParser#parse.
	VisitParse(ctx *ParseContext) interface{}

	// Visit a parse tree produced by ExcellentParser#decimalLiteral.
	VisitDecimalLiteral(ctx *DecimalLiteralContext) interface{}

	// Visit a parse tree produced by ExcellentParser#dotLookup.
	VisitDotLookup(ctx *DotLookupContext) interface{}

	// Visit a parse tree produced by ExcellentParser#stringLiteral.
	VisitStringLiteral(ctx *StringLiteralContext) interface{}

	// Visit a parse tree produced by ExcellentParser#functionCall.
	VisitFunctionCall(ctx *FunctionCallContext) interface{}

	// Visit a parse tree produced by ExcellentParser#true.
	VisitTrue(ctx *TrueContext) interface{}

	// Visit a parse tree produced by ExcellentParser#false.
	VisitFalse(ctx *FalseContext) interface{}

	// Visit a parse tree produced by ExcellentParser#arrayLookup.
	VisitArrayLookup(ctx *ArrayLookupContext) interface{}

	// Visit a parse tree produced by ExcellentParser#contextReference.
	VisitContextReference(ctx *ContextReferenceContext) interface{}

	// Visit a parse tree produced by ExcellentParser#parentheses.
	VisitParentheses(ctx *ParenthesesContext) interface{}

	// Visit a parse tree produced by ExcellentParser#negation.
	VisitNegation(ctx *NegationContext) interface{}

	// Visit a parse tree produced by ExcellentParser#comparison.
	VisitComparison(ctx *ComparisonContext) interface{}

	// Visit a parse tree produced by ExcellentParser#concatenation.
	VisitConcatenation(ctx *ConcatenationContext) interface{}

	// Visit a parse tree produced by ExcellentParser#multiplicationOrDivision.
	VisitMultiplicationOrDivision(ctx *MultiplicationOrDivisionContext) interface{}

	// Visit a parse tree produced by ExcellentParser#atomReference.
	VisitAtomReference(ctx *AtomReferenceContext) interface{}

	// Visit a parse tree produced by ExcellentParser#additionOrSubtraction.
	VisitAdditionOrSubtraction(ctx *AdditionOrSubtractionContext) interface{}

	// Visit a parse tree produced by ExcellentParser#equality.
	VisitEquality(ctx *EqualityContext) interface{}

	// Visit a parse tree produced by ExcellentParser#exponent.
	VisitExponent(ctx *ExponentContext) interface{}

	// Visit a parse tree produced by ExcellentParser#fnname.
	VisitFnname(ctx *FnnameContext) interface{}

	// Visit a parse tree produced by ExcellentParser#functionParameters.
	VisitFunctionParameters(ctx *FunctionParametersContext) interface{}
}
