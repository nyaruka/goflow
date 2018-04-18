package legacy

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/functions"
	"github.com/nyaruka/goflow/excellent/gen"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows/runs"
	"github.com/nyaruka/goflow/utils"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// allowed top-level identifiers in legacy expressions, i.e. @contact.bar is valid but @foo.bar isn't
var legacyContextTopLevels = []string{"channel", "child", "contact", "date", "extra", "flow", "parent", "step"}

// ExtraVarsMapping defines how @extra.* variables should be migrated
type ExtraVarsMapping int

// different ways of mapping @extra in legacy flows
const (
	ExtraAsWebhookJSON ExtraVarsMapping = iota
	ExtraAsTriggerParams
	ExtraAsFunction
)

type varMapper struct {
	// subitems that should be replaced completely with the given strings
	substitutions map[string]string

	// base for fixed subitems, e.g. "contact"
	base string

	// recognized fixed subitems, e.g. "name" or "uuid"
	baseVars map[string]interface{}

	// nesting for arbitrary subitems, e.g. contact fields or run results
	arbitraryNesting string

	// mapper for each arbitrary item
	arbitraryVars map[string]interface{}
}

// returns a copy of this mapper with a prefix applied to the previous base
func (v *varMapper) rebase(prefix string) *varMapper {
	var newBase string
	if prefix != "" {
		newBase = fmt.Sprintf("%s.%s", prefix, v.base)
	} else {
		newBase = v.base
	}
	return &varMapper{
		substitutions:    v.substitutions,
		base:             newBase,
		baseVars:         v.baseVars,
		arbitraryNesting: v.arbitraryNesting,
		arbitraryVars:    v.arbitraryVars,
	}
}

// Resolve resolves the given key to a mapped expression
func (v *varMapper) Resolve(key string) types.XValue {

	// is this a complete substitution?
	if substitute, ok := v.substitutions[key]; ok {
		return types.NewXText(substitute)
	}

	newPath := make([]string, 0, 1)

	if v.base != "" {
		newPath = append(newPath, v.base)
	}

	// is it a fixed base item?
	value, ok := v.baseVars[key]
	if ok {
		// subitem may be a mapper itself
		asVarMapper, isVarMapper := value.(*varMapper)
		if isVarMapper {
			if len(newPath) > 0 {
				return asVarMapper.rebase(strings.Join(newPath, "."))
			}
			return asVarMapper
		}

		asExtraMapper, isExtraMapper := value.(*extraMapper)
		if isExtraMapper {
			return asExtraMapper
		}

		// or a simple string in which case we add to the end of the path and return that
		newPath = append(newPath, value.(string))
		return types.NewXText(strings.Join(newPath, "."))
	}

	// then it must be an arbitrary item
	if v.arbitraryNesting != "" {
		newPath = append(newPath, v.arbitraryNesting)
	}

	newPath = append(newPath, key)

	if v.arbitraryVars != nil {
		return &varMapper{
			base:     strings.Join(newPath, "."),
			baseVars: v.arbitraryVars,
		}
	}

	return types.NewXText(strings.Join(newPath, "."))
}

// Reduce is called when this object needs to be reduced to a primitive
func (v *varMapper) Reduce() types.XPrimitive {
	return types.NewXText(v.String())
}

// ToXJSON won't be called on this but needs to be defined
func (v *varMapper) ToXJSON() types.XText { return types.XTextEmpty }

func (v *varMapper) String() string {
	sub, exists := v.substitutions["__default__"]
	if exists {
		return sub
	}
	return v.base
}

var _ types.XValue = (*varMapper)(nil)
var _ types.XResolvable = (*varMapper)(nil)

// Migration of @extra requires its own mapper because it can map differently depending on the containing flow
type extraMapper struct {
	varMapper

	path    string
	extraAs ExtraVarsMapping
}

// Resolve resolves the given key to a new expression
func (m *extraMapper) Resolve(key string) types.XValue {
	newPath := []string{}
	if m.path != "" {
		newPath = append(newPath, m.path)
	}
	newPath = append(newPath, key)
	return &extraMapper{extraAs: m.extraAs, path: strings.Join(newPath, ".")}
}

// Reduce is called when this object needs to be reduced to a primitive
func (m *extraMapper) Reduce() types.XPrimitive {
	switch m.extraAs {
	case ExtraAsWebhookJSON:
		return types.NewXText(fmt.Sprintf("run.webhook.json.%s", m.path))
	case ExtraAsTriggerParams:
		return types.NewXText(fmt.Sprintf("trigger.params.%s", m.path))
	case ExtraAsFunction:
		return types.NewXText(fmt.Sprintf("if(is_error(run.webhook.json.%s), trigger.params.%s, run.webhook.json.%s)", m.path, m.path, m.path))
	}
	return types.XTextEmpty
}

var _ types.XValue = (*extraMapper)(nil)
var _ types.XResolvable = (*extraMapper)(nil)

type functionTemplate struct {
	name   string
	params string
	join   string
	two    string
	three  string
	four   string
}

var functionTemplates = map[string]functionTemplate{
	"first_word": {name: "word", params: "(%s, 0)"},
	"datevalue":  {name: "date"},
	"edate":      {name: "date_add", params: "(%s, %s, \"M\")"},
	"word":       {name: "word", params: "(%s, %s - 1)"},
	"word_slice": {name: "word_slice", params: "(%s, %s - 1)", three: "(%s, %s - 1, %s - 1)"},
	"field":      {name: "field", params: "(%s, %s - 1, %s)"},
	"datedif":    {name: "date_diff"},
	"date":       {name: "date", params: "(\"%s-%s-%s\")"},
	"days":       {name: "date_diff", params: "(%s, %s, \"D\")"},
	"now":        {name: "now", params: "()"},
	"average":    {name: "mean"},
	"fixed":      {name: "format_num", params: "(%s)", two: "(%s, %s)", three: "(%s, %s, %v)"},

	"roundup":     {name: "round_up"},
	"int":         {name: "round_down"},
	"rounddown":   {name: "round_down"},
	"randbetween": {name: "rand_between"},
	"rept":        {name: "repeat"},

	"year":   {name: "format_date", params: `(%s, "YYYY")`},
	"month":  {name: "format_date", params: `(%s, "M")`},
	"day":    {name: "format_date", params: `(%s, "D")`},
	"hour":   {name: "format_date", params: `(%s, "h")`},
	"minute": {name: "format_date", params: `(%s, "m")`},
	"second": {name: "format_date", params: `(%s, "s")`},

	"proper": {name: "title"},

	// we drop this function, instead joining with the cat operator
	"concatenate": {join: " & "},
	"len":         {name: "length"},

	// translate to maths
	"power": {params: "%s ^ %s"},
	"sum":   {params: "%s + %s"},

	// this one is a special case format, we sum these parts into seconds for date_add
	"time": {name: "time", params: "(%s %s %s)"},
}

func newMigrationBaseVars() map[string]interface{} {
	contact := &varMapper{
		base: "contact",
		baseVars: map[string]interface{}{
			"uuid":       "uuid",
			"name":       "name",
			"first_name": "first_name",
			"language":   "language",
			"groups":     "groups",
			"tel_e164":   "urns.tel.0.path",
		},
		arbitraryNesting: "fields",
	}

	for scheme := range urns.ValidSchemes {
		contact.baseVars[scheme] = &varMapper{
			substitutions: map[string]string{
				"__default__": fmt.Sprintf("format_urn(contact.urns.%s.0)", scheme),
				"display":     fmt.Sprintf("format_urn(contact.urns.%s.0)", scheme),
				"scheme":      fmt.Sprintf("contact.urns.%s.0.scheme", scheme),
				"path":        fmt.Sprintf("contact.urns.%s.0.path", scheme),
				"urn":         fmt.Sprintf("contact.urns.%s.0", scheme),
			},
			base: fmt.Sprintf("urns.%s", scheme),
		}
	}

	return map[string]interface{}{
		"contact": contact,
		"flow": &varMapper{
			base: "run.results",
			arbitraryVars: map[string]interface{}{
				"category": "category_localized",
			},
		},
		"parent": &varMapper{
			base: "parent",
			baseVars: map[string]interface{}{
				"contact": contact,
			},
			arbitraryNesting: "results",
			arbitraryVars: map[string]interface{}{
				"category": "category_localized",
			},
		},
		"child": &varMapper{
			base: "child",
			baseVars: map[string]interface{}{
				"contact": contact,
			},
			arbitraryNesting: "results",
			arbitraryVars: map[string]interface{}{
				"category": "category_localized",
			},
		},
		"step": &varMapper{
			substitutions: map[string]string{
				"__default__": "run.input",
				"value":       "run.input",
				"text":        "run.input.text",
				"attachments": "run.input.attachments",
				"time":        "run.input.created_on",
			},
			baseVars: map[string]interface{}{
				"contact": contact,
			},
		},
		"channel": &varMapper{
			substitutions: map[string]string{
				"__default__": "contact.channel.address",
				"name":        "contact.channel.name",
				"tel":         "contact.channel.address",
				"tel_e164":    "contact.channel.address",
			},
		},
		"date": &varMapper{
			substitutions: map[string]string{
				"__default__": `now()`,
				"now":         `now()`,
				"today":       `today()`,
				"tomorrow":    `date_add(today(), 1, "D")`,
				"yesterday":   `date_add(today(), -1, "D")`,
			},
		},
	}
}

var migrationBaseVars = newMigrationBaseVars()

// creates a new var mapper for migrating expressions
func newMigrationVarMapper(extraAs ExtraVarsMapping) *varMapper {
	// copy the base migration vars
	baseVars := make(map[string]interface{})
	for k, v := range migrationBaseVars {
		baseVars[k] = v
	}

	// add a mapper for extra
	baseVars["extra"] = &extraMapper{extraAs: extraAs}

	return &varMapper{baseVars: baseVars}
}

var datePrefixes = []string{
	"today()",
	"yesterday()",
	"now()",
	"time",
	"timevalue",
}

var ignoredFunctions = map[string]bool{
	"time": true,
	"sum":  true,

	// in some cases we actually remove function names
	// such add switching CONCAT to a simple operator expression
	"": true,
}

func isDate(operand string) bool {
	for i := range datePrefixes {
		if strings.HasPrefix(operand, datePrefixes[i]) {
			return true
		}
	}
	return false
}

func wrapRawExpression(raw string) string {
	for _, topLevel := range runs.RunContextTopLevels {
		if strings.HasPrefix(raw, topLevel+".") || raw == topLevel {
			return "@" + raw
		}
	}
	return fmt.Sprintf("@(%s)", raw)
}

// MigrateTemplate will take a legacy expression and translate it to the new syntax
func MigrateTemplate(template string, extraAs ExtraVarsMapping) (string, error) {
	migrationVarMapper := newMigrationVarMapper(extraAs)

	return migrateLegacyTemplateAsString(migrationVarMapper, template)
}

func migrateLegacyTemplateAsString(resolver types.XValue, template string) (string, error) {
	var buf bytes.Buffer
	var errors excellent.TemplateErrors
	scanner := excellent.NewXScanner(strings.NewReader(template), legacyContextTopLevels)

	for tokenType, token := scanner.Scan(); tokenType != excellent.EOF; tokenType, token = scanner.Scan() {
		switch tokenType {
		case excellent.BODY:
			buf.WriteString(token)
		case excellent.IDENTIFIER:
			value := excellent.ResolveValue(nil, resolver, token)
			if value == nil {
				errors = append(errors, fmt.Errorf("Invalid key: '%s'", token))
				buf.WriteString("@")
				buf.WriteString(token)
			} else {
				strValue, _ := toString(value)

				if strValue == token {
					buf.WriteString("@" + token)
				} else {
					// if expression has been changed, then it might need to be wrapped in @(...)
					buf.WriteString(wrapRawExpression(strValue))
				}
			}

		case excellent.EXPRESSION:
			value, err := translateExpression(nil, resolver, token)
			buf.WriteString("@(")
			if err != nil {
				errors = append(errors, err)
			} else {
				strValue, err := toString(value)
				if err != nil {
					buf.WriteString(token)
					errors = append(errors, err)
				} else {
					buf.WriteString(strValue)
				}
			}
			buf.WriteString(")")
		}
	}

	if len(errors) > 0 {
		return buf.String(), errors
	}
	return buf.String(), nil
}

func toString(params interface{}) (string, error) {
	switch typed := params.(type) {
	case types.XValue:
		str, xerr := types.ToXText(typed)
		return str.Native(), xerr
	case string:
		return typed, nil

	case []interface{}:
		strArr := make([]string, len(typed))
		for i := range strArr {
			str, err := toString(typed[i])
			if err != nil {
				return "", err
			}
			strArr[i] = str
		}
		return strings.Join(strArr, ", "), nil
	default:
		panic(fmt.Sprintf("can't toString a %s", typed))
	}
}

// translateExpression will turn an old expression into a new format expression
func translateExpression(env utils.Environment, resolver types.XValue, template string) (interface{}, error) {
	errors := excellent.NewErrorListener()

	input := antlr.NewInputStream(template)
	lexer := gen.NewExcellent2Lexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := gen.NewExcellent2Parser(stream)
	p.AddErrorListener(errors)

	// speed up parsing
	p.SetErrorHandler(antlr.NewBailErrorStrategy())
	p.GetInterpreter().SetPredictionMode(antlr.PredictionModeSLL)
	// TODO: add second stage - https://github.com/antlr/antlr4/issues/192

	// this is still super slow, awaiting fixes in golang antlr
	// leaving this debug in until then
	// start := time.Now()
	tree := p.Parse()
	// timeTrack(start, "Parsing")

	// if we ran into errors parsing, bail
	if errors.HasErrors() {
		return nil, fmt.Errorf(errors.Errors())
	}

	visitor := newLegacyVisitor(env, resolver)
	value := visitor.Visit(tree)
	err, isErr := value.(error)

	// did our evaluation result in an error? return that
	if isErr {
		return nil, err
	}

	// all is good, return our value
	return value, nil
}

// ---------------------------------------------------------------
// Visitors
// ---------------------------------------------------------------

type legacyVisitor struct {
	gen.BaseExcellent2Visitor
	env      utils.Environment
	resolver types.XValue
}

func newLegacyVisitor(env utils.Environment, resolver types.XValue) *legacyVisitor {
	return &legacyVisitor{env: env, resolver: resolver}
}

// ---------------------------------------------------------------

// Visit the top level parse tree
func (v *legacyVisitor) Visit(tree antlr.ParseTree) interface{} {
	return tree.Accept(v)
}

// VisitParse handles our top level parser
func (v *legacyVisitor) VisitParse(ctx *gen.ParseContext) interface{} {
	return v.Visit(ctx.Expression())
}

// VisitDecimalLiteral deals with decimals like 1.5
func (v *legacyVisitor) VisitDecimalLiteral(ctx *gen.DecimalLiteralContext) interface{} {
	dec, _ := toString(ctx.GetText())
	return dec
}

// VisitDotLookup deals with lookups like foo.0 or foo.bar
func (v *legacyVisitor) VisitDotLookup(ctx *gen.DotLookupContext) interface{} {
	value := v.Visit(ctx.Atom(0)).(types.XValue)
	expression := v.Visit(ctx.Atom(1)).(types.XValue)
	lookup, err := types.ToXText(expression)
	if err != nil {
		return err
	}
	return excellent.ResolveValue(v.env, value, lookup.Native())
}

// VisitStringLiteral deals with string literals such as "asdf"
func (v *legacyVisitor) VisitStringLiteral(ctx *gen.StringLiteralContext) interface{} {
	return ctx.GetText()
}

// VisitFunctionCall deals with function calls like TITLE(foo.bar)
func (v *legacyVisitor) VisitFunctionCall(ctx *gen.FunctionCallContext) interface{} {
	functionName := strings.ToLower(ctx.Fnname().GetText())
	template, found := functionTemplates[functionName]
	if !found {
		template = functionTemplate{name: functionName, params: "(%s)"}
	} else {
		if template.params == "" {
			template.params = "(%s)"
		}
	}

	_, ignored := ignoredFunctions[template.name]
	if !ignored {
		_, found = functions.XFUNCTIONS[template.name]
		if !found {
			return fmt.Errorf("No function with name '%s'", template.name)
		}
	}

	var params []interface{}
	if ctx.Parameters() != nil {
		funcParams := v.Visit(ctx.Parameters())
		switch funcParams.(type) {
		case error:
			return funcParams
		default:
			params = funcParams.([]interface{})
		}
	}

	// special case options for 3 or 4 parameters
	paramTemplate := template.params
	if len(params) == 3 && template.three != "" {
		paramTemplate = template.three
	}

	if len(params) == 4 && template.four != "" {
		paramTemplate = template.four
	}

	if template.join != "" {
		// if our template wants a join, do that instead
		toJoin := make([]string, len(params))
		for i := range params {
			p, err := toString(params[i])
			if err == nil {
				toJoin[i] = p
			}
		}

		paramTemplate = "%s"
		params = make([]interface{}, 1)
		params[0] = strings.Join(toJoin, template.join)
	} else {
		// how many replacements we are expecting
		replacementCount := strings.Count(paramTemplate, "%s") + strings.Count(paramTemplate, "%v")

		if replacementCount != len(params) {
			// if our params don't match our template, turn stringify it
			p, err := toString(params)
			if err != nil {
				return err
			}
			params = make([]interface{}, 1)
			params[0] = p
		}
	}

	return fmt.Sprintf("%s%s", template.name, fmt.Sprintf(paramTemplate, params...))
}

// VisitTrue deals with the "true" literal
func (v *legacyVisitor) VisitTrue(ctx *gen.TrueContext) interface{} {
	return true
}

// VisitFalse deals with the "false" literal
func (v *legacyVisitor) VisitFalse(ctx *gen.FalseContext) interface{} {
	return false
}

// VisitArrayLookup deals with lookups such as foo[5]
func (v *legacyVisitor) VisitArrayLookup(ctx *gen.ArrayLookupContext) interface{} {
	value := v.Visit(ctx.Atom()).(types.XValue)
	expression := v.Visit(ctx.Expression()).(types.XValue)
	lookup, err := types.ToXText(expression)
	if err != nil {
		return err
	}
	return excellent.ResolveValue(v.env, value, lookup.Native())
}

// VisitContextReference deals with references to variables in the context such as "foo"
func (v *legacyVisitor) VisitContextReference(ctx *gen.ContextReferenceContext) interface{} {
	key := strings.ToLower(ctx.GetText())
	val := excellent.ResolveValue(v.env, v.resolver, key)
	if val == nil {
		return fmt.Errorf("Invalid key: '%s'", key)
	}

	err, isErr := val.(error)
	if isErr {
		return err
	}

	return val
}

// VisitParentheses deals with expressions in parentheses such as (1+2)
func (v *legacyVisitor) VisitParentheses(ctx *gen.ParenthesesContext) interface{} {
	return fmt.Sprintf("(%s)", v.Visit(ctx.Expression()))
}

// VisitNegation deals with negations such as -5
func (v *legacyVisitor) VisitNegation(ctx *gen.NegationContext) interface{} {
	dec, err := toString(v.Visit(ctx.Expression()))
	if err != nil {
		return err
	}
	return "-" + dec
}

// VisitExponent deals with exponenets such as 5^5
func (v *legacyVisitor) VisitExponent(ctx *gen.ExponentContext) interface{} {
	arg1, err := toString(v.Visit(ctx.Expression(0)))
	if err != nil {
		return err
	}

	arg2, err := toString(v.Visit(ctx.Expression(1)))
	if err != nil {
		return err
	}

	return fmt.Sprintf("%s ^ %s", arg1, arg2)
}

// VisitConcatenation deals with string concatenations like "foo" & "bar"
func (v *legacyVisitor) VisitConcatenation(ctx *gen.ConcatenationContext) interface{} {
	arg1, err := toString(v.Visit(ctx.Expression(0)))
	if err != nil {
		return err
	}

	arg2, err := toString(v.Visit(ctx.Expression(1)))
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	buffer.WriteString(arg1)
	buffer.WriteString(" & ")
	buffer.WriteString(arg2)

	return buffer.String()
}

// VisitAdditionOrSubtraction deals with addition and subtraction like 5+5 and 5-3
func (v *legacyVisitor) VisitAdditionOrSubtraction(ctx *gen.AdditionOrSubtractionContext) interface{} {
	value, err := toString(v.Visit(ctx.Expression(0)))
	if err != nil {
		return err
	}

	dateUnit := "D"
	firstIsDate := isDate(value)
	if firstIsDate {
		firstSeconds, ok := convertTimeToSeconds(value)
		if ok {
			value = firstSeconds
			dateUnit = "s"
		}
	}

	// see if our first param is an int
	_, firstNumberErr := strconv.Atoi(value)

	next, err := toString(v.Visit(ctx.Expression(1)))
	if err != nil {
		return err
	}

	op := "+"
	if ctx.MINUS() != nil {
		op = "-"
	}

	secondIsDate := isDate(next)
	if secondIsDate {
		secondSeconds, ok := convertTimeToSeconds(next)
		if ok {
			next = secondSeconds
			dateUnit = "s"
		}
	}

	// see if our second param is an int
	_, secondNumberErr := strconv.Atoi(next)
	if (firstIsDate || secondIsDate) && (firstNumberErr != nil || secondNumberErr != nil) {

		// we are adding two values where we know at least one side is a date
		template := "date_add(%s, %s, \"%s\")"
		if op == "-" {
			template = "date_add(%s, -%s, \"%s\")"
		}

		// determine the order of our parameters
		replacements := []interface{}{value, next, dateUnit}
		if firstNumberErr == nil {
			replacements = []interface{}{next, value, dateUnit}
		}

		value = fmt.Sprintf(template, replacements...)

	} else if firstNumberErr == nil && secondNumberErr == nil {
		// we are adding two numbers
		if op == "+" {
			value = fmt.Sprintf("%s + %s", value, next)
		} else {
			value = fmt.Sprintf("%s - %s", value, next)
		}
	} else {
		// we are adding a field of unknown type with an integer
		if op == "+" {
			value = fmt.Sprintf("legacy_add(%s, %s)", value, next)
		} else {
			value = fmt.Sprintf("legacy_add(%s, -%s)", value, next)
		}
	}

	return value
}

// VisitEquality deals with equality or inequality tests 5 = 5 and 5 != 5
func (v *legacyVisitor) VisitEquality(ctx *gen.EqualityContext) interface{} {
	arg1 := v.Visit(ctx.Expression(0))
	err, isErr := arg1.(error)
	if isErr {
		return err
	}

	arg2 := v.Visit(ctx.Expression(1))
	err, isErr = arg2.(error)
	if isErr {
		return err
	}

	if ctx.EQ() != nil {
		return fmt.Sprintf("%s = %s", arg1, arg2)
	}

	return fmt.Sprintf("%s != %s", arg1, arg2)
}

// VisitAtomReference deals with visiting a single atom in our expression
func (v *legacyVisitor) VisitAtomReference(ctx *gen.AtomReferenceContext) interface{} {
	return v.Visit(ctx.Atom())
}

// VisitMultiplicationOrDivision deals with division and multiplication such as 5*5 or 5/2
func (v *legacyVisitor) VisitMultiplicationOrDivision(ctx *gen.MultiplicationOrDivisionContext) interface{} {
	arg1 := v.Visit(ctx.Expression(0))
	str1, err := toString(arg1)
	if err != nil {
		return err
	}

	arg2 := v.Visit(ctx.Expression(1))
	str2, err := toString(arg2)
	if err != nil {
		return err
	}

	if ctx.TIMES() != nil {
		return fmt.Sprintf("%s * %s", str1, str2)
	}

	return fmt.Sprintf("%s / %s", str1, str2)
}

// VisitComparison deals with visiting a comparison between two values, such as 5<3 or 3>5
func (v *legacyVisitor) VisitComparison(ctx *gen.ComparisonContext) interface{} {
	arg1 := v.Visit(ctx.Expression(0))
	arg2 := v.Visit(ctx.Expression(1))

	err, isErr := arg1.(error)
	if isErr {
		return err
	}

	err, isErr = arg2.(error)
	if isErr {
		return err
	}

	return fmt.Sprintf("%s %s %s", arg1, ctx.GetOp().GetText(), arg2)
}

// VisitFunctionParameters deals with the parameters to a function call
func (v *legacyVisitor) VisitFunctionParameters(ctx *gen.FunctionParametersContext) interface{} {
	expressions := ctx.AllExpression()
	params := make([]interface{}, len(expressions))

	for i := range expressions {
		params[i] = v.Visit(expressions[i])
		error, isError := params[i].(error)
		if isError {
			return error
		}
	}
	return params
}

// convertTimeToSeconds takes a old TIME(0,2,5) like expression
// and returns the numeric value in seconds
func convertTimeToSeconds(operand string) (string, bool) {
	converted := false
	if strings.HasPrefix(operand, "time(") {
		var hours, minutes, seconds int
		parsed, _ := fmt.Sscanf(operand, "time(%d %d %d)", &hours, &minutes, &seconds)
		if parsed == 3 {
			operand = strconv.Itoa(seconds + (minutes * 60) + (hours * 3600))
			converted = true
		}
	}
	return operand, converted
}
