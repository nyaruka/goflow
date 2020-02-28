// Code generated from ContactQL.g4 by ANTLR 4.8. DO NOT EDIT.

package gen // ContactQL
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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 11, 40, 4,
	2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 3, 2, 3, 2, 3, 2, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 5, 3, 21, 10, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 7, 3, 31, 10, 3, 12, 3, 14, 3, 34, 11, 3, 3,
	4, 3, 4, 5, 4, 38, 10, 4, 3, 4, 2, 3, 4, 5, 2, 4, 6, 2, 2, 2, 42, 2, 8,
	3, 2, 2, 2, 4, 20, 3, 2, 2, 2, 6, 37, 3, 2, 2, 2, 8, 9, 5, 4, 3, 2, 9,
	10, 7, 2, 2, 3, 10, 3, 3, 2, 2, 2, 11, 12, 8, 3, 1, 2, 12, 13, 7, 3, 2,
	2, 13, 14, 5, 4, 3, 2, 14, 15, 7, 4, 2, 2, 15, 21, 3, 2, 2, 2, 16, 17,
	7, 8, 2, 2, 17, 18, 7, 7, 2, 2, 18, 21, 5, 6, 4, 2, 19, 21, 5, 6, 4, 2,
	20, 11, 3, 2, 2, 2, 20, 16, 3, 2, 2, 2, 20, 19, 3, 2, 2, 2, 21, 32, 3,
	2, 2, 2, 22, 23, 12, 8, 2, 2, 23, 24, 7, 5, 2, 2, 24, 31, 5, 4, 3, 9, 25,
	26, 12, 7, 2, 2, 26, 31, 5, 4, 3, 8, 27, 28, 12, 6, 2, 2, 28, 29, 7, 6,
	2, 2, 29, 31, 5, 4, 3, 7, 30, 22, 3, 2, 2, 2, 30, 25, 3, 2, 2, 2, 30, 27,
	3, 2, 2, 2, 31, 34, 3, 2, 2, 2, 32, 30, 3, 2, 2, 2, 32, 33, 3, 2, 2, 2,
	33, 5, 3, 2, 2, 2, 34, 32, 3, 2, 2, 2, 35, 38, 7, 8, 2, 2, 36, 38, 7, 9,
	2, 2, 37, 35, 3, 2, 2, 2, 37, 36, 3, 2, 2, 2, 38, 7, 3, 2, 2, 2, 6, 20,
	30, 32, 37,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "'('", "')'",
}
var symbolicNames = []string{
	"", "LPAREN", "RPAREN", "AND", "OR", "COMPARATOR", "TEXT", "STRING", "WS",
	"ERROR",
}

var ruleNames = []string{
	"parse", "expression", "literal",
}
var decisionToDFA = make([]*antlr.DFA, len(deserializedATN.DecisionToState))

func init() {
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

type ContactQLParser struct {
	*antlr.BaseParser
}

func NewContactQLParser(input antlr.TokenStream) *ContactQLParser {
	this := new(ContactQLParser)

	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "ContactQL.g4"

	return this
}

// ContactQLParser tokens.
const (
	ContactQLParserEOF        = antlr.TokenEOF
	ContactQLParserLPAREN     = 1
	ContactQLParserRPAREN     = 2
	ContactQLParserAND        = 3
	ContactQLParserOR         = 4
	ContactQLParserCOMPARATOR = 5
	ContactQLParserTEXT       = 6
	ContactQLParserSTRING     = 7
	ContactQLParserWS         = 8
	ContactQLParserERROR      = 9
)

// ContactQLParser rules.
const (
	ContactQLParserRULE_parse      = 0
	ContactQLParserRULE_expression = 1
	ContactQLParserRULE_literal    = 2
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
	p.RuleIndex = ContactQLParserRULE_parse
	return p
}

func (*ParseContext) IsParseContext() {}

func NewParseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParseContext {
	var p = new(ParseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ContactQLParserRULE_parse

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
	return s.GetToken(ContactQLParserEOF, 0)
}

func (s *ParseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ParseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ContactQLListener); ok {
		listenerT.EnterParse(s)
	}
}

func (s *ParseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ContactQLListener); ok {
		listenerT.ExitParse(s)
	}
}

func (s *ParseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ContactQLVisitor:
		return t.VisitParse(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ContactQLParser) Parse() (localctx IParseContext) {
	localctx = NewParseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, ContactQLParserRULE_parse)

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
		p.SetState(6)
		p.expression(0)
	}
	{
		p.SetState(7)
		p.Match(ContactQLParserEOF)
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
	p.RuleIndex = ContactQLParserRULE_expression
	return p
}

func (*ExpressionContext) IsExpressionContext() {}

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext {
	var p = new(ExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ContactQLParserRULE_expression

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

type ImplicitConditionContext struct {
	*ExpressionContext
}

func NewImplicitConditionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ImplicitConditionContext {
	var p = new(ImplicitConditionContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *ImplicitConditionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ImplicitConditionContext) Literal() ILiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILiteralContext)
}

func (s *ImplicitConditionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ContactQLListener); ok {
		listenerT.EnterImplicitCondition(s)
	}
}

func (s *ImplicitConditionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ContactQLListener); ok {
		listenerT.ExitImplicitCondition(s)
	}
}

func (s *ImplicitConditionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ContactQLVisitor:
		return t.VisitImplicitCondition(s)

	default:
		return t.VisitChildren(s)
	}
}

type ConditionContext struct {
	*ExpressionContext
}

func NewConditionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ConditionContext {
	var p = new(ConditionContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *ConditionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConditionContext) TEXT() antlr.TerminalNode {
	return s.GetToken(ContactQLParserTEXT, 0)
}

func (s *ConditionContext) COMPARATOR() antlr.TerminalNode {
	return s.GetToken(ContactQLParserCOMPARATOR, 0)
}

func (s *ConditionContext) Literal() ILiteralContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILiteralContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILiteralContext)
}

func (s *ConditionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ContactQLListener); ok {
		listenerT.EnterCondition(s)
	}
}

func (s *ConditionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ContactQLListener); ok {
		listenerT.ExitCondition(s)
	}
}

func (s *ConditionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ContactQLVisitor:
		return t.VisitCondition(s)

	default:
		return t.VisitChildren(s)
	}
}

type CombinationAndContext struct {
	*ExpressionContext
}

func NewCombinationAndContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CombinationAndContext {
	var p = new(CombinationAndContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *CombinationAndContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CombinationAndContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *CombinationAndContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *CombinationAndContext) AND() antlr.TerminalNode {
	return s.GetToken(ContactQLParserAND, 0)
}

func (s *CombinationAndContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ContactQLListener); ok {
		listenerT.EnterCombinationAnd(s)
	}
}

func (s *CombinationAndContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ContactQLListener); ok {
		listenerT.ExitCombinationAnd(s)
	}
}

func (s *CombinationAndContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ContactQLVisitor:
		return t.VisitCombinationAnd(s)

	default:
		return t.VisitChildren(s)
	}
}

type CombinationImpicitAndContext struct {
	*ExpressionContext
}

func NewCombinationImpicitAndContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CombinationImpicitAndContext {
	var p = new(CombinationImpicitAndContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *CombinationImpicitAndContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CombinationImpicitAndContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *CombinationImpicitAndContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *CombinationImpicitAndContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ContactQLListener); ok {
		listenerT.EnterCombinationImpicitAnd(s)
	}
}

func (s *CombinationImpicitAndContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ContactQLListener); ok {
		listenerT.ExitCombinationImpicitAnd(s)
	}
}

func (s *CombinationImpicitAndContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ContactQLVisitor:
		return t.VisitCombinationImpicitAnd(s)

	default:
		return t.VisitChildren(s)
	}
}

type CombinationOrContext struct {
	*ExpressionContext
}

func NewCombinationOrContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CombinationOrContext {
	var p = new(CombinationOrContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *CombinationOrContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CombinationOrContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *CombinationOrContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *CombinationOrContext) OR() antlr.TerminalNode {
	return s.GetToken(ContactQLParserOR, 0)
}

func (s *CombinationOrContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ContactQLListener); ok {
		listenerT.EnterCombinationOr(s)
	}
}

func (s *CombinationOrContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ContactQLListener); ok {
		listenerT.ExitCombinationOr(s)
	}
}

func (s *CombinationOrContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ContactQLVisitor:
		return t.VisitCombinationOr(s)

	default:
		return t.VisitChildren(s)
	}
}

type ExpressionGroupingContext struct {
	*ExpressionContext
}

func NewExpressionGroupingContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ExpressionGroupingContext {
	var p = new(ExpressionGroupingContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *ExpressionGroupingContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionGroupingContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(ContactQLParserLPAREN, 0)
}

func (s *ExpressionGroupingContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ExpressionGroupingContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(ContactQLParserRPAREN, 0)
}

func (s *ExpressionGroupingContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ContactQLListener); ok {
		listenerT.EnterExpressionGrouping(s)
	}
}

func (s *ExpressionGroupingContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ContactQLListener); ok {
		listenerT.ExitExpressionGrouping(s)
	}
}

func (s *ExpressionGroupingContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ContactQLVisitor:
		return t.VisitExpressionGrouping(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ContactQLParser) Expression() (localctx IExpressionContext) {
	return p.expression(0)
}

func (p *ContactQLParser) expression(_p int) (localctx IExpressionContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewExpressionContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IExpressionContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 2
	p.EnterRecursionRule(localctx, 2, ContactQLParserRULE_expression, _p)

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
	p.SetState(18)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext()) {
	case 1:
		localctx = NewExpressionGroupingContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx

		{
			p.SetState(10)
			p.Match(ContactQLParserLPAREN)
		}
		{
			p.SetState(11)
			p.expression(0)
		}
		{
			p.SetState(12)
			p.Match(ContactQLParserRPAREN)
		}

	case 2:
		localctx = NewConditionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(14)
			p.Match(ContactQLParserTEXT)
		}
		{
			p.SetState(15)
			p.Match(ContactQLParserCOMPARATOR)
		}
		{
			p.SetState(16)
			p.Literal()
		}

	case 3:
		localctx = NewImplicitConditionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(17)
			p.Literal()
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(30)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 2, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(28)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 1, p.GetParserRuleContext()) {
			case 1:
				localctx = NewCombinationAndContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ContactQLParserRULE_expression)
				p.SetState(20)

				if !(p.Precpred(p.GetParserRuleContext(), 6)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 6)", ""))
				}
				{
					p.SetState(21)
					p.Match(ContactQLParserAND)
				}
				{
					p.SetState(22)
					p.expression(7)
				}

			case 2:
				localctx = NewCombinationImpicitAndContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ContactQLParserRULE_expression)
				p.SetState(23)

				if !(p.Precpred(p.GetParserRuleContext(), 5)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 5)", ""))
				}
				{
					p.SetState(24)
					p.expression(6)
				}

			case 3:
				localctx = NewCombinationOrContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, ContactQLParserRULE_expression)
				p.SetState(25)

				if !(p.Precpred(p.GetParserRuleContext(), 4)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 4)", ""))
				}
				{
					p.SetState(26)
					p.Match(ContactQLParserOR)
				}
				{
					p.SetState(27)
					p.expression(5)
				}

			}

		}
		p.SetState(32)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 2, p.GetParserRuleContext())
	}

	return localctx
}

// ILiteralContext is an interface to support dynamic dispatch.
type ILiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLiteralContext differentiates from other interfaces.
	IsLiteralContext()
}

type LiteralContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLiteralContext() *LiteralContext {
	var p = new(LiteralContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ContactQLParserRULE_literal
	return p
}

func (*LiteralContext) IsLiteralContext() {}

func NewLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LiteralContext {
	var p = new(LiteralContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ContactQLParserRULE_literal

	return p
}

func (s *LiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *LiteralContext) CopyFrom(ctx *LiteralContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *LiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type StringLiteralContext struct {
	*LiteralContext
}

func NewStringLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *StringLiteralContext {
	var p = new(StringLiteralContext)

	p.LiteralContext = NewEmptyLiteralContext()
	p.parser = parser
	p.CopyFrom(ctx.(*LiteralContext))

	return p
}

func (s *StringLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StringLiteralContext) STRING() antlr.TerminalNode {
	return s.GetToken(ContactQLParserSTRING, 0)
}

func (s *StringLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ContactQLListener); ok {
		listenerT.EnterStringLiteral(s)
	}
}

func (s *StringLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ContactQLListener); ok {
		listenerT.ExitStringLiteral(s)
	}
}

func (s *StringLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ContactQLVisitor:
		return t.VisitStringLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

type TextLiteralContext struct {
	*LiteralContext
}

func NewTextLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TextLiteralContext {
	var p = new(TextLiteralContext)

	p.LiteralContext = NewEmptyLiteralContext()
	p.parser = parser
	p.CopyFrom(ctx.(*LiteralContext))

	return p
}

func (s *TextLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TextLiteralContext) TEXT() antlr.TerminalNode {
	return s.GetToken(ContactQLParserTEXT, 0)
}

func (s *TextLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ContactQLListener); ok {
		listenerT.EnterTextLiteral(s)
	}
}

func (s *TextLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ContactQLListener); ok {
		listenerT.ExitTextLiteral(s)
	}
}

func (s *TextLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ContactQLVisitor:
		return t.VisitTextLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ContactQLParser) Literal() (localctx ILiteralContext) {
	localctx = NewLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, ContactQLParserRULE_literal)

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

	p.SetState(35)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case ContactQLParserTEXT:
		localctx = NewTextLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(33)
			p.Match(ContactQLParserTEXT)
		}

	case ContactQLParserSTRING:
		localctx = NewStringLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(34)
			p.Match(ContactQLParserSTRING)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

func (p *ContactQLParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
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

func (p *ContactQLParser) Expression_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 6)

	case 1:
		return p.Precpred(p.GetParserRuleContext(), 5)

	case 2:
		return p.Precpred(p.GetParserRuleContext(), 4)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
