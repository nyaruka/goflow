package excellent

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/excellent/gen"
	"github.com/nyaruka/goflow/utils"
)

var topLevelScopes = []string{"contact", "child", "parent", "run", "trigger"}

// ExtraVarsMapping defines how @extra.* variables should be migrated
type ExtraVarsMapping string

// different ways of mapping @extra in legacy flows
const (
	ExtraAsWebhookJSON   ExtraVarsMapping = "run.webhook.json"
	ExtraAsTriggerParams ExtraVarsMapping = "trigger.params"
	ExtraAsFunction      ExtraVarsMapping = "IF(trigger.params.%s, trigger.params.%s, run.webhook.json.%s)"
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
func (v *varMapper) Resolve(key string) interface{} {

	// is this a complete substitution?
	if substitute, ok := v.substitutions[key]; ok {
		return substitute
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
		return strings.Join(newPath, ".")
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

	return strings.Join(newPath, ".")
}

// Default returns the value of this mapper when it is the result of an expression
func (v *varMapper) Default() interface{} {
	return v.base
}

func (v *varMapper) String() string {
	sub, exists := v.substitutions["__default__"]
	if exists {
		return sub
	}
	return v.base
}

// Migration of @extra requires its own mapper because it can map differently depending on the containing flow
type extraMapper struct {
	varMapper

	path    string
	extraAs ExtraVarsMapping
}

// Resolve resolves the given key to a new expression
func (m *extraMapper) Resolve(key string) interface{} {
	newPath := []string{}
	if m.path != "" {
		newPath = append(newPath, m.path)
	}
	newPath = append(newPath, key)
	return &extraMapper{extraAs: m.extraAs, path: strings.Join(newPath, ".")}
}

// Default returns the value of this extra mapper when it is the result of an expression
func (m *extraMapper) Default() interface{} {
	return m
}

func (m *extraMapper) String() string {
	switch m.extraAs {
	case ExtraAsWebhookJSON:
		return fmt.Sprintf("run.webhook.json.%s", m.path)
	case ExtraAsTriggerParams:
		return fmt.Sprintf("trigger.params.%s", m.path)
	case ExtraAsFunction:
		return fmt.Sprintf("if(has_error(run.webhook.json.%s), trigger.params.%s, run.webhook.json.%s)", m.path, m.path, m.path)
	}
	return ""
}

type functionTemplate struct {
	name   string
	params string
	join   string
	two    string
	three  string
	four   string
}

var functionTemplates = map[string]functionTemplate{
	"first_word": {name: "split", params: "(%s, \" \")[0]"},
	"datevalue":  {name: "date"},
	"edate":      {name: "date_add", params: "(%s, \"m\", %s)"},
	"word_slice": {name: "word_slice", params: "(%s, %s)", three: "(%s, %s, %s)", four: "(%s, %s, %s, %v)"},
	"field":      {name: "field", params: "(%s, %s - 1, %s)"},
	"datedif":    {name: "date_diff"},
	"date":       {name: "date", params: "(\"%s-%s-%s\")"},
	"days":       {name: "date_diff", params: "(%s, %s, \"D\")"},
	"now":        {name: "now", params: "()"},
	"average":    {name: "mean"},
	"fixed":      {name: "format_num", params: "(%s)", two: "(%s, %s)", three: "(%s, %s, %v)"},

	"roundup":     {name: "round_up"},
	"rounddown":   {name: "round_down"},
	"randbetween": {name: "rand"},
	"rept":        {name: "repeat"},

	"year":   {name: "format_date", params: `(%s, "yyyy")`},
	"month":  {name: "format_date", params: `(%s, "M")`},
	"day":    {name: "format_date", params: `(%s, "d")`},
	"hour":   {name: "format_date", params: `(%s, "h")`},
	"minute": {name: "format_date", params: `(%s, "m")`},
	"second": {name: "format_date", params: `(%s, "s")`},

	"proper": {name: "title"},

	// we drop this function, instead joining with the cat operator
	"concatenate": {join: " & "},
	"len":         {name: "string_length"},

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
				"__default__": fmt.Sprintf("format_urn(contact.urns.%s)", scheme),
				"display":     fmt.Sprintf("format_urn(contact.urns.%s)", scheme),
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
				"__default__": "now()",
				"now":         "now()",
				"today":       "today()",
				"tomorrow":    "tomorrow()",
				"yesterday":   "yesterday()",
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
	for _, topLevel := range topLevelScopes {
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

func migrateLegacyTemplateAsString(resolver utils.VariableResolver, template string) (string, error) {
	var buf bytes.Buffer
	var errors TemplateErrors
	scanner := newXScanner(strings.NewReader(template))

	for tokenType, token := scanner.Scan(); tokenType != EOF; tokenType, token = scanner.Scan() {
		switch tokenType {
		case BODY:
			buf.WriteString(token)
		case IDENTIFIER:
			value := utils.ResolveVariable(nil, resolver, token)
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

		case EXPRESSION:
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

// toString defers to main ToString, but will also stringify arrays
func toString(params interface{}) (string, error) {
	switch params := params.(type) {
	case []interface{}:
		strArr, err := utils.ToStringArray(nil, params)
		if err == nil {
			return strings.Join(strArr, ", "), nil
		}
	}
	return utils.ToString(nil, params)
}

// translateExpression will turn an old expression into a new format expression
func translateExpression(env utils.Environment, resolver utils.VariableResolver, template string) (interface{}, error) {
	errors := newErrorListener()

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
	resolver utils.VariableResolver
}

func newLegacyVisitor(env utils.Environment, resolver utils.VariableResolver) *legacyVisitor {
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
	value := v.Visit(ctx.Atom(0))
	lookup, err := utils.ToString(v.env, v.Visit(ctx.Atom(1)))
	if err != nil {
		return err
	}
	return utils.ResolveVariable(v.env, value, lookup)
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
		_, found = XFUNCTIONS[template.name]
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
	value := v.Visit(ctx.Atom())
	lookup, err := utils.ToString(v.env, v.Visit(ctx.Expression()))
	if err != nil {
		return err
	}
	return utils.ResolveVariable(v.env, value, lookup)
}

// VisitContextReference deals with references to variables in the context such as "foo"
func (v *legacyVisitor) VisitContextReference(ctx *gen.ContextReferenceContext) interface{} {
	key := strings.ToLower(ctx.GetText())
	val := utils.ResolveVariable(v.env, v.resolver, key)
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

	dateUnit := "d"
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
		template := "date_add(%s, \"%s\", %s)"
		if op == "-" {
			template = "date_add(%s, \"%s\", -%s)"
		}

		// determine the order of our parameters
		replacements := []interface{}{value, dateUnit, next}
		if firstNumberErr == nil {
			replacements = []interface{}{next, dateUnit, value}
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
		return fmt.Sprintf("%s == %s", arg1, arg2)
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
