// Generated from src/github.com/nyaruka/goflow/excellent/gen/Excellent.g4 by ANTLR 4.7.

package gen

import (
	"fmt"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = unicode.IsLetter

var serializedLexerAtn = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 2, 27, 143,
	8, 1, 4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7,
	9, 7, 4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12,
	4, 13, 9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4,
	18, 9, 18, 4, 19, 9, 19, 4, 20, 9, 20, 4, 21, 9, 21, 4, 22, 9, 22, 4, 23,
	9, 23, 4, 24, 9, 24, 4, 25, 9, 25, 4, 26, 9, 26, 3, 2, 3, 2, 3, 3, 3, 3,
	3, 4, 3, 4, 3, 5, 3, 5, 3, 6, 3, 6, 3, 7, 3, 7, 3, 8, 3, 8, 3, 9, 3, 9,
	3, 10, 3, 10, 3, 11, 3, 11, 3, 12, 3, 12, 3, 13, 3, 13, 3, 14, 3, 14, 3,
	14, 3, 15, 3, 15, 3, 15, 3, 16, 3, 16, 3, 17, 3, 17, 3, 17, 3, 18, 3, 18,
	3, 19, 3, 19, 3, 20, 6, 20, 94, 10, 20, 13, 20, 14, 20, 95, 3, 20, 3, 20,
	6, 20, 100, 10, 20, 13, 20, 14, 20, 101, 5, 20, 104, 10, 20, 3, 21, 3,
	21, 3, 21, 3, 21, 7, 21, 110, 10, 21, 12, 21, 14, 21, 113, 11, 21, 3, 21,
	3, 21, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 23, 3, 23, 3, 23, 3, 23, 3,
	23, 3, 23, 3, 24, 3, 24, 7, 24, 130, 10, 24, 12, 24, 14, 24, 133, 11, 24,
	3, 25, 6, 25, 136, 10, 25, 13, 25, 14, 25, 137, 3, 25, 3, 25, 3, 26, 3,
	26, 2, 2, 27, 3, 3, 5, 4, 7, 5, 9, 6, 11, 7, 13, 8, 15, 9, 17, 10, 19,
	11, 21, 12, 23, 13, 25, 14, 27, 15, 29, 16, 31, 17, 33, 18, 35, 19, 37,
	20, 39, 21, 41, 22, 43, 23, 45, 24, 47, 25, 49, 26, 51, 27, 3, 2, 15, 3,
	2, 50, 59, 3, 2, 36, 36, 4, 2, 86, 86, 118, 118, 4, 2, 84, 84, 116, 116,
	4, 2, 87, 87, 119, 119, 4, 2, 71, 71, 103, 103, 4, 2, 72, 72, 104, 104,
	4, 2, 67, 67, 99, 99, 4, 2, 78, 78, 110, 110, 4, 2, 85, 85, 117, 117, 4,
	2, 67, 92, 99, 124, 7, 2, 48, 48, 50, 59, 67, 92, 97, 97, 99, 124, 5, 2,
	11, 12, 15, 15, 34, 34, 2, 149, 2, 3, 3, 2, 2, 2, 2, 5, 3, 2, 2, 2, 2,
	7, 3, 2, 2, 2, 2, 9, 3, 2, 2, 2, 2, 11, 3, 2, 2, 2, 2, 13, 3, 2, 2, 2,
	2, 15, 3, 2, 2, 2, 2, 17, 3, 2, 2, 2, 2, 19, 3, 2, 2, 2, 2, 21, 3, 2, 2,
	2, 2, 23, 3, 2, 2, 2, 2, 25, 3, 2, 2, 2, 2, 27, 3, 2, 2, 2, 2, 29, 3, 2,
	2, 2, 2, 31, 3, 2, 2, 2, 2, 33, 3, 2, 2, 2, 2, 35, 3, 2, 2, 2, 2, 37, 3,
	2, 2, 2, 2, 39, 3, 2, 2, 2, 2, 41, 3, 2, 2, 2, 2, 43, 3, 2, 2, 2, 2, 45,
	3, 2, 2, 2, 2, 47, 3, 2, 2, 2, 2, 49, 3, 2, 2, 2, 2, 51, 3, 2, 2, 2, 3,
	53, 3, 2, 2, 2, 5, 55, 3, 2, 2, 2, 7, 57, 3, 2, 2, 2, 9, 59, 3, 2, 2, 2,
	11, 61, 3, 2, 2, 2, 13, 63, 3, 2, 2, 2, 15, 65, 3, 2, 2, 2, 17, 67, 3,
	2, 2, 2, 19, 69, 3, 2, 2, 2, 21, 71, 3, 2, 2, 2, 23, 73, 3, 2, 2, 2, 25,
	75, 3, 2, 2, 2, 27, 77, 3, 2, 2, 2, 29, 80, 3, 2, 2, 2, 31, 83, 3, 2, 2,
	2, 33, 85, 3, 2, 2, 2, 35, 88, 3, 2, 2, 2, 37, 90, 3, 2, 2, 2, 39, 93,
	3, 2, 2, 2, 41, 105, 3, 2, 2, 2, 43, 116, 3, 2, 2, 2, 45, 121, 3, 2, 2,
	2, 47, 127, 3, 2, 2, 2, 49, 135, 3, 2, 2, 2, 51, 141, 3, 2, 2, 2, 53, 54,
	7, 46, 2, 2, 54, 4, 3, 2, 2, 2, 55, 56, 7, 42, 2, 2, 56, 6, 3, 2, 2, 2,
	57, 58, 7, 43, 2, 2, 58, 8, 3, 2, 2, 2, 59, 60, 7, 93, 2, 2, 60, 10, 3,
	2, 2, 2, 61, 62, 7, 95, 2, 2, 62, 12, 3, 2, 2, 2, 63, 64, 7, 48, 2, 2,
	64, 14, 3, 2, 2, 2, 65, 66, 7, 45, 2, 2, 66, 16, 3, 2, 2, 2, 67, 68, 7,
	47, 2, 2, 68, 18, 3, 2, 2, 2, 69, 70, 7, 44, 2, 2, 70, 20, 3, 2, 2, 2,
	71, 72, 7, 49, 2, 2, 72, 22, 3, 2, 2, 2, 73, 74, 7, 96, 2, 2, 74, 24, 3,
	2, 2, 2, 75, 76, 7, 63, 2, 2, 76, 26, 3, 2, 2, 2, 77, 78, 7, 35, 2, 2,
	78, 79, 7, 63, 2, 2, 79, 28, 3, 2, 2, 2, 80, 81, 7, 62, 2, 2, 81, 82, 7,
	63, 2, 2, 82, 30, 3, 2, 2, 2, 83, 84, 7, 62, 2, 2, 84, 32, 3, 2, 2, 2,
	85, 86, 7, 64, 2, 2, 86, 87, 7, 63, 2, 2, 87, 34, 3, 2, 2, 2, 88, 89, 7,
	64, 2, 2, 89, 36, 3, 2, 2, 2, 90, 91, 7, 40, 2, 2, 91, 38, 3, 2, 2, 2,
	92, 94, 9, 2, 2, 2, 93, 92, 3, 2, 2, 2, 94, 95, 3, 2, 2, 2, 95, 93, 3,
	2, 2, 2, 95, 96, 3, 2, 2, 2, 96, 103, 3, 2, 2, 2, 97, 99, 7, 48, 2, 2,
	98, 100, 9, 2, 2, 2, 99, 98, 3, 2, 2, 2, 100, 101, 3, 2, 2, 2, 101, 99,
	3, 2, 2, 2, 101, 102, 3, 2, 2, 2, 102, 104, 3, 2, 2, 2, 103, 97, 3, 2,
	2, 2, 103, 104, 3, 2, 2, 2, 104, 40, 3, 2, 2, 2, 105, 111, 7, 36, 2, 2,
	106, 110, 10, 3, 2, 2, 107, 108, 7, 36, 2, 2, 108, 110, 7, 36, 2, 2, 109,
	106, 3, 2, 2, 2, 109, 107, 3, 2, 2, 2, 110, 113, 3, 2, 2, 2, 111, 109,
	3, 2, 2, 2, 111, 112, 3, 2, 2, 2, 112, 114, 3, 2, 2, 2, 113, 111, 3, 2,
	2, 2, 114, 115, 7, 36, 2, 2, 115, 42, 3, 2, 2, 2, 116, 117, 9, 4, 2, 2,
	117, 118, 9, 5, 2, 2, 118, 119, 9, 6, 2, 2, 119, 120, 9, 7, 2, 2, 120,
	44, 3, 2, 2, 2, 121, 122, 9, 8, 2, 2, 122, 123, 9, 9, 2, 2, 123, 124, 9,
	10, 2, 2, 124, 125, 9, 11, 2, 2, 125, 126, 9, 7, 2, 2, 126, 46, 3, 2, 2,
	2, 127, 131, 9, 12, 2, 2, 128, 130, 9, 13, 2, 2, 129, 128, 3, 2, 2, 2,
	130, 133, 3, 2, 2, 2, 131, 129, 3, 2, 2, 2, 131, 132, 3, 2, 2, 2, 132,
	48, 3, 2, 2, 2, 133, 131, 3, 2, 2, 2, 134, 136, 9, 14, 2, 2, 135, 134,
	3, 2, 2, 2, 136, 137, 3, 2, 2, 2, 137, 135, 3, 2, 2, 2, 137, 138, 3, 2,
	2, 2, 138, 139, 3, 2, 2, 2, 139, 140, 8, 25, 2, 2, 140, 50, 3, 2, 2, 2,
	141, 142, 11, 2, 2, 2, 142, 52, 3, 2, 2, 2, 10, 2, 95, 101, 103, 109, 111,
	131, 137, 3, 8, 2, 2,
}

var lexerDeserializer = antlr.NewATNDeserializer(nil)
var lexerAtn = lexerDeserializer.DeserializeFromUInt16(serializedLexerAtn)

var lexerChannelNames = []string{
	"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
}

var lexerModeNames = []string{
	"DEFAULT_MODE",
}

var lexerLiteralNames = []string{
	"", "','", "'('", "')'", "'['", "']'", "'.'", "'+'", "'-'", "'*'", "'/'",
	"'^'", "'='", "'!='", "'<='", "'<'", "'>='", "'>'", "'&'",
}

var lexerSymbolicNames = []string{
	"", "COMMA", "LPAREN", "RPAREN", "LBRACK", "RBRACK", "DOT", "PLUS", "MINUS",
	"TIMES", "DIVIDE", "EXPONENT", "EQ", "NEQ", "LTE", "LT", "GTE", "GT", "AMPERSAND",
	"DECIMAL", "STRING", "TRUE", "FALSE", "NAME", "WS", "ERROR",
}

var lexerRuleNames = []string{
	"COMMA", "LPAREN", "RPAREN", "LBRACK", "RBRACK", "DOT", "PLUS", "MINUS",
	"TIMES", "DIVIDE", "EXPONENT", "EQ", "NEQ", "LTE", "LT", "GTE", "GT", "AMPERSAND",
	"DECIMAL", "STRING", "TRUE", "FALSE", "NAME", "WS", "ERROR",
}

type ExcellentLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var lexerDecisionToDFA = make([]*antlr.DFA, len(lexerAtn.DecisionToState))

func init() {
	for index, ds := range lexerAtn.DecisionToState {
		lexerDecisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

func NewExcellentLexer(input antlr.CharStream) *ExcellentLexer {

	l := new(ExcellentLexer)

	l.BaseLexer = antlr.NewBaseLexer(input)
	l.Interpreter = antlr.NewLexerATNSimulator(l, lexerAtn, lexerDecisionToDFA, antlr.NewPredictionContextCache())

	l.channelNames = lexerChannelNames
	l.modeNames = lexerModeNames
	l.RuleNames = lexerRuleNames
	l.LiteralNames = lexerLiteralNames
	l.SymbolicNames = lexerSymbolicNames
	l.GrammarFileName = "Excellent.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// ExcellentLexer tokens.
const (
	ExcellentLexerCOMMA     = 1
	ExcellentLexerLPAREN    = 2
	ExcellentLexerRPAREN    = 3
	ExcellentLexerLBRACK    = 4
	ExcellentLexerRBRACK    = 5
	ExcellentLexerDOT       = 6
	ExcellentLexerPLUS      = 7
	ExcellentLexerMINUS     = 8
	ExcellentLexerTIMES     = 9
	ExcellentLexerDIVIDE    = 10
	ExcellentLexerEXPONENT  = 11
	ExcellentLexerEQ        = 12
	ExcellentLexerNEQ       = 13
	ExcellentLexerLTE       = 14
	ExcellentLexerLT        = 15
	ExcellentLexerGTE       = 16
	ExcellentLexerGT        = 17
	ExcellentLexerAMPERSAND = 18
	ExcellentLexerDECIMAL   = 19
	ExcellentLexerSTRING    = 20
	ExcellentLexerTRUE      = 21
	ExcellentLexerFALSE     = 22
	ExcellentLexerNAME      = 23
	ExcellentLexerWS        = 24
	ExcellentLexerERROR     = 25
)
