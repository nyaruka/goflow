package expressions

type functionTemplate struct {
	name   string
	params string
	join   string
	two    string
	three  string
	four   string
}

var functionTemplates = map[string]functionTemplate{
	"average":     {name: "mean"},
	"date":        {name: "datetime", params: "(\"%s-%s-%s\")"},
	"datedif":     {name: "datetime_diff"},
	"datevalue":   {name: "datetime"},
	"day":         {name: "format_datetime", params: `(%s, "D")`},
	"days":        {name: "datetime_diff", params: "(%s, %s, \"D\")"},
	"edate":       {name: "datetime_add", params: "(%s, %s, \"M\")"},
	"field":       {name: "field", params: "(%s, %s - 1, %s)"},
	"first_word":  {name: "word", params: "(%s, 0)"},
	"fixed":       {name: "format_number", params: "(%s)", two: "(%s, %s)", three: "(%s, %s, %v)"},
	"format_date": {name: "format_datetime"},
	"hour":        {name: "format_datetime", params: `(%s, "h")`},
	"int":         {name: "round_down"},
	"len":         {name: "length"},
	"minute":      {name: "format_datetime", params: `(%s, "m")`},
	"month":       {name: "format_datetime", params: `(%s, "M")`},
	"now":         {name: "now"},
	"proper":      {name: "title"},
	"randbetween": {name: "rand_between"},
	"read_digits": {name: "read_chars"},
	"rept":        {name: "repeat"},
	"rounddown":   {name: "round_down"},
	"roundup":     {name: "round_up"},
	"second":      {name: "format_datetime", params: `(%s, "s")`},
	"substitute":  {name: "replace"},
	"timevalue":   {name: "parse_datetime"},
	"trunc":       {name: "round_down"},
	"unichar":     {name: "char"},
	"unicode":     {name: "code"},
	"word_slice":  {name: "word_slice", params: "(%s, %s - 1)", three: "(%s, %s - 1, %s - 1)"},
	"word":        {name: "word", params: "(%s, %s - 1)"},
	"year":        {name: "format_datetime", params: `(%s, "YYYY")`},

	// we drop this function, instead joining with the cat operator
	"concatenate": {join: " & "},

	// translate to maths
	"power": {params: "%s ^ %s"},
	"exp":   {params: "2.718281828459045 ^ %s"},
	"sum":   {params: "%s + %s"},

	// this one is a special case format, we sum these parts into seconds for datetime_add
	"time": {name: "time", params: "(%s %s %s)"},
}

var ignoredFunctions = map[string]bool{
	"time": true,
	"sum":  true,

	// in some cases we actually remove function names
	// such add switching CONCAT to a simple operator expression
	"": true,
}
