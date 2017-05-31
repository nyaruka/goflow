// Generated from excellent/gen/Excellent.g4 by ANTLR 4.7.

package gen // Excellent
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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 27, 87, 4,
	2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 3, 2, 3, 2, 3,
	2, 3, 3, 3, 3, 3, 3, 3, 3, 5, 3, 20, 10, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 5, 3, 29, 10, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 7, 3, 39, 10, 3, 12, 3, 14, 3, 42, 11, 3, 3, 4, 3, 4, 3, 4, 3,
	4, 3, 4, 3, 4, 3, 4, 3, 4, 5, 4, 52, 10, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3,
	4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3,
	4, 3, 4, 7, 4, 72, 10, 4, 12, 4, 14, 4, 75, 11, 4, 3, 5, 3, 5, 3, 6, 3,
	6, 3, 6, 7, 6, 82, 10, 6, 12, 6, 14, 6, 85, 11, 6, 3, 6, 2, 4, 4, 6, 7,
	2, 4, 6, 8, 10, 2, 7, 3, 2, 11, 12, 3, 2, 9, 10, 3, 2, 16, 19, 3, 2, 14,
	15, 3, 2, 23, 25, 2, 98, 2, 12, 3, 2, 2, 2, 4, 28, 3, 2, 2, 2, 6, 51, 3,
	2, 2, 2, 8, 76, 3, 2, 2, 2, 10, 78, 3, 2, 2, 2, 12, 13, 5, 6, 4, 2, 13,
	14, 7, 2, 2, 3, 14, 3, 3, 2, 2, 2, 15, 16, 8, 3, 1, 2, 16, 17, 5, 8, 5,
	2, 17, 19, 7, 4, 2, 2, 18, 20, 5, 10, 6, 2, 19, 18, 3, 2, 2, 2, 19, 20,
	3, 2, 2, 2, 20, 21, 3, 2, 2, 2, 21, 22, 7, 5, 2, 2, 22, 29, 3, 2, 2, 2,
	23, 29, 7, 25, 2, 2, 24, 29, 7, 22, 2, 2, 25, 29, 7, 21, 2, 2, 26, 29,
	7, 23, 2, 2, 27, 29, 7, 24, 2, 2, 28, 15, 3, 2, 2, 2, 28, 23, 3, 2, 2,
	2, 28, 24, 3, 2, 2, 2, 28, 25, 3, 2, 2, 2, 28, 26, 3, 2, 2, 2, 28, 27,
	3, 2, 2, 2, 29, 40, 3, 2, 2, 2, 30, 31, 12, 9, 2, 2, 31, 32, 7, 8, 2, 2,
	32, 39, 5, 4, 3, 10, 33, 34, 12, 8, 2, 2, 34, 35, 7, 6, 2, 2, 35, 36, 5,
	6, 4, 2, 36, 37, 7, 7, 2, 2, 37, 39, 3, 2, 2, 2, 38, 30, 3, 2, 2, 2, 38,
	33, 3, 2, 2, 2, 39, 42, 3, 2, 2, 2, 40, 38, 3, 2, 2, 2, 40, 41, 3, 2, 2,
	2, 41, 5, 3, 2, 2, 2, 42, 40, 3, 2, 2, 2, 43, 44, 8, 4, 1, 2, 44, 52, 5,
	4, 3, 2, 45, 46, 7, 10, 2, 2, 46, 52, 5, 6, 4, 10, 47, 48, 7, 4, 2, 2,
	48, 49, 5, 6, 4, 2, 49, 50, 7, 5, 2, 2, 50, 52, 3, 2, 2, 2, 51, 43, 3,
	2, 2, 2, 51, 45, 3, 2, 2, 2, 51, 47, 3, 2, 2, 2, 52, 73, 3, 2, 2, 2, 53,
	54, 12, 9, 2, 2, 54, 55, 7, 13, 2, 2, 55, 72, 5, 6, 4, 10, 56, 57, 12,
	8, 2, 2, 57, 58, 9, 2, 2, 2, 58, 72, 5, 6, 4, 9, 59, 60, 12, 7, 2, 2, 60,
	61, 9, 3, 2, 2, 61, 72, 5, 6, 4, 8, 62, 63, 12, 6, 2, 2, 63, 64, 9, 4,
	2, 2, 64, 72, 5, 6, 4, 7, 65, 66, 12, 5, 2, 2, 66, 67, 9, 5, 2, 2, 67,
	72, 5, 6, 4, 6, 68, 69, 12, 4, 2, 2, 69, 70, 7, 20, 2, 2, 70, 72, 5, 6,
	4, 5, 71, 53, 3, 2, 2, 2, 71, 56, 3, 2, 2, 2, 71, 59, 3, 2, 2, 2, 71, 62,
	3, 2, 2, 2, 71, 65, 3, 2, 2, 2, 71, 68, 3, 2, 2, 2, 72, 75, 3, 2, 2, 2,
	73, 71, 3, 2, 2, 2, 73, 74, 3, 2, 2, 2, 74, 7, 3, 2, 2, 2, 75, 73, 3, 2,
	2, 2, 76, 77, 9, 6, 2, 2, 77, 9, 3, 2, 2, 2, 78, 83, 5, 6, 4, 2, 79, 80,
	7, 3, 2, 2, 80, 82, 5, 6, 4, 2, 81, 79, 3, 2, 2, 2, 82, 85, 3, 2, 2, 2,
	83, 81, 3, 2, 2, 2, 83, 84, 3, 2, 2, 2, 84, 11, 3, 2, 2, 2, 85, 83, 3,
	2, 2, 2, 10, 19, 28, 38, 40, 51, 71, 73, 83,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "','", "'('", "')'", "'['", "']'", "'.'", "'+'", "'-'", "'*'", "'/'",
	"'^'", "'='", "'!='", "'<='", "'<'", "'>='", "'>'", "'&'",
}
var symbolicNames = []string{
	"", "COMMA", "LPAREN", "RPAREN", "LBRACK", "RBRACK", "DOT", "PLUS", "MINUS",
	"TIMES", "DIVIDE", "EXPONENT", "EQ", "NEQ", "LTE", "LT", "GTE", "GT", "AMPERSAND",
	"DECIMAL", "STRING", "TRUE", "FALSE", "NAME", "WS", "ERROR",
}

var ruleNames = []string{
	"parse", "atom", "expression", "fnname", "parameters",
}
var decisionToDFA = make([]*antlr.DFA, len(deserializedATN.DecisionToState))

func init() {
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

type ExcellentParser struct {
	*antlr.BaseParser
}

func NewExcellentParser(input antlr.TokenStream) *ExcellentParser {
	this := new(ExcellentParser)

	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "Excellent.g4"

	return this
}

// ExcellentParser tokens.
const (
	ExcellentParserEOF       = antlr.TokenEOF
	ExcellentParserCOMMA     = 1
	ExcellentParserLPAREN    = 2
	ExcellentParserRPAREN    = 3
	ExcellentParserLBRACK    = 4
	ExcellentParserRBRACK    = 5
	ExcellentParserDOT       = 6
	ExcellentParserPLUS      = 7
	ExcellentParserMINUS     = 8
	ExcellentParserTIMES     = 9
	ExcellentParserDIVIDE    = 10
	ExcellentParserEXPONENT  = 11
	ExcellentParserEQ        = 12
	ExcellentParserNEQ       = 13
	ExcellentParserLTE       = 14
	ExcellentParserLT        = 15
	ExcellentParserGTE       = 16
	ExcellentParserGT        = 17
	ExcellentParserAMPERSAND = 18
	ExcellentParserDECIMAL   = 19
	ExcellentParserSTRING    = 20
	ExcellentParserTRUE      = 21
	ExcellentParserFALSE     = 22
	ExcellentParserNAME      = 23
	ExcellentParserWS        = 24
	ExcellentParserERROR     = 25
)

// ExcellentParser rules.
const (
	ExcellentParserRULE_parse      = 0
	ExcellentParserRULE_atom       = 1
	ExcellentParserRULE_expression = 2
	ExcellentParserRULE_fnname     = 3
	ExcellentParserRULE_parameters = 4
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
	p.RuleIndex = ExcellentParserRULE_parse
	return p
}

func (*ParseContext) IsParseContext() {}

func NewParseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParseContext {
	var p = new(ParseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ExcellentParserRULE_parse

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
	return s.GetToken(ExcellentParserEOF, 0)
}

func (s *ParseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ParseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.EnterParse(s)
	}
}

func (s *ParseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.ExitParse(s)
	}
}

func (s *ParseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExcellentVisitor:
		return t.VisitParse(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ExcellentParser) Parse() (localctx IParseContext) {
	localctx = NewParseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, ExcellentParserRULE_parse)

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
		p.Match(ExcellentParserEOF)
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
	p.RuleIndex = ExcellentParserRULE_atom
	return p
}

func (*AtomContext) IsAtomContext() {}

func NewAtomContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AtomContext {
	var p = new(AtomContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ExcellentParserRULE_atom

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

type DecimalLiteralContext struct {
	*AtomContext
}

func NewDecimalLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DecimalLiteralContext {
	var p = new(DecimalLiteralContext)

	p.AtomContext = NewEmptyAtomContext()
	p.parser = parser
	p.CopyFrom(ctx.(*AtomContext))

	return p
}

func (s *DecimalLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DecimalLiteralContext) DECIMAL() antlr.TerminalNode {
	return s.GetToken(ExcellentParserDECIMAL, 0)
}

func (s *DecimalLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.EnterDecimalLiteral(s)
	}
}

func (s *DecimalLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.ExitDecimalLiteral(s)
	}
}

func (s *DecimalLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExcellentVisitor:
		return t.VisitDecimalLiteral(s)

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

func (s *DotLookupContext) AllAtom() []IAtomContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IAtomContext)(nil)).Elem())
	var tst = make([]IAtomContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IAtomContext)
		}
	}

	return tst
}

func (s *DotLookupContext) Atom(i int) IAtomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAtomContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IAtomContext)
}

func (s *DotLookupContext) DOT() antlr.TerminalNode {
	return s.GetToken(ExcellentParserDOT, 0)
}

func (s *DotLookupContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.EnterDotLookup(s)
	}
}

func (s *DotLookupContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.ExitDotLookup(s)
	}
}

func (s *DotLookupContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExcellentVisitor:
		return t.VisitDotLookup(s)

	default:
		return t.VisitChildren(s)
	}
}

type StringLiteralContext struct {
	*AtomContext
}

func NewStringLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *StringLiteralContext {
	var p = new(StringLiteralContext)

	p.AtomContext = NewEmptyAtomContext()
	p.parser = parser
	p.CopyFrom(ctx.(*AtomContext))

	return p
}

func (s *StringLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StringLiteralContext) STRING() antlr.TerminalNode {
	return s.GetToken(ExcellentParserSTRING, 0)
}

func (s *StringLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.EnterStringLiteral(s)
	}
}

func (s *StringLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.ExitStringLiteral(s)
	}
}

func (s *StringLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExcellentVisitor:
		return t.VisitStringLiteral(s)

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

func (s *FunctionCallContext) Fnname() IFnnameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFnnameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFnnameContext)
}

func (s *FunctionCallContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(ExcellentParserLPAREN, 0)
}

func (s *FunctionCallContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(ExcellentParserRPAREN, 0)
}

func (s *FunctionCallContext) Parameters() IParametersContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IParametersContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IParametersContext)
}

func (s *FunctionCallContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.EnterFunctionCall(s)
	}
}

func (s *FunctionCallContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.ExitFunctionCall(s)
	}
}

func (s *FunctionCallContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExcellentVisitor:
		return t.VisitFunctionCall(s)

	default:
		return t.VisitChildren(s)
	}
}

type TrueContext struct {
	*AtomContext
}

func NewTrueContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TrueContext {
	var p = new(TrueContext)

	p.AtomContext = NewEmptyAtomContext()
	p.parser = parser
	p.CopyFrom(ctx.(*AtomContext))

	return p
}

func (s *TrueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TrueContext) TRUE() antlr.TerminalNode {
	return s.GetToken(ExcellentParserTRUE, 0)
}

func (s *TrueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.EnterTrue(s)
	}
}

func (s *TrueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.ExitTrue(s)
	}
}

func (s *TrueContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExcellentVisitor:
		return t.VisitTrue(s)

	default:
		return t.VisitChildren(s)
	}
}

type FalseContext struct {
	*AtomContext
}

func NewFalseContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FalseContext {
	var p = new(FalseContext)

	p.AtomContext = NewEmptyAtomContext()
	p.parser = parser
	p.CopyFrom(ctx.(*AtomContext))

	return p
}

func (s *FalseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FalseContext) FALSE() antlr.TerminalNode {
	return s.GetToken(ExcellentParserFALSE, 0)
}

func (s *FalseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.EnterFalse(s)
	}
}

func (s *FalseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.ExitFalse(s)
	}
}

func (s *FalseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExcellentVisitor:
		return t.VisitFalse(s)

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
	return s.GetToken(ExcellentParserLBRACK, 0)
}

func (s *ArrayLookupContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ArrayLookupContext) RBRACK() antlr.TerminalNode {
	return s.GetToken(ExcellentParserRBRACK, 0)
}

func (s *ArrayLookupContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.EnterArrayLookup(s)
	}
}

func (s *ArrayLookupContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.ExitArrayLookup(s)
	}
}

func (s *ArrayLookupContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExcellentVisitor:
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
	return s.GetToken(ExcellentParserNAME, 0)
}

func (s *ContextReferenceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.EnterContextReference(s)
	}
}

func (s *ContextReferenceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.ExitContextReference(s)
	}
}

func (s *ContextReferenceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExcellentVisitor:
		return t.VisitContextReference(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ExcellentParser) Atom() (localctx IAtomContext) {
	return p.atom(0)
}

func (p *ExcellentParser) atom(_p int) (localctx IAtomContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewAtomContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IAtomContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 2
	p.EnterRecursionRule(localctx, 2, ExcellentParserRULE_atom, _p)
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
	p.SetState(26)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 1, p.GetParserRuleContext()) {
	case 1:
		localctx = NewFunctionCallContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx

		{
			p.SetState(14)
			p.Fnname()
		}
		{
			p.SetState(15)
			p.Match(ExcellentParserLPAREN)
		}
		p.SetState(17)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if ((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<ExcellentParserLPAREN)|(1<<ExcellentParserMINUS)|(1<<ExcellentParserDECIMAL)|(1<<ExcellentParserSTRING)|(1<<ExcellentParserTRUE)|(1<<ExcellentParserFALSE)|(1<<ExcellentParserNAME))) != 0 {
			{
				p.SetState(16)
				p.Parameters()
			}

		}
		{
			p.SetState(19)
			p.Match(ExcellentParserRPAREN)
		}

	case 2:
		localctx = NewContextReferenceContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(21)
			p.Match(ExcellentParserNAME)
		}

	case 3:
		localctx = NewStringLiteralContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(22)
			p.Match(ExcellentParserSTRING)
		}

	case 4:
		localctx = NewDecimalLiteralContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(23)
			p.Match(ExcellentParserDECIMAL)
		}

	case 5:
		localctx = NewTrueContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(24)
			p.Match(ExcellentParserTRUE)
		}

	case 6:
		localctx = NewFalseContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(25)
			p.Match(ExcellentParserFALSE)
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(38)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 3, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(36)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 2, p.GetParserRuleContext()) {
			case 1:
				localctx = NewDotLookupContext(p, NewAtomContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExcellentParserRULE_atom)
				p.SetState(28)

				if !(p.Precpred(p.GetParserRuleContext(), 7)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 7)", ""))
				}
				{
					p.SetState(29)
					p.Match(ExcellentParserDOT)
				}
				{
					p.SetState(30)
					p.atom(8)
				}

			case 2:
				localctx = NewArrayLookupContext(p, NewAtomContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExcellentParserRULE_atom)
				p.SetState(31)

				if !(p.Precpred(p.GetParserRuleContext(), 6)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 6)", ""))
				}
				{
					p.SetState(32)
					p.Match(ExcellentParserLBRACK)
				}
				{
					p.SetState(33)
					p.expression(0)
				}
				{
					p.SetState(34)
					p.Match(ExcellentParserRBRACK)
				}

			}

		}
		p.SetState(40)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 3, p.GetParserRuleContext())
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
	p.RuleIndex = ExcellentParserRULE_expression
	return p
}

func (*ExpressionContext) IsExpressionContext() {}

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext {
	var p = new(ExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ExcellentParserRULE_expression

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
	return s.GetToken(ExcellentParserLPAREN, 0)
}

func (s *ParenthesesContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ParenthesesContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(ExcellentParserRPAREN, 0)
}

func (s *ParenthesesContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.EnterParentheses(s)
	}
}

func (s *ParenthesesContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.ExitParentheses(s)
	}
}

func (s *ParenthesesContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExcellentVisitor:
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
	return s.GetToken(ExcellentParserMINUS, 0)
}

func (s *NegationContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *NegationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.EnterNegation(s)
	}
}

func (s *NegationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.ExitNegation(s)
	}
}

func (s *NegationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExcellentVisitor:
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
	return s.GetToken(ExcellentParserLTE, 0)
}

func (s *ComparisonContext) LT() antlr.TerminalNode {
	return s.GetToken(ExcellentParserLT, 0)
}

func (s *ComparisonContext) GTE() antlr.TerminalNode {
	return s.GetToken(ExcellentParserGTE, 0)
}

func (s *ComparisonContext) GT() antlr.TerminalNode {
	return s.GetToken(ExcellentParserGT, 0)
}

func (s *ComparisonContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.EnterComparison(s)
	}
}

func (s *ComparisonContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.ExitComparison(s)
	}
}

func (s *ComparisonContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExcellentVisitor:
		return t.VisitComparison(s)

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
	return s.GetToken(ExcellentParserAMPERSAND, 0)
}

func (s *ConcatenationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.EnterConcatenation(s)
	}
}

func (s *ConcatenationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.ExitConcatenation(s)
	}
}

func (s *ConcatenationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExcellentVisitor:
		return t.VisitConcatenation(s)

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
	return s.GetToken(ExcellentParserTIMES, 0)
}

func (s *MultiplicationOrDivisionContext) DIVIDE() antlr.TerminalNode {
	return s.GetToken(ExcellentParserDIVIDE, 0)
}

func (s *MultiplicationOrDivisionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.EnterMultiplicationOrDivision(s)
	}
}

func (s *MultiplicationOrDivisionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.ExitMultiplicationOrDivision(s)
	}
}

func (s *MultiplicationOrDivisionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExcellentVisitor:
		return t.VisitMultiplicationOrDivision(s)

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
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.EnterAtomReference(s)
	}
}

func (s *AtomReferenceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.ExitAtomReference(s)
	}
}

func (s *AtomReferenceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExcellentVisitor:
		return t.VisitAtomReference(s)

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
	return s.GetToken(ExcellentParserPLUS, 0)
}

func (s *AdditionOrSubtractionContext) MINUS() antlr.TerminalNode {
	return s.GetToken(ExcellentParserMINUS, 0)
}

func (s *AdditionOrSubtractionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.EnterAdditionOrSubtraction(s)
	}
}

func (s *AdditionOrSubtractionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.ExitAdditionOrSubtraction(s)
	}
}

func (s *AdditionOrSubtractionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExcellentVisitor:
		return t.VisitAdditionOrSubtraction(s)

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
	return s.GetToken(ExcellentParserEQ, 0)
}

func (s *EqualityContext) NEQ() antlr.TerminalNode {
	return s.GetToken(ExcellentParserNEQ, 0)
}

func (s *EqualityContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.EnterEquality(s)
	}
}

func (s *EqualityContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.ExitEquality(s)
	}
}

func (s *EqualityContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExcellentVisitor:
		return t.VisitEquality(s)

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
	return s.GetToken(ExcellentParserEXPONENT, 0)
}

func (s *ExponentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.EnterExponent(s)
	}
}

func (s *ExponentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.ExitExponent(s)
	}
}

func (s *ExponentContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExcellentVisitor:
		return t.VisitExponent(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ExcellentParser) Expression() (localctx IExpressionContext) {
	return p.expression(0)
}

func (p *ExcellentParser) expression(_p int) (localctx IExpressionContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewExpressionContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IExpressionContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 4
	p.EnterRecursionRule(localctx, 4, ExcellentParserRULE_expression, _p)
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
	p.SetState(49)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case ExcellentParserDECIMAL, ExcellentParserSTRING, ExcellentParserTRUE, ExcellentParserFALSE, ExcellentParserNAME:
		localctx = NewAtomReferenceContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx

		{
			p.SetState(42)
			p.atom(0)
		}

	case ExcellentParserMINUS:
		localctx = NewNegationContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(43)
			p.Match(ExcellentParserMINUS)
		}
		{
			p.SetState(44)
			p.expression(8)
		}

	case ExcellentParserLPAREN:
		localctx = NewParenthesesContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(45)
			p.Match(ExcellentParserLPAREN)
		}
		{
			p.SetState(46)
			p.expression(0)
		}
		{
			p.SetState(47)
			p.Match(ExcellentParserRPAREN)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(71)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 6, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(69)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 5, p.GetParserRuleContext()) {
			case 1:
				localctx = NewExponentContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExcellentParserRULE_expression)
				p.SetState(51)

				if !(p.Precpred(p.GetParserRuleContext(), 7)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 7)", ""))
				}
				{
					p.SetState(52)
					p.Match(ExcellentParserEXPONENT)
				}
				{
					p.SetState(53)
					p.expression(8)
				}

			case 2:
				localctx = NewMultiplicationOrDivisionContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExcellentParserRULE_expression)
				p.SetState(54)

				if !(p.Precpred(p.GetParserRuleContext(), 6)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 6)", ""))
				}
				p.SetState(55)

				var _lt = p.GetTokenStream().LT(1)

				localctx.(*MultiplicationOrDivisionContext).op = _lt

				_la = p.GetTokenStream().LA(1)

				if !(_la == ExcellentParserTIMES || _la == ExcellentParserDIVIDE) {
					var _ri = p.GetErrorHandler().RecoverInline(p)

					localctx.(*MultiplicationOrDivisionContext).op = _ri
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
				{
					p.SetState(56)
					p.expression(7)
				}

			case 3:
				localctx = NewAdditionOrSubtractionContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExcellentParserRULE_expression)
				p.SetState(57)

				if !(p.Precpred(p.GetParserRuleContext(), 5)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 5)", ""))
				}
				p.SetState(58)

				var _lt = p.GetTokenStream().LT(1)

				localctx.(*AdditionOrSubtractionContext).op = _lt

				_la = p.GetTokenStream().LA(1)

				if !(_la == ExcellentParserPLUS || _la == ExcellentParserMINUS) {
					var _ri = p.GetErrorHandler().RecoverInline(p)

					localctx.(*AdditionOrSubtractionContext).op = _ri
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
				{
					p.SetState(59)
					p.expression(6)
				}

			case 4:
				localctx = NewComparisonContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExcellentParserRULE_expression)
				p.SetState(60)

				if !(p.Precpred(p.GetParserRuleContext(), 4)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 4)", ""))
				}
				p.SetState(61)

				var _lt = p.GetTokenStream().LT(1)

				localctx.(*ComparisonContext).op = _lt

				_la = p.GetTokenStream().LA(1)

				if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<ExcellentParserLTE)|(1<<ExcellentParserLT)|(1<<ExcellentParserGTE)|(1<<ExcellentParserGT))) != 0) {
					var _ri = p.GetErrorHandler().RecoverInline(p)

					localctx.(*ComparisonContext).op = _ri
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
				{
					p.SetState(62)
					p.expression(5)
				}

			case 5:
				localctx = NewEqualityContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExcellentParserRULE_expression)
				p.SetState(63)

				if !(p.Precpred(p.GetParserRuleContext(), 3)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 3)", ""))
				}
				p.SetState(64)

				var _lt = p.GetTokenStream().LT(1)

				localctx.(*EqualityContext).op = _lt

				_la = p.GetTokenStream().LA(1)

				if !(_la == ExcellentParserEQ || _la == ExcellentParserNEQ) {
					var _ri = p.GetErrorHandler().RecoverInline(p)

					localctx.(*EqualityContext).op = _ri
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
				{
					p.SetState(65)
					p.expression(4)
				}

			case 6:
				localctx = NewConcatenationContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ExcellentParserRULE_expression)
				p.SetState(66)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
				}
				{
					p.SetState(67)
					p.Match(ExcellentParserAMPERSAND)
				}
				{
					p.SetState(68)
					p.expression(3)
				}

			}

		}
		p.SetState(73)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 6, p.GetParserRuleContext())
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
	p.RuleIndex = ExcellentParserRULE_fnname
	return p
}

func (*FnnameContext) IsFnnameContext() {}

func NewFnnameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FnnameContext {
	var p = new(FnnameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ExcellentParserRULE_fnname

	return p
}

func (s *FnnameContext) GetParser() antlr.Parser { return s.parser }

func (s *FnnameContext) NAME() antlr.TerminalNode {
	return s.GetToken(ExcellentParserNAME, 0)
}

func (s *FnnameContext) TRUE() antlr.TerminalNode {
	return s.GetToken(ExcellentParserTRUE, 0)
}

func (s *FnnameContext) FALSE() antlr.TerminalNode {
	return s.GetToken(ExcellentParserFALSE, 0)
}

func (s *FnnameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FnnameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FnnameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.EnterFnname(s)
	}
}

func (s *FnnameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.ExitFnname(s)
	}
}

func (s *FnnameContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExcellentVisitor:
		return t.VisitFnname(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ExcellentParser) Fnname() (localctx IFnnameContext) {
	localctx = NewFnnameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, ExcellentParserRULE_fnname)
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
	p.SetState(74)
	_la = p.GetTokenStream().LA(1)

	if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<ExcellentParserTRUE)|(1<<ExcellentParserFALSE)|(1<<ExcellentParserNAME))) != 0) {
		p.GetErrorHandler().RecoverInline(p)
	} else {
		p.GetErrorHandler().ReportMatch(p)
		p.Consume()
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
	p.RuleIndex = ExcellentParserRULE_parameters
	return p
}

func (*ParametersContext) IsParametersContext() {}

func NewParametersContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParametersContext {
	var p = new(ParametersContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ExcellentParserRULE_parameters

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
	return s.GetTokens(ExcellentParserCOMMA)
}

func (s *FunctionParametersContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(ExcellentParserCOMMA, i)
}

func (s *FunctionParametersContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.EnterFunctionParameters(s)
	}
}

func (s *FunctionParametersContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExcellentListener); ok {
		listenerT.ExitFunctionParameters(s)
	}
}

func (s *FunctionParametersContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ExcellentVisitor:
		return t.VisitFunctionParameters(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ExcellentParser) Parameters() (localctx IParametersContext) {
	localctx = NewParametersContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, ExcellentParserRULE_parameters)
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
		p.SetState(76)
		p.expression(0)
	}
	p.SetState(81)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == ExcellentParserCOMMA {
		{
			p.SetState(77)
			p.Match(ExcellentParserCOMMA)
		}
		{
			p.SetState(78)
			p.expression(0)
		}

		p.SetState(83)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

func (p *ExcellentParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 1:
		var t *AtomContext = nil
		if localctx != nil {
			t = localctx.(*AtomContext)
		}
		return p.Atom_Sempred(t, predIndex)

	case 2:
		var t *ExpressionContext = nil
		if localctx != nil {
			t = localctx.(*ExpressionContext)
		}
		return p.Expression_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *ExcellentParser) Atom_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 7)

	case 1:
		return p.Precpred(p.GetParserRuleContext(), 6)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *ExcellentParser) Expression_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 2:
		return p.Precpred(p.GetParserRuleContext(), 7)

	case 3:
		return p.Precpred(p.GetParserRuleContext(), 6)

	case 4:
		return p.Precpred(p.GetParserRuleContext(), 5)

	case 5:
		return p.Precpred(p.GetParserRuleContext(), 4)

	case 6:
		return p.Precpred(p.GetParserRuleContext(), 3)

	case 7:
		return p.Precpred(p.GetParserRuleContext(), 2)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
