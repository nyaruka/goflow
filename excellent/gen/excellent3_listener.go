// Code generated from Excellent3.g4 by ANTLR 4.10.1. DO NOT EDIT.

package gen // Excellent3
import "github.com/antlr/antlr4/runtime/Go/antlr"

// Excellent3Listener is a complete listener for a parse tree produced by Excellent3Parser.
type Excellent3Listener interface {
	antlr.ParseTreeListener

	// EnterParse is called when entering the parse production.
	EnterParse(c *ParseContext)

	// EnterNegation is called when entering the negation production.
	EnterNegation(c *NegationContext)

	// EnterComparison is called when entering the comparison production.
	EnterComparison(c *ComparisonContext)

	// EnterFalse is called when entering the false production.
	EnterFalse(c *FalseContext)

	// EnterAdditionOrSubtraction is called when entering the additionOrSubtraction production.
	EnterAdditionOrSubtraction(c *AdditionOrSubtractionContext)

	// EnterTextLiteral is called when entering the textLiteral production.
	EnterTextLiteral(c *TextLiteralContext)

	// EnterConcatenation is called when entering the concatenation production.
	EnterConcatenation(c *ConcatenationContext)

	// EnterNull is called when entering the null production.
	EnterNull(c *NullContext)

	// EnterMultiplicationOrDivision is called when entering the multiplicationOrDivision production.
	EnterMultiplicationOrDivision(c *MultiplicationOrDivisionContext)

	// EnterTrue is called when entering the true production.
	EnterTrue(c *TrueContext)

	// EnterAtomReference is called when entering the atomReference production.
	EnterAtomReference(c *AtomReferenceContext)

	// EnterAnonFunction is called when entering the anonFunction production.
	EnterAnonFunction(c *AnonFunctionContext)

	// EnterEquality is called when entering the equality production.
	EnterEquality(c *EqualityContext)

	// EnterNumberLiteral is called when entering the numberLiteral production.
	EnterNumberLiteral(c *NumberLiteralContext)

	// EnterExponent is called when entering the exponent production.
	EnterExponent(c *ExponentContext)

	// EnterParentheses is called when entering the parentheses production.
	EnterParentheses(c *ParenthesesContext)

	// EnterDotLookup is called when entering the dotLookup production.
	EnterDotLookup(c *DotLookupContext)

	// EnterFunctionCall is called when entering the functionCall production.
	EnterFunctionCall(c *FunctionCallContext)

	// EnterArrayLookup is called when entering the arrayLookup production.
	EnterArrayLookup(c *ArrayLookupContext)

	// EnterContextReference is called when entering the contextReference production.
	EnterContextReference(c *ContextReferenceContext)

	// EnterFunctionParameters is called when entering the functionParameters production.
	EnterFunctionParameters(c *FunctionParametersContext)

	// EnterNameList is called when entering the nameList production.
	EnterNameList(c *NameListContext)

	// ExitParse is called when exiting the parse production.
	ExitParse(c *ParseContext)

	// ExitNegation is called when exiting the negation production.
	ExitNegation(c *NegationContext)

	// ExitComparison is called when exiting the comparison production.
	ExitComparison(c *ComparisonContext)

	// ExitFalse is called when exiting the false production.
	ExitFalse(c *FalseContext)

	// ExitAdditionOrSubtraction is called when exiting the additionOrSubtraction production.
	ExitAdditionOrSubtraction(c *AdditionOrSubtractionContext)

	// ExitTextLiteral is called when exiting the textLiteral production.
	ExitTextLiteral(c *TextLiteralContext)

	// ExitConcatenation is called when exiting the concatenation production.
	ExitConcatenation(c *ConcatenationContext)

	// ExitNull is called when exiting the null production.
	ExitNull(c *NullContext)

	// ExitMultiplicationOrDivision is called when exiting the multiplicationOrDivision production.
	ExitMultiplicationOrDivision(c *MultiplicationOrDivisionContext)

	// ExitTrue is called when exiting the true production.
	ExitTrue(c *TrueContext)

	// ExitAtomReference is called when exiting the atomReference production.
	ExitAtomReference(c *AtomReferenceContext)

	// ExitAnonFunction is called when exiting the anonFunction production.
	ExitAnonFunction(c *AnonFunctionContext)

	// ExitEquality is called when exiting the equality production.
	ExitEquality(c *EqualityContext)

	// ExitNumberLiteral is called when exiting the numberLiteral production.
	ExitNumberLiteral(c *NumberLiteralContext)

	// ExitExponent is called when exiting the exponent production.
	ExitExponent(c *ExponentContext)

	// ExitParentheses is called when exiting the parentheses production.
	ExitParentheses(c *ParenthesesContext)

	// ExitDotLookup is called when exiting the dotLookup production.
	ExitDotLookup(c *DotLookupContext)

	// ExitFunctionCall is called when exiting the functionCall production.
	ExitFunctionCall(c *FunctionCallContext)

	// ExitArrayLookup is called when exiting the arrayLookup production.
	ExitArrayLookup(c *ArrayLookupContext)

	// ExitContextReference is called when exiting the contextReference production.
	ExitContextReference(c *ContextReferenceContext)

	// ExitFunctionParameters is called when exiting the functionParameters production.
	ExitFunctionParameters(c *FunctionParametersContext)

	// ExitNameList is called when exiting the nameList production.
	ExitNameList(c *NameListContext)
}
