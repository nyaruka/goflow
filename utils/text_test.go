package utils_test

import (
	"testing"

	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestSnakify(t *testing.T) {
	var snakeTests = []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello_world"},
		{"hello_world", "hello_world"},
		{"hello-world", "hello_world"},
		{"_Hello World", "_hello_world"},
		{"   Hello World    ", "hello_world"},
		{"hi😀😃😄😁there", "hi_there"},
		{"昨夜のコ", "昨夜のコ"},
		{"this@isn't@email", "this_isn_t_email"},
	}

	for _, test := range snakeTests {
		assert.Equal(t, test.expected, utils.Snakify(test.input), "unexpected result snakifying '%s'", test.input)
	}
}

func TestTokenizeString(t *testing.T) {
	tokenizerTests := []struct {
		text   string
		result []string
	}{
		{" one ", []string{"one"}},
		{"one   two three", []string{"one", "two", "three"}},
		{"one.two.three", []string{"one", "two", "three"}},
		{"O'Grady can't foo_bar", []string{"O'Grady", "can't", "foo_bar"}}, // single quotes and underscores don't split tokens
		{"öne.βήταa.thé", []string{"öne", "βήταa", "thé"}},                 // non-latin letters allowed in tokens
		{"واحد اثنين ثلاثة", []string{"واحد", "اثنين", "ثلاثة"}},           // RTL scripts
		{"  \t\none(two!*@three ", []string{"one", "two", "three"}},        // other punctuation ignored
		{"spend$£€₠₣₪", []string{"spend", "$", "£", "€", "₠", "₣", "₪"}},   // currency symbols treated as individual tokens
		{"math+=×÷√∊", []string{"math", "+", "=", "×", "÷", "√", "∊"}},     // math symbols treated as individual tokens
		{"emoji😄🏥👪👰😟🧟", []string{"emoji", "😄", "🏥", "👪", "👰", "😟", "🧟"}},   // emojis treated as individual tokens
		{"👍🏿 👨🏼", []string{"👍", "🏿", "👨", "🏼"}},                            // tone modifiers treated as individual tokens
		{"ℹ ℹ️", []string{"ℹ", "ℹ️"}},                                      // variation selectors ignored
		{"ยกเลิก sasa", []string{"ยกเลิก", "sasa"}},                        // Thai word means Cancelled
		{"বাতিল sasa", []string{"বাতিল", "sasa"}},                          // Bangla word means Cancel
		{"ထွက်သွား sasa", []string{"ထွက်သွား", "sasa"}},                    // Burmese word means exit
	}
	for _, test := range tokenizerTests {
		assert.Equal(t, test.result, utils.TokenizeString(test.text), "unexpected result tokenizing '%s'", test.text)
	}
}

func TestTokenizeStringByChars(t *testing.T) {
	tokenizerTests := []struct {
		text   string
		chars  string
		result []string
	}{
		{"one   two three", " ", []string{"one", "two", "three"}},
		{"Jim O'Grady", " ", []string{"Jim", "O'Grady"}},
		{"one.βήταa/three", "./", []string{"one", "βήταa", "three"}},
		{"one😄three", "😄", []string{"one", "three"}},
		{"  one.two.*@three ", " .*@", []string{"one", "two", "three"}},
		{" one ", " ", []string{"one"}},
	}
	for _, test := range tokenizerTests {
		assert.Equal(t, test.result, utils.TokenizeStringByChars(test.text, test.chars), "unexpected result tokenizing '%s'", test.text)
	}
}

func TestTokenizeStringByUnicodeSeg(t *testing.T) {
	tokenizerTests := []struct {
		text   string
		result []string
	}{
		// test cases taken from Apache Lucene (TestStandardAnalyzer.java) to assert equivalency
		{"Վիքիպեդիայի 13 միլիոն հոդվածները (4,600` հայերեն վիքիպեդիայում) գրվել են կամավորների կողմից ու համարյա բոլոր հոդվածները կարող է խմբագրել ցանկաց մարդ ով կարող է բացել Վիքիպեդիայի կայքը։", []string{"Վիքիպեդիայի", "13", "միլիոն", "հոդվածները", "4,600", "հայերեն", "վիքիպեդիայում", "գրվել", "են", "կամավորների", "կողմից", "ու", "համարյա", "բոլոր", "հոդվածները", "կարող", "է", "խմբագրել", "ցանկաց", "մարդ", "ով", "կարող", "է", "բացել", "Վիքիպեդիայի", "կայքը"}},
		{"ዊኪፔድያ የባለ ብዙ ቋንቋ የተሟላ ትክክለኛና ነጻ መዝገበ ዕውቀት (ኢንሳይክሎፒዲያ) ነው። ማንኛውም", []string{"ዊኪፔድያ", "የባለ", "ብዙ", "ቋንቋ", "የተሟላ", "ትክክለኛና", "ነጻ", "መዝገበ", "ዕውቀት", "ኢንሳይክሎፒዲያ", "ነው", "ማንኛውም"}},
		{"الفيلم الوثائقي الأول عن ويكيبيديا يسمى \"الحقيقة بالأرقام: قصة ويكيبيديا\" (بالإنجليزية: Truth in Numbers: The Wikipedia Story)، سيتم إطلاقه في 2008.", []string{"الفيلم", "الوثائقي", "الأول", "عن", "ويكيبيديا", "يسمى", "الحقيقة", "بالأرقام", "قصة", "ويكيبيديا", "بالإنجليزية", "Truth", "in", "Numbers", "The", "Wikipedia", "Story", "سيتم", "إطلاقه", "في", "2008"}},
		{"ܘܝܩܝܦܕܝܐ (ܐܢܓܠܝܐ: Wikipedia) ܗܘ ܐܝܢܣܩܠܘܦܕܝܐ ܚܐܪܬܐ ܕܐܢܛܪܢܛ ܒܠܫܢ̈ܐ ܣܓܝܐ̈ܐ܂ ܫܡܗ ܐܬܐ ܡܢ ܡ̈ܠܬܐ ܕ\"ܘܝܩܝ\" ܘ\"ܐܝܢܣܩܠܘܦܕܝܐ\"܀", []string{"ܘܝܩܝܦܕܝܐ", "ܐܢܓܠܝܐ", "Wikipedia", "ܗܘ", "ܐܝܢܣܩܠܘܦܕܝܐ", "ܚܐܪܬܐ", "ܕܐܢܛܪܢܛ", "ܒܠܫܢ̈ܐ", "ܣܓܝܐ̈ܐ", "ܫܡܗ", "ܐܬܐ", "ܡܢ", "ܡ̈ܠܬܐ", "ܕ", "ܘܝܩܝ", "ܘ", "ܐܝܢܣܩܠܘܦܕܝܐ"}},
		{"এই বিশ্বকোষ পরিচালনা করে উইকিমিডিয়া ফাউন্ডেশন (একটি অলাভজনক সংস্থা)। উইকিপিডিয়ার শুরু ১৫ জানুয়ারি, ২০০১ সালে। এখন পর্যন্ত ২০০টিরও বেশী ভাষায় উইকিপিডিয়া রয়েছে।", []string{"এই", "বিশ্বকোষ", "পরিচালনা", "করে", "উইকিমিডিয়া", "ফাউন্ডেশন", "একটি", "অলাভজনক", "সংস্থা", "উইকিপিডিয়ার", "শুরু", "১৫", "জানুয়ারি", "২০০১", "সালে", "এখন", "পর্যন্ত", "২০০টিরও", "বেশী", "ভাষায়", "উইকিপিডিয়া", "রয়েছে"}},
		{"ویکی پدیای انگلیسی در تاریخ ۲۵ دی ۱۳۷۹ به صورت مکملی برای دانشنامهٔ تخصصی نوپدیا نوشته شد.", []string{"ویکی", "پدیای", "انگلیسی", "در", "تاریخ", "۲۵", "دی", "۱۳۷۹", "به", "صورت", "مکملی", "برای", "دانشنامهٔ", "تخصصی", "نوپدیا", "نوشته", "شد"}},
		{"Γράφεται σε συνεργασία από εθελοντές με το λογισμικό wiki, κάτι που σημαίνει ότι άρθρα μπορεί να προστεθούν ή να αλλάξουν από τον καθένα.", []string{"Γράφεται", "σε", "συνεργασία", "από", "εθελοντές", "με", "το", "λογισμικό", "wiki", "κάτι", "που", "σημαίνει", "ότι", "άρθρα", "μπορεί", "να", "προστεθούν", "ή", "να", "αλλάξουν", "από", "τον", "καθένα"}},
		//{"我是中国人。 １２３４ Ｔｅｓｔｓ ", []string{"我", "是", "中", "国", "人", "１２３４", "Ｔｅｓｔｓ"}},  // gives different result
		{"", []string{}},
		{".", []string{}},
		{" ", []string{}},
		{"B2B", []string{"B2B"}},
		{"2B", []string{"2B"}},
		{"some-dashed-phrase", []string{"some", "dashed", "phrase"}},
		{"dogs,chase,cats", []string{"dogs", "chase", "cats"}},
		{"ac/dc", []string{"ac", "dc"}},
		{"O'Reilly", []string{"O'Reilly"}},
		{"you're", []string{"you're"}},
		{"she's", []string{"she's"}},
		{"Jim's", []string{"Jim's"}},
		{"don't", []string{"don't"}},
		{"O'Reilly's", []string{"O'Reilly's"}},
		{"21.35", []string{"21.35"}},
		{"R2D2 C3PO", []string{"R2D2", "C3PO"}},
		{"216.239.63.104", []string{"216.239.63.104"}},
		{"216.239.63.104", []string{"216.239.63.104"}},
		{"David has 5000 bones", []string{"David", "has", "5000", "bones"}},
		{"C embedded developers wanted", []string{"C", "embedded", "developers", "wanted"}},
		{"foo bar FOO BAR", []string{"foo", "bar", "FOO", "BAR"}},
		{"foo      bar .  FOO <> BAR", []string{"foo", "bar", "FOO", "BAR"}},
		{"\"QUOTED\" word", []string{"QUOTED", "word"}},
		{"안녕하세요 한글입니다", []string{"안녕하세요", "한글입니다"}},
	}
	for _, test := range tokenizerTests {
		assert.Equal(t, test.result, utils.TokenizeStringByUnicodeSeg(test.text), "unexpected result tokenizing '%s'", test.text)
	}
}

func TestPrefixOverlap(t *testing.T) {
	assert.Equal(t, 0, utils.PrefixOverlap("", ""))
	assert.Equal(t, 0, utils.PrefixOverlap("abc", ""))
	assert.Equal(t, 0, utils.PrefixOverlap("", "abc"))
	assert.Equal(t, 0, utils.PrefixOverlap("a", "x"))
	assert.Equal(t, 1, utils.PrefixOverlap("x", "x"))
	assert.Equal(t, 3, utils.PrefixOverlap("xyz", "xyz"))
	assert.Equal(t, 2, utils.PrefixOverlap("xya", "xyz"))
	assert.Equal(t, 2, utils.PrefixOverlap("😄😟👨🏼", "😄😟👰"))
	assert.Equal(t, 4, utils.PrefixOverlap("25078", "25073254252"))
}

func TestStringSlices(t *testing.T) {
	assert.Equal(t, []string{"he", "hello", "world"}, utils.StringSlices("hello world", []int{0, 2, 0, 5, 6, 11}))
}

func TestStringSliceContains(t *testing.T) {
	assert.False(t, utils.StringSliceContains(nil, "a", true))
	assert.False(t, utils.StringSliceContains([]string{}, "a", true))
	assert.False(t, utils.StringSliceContains([]string{"b", "c"}, "a", true))
	assert.True(t, utils.StringSliceContains([]string{"b", "a", "c"}, "a", true))
	assert.False(t, utils.StringSliceContains([]string{"b", "a", "c"}, "A", true))
	assert.True(t, utils.StringSliceContains([]string{"b", "a", "c"}, "A", false))
}

func TestIndent(t *testing.T) {
	assert.Equal(t, "", utils.Indent("", "  "))
	assert.Equal(t, "  x", utils.Indent("x", "  "))
	assert.Equal(t, "  x\n  y", utils.Indent("x\ny", "  "))
	assert.Equal(t, "  x\n\n  y", utils.Indent("x\n\ny", "  "))
	assert.Equal(t, ">>>x", utils.Indent("x", ">>>"))
}

func TestNestingDepthExceeds(t *testing.T) {
	assert.False(t, utils.NestingDepthExceeds("", 5))
	assert.False(t, utils.NestingDepthExceeds("no brackets here", 5))
	assert.False(t, utils.NestingDepthExceeds("(((((", 5))        // depth 5, not > 5
	assert.False(t, utils.NestingDepthExceeds("()()()()()()", 1)) // flat, max depth 1
	assert.False(t, utils.NestingDepthExceeds("a[b].c{d}(e)", 1)) // mixed but shallow
	assert.True(t, utils.NestingDepthExceeds("((((((", 5))        // depth 6
	assert.True(t, utils.NestingDepthExceeds("([{([{", 5))        // mixed opening brackets, depth 6
	// unbalanced closers don't push depth negative
	assert.False(t, utils.NestingDepthExceeds("))))))((", 5))
}
