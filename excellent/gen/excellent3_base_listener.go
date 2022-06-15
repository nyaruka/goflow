// Code generated from Excellent3.g4 by ANTLR 4.10.1. DO NOT EDIT.

package gen // Excellent3
import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseExcellent3Listener is a complete listener for a parse tree produced by Excellent3Parser.
type BaseExcellent3Listener struct{}

var _ Excellent3Listener = &BaseExcellent3Listener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseExcellent3Listener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseExcellent3Listener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseExcellent3Listener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseExcellent3Listener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterParse is called when production parse is entered.
func (s *BaseExcellent3Listener) EnterParse(ctx *ParseContext) {}

// ExitParse is called when production parse is exited.
func (s *BaseExcellent3Listener) ExitParse(ctx *ParseContext) {}

// EnterNegation is called when production negation is entered.
func (s *BaseExcellent3Listener) EnterNegation(ctx *NegationContext) {}

// ExitNegation is called when production negation is exited.
func (s *BaseExcellent3Listener) ExitNegation(ctx *NegationContext) {}

// EnterComparison is called when production comparison is entered.
func (s *BaseExcellent3Listener) EnterComparison(ctx *ComparisonContext) {}

// ExitComparison is called when production comparison is exited.
func (s *BaseExcellent3Listener) ExitComparison(ctx *ComparisonContext) {}

// EnterFalse is called when production false is entered.
func (s *BaseExcellent3Listener) EnterFalse(ctx *FalseContext) {}

// ExitFalse is called when production false is exited.
func (s *BaseExcellent3Listener) ExitFalse(ctx *FalseContext) {}

// EnterAdditionOrSubtraction is called when production additionOrSubtraction is entered.
func (s *BaseExcellent3Listener) EnterAdditionOrSubtraction(ctx *AdditionOrSubtractionContext) {}

// ExitAdditionOrSubtraction is called when production additionOrSubtraction is exited.
func (s *BaseExcellent3Listener) ExitAdditionOrSubtraction(ctx *AdditionOrSubtractionContext) {}

// EnterTextLiteral is called when production textLiteral is entered.
func (s *BaseExcellent3Listener) EnterTextLiteral(ctx *TextLiteralContext) {}

// ExitTextLiteral is called when production textLiteral is exited.
func (s *BaseExcellent3Listener) ExitTextLiteral(ctx *TextLiteralContext) {}

// EnterConcatenation is called when production concatenation is entered.
func (s *BaseExcellent3Listener) EnterConcatenation(ctx *ConcatenationContext) {}

// ExitConcatenation is called when production concatenation is exited.
func (s *BaseExcellent3Listener) ExitConcatenation(ctx *ConcatenationContext) {}

// EnterNull is called when production null is entered.
func (s *BaseExcellent3Listener) EnterNull(ctx *NullContext) {}

// ExitNull is called when production null is exited.
func (s *BaseExcellent3Listener) ExitNull(ctx *NullContext) {}

// EnterMultiplicationOrDivision is called when production multiplicationOrDivision is entered.
func (s *BaseExcellent3Listener) EnterMultiplicationOrDivision(ctx *MultiplicationOrDivisionContext) {
}

// ExitMultiplicationOrDivision is called when production multiplicationOrDivision is exited.
func (s *BaseExcellent3Listener) ExitMultiplicationOrDivision(ctx *MultiplicationOrDivisionContext) {}

// EnterTrue is called when production true is entered.
func (s *BaseExcellent3Listener) EnterTrue(ctx *TrueContext) {}

// ExitTrue is called when production true is exited.
func (s *BaseExcellent3Listener) ExitTrue(ctx *TrueContext) {}

// EnterAtomReference is called when production atomReference is entered.
func (s *BaseExcellent3Listener) EnterAtomReference(ctx *AtomReferenceContext) {}

// ExitAtomReference is called when production atomReference is exited.
func (s *BaseExcellent3Listener) ExitAtomReference(ctx *AtomReferenceContext) {}

// EnterAnonFunction is called when production anonFunction is entered.
func (s *BaseExcellent3Listener) EnterAnonFunction(ctx *AnonFunctionContext) {}

// ExitAnonFunction is called when production anonFunction is exited.
func (s *BaseExcellent3Listener) ExitAnonFunction(ctx *AnonFunctionContext) {}

// EnterEquality is called when production equality is entered.
func (s *BaseExcellent3Listener) EnterEquality(ctx *EqualityContext) {}

// ExitEquality is called when production equality is exited.
func (s *BaseExcellent3Listener) ExitEquality(ctx *EqualityContext) {}

// EnterNumberLiteral is called when production numberLiteral is entered.
func (s *BaseExcellent3Listener) EnterNumberLiteral(ctx *NumberLiteralContext) {}

// ExitNumberLiteral is called when production numberLiteral is exited.
func (s *BaseExcellent3Listener) ExitNumberLiteral(ctx *NumberLiteralContext) {}

// EnterExponent is called when production exponent is entered.
func (s *BaseExcellent3Listener) EnterExponent(ctx *ExponentContext) {}

// ExitExponent is called when production exponent is exited.
func (s *BaseExcellent3Listener) ExitExponent(ctx *ExponentContext) {}

// EnterParentheses is called when production parentheses is entered.
func (s *BaseExcellent3Listener) EnterParentheses(ctx *ParenthesesContext) {}

// ExitParentheses is called when production parentheses is exited.
func (s *BaseExcellent3Listener) ExitParentheses(ctx *ParenthesesContext) {}

// EnterDotLookup is called when production dotLookup is entered.
func (s *BaseExcellent3Listener) EnterDotLookup(ctx *DotLookupContext) {}

// ExitDotLookup is called when production dotLookup is exited.
func (s *BaseExcellent3Listener) ExitDotLookup(ctx *DotLookupContext) {}

// EnterFunctionCall is called when production functionCall is entered.
func (s *BaseExcellent3Listener) EnterFunctionCall(ctx *FunctionCallContext) {}

// ExitFunctionCall is called when production functionCall is exited.
func (s *BaseExcellent3Listener) ExitFunctionCall(ctx *FunctionCallContext) {}

// EnterArrayLookup is called when production arrayLookup is entered.
func (s *BaseExcellent3Listener) EnterArrayLookup(ctx *ArrayLookupContext) {}

// ExitArrayLookup is called when production arrayLookup is exited.
func (s *BaseExcellent3Listener) ExitArrayLookup(ctx *ArrayLookupContext) {}

// EnterContextReference is called when production contextReference is entered.
func (s *BaseExcellent3Listener) EnterContextReference(ctx *ContextReferenceContext) {}

// ExitContextReference is called when production contextReference is exited.
func (s *BaseExcellent3Listener) ExitContextReference(ctx *ContextReferenceContext) {}

// EnterFunctionParameters is called when production functionParameters is entered.
func (s *BaseExcellent3Listener) EnterFunctionParameters(ctx *FunctionParametersContext) {}

// ExitFunctionParameters is called when production functionParameters is exited.
func (s *BaseExcellent3Listener) ExitFunctionParameters(ctx *FunctionParametersContext) {}

// EnterNameList is called when production nameList is entered.
func (s *BaseExcellent3Listener) EnterNameList(ctx *NameListContext) {}

// ExitNameList is called when production nameList is exited.
func (s *BaseExcellent3Listener) ExitNameList(ctx *NameListContext) {}
