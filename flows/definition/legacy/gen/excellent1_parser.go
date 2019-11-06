// Code generated from Excellent1.g4 by ANTLR 4.7.2. DO NOT EDIT.

package gen // Excellent1
import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = reflect.Copy
var _ = strconv.Itoa

var parserATN = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 24, 68, 4,
	2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 3, 2, 3, 2, 3, 2, 3, 3, 3,
	3, 3, 3, 3, 3, 5, 3, 18, 10, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 5, 3, 33, 10, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 7, 3, 53, 10, 3, 12, 3, 14, 3, 56, 11, 3, 3, 4, 3,
	4, 3, 5, 3, 5, 3, 5, 7, 5, 63, 10, 5, 12, 5, 14, 5, 66, 11, 5, 3, 5, 2,
	3, 4, 6, 2, 4, 6, 8, 2, 7, 3, 2, 8, 9, 3, 2, 6, 7, 3, 2, 13, 16, 3, 2,
	11, 12, 3, 2, 20, 22, 2, 78, 2, 10, 3, 2, 2, 2, 4, 32, 3, 2, 2, 2, 6, 57,
	3, 2, 2, 2, 8, 59, 3, 2, 2, 2, 10, 11, 5, 4, 3, 2, 11, 12, 7, 2, 2, 3,
	12, 3, 3, 2, 2, 2, 13, 14, 8, 3, 1, 2, 14, 15, 5, 6, 4, 2, 15, 17, 7, 4,
	2, 2, 16, 18, 5, 8, 5, 2, 17, 16, 3, 2, 2, 2, 17, 18, 3, 2, 2, 2, 18, 19,
	3, 2, 2, 2, 19, 20, 7, 5, 2, 2, 20, 33, 3, 2, 2, 2, 21, 22, 7, 7, 2, 2,
	22, 33, 5, 4, 3, 15, 23, 33, 7, 19, 2, 2, 24, 33, 7, 18, 2, 2, 25, 33,
	7, 20, 2, 2, 26, 33, 7, 21, 2, 2, 27, 33, 7, 22, 2, 2, 28, 29, 7, 4, 2,
	2, 29, 30, 5, 4, 3, 2, 30, 31, 7, 5, 2, 2, 31, 33, 3, 2, 2, 2, 32, 13,
	3, 2, 2, 2, 32, 21, 3, 2, 2, 2, 32, 23, 3, 2, 2, 2, 32, 24, 3, 2, 2, 2,
	32, 25, 3, 2, 2, 2, 32, 26, 3, 2, 2, 2, 32, 27, 3, 2, 2, 2, 32, 28, 3,
	2, 2, 2, 33, 54, 3, 2, 2, 2, 34, 35, 12, 14, 2, 2, 35, 36, 7, 10, 2, 2,
	36, 53, 5, 4, 3, 15, 37, 38, 12, 13, 2, 2, 38, 39, 9, 2, 2, 2, 39, 53,
	5, 4, 3, 14, 40, 41, 12, 12, 2, 2, 41, 42, 9, 3, 2, 2, 42, 53, 5, 4, 3,
	13, 43, 44, 12, 11, 2, 2, 44, 45, 9, 4, 2, 2, 45, 53, 5, 4, 3, 12, 46,
	47, 12, 10, 2, 2, 47, 48, 9, 5, 2, 2, 48, 53, 5, 4, 3, 11, 49, 50, 12,
	9, 2, 2, 50, 51, 7, 17, 2, 2, 51, 53, 5, 4, 3, 10, 52, 34, 3, 2, 2, 2,
	52, 37, 3, 2, 2, 2, 52, 40, 3, 2, 2, 2, 52, 43, 3, 2, 2, 2, 52, 46, 3,
	2, 2, 2, 52, 49, 3, 2, 2, 2, 53, 56, 3, 2, 2, 2, 54, 52, 3, 2, 2, 2, 54,
	55, 3, 2, 2, 2, 55, 5, 3, 2, 2, 2, 56, 54, 3, 2, 2, 2, 57, 58, 9, 6, 2,
	2, 58, 7, 3, 2, 2, 2, 59, 64, 5, 4, 3, 2, 60, 61, 7, 3, 2, 2, 61, 63, 5,
	4, 3, 2, 62, 60, 3, 2, 2, 2, 63, 66, 3, 2, 2, 2, 64, 62, 3, 2, 2, 2, 64,
	65, 3, 2, 2, 2, 65, 9, 3, 2, 2, 2, 66, 64, 3, 2, 2, 2, 7, 17, 32, 52, 54,
	64,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "','", "'('", "')'", "'+'", "'-'", "'*'", "'/'", "'^'", "'='", "'<>'",
	"'<='", "'<'", "'>='", "'>'", "'&'",
}
var symbolicNames = []string{
	"", "COMMA", "LPAREN", "RPAREN", "PLUS", "MINUS", "TIMES", "DIVIDE", "EXPONENT",
	"EQ", "NEQ", "LTE", "LT", "GTE", "GT", "AMPERSAND", "DECIMAL", "STRING",
	"TRUE", "FALSE", "NAME", "WS", "ERROR",
}

var ruleNames = []string{
	"parse", "expression", "fnname", "parameters",
}
var decisionToDFA = make([]*antlr.DFA, len(deserializedATN.DecisionToState))

func init() {
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

type Excellent1Parser struct {
	*antlr.BaseParser
}

func NewExcellent1Parser(input antlr.TokenStream) *Excellent1Parser {
	this := new(Excellent1Parser)

	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "Excellent1.g4"

	return this
}

// Excellent1Parser tokens.
const (
	Excellent1ParserEOF       = antlr.TokenEOF
	Excellent1ParserCOMMA     = 1
	Excellent1ParserLPAREN    = 2
	Excellent1ParserRPAREN    = 3
	Excellent1ParserPLUS      = 4
	Excellent1ParserMINUS     = 5
	Excellent1ParserTIMES     = 6
	Excellent1ParserDIVIDE    = 7
	Excellent1ParserEXPONENT  = 8
	Excellent1ParserEQ        = 9
	Excellent1ParserNEQ       = 10
	Excellent1ParserLTE       = 11
	Excellent1ParserLT        = 12
	Excellent1ParserGTE       = 13
	Excellent1ParserGT        = 14
	Excellent1ParserAMPERSAND = 15
	Excellent1ParserDECIMAL   = 16
	Excellent1ParserSTRING    = 17
	Excellent1ParserTRUE      = 18
	Excellent1ParserFALSE     = 19
	Excellent1ParserNAME      = 20
	Excellent1ParserWS        = 21
	Excellent1ParserERROR     = 22
)

// Excellent1Parser rules.
const (
	Excellent1ParserRULE_parse      = 0
	Excellent1ParserRULE_expression = 1
	Excellent1ParserRULE_fnname     = 2
	Excellent1ParserRULE_parameters = 3
)

// IParseContext is an interface to support dynamic dispatch.
type IParseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsParseContext differentiates from other interfaces.
	IsParseContext()
}

type ParseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyParseContext() *ParseContext {
	var p = new(ParseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = Excellent1ParserRULE_parse
	return p
}

func (*ParseContext) IsParseContext() {}

func NewParseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParseContext {
	var p = new(ParseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = Excellent1ParserRULE_parse

	return p
}

func (s *ParseContext) GetParser() antlr.Parser { return s.parser }

func (s *ParseContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ParseContext) EOF() antlr.TerminalNode {
	return s.GetToken(Excellent1ParserEOF, 0)
}

func (s *ParseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ParseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.EnterParse(s)
	}
}

func (s *ParseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.ExitParse(s)
	}
}

func (s *ParseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent1Visitor:
		return t.VisitParse(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *Excellent1Parser) Parse() (localctx IParseContext) {
	localctx = NewParseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, Excellent1ParserRULE_parse)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(8)
		p.expression(0)
	}
	{
		p.SetState(9)
		p.Match(Excellent1ParserEOF)
	}

	return localctx
}

// IExpressionContext is an interface to support dynamic dispatch.
type IExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsExpressionContext differentiates from other interfaces.
	IsExpressionContext()
}

type ExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpressionContext() *ExpressionContext {
	var p = new(ExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = Excellent1ParserRULE_expression
	return p
}

func (*ExpressionContext) IsExpressionContext() {}

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext {
	var p = new(ExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = Excellent1ParserRULE_expression

	return p
}

func (s *ExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionContext) CopyFrom(ctx *ExpressionContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *ExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type DecimalLiteralContext struct {
	*ExpressionContext
}

func NewDecimalLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DecimalLiteralContext {
	var p = new(DecimalLiteralContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *DecimalLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DecimalLiteralContext) DECIMAL() antlr.TerminalNode {
	return s.GetToken(Excellent1ParserDECIMAL, 0)
}

func (s *DecimalLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.EnterDecimalLiteral(s)
	}
}

func (s *DecimalLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.ExitDecimalLiteral(s)
	}
}

func (s *DecimalLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent1Visitor:
		return t.VisitDecimalLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

type ParenthesesContext struct {
	*ExpressionContext
}

func NewParenthesesContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ParenthesesContext {
	var p = new(ParenthesesContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *ParenthesesContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParenthesesContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(Excellent1ParserLPAREN, 0)
}

func (s *ParenthesesContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ParenthesesContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(Excellent1ParserRPAREN, 0)
}

func (s *ParenthesesContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.EnterParentheses(s)
	}
}

func (s *ParenthesesContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.ExitParentheses(s)
	}
}

func (s *ParenthesesContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent1Visitor:
		return t.VisitParentheses(s)

	default:
		return t.VisitChildren(s)
	}
}

type NegationContext struct {
	*ExpressionContext
}

func NewNegationContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NegationContext {
	var p = new(NegationContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *NegationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NegationContext) MINUS() antlr.TerminalNode {
	return s.GetToken(Excellent1ParserMINUS, 0)
}

func (s *NegationContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *NegationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.EnterNegation(s)
	}
}

func (s *NegationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.ExitNegation(s)
	}
}

func (s *NegationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent1Visitor:
		return t.VisitNegation(s)

	default:
		return t.VisitChildren(s)
	}
}

type ExponentExpressionContext struct {
	*ExpressionContext
}

func NewExponentExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ExponentExpressionContext {
	var p = new(ExponentExpressionContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *ExponentExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExponentExpressionContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *ExponentExpressionContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ExponentExpressionContext) EXPONENT() antlr.TerminalNode {
	return s.GetToken(Excellent1ParserEXPONENT, 0)
}

func (s *ExponentExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.EnterExponentExpression(s)
	}
}

func (s *ExponentExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.ExitExponentExpression(s)
	}
}

func (s *ExponentExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent1Visitor:
		return t.VisitExponentExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type AdditionOrSubtractionExpressionContext struct {
	*ExpressionContext
	op antlr.Token
}

func NewAdditionOrSubtractionExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *AdditionOrSubtractionExpressionContext {
	var p = new(AdditionOrSubtractionExpressionContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *AdditionOrSubtractionExpressionContext) GetOp() antlr.Token { return s.op }

func (s *AdditionOrSubtractionExpressionContext) SetOp(v antlr.Token) { s.op = v }

func (s *AdditionOrSubtractionExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AdditionOrSubtractionExpressionContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *AdditionOrSubtractionExpressionContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *AdditionOrSubtractionExpressionContext) PLUS() antlr.TerminalNode {
	return s.GetToken(Excellent1ParserPLUS, 0)
}

func (s *AdditionOrSubtractionExpressionContext) MINUS() antlr.TerminalNode {
	return s.GetToken(Excellent1ParserMINUS, 0)
}

func (s *AdditionOrSubtractionExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.EnterAdditionOrSubtractionExpression(s)
	}
}

func (s *AdditionOrSubtractionExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.ExitAdditionOrSubtractionExpression(s)
	}
}

func (s *AdditionOrSubtractionExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent1Visitor:
		return t.VisitAdditionOrSubtractionExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type FalseContext struct {
	*ExpressionContext
}

func NewFalseContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FalseContext {
	var p = new(FalseContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *FalseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FalseContext) FALSE() antlr.TerminalNode {
	return s.GetToken(Excellent1ParserFALSE, 0)
}

func (s *FalseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.EnterFalse(s)
	}
}

func (s *FalseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.ExitFalse(s)
	}
}

func (s *FalseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent1Visitor:
		return t.VisitFalse(s)

	default:
		return t.VisitChildren(s)
	}
}

type ContextReferenceContext struct {
	*ExpressionContext
}

func NewContextReferenceContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ContextReferenceContext {
	var p = new(ContextReferenceContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *ContextReferenceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ContextReferenceContext) NAME() antlr.TerminalNode {
	return s.GetToken(Excellent1ParserNAME, 0)
}

func (s *ContextReferenceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.EnterContextReference(s)
	}
}

func (s *ContextReferenceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.ExitContextReference(s)
	}
}

func (s *ContextReferenceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent1Visitor:
		return t.VisitContextReference(s)

	default:
		return t.VisitChildren(s)
	}
}

type ComparisonExpressionContext struct {
	*ExpressionContext
	op antlr.Token
}

func NewComparisonExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ComparisonExpressionContext {
	var p = new(ComparisonExpressionContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *ComparisonExpressionContext) GetOp() antlr.Token { return s.op }

func (s *ComparisonExpressionContext) SetOp(v antlr.Token) { s.op = v }

func (s *ComparisonExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ComparisonExpressionContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *ComparisonExpressionContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ComparisonExpressionContext) LTE() antlr.TerminalNode {
	return s.GetToken(Excellent1ParserLTE, 0)
}

func (s *ComparisonExpressionContext) LT() antlr.TerminalNode {
	return s.GetToken(Excellent1ParserLT, 0)
}

func (s *ComparisonExpressionContext) GTE() antlr.TerminalNode {
	return s.GetToken(Excellent1ParserGTE, 0)
}

func (s *ComparisonExpressionContext) GT() antlr.TerminalNode {
	return s.GetToken(Excellent1ParserGT, 0)
}

func (s *ComparisonExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.EnterComparisonExpression(s)
	}
}

func (s *ComparisonExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.ExitComparisonExpression(s)
	}
}

func (s *ComparisonExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent1Visitor:
		return t.VisitComparisonExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type ConcatenationContext struct {
	*ExpressionContext
}

func NewConcatenationContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ConcatenationContext {
	var p = new(ConcatenationContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *ConcatenationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConcatenationContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *ConcatenationContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ConcatenationContext) AMPERSAND() antlr.TerminalNode {
	return s.GetToken(Excellent1ParserAMPERSAND, 0)
}

func (s *ConcatenationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.EnterConcatenation(s)
	}
}

func (s *ConcatenationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.ExitConcatenation(s)
	}
}

func (s *ConcatenationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent1Visitor:
		return t.VisitConcatenation(s)

	default:
		return t.VisitChildren(s)
	}
}

type StringLiteralContext struct {
	*ExpressionContext
}

func NewStringLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *StringLiteralContext {
	var p = new(StringLiteralContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *StringLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StringLiteralContext) STRING() antlr.TerminalNode {
	return s.GetToken(Excellent1ParserSTRING, 0)
}

func (s *StringLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.EnterStringLiteral(s)
	}
}

func (s *StringLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.ExitStringLiteral(s)
	}
}

func (s *StringLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent1Visitor:
		return t.VisitStringLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

type FunctionCallContext struct {
	*ExpressionContext
}

func NewFunctionCallContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FunctionCallContext {
	var p = new(FunctionCallContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *FunctionCallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionCallContext) Fnname() IFnnameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFnnameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFnnameContext)
}

func (s *FunctionCallContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(Excellent1ParserLPAREN, 0)
}

func (s *FunctionCallContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(Excellent1ParserRPAREN, 0)
}

func (s *FunctionCallContext) Parameters() IParametersContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IParametersContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IParametersContext)
}

func (s *FunctionCallContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.EnterFunctionCall(s)
	}
}

func (s *FunctionCallContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.ExitFunctionCall(s)
	}
}

func (s *FunctionCallContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent1Visitor:
		return t.VisitFunctionCall(s)

	default:
		return t.VisitChildren(s)
	}
}

type TrueContext struct {
	*ExpressionContext
}

func NewTrueContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TrueContext {
	var p = new(TrueContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *TrueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TrueContext) TRUE() antlr.TerminalNode {
	return s.GetToken(Excellent1ParserTRUE, 0)
}

func (s *TrueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.EnterTrue(s)
	}
}

func (s *TrueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.ExitTrue(s)
	}
}

func (s *TrueContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent1Visitor:
		return t.VisitTrue(s)

	default:
		return t.VisitChildren(s)
	}
}

type EqualityExpressionContext struct {
	*ExpressionContext
	op antlr.Token
}

func NewEqualityExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *EqualityExpressionContext {
	var p = new(EqualityExpressionContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *EqualityExpressionContext) GetOp() antlr.Token { return s.op }

func (s *EqualityExpressionContext) SetOp(v antlr.Token) { s.op = v }

func (s *EqualityExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EqualityExpressionContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *EqualityExpressionContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *EqualityExpressionContext) EQ() antlr.TerminalNode {
	return s.GetToken(Excellent1ParserEQ, 0)
}

func (s *EqualityExpressionContext) NEQ() antlr.TerminalNode {
	return s.GetToken(Excellent1ParserNEQ, 0)
}

func (s *EqualityExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.EnterEqualityExpression(s)
	}
}

func (s *EqualityExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.ExitEqualityExpression(s)
	}
}

func (s *EqualityExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent1Visitor:
		return t.VisitEqualityExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type MultiplicationOrDivisionExpressionContext struct {
	*ExpressionContext
	op antlr.Token
}

func NewMultiplicationOrDivisionExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MultiplicationOrDivisionExpressionContext {
	var p = new(MultiplicationOrDivisionExpressionContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *MultiplicationOrDivisionExpressionContext) GetOp() antlr.Token { return s.op }

func (s *MultiplicationOrDivisionExpressionContext) SetOp(v antlr.Token) { s.op = v }

func (s *MultiplicationOrDivisionExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MultiplicationOrDivisionExpressionContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *MultiplicationOrDivisionExpressionContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *MultiplicationOrDivisionExpressionContext) TIMES() antlr.TerminalNode {
	return s.GetToken(Excellent1ParserTIMES, 0)
}

func (s *MultiplicationOrDivisionExpressionContext) DIVIDE() antlr.TerminalNode {
	return s.GetToken(Excellent1ParserDIVIDE, 0)
}

func (s *MultiplicationOrDivisionExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.EnterMultiplicationOrDivisionExpression(s)
	}
}

func (s *MultiplicationOrDivisionExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.ExitMultiplicationOrDivisionExpression(s)
	}
}

func (s *MultiplicationOrDivisionExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent1Visitor:
		return t.VisitMultiplicationOrDivisionExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *Excellent1Parser) Expression() (localctx IExpressionContext) {
	return p.expression(0)
}

func (p *Excellent1Parser) expression(_p int) (localctx IExpressionContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewExpressionContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IExpressionContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 2
	p.EnterRecursionRule(localctx, 2, Excellent1ParserRULE_expression, _p)
	var _la int

	defer func() {
		p.UnrollRecursionContexts(_parentctx)
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(30)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 1, p.GetParserRuleContext()) {
	case 1:
		localctx = NewFunctionCallContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx

		{
			p.SetState(12)
			p.Fnname()
		}
		{
			p.SetState(13)
			p.Match(Excellent1ParserLPAREN)
		}
		p.SetState(15)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if ((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<Excellent1ParserLPAREN)|(1<<Excellent1ParserMINUS)|(1<<Excellent1ParserDECIMAL)|(1<<Excellent1ParserSTRING)|(1<<Excellent1ParserTRUE)|(1<<Excellent1ParserFALSE)|(1<<Excellent1ParserNAME))) != 0 {
			{
				p.SetState(14)
				p.Parameters()
			}

		}
		{
			p.SetState(17)
			p.Match(Excellent1ParserRPAREN)
		}

	case 2:
		localctx = NewNegationContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(19)
			p.Match(Excellent1ParserMINUS)
		}
		{
			p.SetState(20)
			p.expression(13)
		}

	case 3:
		localctx = NewStringLiteralContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(21)
			p.Match(Excellent1ParserSTRING)
		}

	case 4:
		localctx = NewDecimalLiteralContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(22)
			p.Match(Excellent1ParserDECIMAL)
		}

	case 5:
		localctx = NewTrueContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(23)
			p.Match(Excellent1ParserTRUE)
		}

	case 6:
		localctx = NewFalseContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(24)
			p.Match(Excellent1ParserFALSE)
		}

	case 7:
		localctx = NewContextReferenceContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(25)
			p.Match(Excellent1ParserNAME)
		}

	case 8:
		localctx = NewParenthesesContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(26)
			p.Match(Excellent1ParserLPAREN)
		}
		{
			p.SetState(27)
			p.expression(0)
		}
		{
			p.SetState(28)
			p.Match(Excellent1ParserRPAREN)
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(52)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 3, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(50)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 2, p.GetParserRuleContext()) {
			case 1:
				localctx = NewExponentExpressionContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, Excellent1ParserRULE_expression)
				p.SetState(32)

				if !(p.Precpred(p.GetParserRuleContext(), 12)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 12)", ""))
				}
				{
					p.SetState(33)
					p.Match(Excellent1ParserEXPONENT)
				}
				{
					p.SetState(34)
					p.expression(13)
				}

			case 2:
				localctx = NewMultiplicationOrDivisionExpressionContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, Excellent1ParserRULE_expression)
				p.SetState(35)

				if !(p.Precpred(p.GetParserRuleContext(), 11)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 11)", ""))
				}
				{
					p.SetState(36)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*MultiplicationOrDivisionExpressionContext).op = _lt

					_la = p.GetTokenStream().LA(1)

					if !(_la == Excellent1ParserTIMES || _la == Excellent1ParserDIVIDE) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*MultiplicationOrDivisionExpressionContext).op = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(37)
					p.expression(12)
				}

			case 3:
				localctx = NewAdditionOrSubtractionExpressionContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, Excellent1ParserRULE_expression)
				p.SetState(38)

				if !(p.Precpred(p.GetParserRuleContext(), 10)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 10)", ""))
				}
				{
					p.SetState(39)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*AdditionOrSubtractionExpressionContext).op = _lt

					_la = p.GetTokenStream().LA(1)

					if !(_la == Excellent1ParserPLUS || _la == Excellent1ParserMINUS) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*AdditionOrSubtractionExpressionContext).op = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(40)
					p.expression(11)
				}

			case 4:
				localctx = NewComparisonExpressionContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, Excellent1ParserRULE_expression)
				p.SetState(41)

				if !(p.Precpred(p.GetParserRuleContext(), 9)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 9)", ""))
				}
				{
					p.SetState(42)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*ComparisonExpressionContext).op = _lt

					_la = p.GetTokenStream().LA(1)

					if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<Excellent1ParserLTE)|(1<<Excellent1ParserLT)|(1<<Excellent1ParserGTE)|(1<<Excellent1ParserGT))) != 0) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*ComparisonExpressionContext).op = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(43)
					p.expression(10)
				}

			case 5:
				localctx = NewEqualityExpressionContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, Excellent1ParserRULE_expression)
				p.SetState(44)

				if !(p.Precpred(p.GetParserRuleContext(), 8)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 8)", ""))
				}
				{
					p.SetState(45)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*EqualityExpressionContext).op = _lt

					_la = p.GetTokenStream().LA(1)

					if !(_la == Excellent1ParserEQ || _la == Excellent1ParserNEQ) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*EqualityExpressionContext).op = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(46)
					p.expression(9)
				}

			case 6:
				localctx = NewConcatenationContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, Excellent1ParserRULE_expression)
				p.SetState(47)

				if !(p.Precpred(p.GetParserRuleContext(), 7)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 7)", ""))
				}
				{
					p.SetState(48)
					p.Match(Excellent1ParserAMPERSAND)
				}
				{
					p.SetState(49)
					p.expression(8)
				}

			}

		}
		p.SetState(54)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 3, p.GetParserRuleContext())
	}

	return localctx
}

// IFnnameContext is an interface to support dynamic dispatch.
type IFnnameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFnnameContext differentiates from other interfaces.
	IsFnnameContext()
}

type FnnameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFnnameContext() *FnnameContext {
	var p = new(FnnameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = Excellent1ParserRULE_fnname
	return p
}

func (*FnnameContext) IsFnnameContext() {}

func NewFnnameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FnnameContext {
	var p = new(FnnameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = Excellent1ParserRULE_fnname

	return p
}

func (s *FnnameContext) GetParser() antlr.Parser { return s.parser }

func (s *FnnameContext) NAME() antlr.TerminalNode {
	return s.GetToken(Excellent1ParserNAME, 0)
}

func (s *FnnameContext) TRUE() antlr.TerminalNode {
	return s.GetToken(Excellent1ParserTRUE, 0)
}

func (s *FnnameContext) FALSE() antlr.TerminalNode {
	return s.GetToken(Excellent1ParserFALSE, 0)
}

func (s *FnnameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FnnameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FnnameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.EnterFnname(s)
	}
}

func (s *FnnameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.ExitFnname(s)
	}
}

func (s *FnnameContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent1Visitor:
		return t.VisitFnname(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *Excellent1Parser) Fnname() (localctx IFnnameContext) {
	localctx = NewFnnameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, Excellent1ParserRULE_fnname)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(55)
		_la = p.GetTokenStream().LA(1)

		if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<Excellent1ParserTRUE)|(1<<Excellent1ParserFALSE)|(1<<Excellent1ParserNAME))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IParametersContext is an interface to support dynamic dispatch.
type IParametersContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsParametersContext differentiates from other interfaces.
	IsParametersContext()
}

type ParametersContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyParametersContext() *ParametersContext {
	var p = new(ParametersContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = Excellent1ParserRULE_parameters
	return p
}

func (*ParametersContext) IsParametersContext() {}

func NewParametersContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParametersContext {
	var p = new(ParametersContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = Excellent1ParserRULE_parameters

	return p
}

func (s *ParametersContext) GetParser() antlr.Parser { return s.parser }

func (s *ParametersContext) CopyFrom(ctx *ParametersContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *ParametersContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParametersContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type FunctionParametersContext struct {
	*ParametersContext
}

func NewFunctionParametersContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FunctionParametersContext {
	var p = new(FunctionParametersContext)

	p.ParametersContext = NewEmptyParametersContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ParametersContext))

	return p
}

func (s *FunctionParametersContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionParametersContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *FunctionParametersContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *FunctionParametersContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(Excellent1ParserCOMMA)
}

func (s *FunctionParametersContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(Excellent1ParserCOMMA, i)
}

func (s *FunctionParametersContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.EnterFunctionParameters(s)
	}
}

func (s *FunctionParametersContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent1Listener); ok {
		listenerT.ExitFunctionParameters(s)
	}
}

func (s *FunctionParametersContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent1Visitor:
		return t.VisitFunctionParameters(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *Excellent1Parser) Parameters() (localctx IParametersContext) {
	localctx = NewParametersContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, Excellent1ParserRULE_parameters)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	localctx = NewFunctionParametersContext(p, localctx)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(57)
		p.expression(0)
	}
	p.SetState(62)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == Excellent1ParserCOMMA {
		{
			p.SetState(58)
			p.Match(Excellent1ParserCOMMA)
		}
		{
			p.SetState(59)
			p.expression(0)
		}

		p.SetState(64)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

func (p *Excellent1Parser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 1:
		var t *ExpressionContext = nil
		if localctx != nil {
			t = localctx.(*ExpressionContext)
		}
		return p.Expression_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *Excellent1Parser) Expression_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 12)

	case 1:
		return p.Precpred(p.GetParserRuleContext(), 11)

	case 2:
		return p.Precpred(p.GetParserRuleContext(), 10)

	case 3:
		return p.Precpred(p.GetParserRuleContext(), 9)

	case 4:
		return p.Precpred(p.GetParserRuleContext(), 8)

	case 5:
		return p.Precpred(p.GetParserRuleContext(), 7)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
