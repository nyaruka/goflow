// Code generated from ContactQL.g4 by ANTLR 4.10.1. DO NOT EDIT.

package gen // ContactQL
import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type ContactQLParser struct {
	*antlr.BaseParser
}

var contactqlParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	literalNames           []string
	symbolicNames          []string
	ruleNames              []string
	predictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func contactqlParserInit() {
	staticData := &contactqlParserStaticData
	staticData.literalNames = []string{
		"", "'('", "')'",
	}
	staticData.symbolicNames = []string{
		"", "LPAREN", "RPAREN", "AND", "OR", "COMPARATOR", "TEXT", "STRING",
		"WS", "ERROR",
	}
	staticData.ruleNames = []string{
		"parse", "expression", "literal",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 9, 38, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 1, 0, 1, 0, 1, 0, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 19, 8, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 5, 1, 29, 8, 1, 10, 1, 12, 1, 32,
		9, 1, 1, 2, 1, 2, 3, 2, 36, 8, 2, 1, 2, 0, 1, 2, 3, 0, 2, 4, 0, 0, 40,
		0, 6, 1, 0, 0, 0, 2, 18, 1, 0, 0, 0, 4, 35, 1, 0, 0, 0, 6, 7, 3, 2, 1,
		0, 7, 8, 5, 0, 0, 1, 8, 1, 1, 0, 0, 0, 9, 10, 6, 1, -1, 0, 10, 11, 5, 1,
		0, 0, 11, 12, 3, 2, 1, 0, 12, 13, 5, 2, 0, 0, 13, 19, 1, 0, 0, 0, 14, 15,
		5, 6, 0, 0, 15, 16, 5, 5, 0, 0, 16, 19, 3, 4, 2, 0, 17, 19, 3, 4, 2, 0,
		18, 9, 1, 0, 0, 0, 18, 14, 1, 0, 0, 0, 18, 17, 1, 0, 0, 0, 19, 30, 1, 0,
		0, 0, 20, 21, 10, 6, 0, 0, 21, 22, 5, 3, 0, 0, 22, 29, 3, 2, 1, 7, 23,
		24, 10, 5, 0, 0, 24, 29, 3, 2, 1, 6, 25, 26, 10, 4, 0, 0, 26, 27, 5, 4,
		0, 0, 27, 29, 3, 2, 1, 5, 28, 20, 1, 0, 0, 0, 28, 23, 1, 0, 0, 0, 28, 25,
		1, 0, 0, 0, 29, 32, 1, 0, 0, 0, 30, 28, 1, 0, 0, 0, 30, 31, 1, 0, 0, 0,
		31, 3, 1, 0, 0, 0, 32, 30, 1, 0, 0, 0, 33, 36, 5, 6, 0, 0, 34, 36, 5, 7,
		0, 0, 35, 33, 1, 0, 0, 0, 35, 34, 1, 0, 0, 0, 36, 5, 1, 0, 0, 0, 4, 18,
		28, 30, 35,
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

// ContactQLParserInit initializes any static state used to implement ContactQLParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewContactQLParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func ContactQLParserInit() {
	staticData := &contactqlParserStaticData
	staticData.once.Do(contactqlParserInit)
}

// NewContactQLParser produces a new parser instance for the optional input antlr.TokenStream.
func NewContactQLParser(input antlr.TokenStream) *ContactQLParser {
	ContactQLParserInit()
	this := new(ContactQLParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &contactqlParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	this.RuleNames = staticData.ruleNames
	this.LiteralNames = staticData.literalNames
	this.SymbolicNames = staticData.symbolicNames
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
	this := p
	_ = this

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
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

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
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

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

func (s *CombinationAndContext) Expression(i int) IExpressionContext {
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

func (s *CombinationImpicitAndContext) Expression(i int) IExpressionContext {
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

func (s *CombinationOrContext) Expression(i int) IExpressionContext {
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
	this := p
	_ = this

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
	this := p
	_ = this

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
	this := p
	_ = this

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
