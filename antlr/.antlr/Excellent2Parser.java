// Generated from /Users/rowan/Nyaruka/go/src/github.com/nyaruka/goflow/antlr/Excellent2.g4 by ANTLR 4.7.1
import org.antlr.v4.runtime.atn.*;
import org.antlr.v4.runtime.dfa.DFA;
import org.antlr.v4.runtime.*;
import org.antlr.v4.runtime.misc.*;
import org.antlr.v4.runtime.tree.*;
import java.util.List;
import java.util.Iterator;
import java.util.ArrayList;

@SuppressWarnings({"all", "warnings", "unchecked", "unused", "cast"})
public class Excellent2Parser extends Parser {
	static { RuntimeMetaData.checkVersion("4.7.1", RuntimeMetaData.VERSION); }

	protected static final DFA[] _decisionToDFA;
	protected static final PredictionContextCache _sharedContextCache =
		new PredictionContextCache();
	public static final int
		COMMA=1, LPAREN=2, RPAREN=3, LBRACK=4, RBRACK=5, DOT=6, PLUS=7, MINUS=8, 
		TIMES=9, DIVIDE=10, EXPONENT=11, EQ=12, NEQ=13, LTE=14, LT=15, GTE=16, 
		GT=17, AMPERSAND=18, DECIMAL=19, STRING=20, TRUE=21, FALSE=22, NULL=23, 
		NAME=24, WS=25, ERROR=26;
	public static final int
		RULE_parse = 0, RULE_atom = 1, RULE_expression = 2, RULE_fnname = 3, RULE_parameters = 4;
	public static final String[] ruleNames = {
		"parse", "atom", "expression", "fnname", "parameters"
	};

	private static final String[] _LITERAL_NAMES = {
		null, "','", "'('", "')'", "'['", "']'", "'.'", "'+'", "'-'", "'*'", "'/'", 
		"'^'", "'='", "'!='", "'<='", "'<'", "'>='", "'>'", "'&'"
	};
	private static final String[] _SYMBOLIC_NAMES = {
		null, "COMMA", "LPAREN", "RPAREN", "LBRACK", "RBRACK", "DOT", "PLUS", 
		"MINUS", "TIMES", "DIVIDE", "EXPONENT", "EQ", "NEQ", "LTE", "LT", "GTE", 
		"GT", "AMPERSAND", "DECIMAL", "STRING", "TRUE", "FALSE", "NULL", "NAME", 
		"WS", "ERROR"
	};
	public static final Vocabulary VOCABULARY = new VocabularyImpl(_LITERAL_NAMES, _SYMBOLIC_NAMES);

	/**
	 * @deprecated Use {@link #VOCABULARY} instead.
	 */
	@Deprecated
	public static final String[] tokenNames;
	static {
		tokenNames = new String[_SYMBOLIC_NAMES.length];
		for (int i = 0; i < tokenNames.length; i++) {
			tokenNames[i] = VOCABULARY.getLiteralName(i);
			if (tokenNames[i] == null) {
				tokenNames[i] = VOCABULARY.getSymbolicName(i);
			}

			if (tokenNames[i] == null) {
				tokenNames[i] = "<INVALID>";
			}
		}
	}

	@Override
	@Deprecated
	public String[] getTokenNames() {
		return tokenNames;
	}

	@Override

	public Vocabulary getVocabulary() {
		return VOCABULARY;
	}

	@Override
	public String getGrammarFileName() { return "Excellent2.g4"; }

	@Override
	public String[] getRuleNames() { return ruleNames; }

	@Override
	public String getSerializedATN() { return _serializedATN; }

	@Override
	public ATN getATN() { return _ATN; }

	public Excellent2Parser(TokenStream input) {
		super(input);
		_interp = new ParserATNSimulator(this,_ATN,_decisionToDFA,_sharedContextCache);
	}
	public static class ParseContext extends ParserRuleContext {
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode EOF() { return getToken(Excellent2Parser.EOF, 0); }
		public ParseContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_parse; }
	}

	public final ParseContext parse() throws RecognitionException {
		ParseContext _localctx = new ParseContext(_ctx, getState());
		enterRule(_localctx, 0, RULE_parse);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(10);
			expression(0);
			setState(11);
			match(EOF);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class AtomContext extends ParserRuleContext {
		public AtomContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_atom; }
	 
		public AtomContext() { }
		public void copyFrom(AtomContext ctx) {
			super.copyFrom(ctx);
		}
	}
	public static class DecimalLiteralContext extends AtomContext {
		public TerminalNode DECIMAL() { return getToken(Excellent2Parser.DECIMAL, 0); }
		public DecimalLiteralContext(AtomContext ctx) { copyFrom(ctx); }
	}
	public static class DotLookupContext extends AtomContext {
		public List<AtomContext> atom() {
			return getRuleContexts(AtomContext.class);
		}
		public AtomContext atom(int i) {
			return getRuleContext(AtomContext.class,i);
		}
		public TerminalNode DOT() { return getToken(Excellent2Parser.DOT, 0); }
		public DotLookupContext(AtomContext ctx) { copyFrom(ctx); }
	}
	public static class NullContext extends AtomContext {
		public TerminalNode NULL() { return getToken(Excellent2Parser.NULL, 0); }
		public NullContext(AtomContext ctx) { copyFrom(ctx); }
	}
	public static class StringLiteralContext extends AtomContext {
		public TerminalNode STRING() { return getToken(Excellent2Parser.STRING, 0); }
		public StringLiteralContext(AtomContext ctx) { copyFrom(ctx); }
	}
	public static class FunctionCallContext extends AtomContext {
		public FnnameContext fnname() {
			return getRuleContext(FnnameContext.class,0);
		}
		public TerminalNode LPAREN() { return getToken(Excellent2Parser.LPAREN, 0); }
		public TerminalNode RPAREN() { return getToken(Excellent2Parser.RPAREN, 0); }
		public ParametersContext parameters() {
			return getRuleContext(ParametersContext.class,0);
		}
		public FunctionCallContext(AtomContext ctx) { copyFrom(ctx); }
	}
	public static class TrueContext extends AtomContext {
		public TerminalNode TRUE() { return getToken(Excellent2Parser.TRUE, 0); }
		public TrueContext(AtomContext ctx) { copyFrom(ctx); }
	}
	public static class FalseContext extends AtomContext {
		public TerminalNode FALSE() { return getToken(Excellent2Parser.FALSE, 0); }
		public FalseContext(AtomContext ctx) { copyFrom(ctx); }
	}
	public static class ArrayLookupContext extends AtomContext {
		public AtomContext atom() {
			return getRuleContext(AtomContext.class,0);
		}
		public TerminalNode LBRACK() { return getToken(Excellent2Parser.LBRACK, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode RBRACK() { return getToken(Excellent2Parser.RBRACK, 0); }
		public ArrayLookupContext(AtomContext ctx) { copyFrom(ctx); }
	}
	public static class ContextReferenceContext extends AtomContext {
		public TerminalNode NAME() { return getToken(Excellent2Parser.NAME, 0); }
		public ContextReferenceContext(AtomContext ctx) { copyFrom(ctx); }
	}

	public final AtomContext atom() throws RecognitionException {
		return atom(0);
	}

	private AtomContext atom(int _p) throws RecognitionException {
		ParserRuleContext _parentctx = _ctx;
		int _parentState = getState();
		AtomContext _localctx = new AtomContext(_ctx, _parentState);
		AtomContext _prevctx = _localctx;
		int _startState = 2;
		enterRecursionRule(_localctx, 2, RULE_atom, _p);
		int _la;
		try {
			int _alt;
			enterOuterAlt(_localctx, 1);
			{
			setState(27);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,1,_ctx) ) {
			case 1:
				{
				_localctx = new FunctionCallContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;

				setState(14);
				fnname();
				setState(15);
				match(LPAREN);
				setState(17);
				_errHandler.sync(this);
				_la = _input.LA(1);
				if ((((_la) & ~0x3f) == 0 && ((1L << _la) & ((1L << LPAREN) | (1L << MINUS) | (1L << DECIMAL) | (1L << STRING) | (1L << TRUE) | (1L << FALSE) | (1L << NULL) | (1L << NAME))) != 0)) {
					{
					setState(16);
					parameters();
					}
				}

				setState(19);
				match(RPAREN);
				}
				break;
			case 2:
				{
				_localctx = new ContextReferenceContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;
				setState(21);
				match(NAME);
				}
				break;
			case 3:
				{
				_localctx = new StringLiteralContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;
				setState(22);
				match(STRING);
				}
				break;
			case 4:
				{
				_localctx = new DecimalLiteralContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;
				setState(23);
				match(DECIMAL);
				}
				break;
			case 5:
				{
				_localctx = new TrueContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;
				setState(24);
				match(TRUE);
				}
				break;
			case 6:
				{
				_localctx = new FalseContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;
				setState(25);
				match(FALSE);
				}
				break;
			case 7:
				{
				_localctx = new NullContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;
				setState(26);
				match(NULL);
				}
				break;
			}
			_ctx.stop = _input.LT(-1);
			setState(39);
			_errHandler.sync(this);
			_alt = getInterpreter().adaptivePredict(_input,3,_ctx);
			while ( _alt!=2 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER ) {
				if ( _alt==1 ) {
					if ( _parseListeners!=null ) triggerExitRuleEvent();
					_prevctx = _localctx;
					{
					setState(37);
					_errHandler.sync(this);
					switch ( getInterpreter().adaptivePredict(_input,2,_ctx) ) {
					case 1:
						{
						_localctx = new DotLookupContext(new AtomContext(_parentctx, _parentState));
						pushNewRecursionContext(_localctx, _startState, RULE_atom);
						setState(29);
						if (!(precpred(_ctx, 8))) throw new FailedPredicateException(this, "precpred(_ctx, 8)");
						setState(30);
						match(DOT);
						setState(31);
						atom(9);
						}
						break;
					case 2:
						{
						_localctx = new ArrayLookupContext(new AtomContext(_parentctx, _parentState));
						pushNewRecursionContext(_localctx, _startState, RULE_atom);
						setState(32);
						if (!(precpred(_ctx, 7))) throw new FailedPredicateException(this, "precpred(_ctx, 7)");
						setState(33);
						match(LBRACK);
						setState(34);
						expression(0);
						setState(35);
						match(RBRACK);
						}
						break;
					}
					} 
				}
				setState(41);
				_errHandler.sync(this);
				_alt = getInterpreter().adaptivePredict(_input,3,_ctx);
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			unrollRecursionContexts(_parentctx);
		}
		return _localctx;
	}

	public static class ExpressionContext extends ParserRuleContext {
		public ExpressionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_expression; }
	 
		public ExpressionContext() { }
		public void copyFrom(ExpressionContext ctx) {
			super.copyFrom(ctx);
		}
	}
	public static class ParenthesesContext extends ExpressionContext {
		public TerminalNode LPAREN() { return getToken(Excellent2Parser.LPAREN, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode RPAREN() { return getToken(Excellent2Parser.RPAREN, 0); }
		public ParenthesesContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	public static class NegationContext extends ExpressionContext {
		public TerminalNode MINUS() { return getToken(Excellent2Parser.MINUS, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public NegationContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	public static class ComparisonContext extends ExpressionContext {
		public Token op;
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public TerminalNode LTE() { return getToken(Excellent2Parser.LTE, 0); }
		public TerminalNode LT() { return getToken(Excellent2Parser.LT, 0); }
		public TerminalNode GTE() { return getToken(Excellent2Parser.GTE, 0); }
		public TerminalNode GT() { return getToken(Excellent2Parser.GT, 0); }
		public ComparisonContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	public static class ConcatenationContext extends ExpressionContext {
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public TerminalNode AMPERSAND() { return getToken(Excellent2Parser.AMPERSAND, 0); }
		public ConcatenationContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	public static class MultiplicationOrDivisionContext extends ExpressionContext {
		public Token op;
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public TerminalNode TIMES() { return getToken(Excellent2Parser.TIMES, 0); }
		public TerminalNode DIVIDE() { return getToken(Excellent2Parser.DIVIDE, 0); }
		public MultiplicationOrDivisionContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	public static class AtomReferenceContext extends ExpressionContext {
		public AtomContext atom() {
			return getRuleContext(AtomContext.class,0);
		}
		public AtomReferenceContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	public static class AdditionOrSubtractionContext extends ExpressionContext {
		public Token op;
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public TerminalNode PLUS() { return getToken(Excellent2Parser.PLUS, 0); }
		public TerminalNode MINUS() { return getToken(Excellent2Parser.MINUS, 0); }
		public AdditionOrSubtractionContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	public static class EqualityContext extends ExpressionContext {
		public Token op;
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public TerminalNode EQ() { return getToken(Excellent2Parser.EQ, 0); }
		public TerminalNode NEQ() { return getToken(Excellent2Parser.NEQ, 0); }
		public EqualityContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	public static class ExponentContext extends ExpressionContext {
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public TerminalNode EXPONENT() { return getToken(Excellent2Parser.EXPONENT, 0); }
		public ExponentContext(ExpressionContext ctx) { copyFrom(ctx); }
	}

	public final ExpressionContext expression() throws RecognitionException {
		return expression(0);
	}

	private ExpressionContext expression(int _p) throws RecognitionException {
		ParserRuleContext _parentctx = _ctx;
		int _parentState = getState();
		ExpressionContext _localctx = new ExpressionContext(_ctx, _parentState);
		ExpressionContext _prevctx = _localctx;
		int _startState = 4;
		enterRecursionRule(_localctx, 4, RULE_expression, _p);
		int _la;
		try {
			int _alt;
			enterOuterAlt(_localctx, 1);
			{
			setState(50);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case DECIMAL:
			case STRING:
			case TRUE:
			case FALSE:
			case NULL:
			case NAME:
				{
				_localctx = new AtomReferenceContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;

				setState(43);
				atom(0);
				}
				break;
			case MINUS:
				{
				_localctx = new NegationContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;
				setState(44);
				match(MINUS);
				setState(45);
				expression(8);
				}
				break;
			case LPAREN:
				{
				_localctx = new ParenthesesContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;
				setState(46);
				match(LPAREN);
				setState(47);
				expression(0);
				setState(48);
				match(RPAREN);
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
			_ctx.stop = _input.LT(-1);
			setState(72);
			_errHandler.sync(this);
			_alt = getInterpreter().adaptivePredict(_input,6,_ctx);
			while ( _alt!=2 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER ) {
				if ( _alt==1 ) {
					if ( _parseListeners!=null ) triggerExitRuleEvent();
					_prevctx = _localctx;
					{
					setState(70);
					_errHandler.sync(this);
					switch ( getInterpreter().adaptivePredict(_input,5,_ctx) ) {
					case 1:
						{
						_localctx = new ExponentContext(new ExpressionContext(_parentctx, _parentState));
						pushNewRecursionContext(_localctx, _startState, RULE_expression);
						setState(52);
						if (!(precpred(_ctx, 7))) throw new FailedPredicateException(this, "precpred(_ctx, 7)");
						setState(53);
						match(EXPONENT);
						setState(54);
						expression(8);
						}
						break;
					case 2:
						{
						_localctx = new MultiplicationOrDivisionContext(new ExpressionContext(_parentctx, _parentState));
						pushNewRecursionContext(_localctx, _startState, RULE_expression);
						setState(55);
						if (!(precpred(_ctx, 6))) throw new FailedPredicateException(this, "precpred(_ctx, 6)");
						setState(56);
						((MultiplicationOrDivisionContext)_localctx).op = _input.LT(1);
						_la = _input.LA(1);
						if ( !(_la==TIMES || _la==DIVIDE) ) {
							((MultiplicationOrDivisionContext)_localctx).op = (Token)_errHandler.recoverInline(this);
						}
						else {
							if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
							_errHandler.reportMatch(this);
							consume();
						}
						setState(57);
						expression(7);
						}
						break;
					case 3:
						{
						_localctx = new AdditionOrSubtractionContext(new ExpressionContext(_parentctx, _parentState));
						pushNewRecursionContext(_localctx, _startState, RULE_expression);
						setState(58);
						if (!(precpred(_ctx, 5))) throw new FailedPredicateException(this, "precpred(_ctx, 5)");
						setState(59);
						((AdditionOrSubtractionContext)_localctx).op = _input.LT(1);
						_la = _input.LA(1);
						if ( !(_la==PLUS || _la==MINUS) ) {
							((AdditionOrSubtractionContext)_localctx).op = (Token)_errHandler.recoverInline(this);
						}
						else {
							if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
							_errHandler.reportMatch(this);
							consume();
						}
						setState(60);
						expression(6);
						}
						break;
					case 4:
						{
						_localctx = new ComparisonContext(new ExpressionContext(_parentctx, _parentState));
						pushNewRecursionContext(_localctx, _startState, RULE_expression);
						setState(61);
						if (!(precpred(_ctx, 4))) throw new FailedPredicateException(this, "precpred(_ctx, 4)");
						setState(62);
						((ComparisonContext)_localctx).op = _input.LT(1);
						_la = _input.LA(1);
						if ( !((((_la) & ~0x3f) == 0 && ((1L << _la) & ((1L << LTE) | (1L << LT) | (1L << GTE) | (1L << GT))) != 0)) ) {
							((ComparisonContext)_localctx).op = (Token)_errHandler.recoverInline(this);
						}
						else {
							if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
							_errHandler.reportMatch(this);
							consume();
						}
						setState(63);
						expression(5);
						}
						break;
					case 5:
						{
						_localctx = new EqualityContext(new ExpressionContext(_parentctx, _parentState));
						pushNewRecursionContext(_localctx, _startState, RULE_expression);
						setState(64);
						if (!(precpred(_ctx, 3))) throw new FailedPredicateException(this, "precpred(_ctx, 3)");
						setState(65);
						((EqualityContext)_localctx).op = _input.LT(1);
						_la = _input.LA(1);
						if ( !(_la==EQ || _la==NEQ) ) {
							((EqualityContext)_localctx).op = (Token)_errHandler.recoverInline(this);
						}
						else {
							if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
							_errHandler.reportMatch(this);
							consume();
						}
						setState(66);
						expression(4);
						}
						break;
					case 6:
						{
						_localctx = new ConcatenationContext(new ExpressionContext(_parentctx, _parentState));
						pushNewRecursionContext(_localctx, _startState, RULE_expression);
						setState(67);
						if (!(precpred(_ctx, 2))) throw new FailedPredicateException(this, "precpred(_ctx, 2)");
						setState(68);
						match(AMPERSAND);
						setState(69);
						expression(3);
						}
						break;
					}
					} 
				}
				setState(74);
				_errHandler.sync(this);
				_alt = getInterpreter().adaptivePredict(_input,6,_ctx);
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			unrollRecursionContexts(_parentctx);
		}
		return _localctx;
	}

	public static class FnnameContext extends ParserRuleContext {
		public TerminalNode NAME() { return getToken(Excellent2Parser.NAME, 0); }
		public TerminalNode TRUE() { return getToken(Excellent2Parser.TRUE, 0); }
		public TerminalNode FALSE() { return getToken(Excellent2Parser.FALSE, 0); }
		public FnnameContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_fnname; }
	}

	public final FnnameContext fnname() throws RecognitionException {
		FnnameContext _localctx = new FnnameContext(_ctx, getState());
		enterRule(_localctx, 6, RULE_fnname);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(75);
			_la = _input.LA(1);
			if ( !((((_la) & ~0x3f) == 0 && ((1L << _la) & ((1L << TRUE) | (1L << FALSE) | (1L << NAME))) != 0)) ) {
			_errHandler.recoverInline(this);
			}
			else {
				if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
				_errHandler.reportMatch(this);
				consume();
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class ParametersContext extends ParserRuleContext {
		public ParametersContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_parameters; }
	 
		public ParametersContext() { }
		public void copyFrom(ParametersContext ctx) {
			super.copyFrom(ctx);
		}
	}
	public static class FunctionParametersContext extends ParametersContext {
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public List<TerminalNode> COMMA() { return getTokens(Excellent2Parser.COMMA); }
		public TerminalNode COMMA(int i) {
			return getToken(Excellent2Parser.COMMA, i);
		}
		public FunctionParametersContext(ParametersContext ctx) { copyFrom(ctx); }
	}

	public final ParametersContext parameters() throws RecognitionException {
		ParametersContext _localctx = new ParametersContext(_ctx, getState());
		enterRule(_localctx, 8, RULE_parameters);
		int _la;
		try {
			_localctx = new FunctionParametersContext(_localctx);
			enterOuterAlt(_localctx, 1);
			{
			setState(77);
			expression(0);
			setState(82);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==COMMA) {
				{
				{
				setState(78);
				match(COMMA);
				setState(79);
				expression(0);
				}
				}
				setState(84);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public boolean sempred(RuleContext _localctx, int ruleIndex, int predIndex) {
		switch (ruleIndex) {
		case 1:
			return atom_sempred((AtomContext)_localctx, predIndex);
		case 2:
			return expression_sempred((ExpressionContext)_localctx, predIndex);
		}
		return true;
	}
	private boolean atom_sempred(AtomContext _localctx, int predIndex) {
		switch (predIndex) {
		case 0:
			return precpred(_ctx, 8);
		case 1:
			return precpred(_ctx, 7);
		}
		return true;
	}
	private boolean expression_sempred(ExpressionContext _localctx, int predIndex) {
		switch (predIndex) {
		case 2:
			return precpred(_ctx, 7);
		case 3:
			return precpred(_ctx, 6);
		case 4:
			return precpred(_ctx, 5);
		case 5:
			return precpred(_ctx, 4);
		case 6:
			return precpred(_ctx, 3);
		case 7:
			return precpred(_ctx, 2);
		}
		return true;
	}

	public static final String _serializedATN =
		"\3\u608b\ua72a\u8133\ub9ed\u417c\u3be7\u7786\u5964\3\34X\4\2\t\2\4\3\t"+
		"\3\4\4\t\4\4\5\t\5\4\6\t\6\3\2\3\2\3\2\3\3\3\3\3\3\3\3\5\3\24\n\3\3\3"+
		"\3\3\3\3\3\3\3\3\3\3\3\3\3\3\5\3\36\n\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3"+
		"\3\7\3(\n\3\f\3\16\3+\13\3\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\5\4\65\n\4"+
		"\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3"+
		"\4\7\4I\n\4\f\4\16\4L\13\4\3\5\3\5\3\6\3\6\3\6\7\6S\n\6\f\6\16\6V\13\6"+
		"\3\6\2\4\4\6\7\2\4\6\b\n\2\7\3\2\13\f\3\2\t\n\3\2\20\23\3\2\16\17\4\2"+
		"\27\30\32\32\2d\2\f\3\2\2\2\4\35\3\2\2\2\6\64\3\2\2\2\bM\3\2\2\2\nO\3"+
		"\2\2\2\f\r\5\6\4\2\r\16\7\2\2\3\16\3\3\2\2\2\17\20\b\3\1\2\20\21\5\b\5"+
		"\2\21\23\7\4\2\2\22\24\5\n\6\2\23\22\3\2\2\2\23\24\3\2\2\2\24\25\3\2\2"+
		"\2\25\26\7\5\2\2\26\36\3\2\2\2\27\36\7\32\2\2\30\36\7\26\2\2\31\36\7\25"+
		"\2\2\32\36\7\27\2\2\33\36\7\30\2\2\34\36\7\31\2\2\35\17\3\2\2\2\35\27"+
		"\3\2\2\2\35\30\3\2\2\2\35\31\3\2\2\2\35\32\3\2\2\2\35\33\3\2\2\2\35\34"+
		"\3\2\2\2\36)\3\2\2\2\37 \f\n\2\2 !\7\b\2\2!(\5\4\3\13\"#\f\t\2\2#$\7\6"+
		"\2\2$%\5\6\4\2%&\7\7\2\2&(\3\2\2\2\'\37\3\2\2\2\'\"\3\2\2\2(+\3\2\2\2"+
		")\'\3\2\2\2)*\3\2\2\2*\5\3\2\2\2+)\3\2\2\2,-\b\4\1\2-\65\5\4\3\2./\7\n"+
		"\2\2/\65\5\6\4\n\60\61\7\4\2\2\61\62\5\6\4\2\62\63\7\5\2\2\63\65\3\2\2"+
		"\2\64,\3\2\2\2\64.\3\2\2\2\64\60\3\2\2\2\65J\3\2\2\2\66\67\f\t\2\2\67"+
		"8\7\r\2\28I\5\6\4\n9:\f\b\2\2:;\t\2\2\2;I\5\6\4\t<=\f\7\2\2=>\t\3\2\2"+
		">I\5\6\4\b?@\f\6\2\2@A\t\4\2\2AI\5\6\4\7BC\f\5\2\2CD\t\5\2\2DI\5\6\4\6"+
		"EF\f\4\2\2FG\7\24\2\2GI\5\6\4\5H\66\3\2\2\2H9\3\2\2\2H<\3\2\2\2H?\3\2"+
		"\2\2HB\3\2\2\2HE\3\2\2\2IL\3\2\2\2JH\3\2\2\2JK\3\2\2\2K\7\3\2\2\2LJ\3"+
		"\2\2\2MN\t\6\2\2N\t\3\2\2\2OT\5\6\4\2PQ\7\3\2\2QS\5\6\4\2RP\3\2\2\2SV"+
		"\3\2\2\2TR\3\2\2\2TU\3\2\2\2U\13\3\2\2\2VT\3\2\2\2\n\23\35\')\64HJT";
	public static final ATN _ATN =
		new ATNDeserializer().deserialize(_serializedATN.toCharArray());
	static {
		_decisionToDFA = new DFA[_ATN.getNumberOfDecisions()];
		for (int i = 0; i < _ATN.getNumberOfDecisions(); i++) {
			_decisionToDFA[i] = new DFA(_ATN.getDecisionState(i), i);
		}
	}
}