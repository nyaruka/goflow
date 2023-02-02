// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package gen

import (
	"fmt"
	"sync"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter

type ContactQLLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var contactqllexerLexerStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	channelNames           []string
	modeNames              []string
	literalNames           []string
	symbolicNames          []string
	ruleNames              []string
	predictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func contactqllexerLexerInit() {
	staticData := &contactqllexerLexerStaticData
	staticData.channelNames = []string{
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
	}
	staticData.modeNames = []string{
		"DEFAULT_MODE",
	}
	staticData.literalNames = []string{
		"", "'('", "')'",
	}
	staticData.symbolicNames = []string{
		"", "LPAREN", "RPAREN", "AND", "OR", "COMPARATOR", "TEXT", "STRING",
		"WS", "ERROR",
	}
	staticData.ruleNames = []string{
		"HAS", "IS", "LPAREN", "RPAREN", "AND", "OR", "COMPARATOR", "TEXT",
		"STRING", "WS", "ERROR", "UnicodeLetter", "UnicodeClass_LU", "UnicodeClass_LL",
		"UnicodeClass_LT", "UnicodeClass_LM", "UnicodeClass_LO", "UnicodeDigit",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 9, 114, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2,
		10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15,
		7, 15, 2, 16, 7, 16, 2, 17, 7, 17, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1,
		1, 1, 1, 2, 1, 2, 1, 3, 1, 3, 1, 4, 1, 4, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5,
		1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 3, 6,
		67, 8, 6, 1, 7, 1, 7, 1, 7, 4, 7, 72, 8, 7, 11, 7, 12, 7, 73, 1, 8, 1,
		8, 1, 8, 1, 8, 5, 8, 80, 8, 8, 10, 8, 12, 8, 83, 9, 8, 1, 8, 1, 8, 1, 9,
		4, 9, 88, 8, 9, 11, 9, 12, 9, 89, 1, 9, 1, 9, 1, 10, 1, 10, 1, 11, 1, 11,
		1, 11, 1, 11, 1, 11, 3, 11, 101, 8, 11, 1, 12, 1, 12, 1, 13, 1, 13, 1,
		14, 1, 14, 1, 15, 1, 15, 1, 16, 1, 16, 1, 17, 1, 17, 0, 0, 18, 1, 0, 3,
		0, 5, 1, 7, 2, 9, 3, 11, 4, 13, 5, 15, 6, 17, 7, 19, 8, 21, 9, 23, 0, 25,
		0, 27, 0, 29, 0, 31, 0, 33, 0, 35, 0, 1, 0, 18, 2, 0, 72, 72, 104, 104,
		2, 0, 65, 65, 97, 97, 2, 0, 83, 83, 115, 115, 2, 0, 73, 73, 105, 105, 2,
		0, 78, 78, 110, 110, 2, 0, 68, 68, 100, 100, 2, 0, 79, 79, 111, 111, 2,
		0, 82, 82, 114, 114, 2, 0, 60, 60, 62, 62, 6, 0, 39, 39, 43, 43, 45, 47,
		58, 58, 64, 64, 95, 95, 1, 0, 34, 34, 3, 0, 9, 10, 13, 13, 32, 32, 82,
		0, 65, 90, 192, 214, 216, 222, 256, 310, 313, 327, 330, 381, 385, 386,
		388, 395, 398, 401, 403, 404, 406, 408, 412, 413, 415, 416, 418, 425, 428,
		435, 437, 444, 452, 461, 463, 475, 478, 494, 497, 500, 502, 504, 506, 562,
		570, 571, 573, 574, 577, 582, 584, 590, 880, 882, 886, 895, 902, 906, 908,
		929, 931, 939, 975, 980, 984, 1006, 1012, 1015, 1017, 1018, 1021, 1071,
		1120, 1152, 1162, 1229, 1232, 1326, 1329, 1366, 4256, 4293, 4295, 4301,
		7680, 7828, 7838, 7934, 7944, 7951, 7960, 7965, 7976, 7983, 7992, 7999,
		8008, 8013, 8025, 8031, 8040, 8047, 8120, 8123, 8136, 8139, 8152, 8155,
		8168, 8172, 8184, 8187, 8450, 8455, 8459, 8461, 8464, 8466, 8469, 8477,
		8484, 8493, 8496, 8499, 8510, 8511, 8517, 8579, 11264, 11310, 11360, 11364,
		11367, 11376, 11378, 11381, 11390, 11392, 11394, 11490, 11499, 11501, 11506,
		42560, 42562, 42604, 42624, 42650, 42786, 42798, 42802, 42862, 42873, 42886,
		42891, 42893, 42896, 42898, 42902, 42925, 42928, 42929, 65313, 65338, 81,
		0, 97, 122, 181, 246, 248, 255, 257, 375, 378, 384, 387, 389, 392, 402,
		405, 411, 414, 417, 419, 421, 424, 429, 432, 436, 438, 447, 454, 460, 462,
		499, 501, 505, 507, 569, 572, 578, 583, 659, 661, 687, 881, 883, 887, 893,
		912, 974, 976, 977, 981, 983, 985, 1011, 1013, 1119, 1121, 1153, 1163,
		1215, 1218, 1327, 1377, 1415, 7424, 7467, 7531, 7543, 7545, 7578, 7681,
		7837, 7839, 7943, 7952, 7957, 7968, 7975, 7984, 7991, 8000, 8005, 8016,
		8023, 8032, 8039, 8048, 8061, 8064, 8071, 8080, 8087, 8096, 8103, 8112,
		8116, 8118, 8119, 8126, 8132, 8134, 8135, 8144, 8147, 8150, 8151, 8160,
		8167, 8178, 8180, 8182, 8183, 8458, 8467, 8495, 8505, 8508, 8509, 8518,
		8521, 8526, 8580, 11312, 11358, 11361, 11372, 11377, 11387, 11393, 11500,
		11502, 11507, 11520, 11557, 11559, 11565, 42561, 42605, 42625, 42651, 42787,
		42801, 42803, 42872, 42874, 42876, 42879, 42887, 42892, 42894, 42897, 42901,
		42903, 42921, 43002, 43866, 43876, 43877, 64256, 64262, 64275, 64279, 65345,
		65370, 6, 0, 453, 459, 498, 8079, 8088, 8095, 8104, 8111, 8124, 8140, 8188,
		8188, 33, 0, 688, 705, 710, 721, 736, 740, 748, 750, 884, 890, 1369, 1600,
		1765, 1766, 2036, 2037, 2042, 2074, 2084, 2088, 2417, 3654, 3782, 4348,
		6103, 6211, 6823, 7293, 7468, 7530, 7544, 7615, 8305, 8319, 8336, 8348,
		11388, 11389, 11631, 11823, 12293, 12341, 12347, 12542, 40981, 42237, 42508,
		42623, 42652, 42653, 42775, 42783, 42864, 42888, 43000, 43001, 43471, 43494,
		43632, 43741, 43763, 43764, 43868, 43871, 65392, 65439, 234, 0, 170, 186,
		443, 451, 660, 1514, 1520, 1522, 1568, 1599, 1601, 1610, 1646, 1647, 1649,
		1747, 1749, 1788, 1791, 1808, 1810, 1839, 1869, 1957, 1969, 2026, 2048,
		2069, 2112, 2136, 2208, 2226, 2308, 2361, 2365, 2384, 2392, 2401, 2418,
		2432, 2437, 2444, 2447, 2448, 2451, 2472, 2474, 2480, 2482, 2489, 2493,
		2510, 2524, 2525, 2527, 2529, 2544, 2545, 2565, 2570, 2575, 2576, 2579,
		2600, 2602, 2608, 2610, 2611, 2613, 2614, 2616, 2617, 2649, 2652, 2654,
		2676, 2693, 2701, 2703, 2705, 2707, 2728, 2730, 2736, 2738, 2739, 2741,
		2745, 2749, 2768, 2784, 2785, 2821, 2828, 2831, 2832, 2835, 2856, 2858,
		2864, 2866, 2867, 2869, 2873, 2877, 2913, 2929, 2947, 2949, 2954, 2958,
		2960, 2962, 2965, 2969, 2970, 2972, 2986, 2990, 3001, 3024, 3084, 3086,
		3088, 3090, 3112, 3114, 3129, 3133, 3212, 3214, 3216, 3218, 3240, 3242,
		3251, 3253, 3257, 3261, 3294, 3296, 3297, 3313, 3314, 3333, 3340, 3342,
		3344, 3346, 3386, 3389, 3406, 3424, 3425, 3450, 3455, 3461, 3478, 3482,
		3505, 3507, 3515, 3517, 3526, 3585, 3632, 3634, 3635, 3648, 3653, 3713,
		3714, 3716, 3722, 3725, 3735, 3737, 3743, 3745, 3747, 3749, 3751, 3754,
		3755, 3757, 3760, 3762, 3763, 3773, 3780, 3804, 3807, 3840, 3911, 3913,
		3948, 3976, 3980, 4096, 4138, 4159, 4181, 4186, 4189, 4193, 4208, 4213,
		4225, 4238, 4346, 4349, 4680, 4682, 4685, 4688, 4694, 4696, 4701, 4704,
		4744, 4746, 4749, 4752, 4784, 4786, 4789, 4792, 4798, 4800, 4805, 4808,
		4822, 4824, 4880, 4882, 4885, 4888, 4954, 4992, 5007, 5024, 5108, 5121,
		5740, 5743, 5759, 5761, 5786, 5792, 5866, 5873, 5880, 5888, 5900, 5902,
		5905, 5920, 5937, 5952, 5969, 5984, 5996, 5998, 6000, 6016, 6067, 6108,
		6210, 6212, 6263, 6272, 6312, 6314, 6389, 6400, 6430, 6480, 6509, 6512,
		6516, 6528, 6571, 6593, 6599, 6656, 6678, 6688, 6740, 6917, 6963, 6981,
		6987, 7043, 7072, 7086, 7087, 7098, 7141, 7168, 7203, 7245, 7247, 7258,
		7287, 7401, 7404, 7406, 7409, 7413, 7414, 8501, 8504, 11568, 11623, 11648,
		11670, 11680, 11686, 11688, 11694, 11696, 11702, 11704, 11710, 11712, 11718,
		11720, 11726, 11728, 11734, 11736, 11742, 12294, 12348, 12353, 12438, 12447,
		12538, 12543, 12589, 12593, 12686, 12704, 12730, 12784, 12799, 13312, 19893,
		19968, 40908, 40960, 40980, 40982, 42124, 42192, 42231, 42240, 42507, 42512,
		42527, 42538, 42539, 42606, 42725, 42999, 43009, 43011, 43013, 43015, 43018,
		43020, 43042, 43072, 43123, 43138, 43187, 43250, 43255, 43259, 43301, 43312,
		43334, 43360, 43388, 43396, 43442, 43488, 43492, 43495, 43503, 43514, 43518,
		43520, 43560, 43584, 43586, 43588, 43595, 43616, 43631, 43633, 43638, 43642,
		43695, 43697, 43709, 43712, 43714, 43739, 43740, 43744, 43754, 43762, 43782,
		43785, 43790, 43793, 43798, 43808, 43814, 43816, 43822, 43968, 44002, 44032,
		55203, 55216, 55238, 55243, 55291, 63744, 64109, 64112, 64217, 64285, 64296,
		64298, 64310, 64312, 64316, 64318, 64433, 64467, 64829, 64848, 64911, 64914,
		64967, 65008, 65019, 65136, 65140, 65142, 65276, 65382, 65391, 65393, 65437,
		65440, 65470, 65474, 65479, 65482, 65487, 65490, 65495, 65498, 65500, 37,
		0, 48, 57, 1632, 1641, 1776, 1785, 1984, 1993, 2406, 2415, 2534, 2543,
		2662, 2671, 2790, 2799, 2918, 2927, 3046, 3055, 3174, 3183, 3302, 3311,
		3430, 3439, 3558, 3567, 3664, 3673, 3792, 3801, 3872, 3881, 4160, 4169,
		4240, 4249, 6112, 6121, 6160, 6169, 6470, 6479, 6608, 6617, 6784, 6793,
		6800, 6809, 6992, 7001, 7088, 7097, 7232, 7241, 7248, 7257, 42528, 42537,
		43216, 43225, 43264, 43273, 43472, 43481, 43504, 43513, 43600, 43609, 44016,
		44025, 65296, 65305, 121, 0, 5, 1, 0, 0, 0, 0, 7, 1, 0, 0, 0, 0, 9, 1,
		0, 0, 0, 0, 11, 1, 0, 0, 0, 0, 13, 1, 0, 0, 0, 0, 15, 1, 0, 0, 0, 0, 17,
		1, 0, 0, 0, 0, 19, 1, 0, 0, 0, 0, 21, 1, 0, 0, 0, 1, 37, 1, 0, 0, 0, 3,
		41, 1, 0, 0, 0, 5, 44, 1, 0, 0, 0, 7, 46, 1, 0, 0, 0, 9, 48, 1, 0, 0, 0,
		11, 52, 1, 0, 0, 0, 13, 66, 1, 0, 0, 0, 15, 71, 1, 0, 0, 0, 17, 75, 1,
		0, 0, 0, 19, 87, 1, 0, 0, 0, 21, 93, 1, 0, 0, 0, 23, 100, 1, 0, 0, 0, 25,
		102, 1, 0, 0, 0, 27, 104, 1, 0, 0, 0, 29, 106, 1, 0, 0, 0, 31, 108, 1,
		0, 0, 0, 33, 110, 1, 0, 0, 0, 35, 112, 1, 0, 0, 0, 37, 38, 7, 0, 0, 0,
		38, 39, 7, 1, 0, 0, 39, 40, 7, 2, 0, 0, 40, 2, 1, 0, 0, 0, 41, 42, 7, 3,
		0, 0, 42, 43, 7, 2, 0, 0, 43, 4, 1, 0, 0, 0, 44, 45, 5, 40, 0, 0, 45, 6,
		1, 0, 0, 0, 46, 47, 5, 41, 0, 0, 47, 8, 1, 0, 0, 0, 48, 49, 7, 1, 0, 0,
		49, 50, 7, 4, 0, 0, 50, 51, 7, 5, 0, 0, 51, 10, 1, 0, 0, 0, 52, 53, 7,
		6, 0, 0, 53, 54, 7, 7, 0, 0, 54, 12, 1, 0, 0, 0, 55, 67, 5, 61, 0, 0, 56,
		57, 5, 33, 0, 0, 57, 67, 5, 61, 0, 0, 58, 67, 5, 126, 0, 0, 59, 60, 5,
		62, 0, 0, 60, 67, 5, 61, 0, 0, 61, 62, 5, 60, 0, 0, 62, 67, 5, 61, 0, 0,
		63, 67, 7, 8, 0, 0, 64, 67, 3, 1, 0, 0, 65, 67, 3, 3, 1, 0, 66, 55, 1,
		0, 0, 0, 66, 56, 1, 0, 0, 0, 66, 58, 1, 0, 0, 0, 66, 59, 1, 0, 0, 0, 66,
		61, 1, 0, 0, 0, 66, 63, 1, 0, 0, 0, 66, 64, 1, 0, 0, 0, 66, 65, 1, 0, 0,
		0, 67, 14, 1, 0, 0, 0, 68, 72, 3, 23, 11, 0, 69, 72, 3, 35, 17, 0, 70,
		72, 7, 9, 0, 0, 71, 68, 1, 0, 0, 0, 71, 69, 1, 0, 0, 0, 71, 70, 1, 0, 0,
		0, 72, 73, 1, 0, 0, 0, 73, 71, 1, 0, 0, 0, 73, 74, 1, 0, 0, 0, 74, 16,
		1, 0, 0, 0, 75, 81, 5, 34, 0, 0, 76, 80, 8, 10, 0, 0, 77, 78, 5, 92, 0,
		0, 78, 80, 5, 34, 0, 0, 79, 76, 1, 0, 0, 0, 79, 77, 1, 0, 0, 0, 80, 83,
		1, 0, 0, 0, 81, 79, 1, 0, 0, 0, 81, 82, 1, 0, 0, 0, 82, 84, 1, 0, 0, 0,
		83, 81, 1, 0, 0, 0, 84, 85, 5, 34, 0, 0, 85, 18, 1, 0, 0, 0, 86, 88, 7,
		11, 0, 0, 87, 86, 1, 0, 0, 0, 88, 89, 1, 0, 0, 0, 89, 87, 1, 0, 0, 0, 89,
		90, 1, 0, 0, 0, 90, 91, 1, 0, 0, 0, 91, 92, 6, 9, 0, 0, 92, 20, 1, 0, 0,
		0, 93, 94, 9, 0, 0, 0, 94, 22, 1, 0, 0, 0, 95, 101, 3, 25, 12, 0, 96, 101,
		3, 27, 13, 0, 97, 101, 3, 29, 14, 0, 98, 101, 3, 31, 15, 0, 99, 101, 3,
		33, 16, 0, 100, 95, 1, 0, 0, 0, 100, 96, 1, 0, 0, 0, 100, 97, 1, 0, 0,
		0, 100, 98, 1, 0, 0, 0, 100, 99, 1, 0, 0, 0, 101, 24, 1, 0, 0, 0, 102,
		103, 7, 12, 0, 0, 103, 26, 1, 0, 0, 0, 104, 105, 7, 13, 0, 0, 105, 28,
		1, 0, 0, 0, 106, 107, 7, 14, 0, 0, 107, 30, 1, 0, 0, 0, 108, 109, 7, 15,
		0, 0, 109, 32, 1, 0, 0, 0, 110, 111, 7, 16, 0, 0, 111, 34, 1, 0, 0, 0,
		112, 113, 7, 17, 0, 0, 113, 36, 1, 0, 0, 0, 8, 0, 66, 71, 73, 79, 81, 89,
		100, 1, 6, 0, 0,
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

// ContactQLLexerInit initializes any static state used to implement ContactQLLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewContactQLLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func ContactQLLexerInit() {
	staticData := &contactqllexerLexerStaticData
	staticData.once.Do(contactqllexerLexerInit)
}

// NewContactQLLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewContactQLLexer(input antlr.CharStream) *ContactQLLexer {
	ContactQLLexerInit()
	l := new(ContactQLLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
	staticData := &contactqllexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	l.channelNames = staticData.channelNames
	l.modeNames = staticData.modeNames
	l.RuleNames = staticData.ruleNames
	l.LiteralNames = staticData.literalNames
	l.SymbolicNames = staticData.symbolicNames
	l.GrammarFileName = "ContactQL.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// ContactQLLexer tokens.
const (
	ContactQLLexerLPAREN     = 1
	ContactQLLexerRPAREN     = 2
	ContactQLLexerAND        = 3
	ContactQLLexerOR         = 4
	ContactQLLexerCOMPARATOR = 5
	ContactQLLexerTEXT       = 6
	ContactQLLexerSTRING     = 7
	ContactQLLexerWS         = 8
	ContactQLLexerERROR      = 9
)
