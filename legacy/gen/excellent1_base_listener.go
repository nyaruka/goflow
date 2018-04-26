// Code generated from Excellent1.g4 by ANTLR 4.7.1. DO NOT EDIT.

package gen // Excellent1
import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseExcellent1Listener is a complete listener for a parse tree produced by Excellent1Parser.
type BaseExcellent1Listener struct{}

var _ Excellent1Listener = &BaseExcellent1Listener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseExcellent1Listener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseExcellent1Listener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseExcellent1Listener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseExcellent1Listener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterParse is called when production parse is entered.
func (s *BaseExcellent1Listener) EnterParse(ctx *ParseContext) {}

// ExitParse is called when production parse is exited.
func (s *BaseExcellent1Listener) ExitParse(ctx *ParseContext) {}

// EnterDecimalLiteral is called when production decimalLiteral is entered.
func (s *BaseExcellent1Listener) EnterDecimalLiteral(ctx *DecimalLiteralContext) {}

// ExitDecimalLiteral is called when production decimalLiteral is exited.
func (s *BaseExcellent1Listener) ExitDecimalLiteral(ctx *DecimalLiteralContext) {}

// EnterDotLookup is called when production dotLookup is entered.
func (s *BaseExcellent1Listener) EnterDotLookup(ctx *DotLookupContext) {}

// ExitDotLookup is called when production dotLookup is exited.
func (s *BaseExcellent1Listener) ExitDotLookup(ctx *DotLookupContext) {}

// EnterNull is called when production null is entered.
func (s *BaseExcellent1Listener) EnterNull(ctx *NullContext) {}

// ExitNull is called when production null is exited.
func (s *BaseExcellent1Listener) ExitNull(ctx *NullContext) {}

// EnterStringLiteral is called when production stringLiteral is entered.
func (s *BaseExcellent1Listener) EnterStringLiteral(ctx *StringLiteralContext) {}

// ExitStringLiteral is called when production stringLiteral is exited.
func (s *BaseExcellent1Listener) ExitStringLiteral(ctx *StringLiteralContext) {}

// EnterFunctionCall is called when production functionCall is entered.
func (s *BaseExcellent1Listener) EnterFunctionCall(ctx *FunctionCallContext) {}

// ExitFunctionCall is called when production functionCall is exited.
func (s *BaseExcellent1Listener) ExitFunctionCall(ctx *FunctionCallContext) {}

// EnterTrue is called when production true is entered.
func (s *BaseExcellent1Listener) EnterTrue(ctx *TrueContext) {}

// ExitTrue is called when production true is exited.
func (s *BaseExcellent1Listener) ExitTrue(ctx *TrueContext) {}

// EnterFalse is called when production false is entered.
func (s *BaseExcellent1Listener) EnterFalse(ctx *FalseContext) {}

// ExitFalse is called when production false is exited.
func (s *BaseExcellent1Listener) ExitFalse(ctx *FalseContext) {}

// EnterArrayLookup is called when production arrayLookup is entered.
func (s *BaseExcellent1Listener) EnterArrayLookup(ctx *ArrayLookupContext) {}

// ExitArrayLookup is called when production arrayLookup is exited.
func (s *BaseExcellent1Listener) ExitArrayLookup(ctx *ArrayLookupContext) {}

// EnterContextReference is called when production contextReference is entered.
func (s *BaseExcellent1Listener) EnterContextReference(ctx *ContextReferenceContext) {}

// ExitContextReference is called when production contextReference is exited.
func (s *BaseExcellent1Listener) ExitContextReference(ctx *ContextReferenceContext) {}

// EnterParentheses is called when production parentheses is entered.
func (s *BaseExcellent1Listener) EnterParentheses(ctx *ParenthesesContext) {}

// ExitParentheses is called when production parentheses is exited.
func (s *BaseExcellent1Listener) ExitParentheses(ctx *ParenthesesContext) {}

// EnterNegation is called when production negation is entered.
func (s *BaseExcellent1Listener) EnterNegation(ctx *NegationContext) {}

// ExitNegation is called when production negation is exited.
func (s *BaseExcellent1Listener) ExitNegation(ctx *NegationContext) {}

// EnterComparison is called when production comparison is entered.
func (s *BaseExcellent1Listener) EnterComparison(ctx *ComparisonContext) {}

// ExitComparison is called when production comparison is exited.
func (s *BaseExcellent1Listener) ExitComparison(ctx *ComparisonContext) {}

// EnterConcatenation is called when production concatenation is entered.
func (s *BaseExcellent1Listener) EnterConcatenation(ctx *ConcatenationContext) {}

// ExitConcatenation is called when production concatenation is exited.
func (s *BaseExcellent1Listener) ExitConcatenation(ctx *ConcatenationContext) {}

// EnterMultiplicationOrDivision is called when production multiplicationOrDivision is entered.
func (s *BaseExcellent1Listener) EnterMultiplicationOrDivision(ctx *MultiplicationOrDivisionContext) {}

// ExitMultiplicationOrDivision is called when production multiplicationOrDivision is exited.
func (s *BaseExcellent1Listener) ExitMultiplicationOrDivision(ctx *MultiplicationOrDivisionContext) {}

// EnterAtomReference is called when production atomReference is entered.
func (s *BaseExcellent1Listener) EnterAtomReference(ctx *AtomReferenceContext) {}

// ExitAtomReference is called when production atomReference is exited.
func (s *BaseExcellent1Listener) ExitAtomReference(ctx *AtomReferenceContext) {}

// EnterAdditionOrSubtraction is called when production additionOrSubtraction is entered.
func (s *BaseExcellent1Listener) EnterAdditionOrSubtraction(ctx *AdditionOrSubtractionContext) {}

// ExitAdditionOrSubtraction is called when production additionOrSubtraction is exited.
func (s *BaseExcellent1Listener) ExitAdditionOrSubtraction(ctx *AdditionOrSubtractionContext) {}

// EnterEquality is called when production equality is entered.
func (s *BaseExcellent1Listener) EnterEquality(ctx *EqualityContext) {}

// ExitEquality is called when production equality is exited.
func (s *BaseExcellent1Listener) ExitEquality(ctx *EqualityContext) {}

// EnterExponent is called when production exponent is entered.
func (s *BaseExcellent1Listener) EnterExponent(ctx *ExponentContext) {}

// ExitExponent is called when production exponent is exited.
func (s *BaseExcellent1Listener) ExitExponent(ctx *ExponentContext) {}

// EnterFnname is called when production fnname is entered.
func (s *BaseExcellent1Listener) EnterFnname(ctx *FnnameContext) {}

// ExitFnname is called when production fnname is exited.
func (s *BaseExcellent1Listener) ExitFnname(ctx *FnnameContext) {}

// EnterFunctionParameters is called when production functionParameters is entered.
func (s *BaseExcellent1Listener) EnterFunctionParameters(ctx *FunctionParametersContext) {}

// ExitFunctionParameters is called when production functionParameters is exited.
func (s *BaseExcellent1Listener) ExitFunctionParameters(ctx *FunctionParametersContext) {}
