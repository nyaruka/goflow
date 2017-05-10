// Generated from src/github.com/nyaruka/goflow/excellent/gen/Excellent.g4 by ANTLR 4.7.

package gen // Excellent
import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseExcellentListener is a complete listener for a parse tree produced by ExcellentParser.
type BaseExcellentListener struct{}

var _ ExcellentListener = &BaseExcellentListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseExcellentListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseExcellentListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseExcellentListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseExcellentListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterParse is called when production parse is entered.
func (s *BaseExcellentListener) EnterParse(ctx *ParseContext) {}

// ExitParse is called when production parse is exited.
func (s *BaseExcellentListener) ExitParse(ctx *ParseContext) {}

// EnterDecimalLiteral is called when production decimalLiteral is entered.
func (s *BaseExcellentListener) EnterDecimalLiteral(ctx *DecimalLiteralContext) {}

// ExitDecimalLiteral is called when production decimalLiteral is exited.
func (s *BaseExcellentListener) ExitDecimalLiteral(ctx *DecimalLiteralContext) {}

// EnterDotLookup is called when production dotLookup is entered.
func (s *BaseExcellentListener) EnterDotLookup(ctx *DotLookupContext) {}

// ExitDotLookup is called when production dotLookup is exited.
func (s *BaseExcellentListener) ExitDotLookup(ctx *DotLookupContext) {}

// EnterStringLiteral is called when production stringLiteral is entered.
func (s *BaseExcellentListener) EnterStringLiteral(ctx *StringLiteralContext) {}

// ExitStringLiteral is called when production stringLiteral is exited.
func (s *BaseExcellentListener) ExitStringLiteral(ctx *StringLiteralContext) {}

// EnterFunctionCall is called when production functionCall is entered.
func (s *BaseExcellentListener) EnterFunctionCall(ctx *FunctionCallContext) {}

// ExitFunctionCall is called when production functionCall is exited.
func (s *BaseExcellentListener) ExitFunctionCall(ctx *FunctionCallContext) {}

// EnterTrue is called when production true is entered.
func (s *BaseExcellentListener) EnterTrue(ctx *TrueContext) {}

// ExitTrue is called when production true is exited.
func (s *BaseExcellentListener) ExitTrue(ctx *TrueContext) {}

// EnterFalse is called when production false is entered.
func (s *BaseExcellentListener) EnterFalse(ctx *FalseContext) {}

// ExitFalse is called when production false is exited.
func (s *BaseExcellentListener) ExitFalse(ctx *FalseContext) {}

// EnterArrayLookup is called when production arrayLookup is entered.
func (s *BaseExcellentListener) EnterArrayLookup(ctx *ArrayLookupContext) {}

// ExitArrayLookup is called when production arrayLookup is exited.
func (s *BaseExcellentListener) ExitArrayLookup(ctx *ArrayLookupContext) {}

// EnterContextReference is called when production contextReference is entered.
func (s *BaseExcellentListener) EnterContextReference(ctx *ContextReferenceContext) {}

// ExitContextReference is called when production contextReference is exited.
func (s *BaseExcellentListener) ExitContextReference(ctx *ContextReferenceContext) {}

// EnterParentheses is called when production parentheses is entered.
func (s *BaseExcellentListener) EnterParentheses(ctx *ParenthesesContext) {}

// ExitParentheses is called when production parentheses is exited.
func (s *BaseExcellentListener) ExitParentheses(ctx *ParenthesesContext) {}

// EnterNegation is called when production negation is entered.
func (s *BaseExcellentListener) EnterNegation(ctx *NegationContext) {}

// ExitNegation is called when production negation is exited.
func (s *BaseExcellentListener) ExitNegation(ctx *NegationContext) {}

// EnterComparison is called when production comparison is entered.
func (s *BaseExcellentListener) EnterComparison(ctx *ComparisonContext) {}

// ExitComparison is called when production comparison is exited.
func (s *BaseExcellentListener) ExitComparison(ctx *ComparisonContext) {}

// EnterConcatenation is called when production concatenation is entered.
func (s *BaseExcellentListener) EnterConcatenation(ctx *ConcatenationContext) {}

// ExitConcatenation is called when production concatenation is exited.
func (s *BaseExcellentListener) ExitConcatenation(ctx *ConcatenationContext) {}

// EnterMultiplicationOrDivision is called when production multiplicationOrDivision is entered.
func (s *BaseExcellentListener) EnterMultiplicationOrDivision(ctx *MultiplicationOrDivisionContext) {}

// ExitMultiplicationOrDivision is called when production multiplicationOrDivision is exited.
func (s *BaseExcellentListener) ExitMultiplicationOrDivision(ctx *MultiplicationOrDivisionContext) {}

// EnterAtomReference is called when production atomReference is entered.
func (s *BaseExcellentListener) EnterAtomReference(ctx *AtomReferenceContext) {}

// ExitAtomReference is called when production atomReference is exited.
func (s *BaseExcellentListener) ExitAtomReference(ctx *AtomReferenceContext) {}

// EnterAdditionOrSubtraction is called when production additionOrSubtraction is entered.
func (s *BaseExcellentListener) EnterAdditionOrSubtraction(ctx *AdditionOrSubtractionContext) {}

// ExitAdditionOrSubtraction is called when production additionOrSubtraction is exited.
func (s *BaseExcellentListener) ExitAdditionOrSubtraction(ctx *AdditionOrSubtractionContext) {}

// EnterEquality is called when production equality is entered.
func (s *BaseExcellentListener) EnterEquality(ctx *EqualityContext) {}

// ExitEquality is called when production equality is exited.
func (s *BaseExcellentListener) ExitEquality(ctx *EqualityContext) {}

// EnterExponent is called when production exponent is entered.
func (s *BaseExcellentListener) EnterExponent(ctx *ExponentContext) {}

// ExitExponent is called when production exponent is exited.
func (s *BaseExcellentListener) ExitExponent(ctx *ExponentContext) {}

// EnterFnname is called when production fnname is entered.
func (s *BaseExcellentListener) EnterFnname(ctx *FnnameContext) {}

// ExitFnname is called when production fnname is exited.
func (s *BaseExcellentListener) ExitFnname(ctx *FnnameContext) {}

// EnterFunctionParameters is called when production functionParameters is entered.
func (s *BaseExcellentListener) EnterFunctionParameters(ctx *FunctionParametersContext) {}

// ExitFunctionParameters is called when production functionParameters is exited.
func (s *BaseExcellentListener) ExitFunctionParameters(ctx *FunctionParametersContext) {}
