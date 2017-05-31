// Generated from excellent/gen/Excellent.g4 by ANTLR 4.7.

package gen // Excellent
import "github.com/antlr/antlr4/runtime/Go/antlr"

// ExcellentListener is a complete listener for a parse tree produced by ExcellentParser.
type ExcellentListener interface {
	antlr.ParseTreeListener

	// EnterParse is called when entering the parse production.
	EnterParse(c *ParseContext)

	// EnterDecimalLiteral is called when entering the decimalLiteral production.
	EnterDecimalLiteral(c *DecimalLiteralContext)

	// EnterDotLookup is called when entering the dotLookup production.
	EnterDotLookup(c *DotLookupContext)

	// EnterStringLiteral is called when entering the stringLiteral production.
	EnterStringLiteral(c *StringLiteralContext)

	// EnterFunctionCall is called when entering the functionCall production.
	EnterFunctionCall(c *FunctionCallContext)

	// EnterTrue is called when entering the true production.
	EnterTrue(c *TrueContext)

	// EnterFalse is called when entering the false production.
	EnterFalse(c *FalseContext)

	// EnterArrayLookup is called when entering the arrayLookup production.
	EnterArrayLookup(c *ArrayLookupContext)

	// EnterContextReference is called when entering the contextReference production.
	EnterContextReference(c *ContextReferenceContext)

	// EnterParentheses is called when entering the parentheses production.
	EnterParentheses(c *ParenthesesContext)

	// EnterNegation is called when entering the negation production.
	EnterNegation(c *NegationContext)

	// EnterComparison is called when entering the comparison production.
	EnterComparison(c *ComparisonContext)

	// EnterConcatenation is called when entering the concatenation production.
	EnterConcatenation(c *ConcatenationContext)

	// EnterMultiplicationOrDivision is called when entering the multiplicationOrDivision production.
	EnterMultiplicationOrDivision(c *MultiplicationOrDivisionContext)

	// EnterAtomReference is called when entering the atomReference production.
	EnterAtomReference(c *AtomReferenceContext)

	// EnterAdditionOrSubtraction is called when entering the additionOrSubtraction production.
	EnterAdditionOrSubtraction(c *AdditionOrSubtractionContext)

	// EnterEquality is called when entering the equality production.
	EnterEquality(c *EqualityContext)

	// EnterExponent is called when entering the exponent production.
	EnterExponent(c *ExponentContext)

	// EnterFnname is called when entering the fnname production.
	EnterFnname(c *FnnameContext)

	// EnterFunctionParameters is called when entering the functionParameters production.
	EnterFunctionParameters(c *FunctionParametersContext)

	// ExitParse is called when exiting the parse production.
	ExitParse(c *ParseContext)

	// ExitDecimalLiteral is called when exiting the decimalLiteral production.
	ExitDecimalLiteral(c *DecimalLiteralContext)

	// ExitDotLookup is called when exiting the dotLookup production.
	ExitDotLookup(c *DotLookupContext)

	// ExitStringLiteral is called when exiting the stringLiteral production.
	ExitStringLiteral(c *StringLiteralContext)

	// ExitFunctionCall is called when exiting the functionCall production.
	ExitFunctionCall(c *FunctionCallContext)

	// ExitTrue is called when exiting the true production.
	ExitTrue(c *TrueContext)

	// ExitFalse is called when exiting the false production.
	ExitFalse(c *FalseContext)

	// ExitArrayLookup is called when exiting the arrayLookup production.
	ExitArrayLookup(c *ArrayLookupContext)

	// ExitContextReference is called when exiting the contextReference production.
	ExitContextReference(c *ContextReferenceContext)

	// ExitParentheses is called when exiting the parentheses production.
	ExitParentheses(c *ParenthesesContext)

	// ExitNegation is called when exiting the negation production.
	ExitNegation(c *NegationContext)

	// ExitComparison is called when exiting the comparison production.
	ExitComparison(c *ComparisonContext)

	// ExitConcatenation is called when exiting the concatenation production.
	ExitConcatenation(c *ConcatenationContext)

	// ExitMultiplicationOrDivision is called when exiting the multiplicationOrDivision production.
	ExitMultiplicationOrDivision(c *MultiplicationOrDivisionContext)

	// ExitAtomReference is called when exiting the atomReference production.
	ExitAtomReference(c *AtomReferenceContext)

	// ExitAdditionOrSubtraction is called when exiting the additionOrSubtraction production.
	ExitAdditionOrSubtraction(c *AdditionOrSubtractionContext)

	// ExitEquality is called when exiting the equality production.
	ExitEquality(c *EqualityContext)

	// ExitExponent is called when exiting the exponent production.
	ExitExponent(c *ExponentContext)

	// ExitFnname is called when exiting the fnname production.
	ExitFnname(c *FnnameContext)

	// ExitFunctionParameters is called when exiting the functionParameters production.
	ExitFunctionParameters(c *FunctionParametersContext)
}
