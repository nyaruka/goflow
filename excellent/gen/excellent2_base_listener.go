// Code generated from Excellent2.g4 by ANTLR 4.7.2. DO NOT EDIT.

package gen // Excellent2
import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseExcellent2Listener is a complete listener for a parse tree produced by Excellent2Parser.
type BaseExcellent2Listener struct{}

var _ Excellent2Listener = &BaseExcellent2Listener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseExcellent2Listener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseExcellent2Listener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseExcellent2Listener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseExcellent2Listener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterParse is called when production parse is entered.
func (s *BaseExcellent2Listener) EnterParse(ctx *ParseContext) {}

// ExitParse is called when production parse is exited.
func (s *BaseExcellent2Listener) ExitParse(ctx *ParseContext) {}

// EnterNegation is called when production negation is entered.
func (s *BaseExcellent2Listener) EnterNegation(ctx *NegationContext) {}

// ExitNegation is called when production negation is exited.
func (s *BaseExcellent2Listener) ExitNegation(ctx *NegationContext) {}

// EnterComparison is called when production comparison is entered.
func (s *BaseExcellent2Listener) EnterComparison(ctx *ComparisonContext) {}

// ExitComparison is called when production comparison is exited.
func (s *BaseExcellent2Listener) ExitComparison(ctx *ComparisonContext) {}

// EnterFalse is called when production false is entered.
func (s *BaseExcellent2Listener) EnterFalse(ctx *FalseContext) {}

// ExitFalse is called when production false is exited.
func (s *BaseExcellent2Listener) ExitFalse(ctx *FalseContext) {}

// EnterAdditionOrSubtraction is called when production additionOrSubtraction is entered.
func (s *BaseExcellent2Listener) EnterAdditionOrSubtraction(ctx *AdditionOrSubtractionContext) {}

// ExitAdditionOrSubtraction is called when production additionOrSubtraction is exited.
func (s *BaseExcellent2Listener) ExitAdditionOrSubtraction(ctx *AdditionOrSubtractionContext) {}

// EnterTextLiteral is called when production textLiteral is entered.
func (s *BaseExcellent2Listener) EnterTextLiteral(ctx *TextLiteralContext) {}

// ExitTextLiteral is called when production textLiteral is exited.
func (s *BaseExcellent2Listener) ExitTextLiteral(ctx *TextLiteralContext) {}

// EnterConcatenation is called when production concatenation is entered.
func (s *BaseExcellent2Listener) EnterConcatenation(ctx *ConcatenationContext) {}

// ExitConcatenation is called when production concatenation is exited.
func (s *BaseExcellent2Listener) ExitConcatenation(ctx *ConcatenationContext) {}

// EnterNull is called when production null is entered.
func (s *BaseExcellent2Listener) EnterNull(ctx *NullContext) {}

// ExitNull is called when production null is exited.
func (s *BaseExcellent2Listener) ExitNull(ctx *NullContext) {}

// EnterMultiplicationOrDivision is called when production multiplicationOrDivision is entered.
func (s *BaseExcellent2Listener) EnterMultiplicationOrDivision(ctx *MultiplicationOrDivisionContext) {}

// ExitMultiplicationOrDivision is called when production multiplicationOrDivision is exited.
func (s *BaseExcellent2Listener) ExitMultiplicationOrDivision(ctx *MultiplicationOrDivisionContext) {}

// EnterTrue is called when production true is entered.
func (s *BaseExcellent2Listener) EnterTrue(ctx *TrueContext) {}

// ExitTrue is called when production true is exited.
func (s *BaseExcellent2Listener) ExitTrue(ctx *TrueContext) {}

// EnterAtomReference is called when production atomReference is entered.
func (s *BaseExcellent2Listener) EnterAtomReference(ctx *AtomReferenceContext) {}

// ExitAtomReference is called when production atomReference is exited.
func (s *BaseExcellent2Listener) ExitAtomReference(ctx *AtomReferenceContext) {}

// EnterEquality is called when production equality is entered.
func (s *BaseExcellent2Listener) EnterEquality(ctx *EqualityContext) {}

// ExitEquality is called when production equality is exited.
func (s *BaseExcellent2Listener) ExitEquality(ctx *EqualityContext) {}

// EnterNumberLiteral is called when production numberLiteral is entered.
func (s *BaseExcellent2Listener) EnterNumberLiteral(ctx *NumberLiteralContext) {}

// ExitNumberLiteral is called when production numberLiteral is exited.
func (s *BaseExcellent2Listener) ExitNumberLiteral(ctx *NumberLiteralContext) {}

// EnterExponent is called when production exponent is entered.
func (s *BaseExcellent2Listener) EnterExponent(ctx *ExponentContext) {}

// ExitExponent is called when production exponent is exited.
func (s *BaseExcellent2Listener) ExitExponent(ctx *ExponentContext) {}

// EnterParentheses is called when production parentheses is entered.
func (s *BaseExcellent2Listener) EnterParentheses(ctx *ParenthesesContext) {}

// ExitParentheses is called when production parentheses is exited.
func (s *BaseExcellent2Listener) ExitParentheses(ctx *ParenthesesContext) {}

// EnterDotLookup is called when production dotLookup is entered.
func (s *BaseExcellent2Listener) EnterDotLookup(ctx *DotLookupContext) {}

// ExitDotLookup is called when production dotLookup is exited.
func (s *BaseExcellent2Listener) ExitDotLookup(ctx *DotLookupContext) {}

// EnterFunctionCall is called when production functionCall is entered.
func (s *BaseExcellent2Listener) EnterFunctionCall(ctx *FunctionCallContext) {}

// ExitFunctionCall is called when production functionCall is exited.
func (s *BaseExcellent2Listener) ExitFunctionCall(ctx *FunctionCallContext) {}

// EnterArrayLookup is called when production arrayLookup is entered.
func (s *BaseExcellent2Listener) EnterArrayLookup(ctx *ArrayLookupContext) {}

// ExitArrayLookup is called when production arrayLookup is exited.
func (s *BaseExcellent2Listener) ExitArrayLookup(ctx *ArrayLookupContext) {}

// EnterContextReference is called when production contextReference is entered.
func (s *BaseExcellent2Listener) EnterContextReference(ctx *ContextReferenceContext) {}

// ExitContextReference is called when production contextReference is exited.
func (s *BaseExcellent2Listener) ExitContextReference(ctx *ContextReferenceContext) {}

// EnterFunctionParameters is called when production functionParameters is entered.
func (s *BaseExcellent2Listener) EnterFunctionParameters(ctx *FunctionParametersContext) {}

// ExitFunctionParameters is called when production functionParameters is exited.
func (s *BaseExcellent2Listener) ExitFunctionParameters(ctx *FunctionParametersContext) {}
