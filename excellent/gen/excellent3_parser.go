// Code generated from Excellent3.g4 by ANTLR 4.9.2. DO NOT EDIT.

package gen // Excellent3
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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 30, 99, 4,
	2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 3, 2, 3, 2, 3,
	2, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 5, 3, 31, 10, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 7, 3, 51, 10, 3, 12, 3, 14, 3, 54, 11, 3, 3, 4, 3, 4, 3, 4, 3, 4, 3,
	4, 3, 4, 5, 4, 62, 10, 4, 3, 4, 3, 4, 3, 4, 5, 4, 67, 10, 4, 3, 4, 3, 4,
	3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 7, 4, 78, 10, 4, 12, 4, 14, 4,
	81, 11, 4, 3, 5, 3, 5, 3, 5, 7, 5, 86, 10, 5, 12, 5, 14, 5, 89, 11, 5,
	3, 6, 3, 6, 3, 6, 7, 6, 94, 10, 6, 12, 6, 14, 6, 97, 11, 6, 3, 6, 2, 4,
	4, 6, 7, 2, 4, 6, 8, 10, 2, 8, 3, 2, 23, 24, 3, 2, 12, 13, 3, 2, 10, 11,
	3, 2, 17, 20, 3, 2, 15, 16, 4, 2, 23, 23, 28, 28, 2, 113, 2, 12, 3, 2,
	2, 2, 4, 30, 3, 2, 2, 2, 6, 61, 3, 2, 2, 2, 8, 82, 3, 2, 2, 2, 10, 90,
	3, 2, 2, 2, 12, 13, 5, 4, 3, 2, 13, 14, 7, 2, 2, 3, 14, 3, 3, 2, 2, 2,
	15, 16, 8, 3, 1, 2, 16, 31, 5, 6, 4, 2, 17, 18, 7, 11, 2, 2, 18, 31, 5,
	4, 3, 15, 19, 20, 7, 4, 2, 2, 20, 21, 5, 10, 6, 2, 21, 22, 7, 5, 2, 2,
	22, 23, 7, 9, 2, 2, 23, 24, 5, 4, 3, 8, 24, 31, 3, 2, 2, 2, 25, 31, 7,
	22, 2, 2, 26, 31, 9, 2, 2, 2, 27, 31, 7, 25, 2, 2, 28, 31, 7, 26, 2, 2,
	29, 31, 7, 27, 2, 2, 30, 15, 3, 2, 2, 2, 30, 17, 3, 2, 2, 2, 30, 19, 3,
	2, 2, 2, 30, 25, 3, 2, 2, 2, 30, 26, 3, 2, 2, 2, 30, 27, 3, 2, 2, 2, 30,
	28, 3, 2, 2, 2, 30, 29, 3, 2, 2, 2, 31, 52, 3, 2, 2, 2, 32, 33, 12, 14,
	2, 2, 33, 34, 7, 14, 2, 2, 34, 51, 5, 4, 3, 15, 35, 36, 12, 13, 2, 2, 36,
	37, 9, 3, 2, 2, 37, 51, 5, 4, 3, 14, 38, 39, 12, 12, 2, 2, 39, 40, 9, 4,
	2, 2, 40, 51, 5, 4, 3, 13, 41, 42, 12, 11, 2, 2, 42, 43, 9, 5, 2, 2, 43,
	51, 5, 4, 3, 12, 44, 45, 12, 10, 2, 2, 45, 46, 9, 6, 2, 2, 46, 51, 5, 4,
	3, 11, 47, 48, 12, 9, 2, 2, 48, 49, 7, 21, 2, 2, 49, 51, 5, 4, 3, 10, 50,
	32, 3, 2, 2, 2, 50, 35, 3, 2, 2, 2, 50, 38, 3, 2, 2, 2, 50, 41, 3, 2, 2,
	2, 50, 44, 3, 2, 2, 2, 50, 47, 3, 2, 2, 2, 51, 54, 3, 2, 2, 2, 52, 50,
	3, 2, 2, 2, 52, 53, 3, 2, 2, 2, 53, 5, 3, 2, 2, 2, 54, 52, 3, 2, 2, 2,
	55, 56, 8, 4, 1, 2, 56, 57, 7, 4, 2, 2, 57, 58, 5, 4, 3, 2, 58, 59, 7,
	5, 2, 2, 59, 62, 3, 2, 2, 2, 60, 62, 7, 28, 2, 2, 61, 55, 3, 2, 2, 2, 61,
	60, 3, 2, 2, 2, 62, 79, 3, 2, 2, 2, 63, 64, 12, 7, 2, 2, 64, 66, 7, 4,
	2, 2, 65, 67, 5, 8, 5, 2, 66, 65, 3, 2, 2, 2, 66, 67, 3, 2, 2, 2, 67, 68,
	3, 2, 2, 2, 68, 78, 7, 5, 2, 2, 69, 70, 12, 6, 2, 2, 70, 71, 7, 8, 2, 2,
	71, 78, 9, 7, 2, 2, 72, 73, 12, 5, 2, 2, 73, 74, 7, 6, 2, 2, 74, 75, 5,
	4, 3, 2, 75, 76, 7, 7, 2, 2, 76, 78, 3, 2, 2, 2, 77, 63, 3, 2, 2, 2, 77,
	69, 3, 2, 2, 2, 77, 72, 3, 2, 2, 2, 78, 81, 3, 2, 2, 2, 79, 77, 3, 2, 2,
	2, 79, 80, 3, 2, 2, 2, 80, 7, 3, 2, 2, 2, 81, 79, 3, 2, 2, 2, 82, 87, 5,
	4, 3, 2, 83, 84, 7, 3, 2, 2, 84, 86, 5, 4, 3, 2, 85, 83, 3, 2, 2, 2, 86,
	89, 3, 2, 2, 2, 87, 85, 3, 2, 2, 2, 87, 88, 3, 2, 2, 2, 88, 9, 3, 2, 2,
	2, 89, 87, 3, 2, 2, 2, 90, 95, 7, 28, 2, 2, 91, 92, 7, 3, 2, 2, 92, 94,
	7, 28, 2, 2, 93, 91, 3, 2, 2, 2, 94, 97, 3, 2, 2, 2, 95, 93, 3, 2, 2, 2,
	95, 96, 3, 2, 2, 2, 96, 11, 3, 2, 2, 2, 97, 95, 3, 2, 2, 2, 11, 30, 50,
	52, 61, 66, 77, 79, 87, 95,
}
var literalNames = []string{
	"", "','", "'('", "')'", "'['", "']'", "'.'", "'=>'", "'+'", "'-'", "'*'",
	"'/'", "'^'", "'='", "'!='", "'<='", "'<'", "'>='", "'>'", "'&'",
}
var symbolicNames = []string{
	"", "COMMA", "LPAREN", "RPAREN", "LBRACK", "RBRACK", "DOT", "ARROW", "PLUS",
	"MINUS", "TIMES", "DIVIDE", "EXPONENT", "EQ", "NEQ", "LTE", "LT", "GTE",
	"GT", "AMPERSAND", "TEXT", "INTEGER", "DECIMAL", "TRUE", "FALSE", "NULL",
	"NAME", "WS", "ERROR",
}

var ruleNames = []string{
	"parse", "expression", "atom", "parameters", "nameList",
}

type Excellent3Parser struct {
	*antlr.BaseParser
}

// NewExcellent3Parser produces a new parser instance for the optional input antlr.TokenStream.
//
// The *Excellent3Parser instance produced may be reused by calling the SetInputStream method.
// The initial parser configuration is expensive to construct, and the object is not thread-safe;
// however, if used within a Golang sync.Pool, the construction cost amortizes well and the
// objects can be used in a thread-safe manner.
func NewExcellent3Parser(input antlr.TokenStream) *Excellent3Parser {
	this := new(Excellent3Parser)
	deserializer := antlr.NewATNDeserializer(nil)
	deserializedATN := deserializer.DeserializeFromUInt16(parserATN)
	decisionToDFA := make([]*antlr.DFA, len(deserializedATN.DecisionToState))
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "Excellent3.g4"

	return this
}

// Excellent3Parser tokens.
const (
	Excellent3ParserEOF       = antlr.TokenEOF
	Excellent3ParserCOMMA     = 1
	Excellent3ParserLPAREN    = 2
	Excellent3ParserRPAREN    = 3
	Excellent3ParserLBRACK    = 4
	Excellent3ParserRBRACK    = 5
	Excellent3ParserDOT       = 6
	Excellent3ParserARROW     = 7
	Excellent3ParserPLUS      = 8
	Excellent3ParserMINUS     = 9
	Excellent3ParserTIMES     = 10
	Excellent3ParserDIVIDE    = 11
	Excellent3ParserEXPONENT  = 12
	Excellent3ParserEQ        = 13
	Excellent3ParserNEQ       = 14
	Excellent3ParserLTE       = 15
	Excellent3ParserLT        = 16
	Excellent3ParserGTE       = 17
	Excellent3ParserGT        = 18
	Excellent3ParserAMPERSAND = 19
	Excellent3ParserTEXT      = 20
	Excellent3ParserINTEGER   = 21
	Excellent3ParserDECIMAL   = 22
	Excellent3ParserTRUE      = 23
	Excellent3ParserFALSE     = 24
	Excellent3ParserNULL      = 25
	Excellent3ParserNAME      = 26
	Excellent3ParserWS        = 27
	Excellent3ParserERROR     = 28
)

// Excellent3Parser rules.
const (
	Excellent3ParserRULE_parse      = 0
	Excellent3ParserRULE_expression = 1
	Excellent3ParserRULE_atom       = 2
	Excellent3ParserRULE_parameters = 3
	Excellent3ParserRULE_nameList   = 4
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
	p.RuleIndex = Excellent3ParserRULE_parse
	return p
}

func (*ParseContext) IsParseContext() {}

func NewParseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParseContext {
	var p = new(ParseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = Excellent3ParserRULE_parse

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
	return s.GetToken(Excellent3ParserEOF, 0)
}

func (s *ParseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ParseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.EnterParse(s)
	}
}

func (s *ParseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.ExitParse(s)
	}
}

func (s *ParseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent3Visitor:
		return t.VisitParse(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *Excellent3Parser) Parse() (localctx IParseContext) {
	localctx = NewParseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, Excellent3ParserRULE_parse)

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
		p.SetState(10)
		p.expression(0)
	}
	{
		p.SetState(11)
		p.Match(Excellent3ParserEOF)
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
	p.RuleIndex = Excellent3ParserRULE_expression
	return p
}

func (*ExpressionContext) IsExpressionContext() {}

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext {
	var p = new(ExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = Excellent3ParserRULE_expression

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
	return s.GetToken(Excellent3ParserMINUS, 0)
}

func (s *NegationContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *NegationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.EnterNegation(s)
	}
}

func (s *NegationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.ExitNegation(s)
	}
}

func (s *NegationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent3Visitor:
		return t.VisitNegation(s)

	default:
		return t.VisitChildren(s)
	}
}

type ComparisonContext struct {
	*ExpressionContext
	op antlr.Token
}

func NewComparisonContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ComparisonContext {
	var p = new(ComparisonContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *ComparisonContext) GetOp() antlr.Token { return s.op }

func (s *ComparisonContext) SetOp(v antlr.Token) { s.op = v }

func (s *ComparisonContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ComparisonContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *ComparisonContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ComparisonContext) LTE() antlr.TerminalNode {
	return s.GetToken(Excellent3ParserLTE, 0)
}

func (s *ComparisonContext) LT() antlr.TerminalNode {
	return s.GetToken(Excellent3ParserLT, 0)
}

func (s *ComparisonContext) GTE() antlr.TerminalNode {
	return s.GetToken(Excellent3ParserGTE, 0)
}

func (s *ComparisonContext) GT() antlr.TerminalNode {
	return s.GetToken(Excellent3ParserGT, 0)
}

func (s *ComparisonContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.EnterComparison(s)
	}
}

func (s *ComparisonContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.ExitComparison(s)
	}
}

func (s *ComparisonContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent3Visitor:
		return t.VisitComparison(s)

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
	return s.GetToken(Excellent3ParserFALSE, 0)
}

func (s *FalseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.EnterFalse(s)
	}
}

func (s *FalseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.ExitFalse(s)
	}
}

func (s *FalseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent3Visitor:
		return t.VisitFalse(s)

	default:
		return t.VisitChildren(s)
	}
}

type AdditionOrSubtractionContext struct {
	*ExpressionContext
	op antlr.Token
}

func NewAdditionOrSubtractionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *AdditionOrSubtractionContext {
	var p = new(AdditionOrSubtractionContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *AdditionOrSubtractionContext) GetOp() antlr.Token { return s.op }

func (s *AdditionOrSubtractionContext) SetOp(v antlr.Token) { s.op = v }

func (s *AdditionOrSubtractionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AdditionOrSubtractionContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *AdditionOrSubtractionContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *AdditionOrSubtractionContext) PLUS() antlr.TerminalNode {
	return s.GetToken(Excellent3ParserPLUS, 0)
}

func (s *AdditionOrSubtractionContext) MINUS() antlr.TerminalNode {
	return s.GetToken(Excellent3ParserMINUS, 0)
}

func (s *AdditionOrSubtractionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.EnterAdditionOrSubtraction(s)
	}
}

func (s *AdditionOrSubtractionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.ExitAdditionOrSubtraction(s)
	}
}

func (s *AdditionOrSubtractionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent3Visitor:
		return t.VisitAdditionOrSubtraction(s)

	default:
		return t.VisitChildren(s)
	}
}

type TextLiteralContext struct {
	*ExpressionContext
}

func NewTextLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TextLiteralContext {
	var p = new(TextLiteralContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *TextLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TextLiteralContext) TEXT() antlr.TerminalNode {
	return s.GetToken(Excellent3ParserTEXT, 0)
}

func (s *TextLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.EnterTextLiteral(s)
	}
}

func (s *TextLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.ExitTextLiteral(s)
	}
}

func (s *TextLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent3Visitor:
		return t.VisitTextLiteral(s)

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
	return s.GetToken(Excellent3ParserAMPERSAND, 0)
}

func (s *ConcatenationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.EnterConcatenation(s)
	}
}

func (s *ConcatenationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.ExitConcatenation(s)
	}
}

func (s *ConcatenationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent3Visitor:
		return t.VisitConcatenation(s)

	default:
		return t.VisitChildren(s)
	}
}

type NullContext struct {
	*ExpressionContext
}

func NewNullContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NullContext {
	var p = new(NullContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *NullContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NullContext) NULL() antlr.TerminalNode {
	return s.GetToken(Excellent3ParserNULL, 0)
}

func (s *NullContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.EnterNull(s)
	}
}

func (s *NullContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.ExitNull(s)
	}
}

func (s *NullContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent3Visitor:
		return t.VisitNull(s)

	default:
		return t.VisitChildren(s)
	}
}

type MultiplicationOrDivisionContext struct {
	*ExpressionContext
	op antlr.Token
}

func NewMultiplicationOrDivisionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MultiplicationOrDivisionContext {
	var p = new(MultiplicationOrDivisionContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *MultiplicationOrDivisionContext) GetOp() antlr.Token { return s.op }

func (s *MultiplicationOrDivisionContext) SetOp(v antlr.Token) { s.op = v }

func (s *MultiplicationOrDivisionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MultiplicationOrDivisionContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *MultiplicationOrDivisionContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *MultiplicationOrDivisionContext) TIMES() antlr.TerminalNode {
	return s.GetToken(Excellent3ParserTIMES, 0)
}

func (s *MultiplicationOrDivisionContext) DIVIDE() antlr.TerminalNode {
	return s.GetToken(Excellent3ParserDIVIDE, 0)
}

func (s *MultiplicationOrDivisionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.EnterMultiplicationOrDivision(s)
	}
}

func (s *MultiplicationOrDivisionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.ExitMultiplicationOrDivision(s)
	}
}

func (s *MultiplicationOrDivisionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent3Visitor:
		return t.VisitMultiplicationOrDivision(s)

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
	return s.GetToken(Excellent3ParserTRUE, 0)
}

func (s *TrueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.EnterTrue(s)
	}
}

func (s *TrueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.ExitTrue(s)
	}
}

func (s *TrueContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent3Visitor:
		return t.VisitTrue(s)

	default:
		return t.VisitChildren(s)
	}
}

type AtomReferenceContext struct {
	*ExpressionContext
}

func NewAtomReferenceContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *AtomReferenceContext {
	var p = new(AtomReferenceContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *AtomReferenceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AtomReferenceContext) Atom() IAtomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAtomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAtomContext)
}

func (s *AtomReferenceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.EnterAtomReference(s)
	}
}

func (s *AtomReferenceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.ExitAtomReference(s)
	}
}

func (s *AtomReferenceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent3Visitor:
		return t.VisitAtomReference(s)

	default:
		return t.VisitChildren(s)
	}
}

type AnonFunctionContext struct {
	*ExpressionContext
}

func NewAnonFunctionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *AnonFunctionContext {
	var p = new(AnonFunctionContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *AnonFunctionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AnonFunctionContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(Excellent3ParserLPAREN, 0)
}

func (s *AnonFunctionContext) NameList() INameListContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INameListContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INameListContext)
}

func (s *AnonFunctionContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(Excellent3ParserRPAREN, 0)
}

func (s *AnonFunctionContext) ARROW() antlr.TerminalNode {
	return s.GetToken(Excellent3ParserARROW, 0)
}

func (s *AnonFunctionContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *AnonFunctionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.EnterAnonFunction(s)
	}
}

func (s *AnonFunctionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.ExitAnonFunction(s)
	}
}

func (s *AnonFunctionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent3Visitor:
		return t.VisitAnonFunction(s)

	default:
		return t.VisitChildren(s)
	}
}

type EqualityContext struct {
	*ExpressionContext
	op antlr.Token
}

func NewEqualityContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *EqualityContext {
	var p = new(EqualityContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *EqualityContext) GetOp() antlr.Token { return s.op }

func (s *EqualityContext) SetOp(v antlr.Token) { s.op = v }

func (s *EqualityContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EqualityContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *EqualityContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *EqualityContext) EQ() antlr.TerminalNode {
	return s.GetToken(Excellent3ParserEQ, 0)
}

func (s *EqualityContext) NEQ() antlr.TerminalNode {
	return s.GetToken(Excellent3ParserNEQ, 0)
}

func (s *EqualityContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.EnterEquality(s)
	}
}

func (s *EqualityContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.ExitEquality(s)
	}
}

func (s *EqualityContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent3Visitor:
		return t.VisitEquality(s)

	default:
		return t.VisitChildren(s)
	}
}

type NumberLiteralContext struct {
	*ExpressionContext
}

func NewNumberLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NumberLiteralContext {
	var p = new(NumberLiteralContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *NumberLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NumberLiteralContext) INTEGER() antlr.TerminalNode {
	return s.GetToken(Excellent3ParserINTEGER, 0)
}

func (s *NumberLiteralContext) DECIMAL() antlr.TerminalNode {
	return s.GetToken(Excellent3ParserDECIMAL, 0)
}

func (s *NumberLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.EnterNumberLiteral(s)
	}
}

func (s *NumberLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.ExitNumberLiteral(s)
	}
}

func (s *NumberLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent3Visitor:
		return t.VisitNumberLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

type ExponentContext struct {
	*ExpressionContext
}

func NewExponentContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ExponentContext {
	var p = new(ExponentContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *ExponentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExponentContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *ExponentContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ExponentContext) EXPONENT() antlr.TerminalNode {
	return s.GetToken(Excellent3ParserEXPONENT, 0)
}

func (s *ExponentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.EnterExponent(s)
	}
}

func (s *ExponentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.ExitExponent(s)
	}
}

func (s *ExponentContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent3Visitor:
		return t.VisitExponent(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *Excellent3Parser) Expression() (localctx IExpressionContext) {
	return p.expression(0)
}

func (p *Excellent3Parser) expression(_p int) (localctx IExpressionContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewExpressionContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IExpressionContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 2
	p.EnterRecursionRule(localctx, 2, Excellent3ParserRULE_expression, _p)
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
	p.SetState(28)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext()) {
	case 1:
		localctx = NewAtomReferenceContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx

		{
			p.SetState(14)
			p.atom(0)
		}

	case 2:
		localctx = NewNegationContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(15)
			p.Match(Excellent3ParserMINUS)
		}
		{
			p.SetState(16)
			p.expression(13)
		}

	case 3:
		localctx = NewAnonFunctionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(17)
			p.Match(Excellent3ParserLPAREN)
		}
		{
			p.SetState(18)
			p.NameList()
		}
		{
			p.SetState(19)
			p.Match(Excellent3ParserRPAREN)
		}
		{
			p.SetState(20)
			p.Match(Excellent3ParserARROW)
		}
		{
			p.SetState(21)
			p.expression(6)
		}

	case 4:
		localctx = NewTextLiteralContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(23)
			p.Match(Excellent3ParserTEXT)
		}

	case 5:
		localctx = NewNumberLiteralContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(24)
			_la = p.GetTokenStream().LA(1)

			if !(_la == Excellent3ParserINTEGER || _la == Excellent3ParserDECIMAL) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	case 6:
		localctx = NewTrueContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(25)
			p.Match(Excellent3ParserTRUE)
		}

	case 7:
		localctx = NewFalseContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(26)
			p.Match(Excellent3ParserFALSE)
		}

	case 8:
		localctx = NewNullContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(27)
			p.Match(Excellent3ParserNULL)
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(50)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 2, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(48)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 1, p.GetParserRuleContext()) {
			case 1:
				localctx = NewExponentContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, Excellent3ParserRULE_expression)
				p.SetState(30)

				if !(p.Precpred(p.GetParserRuleContext(), 12)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 12)", ""))
				}
				{
					p.SetState(31)
					p.Match(Excellent3ParserEXPONENT)
				}
				{
					p.SetState(32)
					p.expression(13)
				}

			case 2:
				localctx = NewMultiplicationOrDivisionContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, Excellent3ParserRULE_expression)
				p.SetState(33)

				if !(p.Precpred(p.GetParserRuleContext(), 11)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 11)", ""))
				}
				{
					p.SetState(34)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*MultiplicationOrDivisionContext).op = _lt

					_la = p.GetTokenStream().LA(1)

					if !(_la == Excellent3ParserTIMES || _la == Excellent3ParserDIVIDE) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*MultiplicationOrDivisionContext).op = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(35)
					p.expression(12)
				}

			case 3:
				localctx = NewAdditionOrSubtractionContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, Excellent3ParserRULE_expression)
				p.SetState(36)

				if !(p.Precpred(p.GetParserRuleContext(), 10)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 10)", ""))
				}
				{
					p.SetState(37)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*AdditionOrSubtractionContext).op = _lt

					_la = p.GetTokenStream().LA(1)

					if !(_la == Excellent3ParserPLUS || _la == Excellent3ParserMINUS) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*AdditionOrSubtractionContext).op = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(38)
					p.expression(11)
				}

			case 4:
				localctx = NewComparisonContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, Excellent3ParserRULE_expression)
				p.SetState(39)

				if !(p.Precpred(p.GetParserRuleContext(), 9)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 9)", ""))
				}
				{
					p.SetState(40)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*ComparisonContext).op = _lt

					_la = p.GetTokenStream().LA(1)

					if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<Excellent3ParserLTE)|(1<<Excellent3ParserLT)|(1<<Excellent3ParserGTE)|(1<<Excellent3ParserGT))) != 0) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*ComparisonContext).op = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(41)
					p.expression(10)
				}

			case 5:
				localctx = NewEqualityContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, Excellent3ParserRULE_expression)
				p.SetState(42)

				if !(p.Precpred(p.GetParserRuleContext(), 8)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 8)", ""))
				}
				{
					p.SetState(43)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*EqualityContext).op = _lt

					_la = p.GetTokenStream().LA(1)

					if !(_la == Excellent3ParserEQ || _la == Excellent3ParserNEQ) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*EqualityContext).op = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(44)
					p.expression(9)
				}

			case 6:
				localctx = NewConcatenationContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, Excellent3ParserRULE_expression)
				p.SetState(45)

				if !(p.Precpred(p.GetParserRuleContext(), 7)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 7)", ""))
				}
				{
					p.SetState(46)
					p.Match(Excellent3ParserAMPERSAND)
				}
				{
					p.SetState(47)
					p.expression(8)
				}

			}

		}
		p.SetState(52)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 2, p.GetParserRuleContext())
	}

	return localctx
}

// IAtomContext is an interface to support dynamic dispatch.
type IAtomContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAtomContext differentiates from other interfaces.
	IsAtomContext()
}

type AtomContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAtomContext() *AtomContext {
	var p = new(AtomContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = Excellent3ParserRULE_atom
	return p
}

func (*AtomContext) IsAtomContext() {}

func NewAtomContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AtomContext {
	var p = new(AtomContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = Excellent3ParserRULE_atom

	return p
}

func (s *AtomContext) GetParser() antlr.Parser { return s.parser }

func (s *AtomContext) CopyFrom(ctx *AtomContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *AtomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AtomContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type ParenthesesContext struct {
	*AtomContext
}

func NewParenthesesContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ParenthesesContext {
	var p = new(ParenthesesContext)

	p.AtomContext = NewEmptyAtomContext()
	p.parser = parser
	p.CopyFrom(ctx.(*AtomContext))

	return p
}

func (s *ParenthesesContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParenthesesContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(Excellent3ParserLPAREN, 0)
}

func (s *ParenthesesContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ParenthesesContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(Excellent3ParserRPAREN, 0)
}

func (s *ParenthesesContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.EnterParentheses(s)
	}
}

func (s *ParenthesesContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.ExitParentheses(s)
	}
}

func (s *ParenthesesContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent3Visitor:
		return t.VisitParentheses(s)

	default:
		return t.VisitChildren(s)
	}
}

type DotLookupContext struct {
	*AtomContext
}

func NewDotLookupContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DotLookupContext {
	var p = new(DotLookupContext)

	p.AtomContext = NewEmptyAtomContext()
	p.parser = parser
	p.CopyFrom(ctx.(*AtomContext))

	return p
}

func (s *DotLookupContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DotLookupContext) Atom() IAtomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAtomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAtomContext)
}

func (s *DotLookupContext) DOT() antlr.TerminalNode {
	return s.GetToken(Excellent3ParserDOT, 0)
}

func (s *DotLookupContext) NAME() antlr.TerminalNode {
	return s.GetToken(Excellent3ParserNAME, 0)
}

func (s *DotLookupContext) INTEGER() antlr.TerminalNode {
	return s.GetToken(Excellent3ParserINTEGER, 0)
}

func (s *DotLookupContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.EnterDotLookup(s)
	}
}

func (s *DotLookupContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.ExitDotLookup(s)
	}
}

func (s *DotLookupContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent3Visitor:
		return t.VisitDotLookup(s)

	default:
		return t.VisitChildren(s)
	}
}

type FunctionCallContext struct {
	*AtomContext
}

func NewFunctionCallContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FunctionCallContext {
	var p = new(FunctionCallContext)

	p.AtomContext = NewEmptyAtomContext()
	p.parser = parser
	p.CopyFrom(ctx.(*AtomContext))

	return p
}

func (s *FunctionCallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionCallContext) Atom() IAtomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAtomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAtomContext)
}

func (s *FunctionCallContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(Excellent3ParserLPAREN, 0)
}

func (s *FunctionCallContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(Excellent3ParserRPAREN, 0)
}

func (s *FunctionCallContext) Parameters() IParametersContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IParametersContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IParametersContext)
}

func (s *FunctionCallContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.EnterFunctionCall(s)
	}
}

func (s *FunctionCallContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.ExitFunctionCall(s)
	}
}

func (s *FunctionCallContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent3Visitor:
		return t.VisitFunctionCall(s)

	default:
		return t.VisitChildren(s)
	}
}

type ArrayLookupContext struct {
	*AtomContext
}

func NewArrayLookupContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ArrayLookupContext {
	var p = new(ArrayLookupContext)

	p.AtomContext = NewEmptyAtomContext()
	p.parser = parser
	p.CopyFrom(ctx.(*AtomContext))

	return p
}

func (s *ArrayLookupContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArrayLookupContext) Atom() IAtomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAtomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAtomContext)
}

func (s *ArrayLookupContext) LBRACK() antlr.TerminalNode {
	return s.GetToken(Excellent3ParserLBRACK, 0)
}

func (s *ArrayLookupContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ArrayLookupContext) RBRACK() antlr.TerminalNode {
	return s.GetToken(Excellent3ParserRBRACK, 0)
}

func (s *ArrayLookupContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.EnterArrayLookup(s)
	}
}

func (s *ArrayLookupContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.ExitArrayLookup(s)
	}
}

func (s *ArrayLookupContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent3Visitor:
		return t.VisitArrayLookup(s)

	default:
		return t.VisitChildren(s)
	}
}

type ContextReferenceContext struct {
	*AtomContext
}

func NewContextReferenceContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ContextReferenceContext {
	var p = new(ContextReferenceContext)

	p.AtomContext = NewEmptyAtomContext()
	p.parser = parser
	p.CopyFrom(ctx.(*AtomContext))

	return p
}

func (s *ContextReferenceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ContextReferenceContext) NAME() antlr.TerminalNode {
	return s.GetToken(Excellent3ParserNAME, 0)
}

func (s *ContextReferenceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.EnterContextReference(s)
	}
}

func (s *ContextReferenceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.ExitContextReference(s)
	}
}

func (s *ContextReferenceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent3Visitor:
		return t.VisitContextReference(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *Excellent3Parser) Atom() (localctx IAtomContext) {
	return p.atom(0)
}

func (p *Excellent3Parser) atom(_p int) (localctx IAtomContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewAtomContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IAtomContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 4
	p.EnterRecursionRule(localctx, 4, Excellent3ParserRULE_atom, _p)
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
	p.SetState(59)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case Excellent3ParserLPAREN:
		localctx = NewParenthesesContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx

		{
			p.SetState(54)
			p.Match(Excellent3ParserLPAREN)
		}
		{
			p.SetState(55)
			p.expression(0)
		}
		{
			p.SetState(56)
			p.Match(Excellent3ParserRPAREN)
		}

	case Excellent3ParserNAME:
		localctx = NewContextReferenceContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(58)
			p.Match(Excellent3ParserNAME)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(77)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 6, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(75)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 5, p.GetParserRuleContext()) {
			case 1:
				localctx = NewFunctionCallContext(p, NewAtomContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, Excellent3ParserRULE_atom)
				p.SetState(61)

				if !(p.Precpred(p.GetParserRuleContext(), 5)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 5)", ""))
				}
				{
					p.SetState(62)
					p.Match(Excellent3ParserLPAREN)
				}
				p.SetState(64)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)

				if ((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<Excellent3ParserLPAREN)|(1<<Excellent3ParserMINUS)|(1<<Excellent3ParserTEXT)|(1<<Excellent3ParserINTEGER)|(1<<Excellent3ParserDECIMAL)|(1<<Excellent3ParserTRUE)|(1<<Excellent3ParserFALSE)|(1<<Excellent3ParserNULL)|(1<<Excellent3ParserNAME))) != 0 {
					{
						p.SetState(63)
						p.Parameters()
					}

				}
				{
					p.SetState(66)
					p.Match(Excellent3ParserRPAREN)
				}

			case 2:
				localctx = NewDotLookupContext(p, NewAtomContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, Excellent3ParserRULE_atom)
				p.SetState(67)

				if !(p.Precpred(p.GetParserRuleContext(), 4)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 4)", ""))
				}
				{
					p.SetState(68)
					p.Match(Excellent3ParserDOT)
				}
				{
					p.SetState(69)
					_la = p.GetTokenStream().LA(1)

					if !(_la == Excellent3ParserINTEGER || _la == Excellent3ParserNAME) {
						p.GetErrorHandler().RecoverInline(p)
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}

			case 3:
				localctx = NewArrayLookupContext(p, NewAtomContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, Excellent3ParserRULE_atom)
				p.SetState(70)

				if !(p.Precpred(p.GetParserRuleContext(), 3)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 3)", ""))
				}
				{
					p.SetState(71)
					p.Match(Excellent3ParserLBRACK)
				}
				{
					p.SetState(72)
					p.expression(0)
				}
				{
					p.SetState(73)
					p.Match(Excellent3ParserRBRACK)
				}

			}

		}
		p.SetState(79)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 6, p.GetParserRuleContext())
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
	p.RuleIndex = Excellent3ParserRULE_parameters
	return p
}

func (*ParametersContext) IsParametersContext() {}

func NewParametersContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParametersContext {
	var p = new(ParametersContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = Excellent3ParserRULE_parameters

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
	return s.GetTokens(Excellent3ParserCOMMA)
}

func (s *FunctionParametersContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(Excellent3ParserCOMMA, i)
}

func (s *FunctionParametersContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.EnterFunctionParameters(s)
	}
}

func (s *FunctionParametersContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.ExitFunctionParameters(s)
	}
}

func (s *FunctionParametersContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent3Visitor:
		return t.VisitFunctionParameters(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *Excellent3Parser) Parameters() (localctx IParametersContext) {
	localctx = NewParametersContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, Excellent3ParserRULE_parameters)
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
		p.SetState(80)
		p.expression(0)
	}
	p.SetState(85)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == Excellent3ParserCOMMA {
		{
			p.SetState(81)
			p.Match(Excellent3ParserCOMMA)
		}
		{
			p.SetState(82)
			p.expression(0)
		}

		p.SetState(87)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// INameListContext is an interface to support dynamic dispatch.
type INameListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsNameListContext differentiates from other interfaces.
	IsNameListContext()
}

type NameListContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNameListContext() *NameListContext {
	var p = new(NameListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = Excellent3ParserRULE_nameList
	return p
}

func (*NameListContext) IsNameListContext() {}

func NewNameListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NameListContext {
	var p = new(NameListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = Excellent3ParserRULE_nameList

	return p
}

func (s *NameListContext) GetParser() antlr.Parser { return s.parser }

func (s *NameListContext) AllNAME() []antlr.TerminalNode {
	return s.GetTokens(Excellent3ParserNAME)
}

func (s *NameListContext) NAME(i int) antlr.TerminalNode {
	return s.GetToken(Excellent3ParserNAME, i)
}

func (s *NameListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(Excellent3ParserCOMMA)
}

func (s *NameListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(Excellent3ParserCOMMA, i)
}

func (s *NameListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NameListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NameListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.EnterNameList(s)
	}
}

func (s *NameListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Excellent3Listener); ok {
		listenerT.ExitNameList(s)
	}
}

func (s *NameListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case Excellent3Visitor:
		return t.VisitNameList(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *Excellent3Parser) NameList() (localctx INameListContext) {
	localctx = NewNameListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, Excellent3ParserRULE_nameList)
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
		p.SetState(88)
		p.Match(Excellent3ParserNAME)
	}
	p.SetState(93)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == Excellent3ParserCOMMA {
		{
			p.SetState(89)
			p.Match(Excellent3ParserCOMMA)
		}
		{
			p.SetState(90)
			p.Match(Excellent3ParserNAME)
		}

		p.SetState(95)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

func (p *Excellent3Parser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 1:
		var t *ExpressionContext = nil
		if localctx != nil {
			t = localctx.(*ExpressionContext)
		}
		return p.Expression_Sempred(t, predIndex)

	case 2:
		var t *AtomContext = nil
		if localctx != nil {
			t = localctx.(*AtomContext)
		}
		return p.Atom_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *Excellent3Parser) Expression_Sempred(localctx antlr.RuleContext, predIndex int) bool {
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

func (p *Excellent3Parser) Atom_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 6:
		return p.Precpred(p.GetParserRuleContext(), 5)

	case 7:
		return p.Precpred(p.GetParserRuleContext(), 4)

	case 8:
		return p.Precpred(p.GetParserRuleContext(), 3)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
