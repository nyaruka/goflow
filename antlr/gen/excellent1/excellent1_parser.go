// Code generated from Excellent1.g4 by ANTLR 4.13.1. DO NOT EDIT.

package gen // Excellent1
import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type Excellent1Parser struct {
	*antlr.BaseParser
}

var Excellent1ParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func excellent1ParserInit() {
	staticData := &Excellent1ParserStaticData
	staticData.LiteralNames = []string{
		"", "','", "'('", "')'", "'+'", "'-'", "'*'", "'/'", "'^'", "'='", "'<>'",
		"'<='", "'<'", "'>='", "'>'", "'&'",
	}
	staticData.SymbolicNames = []string{
		"", "COMMA", "LPAREN", "RPAREN", "PLUS", "MINUS", "TIMES", "DIVIDE",
		"EXPONENT", "EQ", "NEQ", "LTE", "LT", "GTE", "GT", "AMPERSAND", "DECIMAL",
		"STRING", "TRUE", "FALSE", "NAME", "WS", "ERROR",
	}
	staticData.RuleNames = []string{
		"parse", "expression", "fnname", "parameters",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 22, 66, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 1, 0, 1,
		0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 16, 8, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 31, 8, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 5, 1, 51, 8, 1, 10, 1, 12, 1, 54, 9, 1,
		1, 2, 1, 2, 1, 3, 1, 3, 1, 3, 5, 3, 61, 8, 3, 10, 3, 12, 3, 64, 9, 3, 1,
		3, 0, 1, 2, 4, 0, 2, 4, 6, 0, 5, 1, 0, 6, 7, 1, 0, 4, 5, 1, 0, 11, 14,
		1, 0, 9, 10, 1, 0, 18, 20, 76, 0, 8, 1, 0, 0, 0, 2, 30, 1, 0, 0, 0, 4,
		55, 1, 0, 0, 0, 6, 57, 1, 0, 0, 0, 8, 9, 3, 2, 1, 0, 9, 10, 5, 0, 0, 1,
		10, 1, 1, 0, 0, 0, 11, 12, 6, 1, -1, 0, 12, 13, 3, 4, 2, 0, 13, 15, 5,
		2, 0, 0, 14, 16, 3, 6, 3, 0, 15, 14, 1, 0, 0, 0, 15, 16, 1, 0, 0, 0, 16,
		17, 1, 0, 0, 0, 17, 18, 5, 3, 0, 0, 18, 31, 1, 0, 0, 0, 19, 20, 5, 5, 0,
		0, 20, 31, 3, 2, 1, 13, 21, 31, 5, 17, 0, 0, 22, 31, 5, 16, 0, 0, 23, 31,
		5, 18, 0, 0, 24, 31, 5, 19, 0, 0, 25, 31, 5, 20, 0, 0, 26, 27, 5, 2, 0,
		0, 27, 28, 3, 2, 1, 0, 28, 29, 5, 3, 0, 0, 29, 31, 1, 0, 0, 0, 30, 11,
		1, 0, 0, 0, 30, 19, 1, 0, 0, 0, 30, 21, 1, 0, 0, 0, 30, 22, 1, 0, 0, 0,
		30, 23, 1, 0, 0, 0, 30, 24, 1, 0, 0, 0, 30, 25, 1, 0, 0, 0, 30, 26, 1,
		0, 0, 0, 31, 52, 1, 0, 0, 0, 32, 33, 10, 12, 0, 0, 33, 34, 5, 8, 0, 0,
		34, 51, 3, 2, 1, 13, 35, 36, 10, 11, 0, 0, 36, 37, 7, 0, 0, 0, 37, 51,
		3, 2, 1, 12, 38, 39, 10, 10, 0, 0, 39, 40, 7, 1, 0, 0, 40, 51, 3, 2, 1,
		11, 41, 42, 10, 9, 0, 0, 42, 43, 7, 2, 0, 0, 43, 51, 3, 2, 1, 10, 44, 45,
		10, 8, 0, 0, 45, 46, 7, 3, 0, 0, 46, 51, 3, 2, 1, 9, 47, 48, 10, 7, 0,
		0, 48, 49, 5, 15, 0, 0, 49, 51, 3, 2, 1, 8, 50, 32, 1, 0, 0, 0, 50, 35,
		1, 0, 0, 0, 50, 38, 1, 0, 0, 0, 50, 41, 1, 0, 0, 0, 50, 44, 1, 0, 0, 0,
		50, 47, 1, 0, 0, 0, 51, 54, 1, 0, 0, 0, 52, 50, 1, 0, 0, 0, 52, 53, 1,
		0, 0, 0, 53, 3, 1, 0, 0, 0, 54, 52, 1, 0, 0, 0, 55, 56, 7, 4, 0, 0, 56,
		5, 1, 0, 0, 0, 57, 62, 3, 2, 1, 0, 58, 59, 5, 1, 0, 0, 59, 61, 3, 2, 1,
		0, 60, 58, 1, 0, 0, 0, 61, 64, 1, 0, 0, 0, 62, 60, 1, 0, 0, 0, 62, 63,
		1, 0, 0, 0, 63, 7, 1, 0, 0, 0, 64, 62, 1, 0, 0, 0, 5, 15, 30, 50, 52, 62,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// Excellent1ParserInit initializes any static state used to implement Excellent1Parser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewExcellent1Parser(). You can call this function if you wish to initialize the static state ahead
// of time.
func Excellent1ParserInit() {
	staticData := &Excellent1ParserStaticData
	staticData.once.Do(excellent1ParserInit)
}

// NewExcellent1Parser produces a new parser instance for the optional input antlr.TokenStream.
func NewExcellent1Parser(input antlr.TokenStream) *Excellent1Parser {
	Excellent1ParserInit()
	this := new(Excellent1Parser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &Excellent1ParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
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

	// Getter signatures
	Expression() IExpressionContext
	EOF() antlr.TerminalNode

	// IsParseContext differentiates from other interfaces.
	IsParseContext()
}

type ParseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyParseContext() *ParseContext {
	var p = new(ParseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = Excellent1ParserRULE_parse
	return p
}

func InitEmptyParseContext(p *ParseContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = Excellent1ParserRULE_parse
}

func (*ParseContext) IsParseContext() {}

func NewParseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParseContext {
	var p = new(ParseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = Excellent1ParserRULE_parse

	return p
}

func (s *ParseContext) GetParser() antlr.Parser { return s.parser }

func (s *ParseContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

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
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(8)
		p.expression(0)
	}
	{
		p.SetState(9)
		p.Match(Excellent1ParserEOF)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
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
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpressionContext() *ExpressionContext {
	var p = new(ExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = Excellent1ParserRULE_expression
	return p
}

func InitEmptyExpressionContext(p *ExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = Excellent1ParserRULE_expression
}

func (*ExpressionContext) IsExpressionContext() {}

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext {
	var p = new(ExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = Excellent1ParserRULE_expression

	return p
}

func (s *ExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionContext) CopyAll(ctx *ExpressionContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *ExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type DecimalLiteralContext struct {
	ExpressionContext
}

func NewDecimalLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DecimalLiteralContext {
	var p = new(DecimalLiteralContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

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
	ExpressionContext
}

func NewParenthesesContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ParenthesesContext {
	var p = new(ParenthesesContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}

func (s *ParenthesesContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParenthesesContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(Excellent1ParserLPAREN, 0)
}

func (s *ParenthesesContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

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
	ExpressionContext
}

func NewNegationContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NegationContext {
	var p = new(NegationContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}

func (s *NegationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NegationContext) MINUS() antlr.TerminalNode {
	return s.GetToken(Excellent1ParserMINUS, 0)
}

func (s *NegationContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

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
	ExpressionContext
}

func NewExponentExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ExponentExpressionContext {
	var p = new(ExponentExpressionContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}

func (s *ExponentExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExponentExpressionContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *ExponentExpressionContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

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
	ExpressionContext
	op antlr.Token
}

func NewAdditionOrSubtractionExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *AdditionOrSubtractionExpressionContext {
	var p = new(AdditionOrSubtractionExpressionContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}

func (s *AdditionOrSubtractionExpressionContext) GetOp() antlr.Token { return s.op }

func (s *AdditionOrSubtractionExpressionContext) SetOp(v antlr.Token) { s.op = v }

func (s *AdditionOrSubtractionExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AdditionOrSubtractionExpressionContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *AdditionOrSubtractionExpressionContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

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
	ExpressionContext
}

func NewFalseContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FalseContext {
	var p = new(FalseContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

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
	ExpressionContext
}

func NewContextReferenceContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ContextReferenceContext {
	var p = new(ContextReferenceContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

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
	ExpressionContext
	op antlr.Token
}

func NewComparisonExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ComparisonExpressionContext {
	var p = new(ComparisonExpressionContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}

func (s *ComparisonExpressionContext) GetOp() antlr.Token { return s.op }

func (s *ComparisonExpressionContext) SetOp(v antlr.Token) { s.op = v }

func (s *ComparisonExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ComparisonExpressionContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *ComparisonExpressionContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

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
	ExpressionContext
}

func NewConcatenationContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ConcatenationContext {
	var p = new(ConcatenationContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}

func (s *ConcatenationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConcatenationContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *ConcatenationContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

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
	ExpressionContext
}

func NewStringLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *StringLiteralContext {
	var p = new(StringLiteralContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

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
	ExpressionContext
}

func NewFunctionCallContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FunctionCallContext {
	var p = new(FunctionCallContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}

func (s *FunctionCallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionCallContext) Fnname() IFnnameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFnnameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

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
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParametersContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

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
	ExpressionContext
}

func NewTrueContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TrueContext {
	var p = new(TrueContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

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
	ExpressionContext
	op antlr.Token
}

func NewEqualityExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *EqualityExpressionContext {
	var p = new(EqualityExpressionContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}

func (s *EqualityExpressionContext) GetOp() antlr.Token { return s.op }

func (s *EqualityExpressionContext) SetOp(v antlr.Token) { s.op = v }

func (s *EqualityExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EqualityExpressionContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *EqualityExpressionContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

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
	ExpressionContext
	op antlr.Token
}

func NewMultiplicationOrDivisionExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MultiplicationOrDivisionExpressionContext {
	var p = new(MultiplicationOrDivisionExpressionContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}

func (s *MultiplicationOrDivisionExpressionContext) GetOp() antlr.Token { return s.op }

func (s *MultiplicationOrDivisionExpressionContext) SetOp(v antlr.Token) { s.op = v }

func (s *MultiplicationOrDivisionExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MultiplicationOrDivisionExpressionContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *MultiplicationOrDivisionExpressionContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

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

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(30)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 1, p.GetParserRuleContext()) {
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
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(15)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&2031652) != 0 {
			{
				p.SetState(14)
				p.Parameters()
			}

		}
		{
			p.SetState(17)
			p.Match(Excellent1ParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		localctx = NewNegationContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(19)
			p.Match(Excellent1ParserMINUS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
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
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 4:
		localctx = NewDecimalLiteralContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(22)
			p.Match(Excellent1ParserDECIMAL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 5:
		localctx = NewTrueContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(23)
			p.Match(Excellent1ParserTRUE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 6:
		localctx = NewFalseContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(24)
			p.Match(Excellent1ParserFALSE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 7:
		localctx = NewContextReferenceContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(25)
			p.Match(Excellent1ParserNAME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 8:
		localctx = NewParenthesesContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(26)
			p.Match(Excellent1ParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(27)
			p.expression(0)
		}
		{
			p.SetState(28)
			p.Match(Excellent1ParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(52)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 3, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(50)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 2, p.GetParserRuleContext()) {
			case 1:
				localctx = NewExponentExpressionContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, Excellent1ParserRULE_expression)
				p.SetState(32)

				if !(p.Precpred(p.GetParserRuleContext(), 12)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 12)", ""))
					goto errorExit
				}
				{
					p.SetState(33)
					p.Match(Excellent1ParserEXPONENT)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
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
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 11)", ""))
					goto errorExit
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
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 10)", ""))
					goto errorExit
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
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 9)", ""))
					goto errorExit
				}
				{
					p.SetState(42)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*ComparisonExpressionContext).op = _lt

					_la = p.GetTokenStream().LA(1)

					if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&30720) != 0) {
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
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 8)", ""))
					goto errorExit
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
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 7)", ""))
					goto errorExit
				}
				{
					p.SetState(48)
					p.Match(Excellent1ParserAMPERSAND)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(49)
					p.expression(8)
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(54)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 3, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.UnrollRecursionContexts(_parentctx)
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFnnameContext is an interface to support dynamic dispatch.
type IFnnameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NAME() antlr.TerminalNode
	TRUE() antlr.TerminalNode
	FALSE() antlr.TerminalNode

	// IsFnnameContext differentiates from other interfaces.
	IsFnnameContext()
}

type FnnameContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFnnameContext() *FnnameContext {
	var p = new(FnnameContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = Excellent1ParserRULE_fnname
	return p
}

func InitEmptyFnnameContext(p *FnnameContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = Excellent1ParserRULE_fnname
}

func (*FnnameContext) IsFnnameContext() {}

func NewFnnameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FnnameContext {
	var p = new(FnnameContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

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

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(55)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&1835008) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
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
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyParametersContext() *ParametersContext {
	var p = new(ParametersContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = Excellent1ParserRULE_parameters
	return p
}

func InitEmptyParametersContext(p *ParametersContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = Excellent1ParserRULE_parameters
}

func (*ParametersContext) IsParametersContext() {}

func NewParametersContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParametersContext {
	var p = new(ParametersContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = Excellent1ParserRULE_parameters

	return p
}

func (s *ParametersContext) GetParser() antlr.Parser { return s.parser }

func (s *ParametersContext) CopyAll(ctx *ParametersContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *ParametersContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParametersContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type FunctionParametersContext struct {
	ParametersContext
}

func NewFunctionParametersContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FunctionParametersContext {
	var p = new(FunctionParametersContext)

	InitEmptyParametersContext(&p.ParametersContext)
	p.parser = parser
	p.CopyAll(ctx.(*ParametersContext))

	return p
}

func (s *FunctionParametersContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionParametersContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *FunctionParametersContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

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

	localctx = NewFunctionParametersContext(p, localctx)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(57)
		p.expression(0)
	}
	p.SetState(62)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == Excellent1ParserCOMMA {
		{
			p.SetState(58)
			p.Match(Excellent1ParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(59)
			p.expression(0)
		}

		p.SetState(64)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
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
