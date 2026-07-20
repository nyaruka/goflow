package excellent

import (
	"context"
	"fmt"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	gen "github.com/nyaruka/goflow/antlr/gen/excellent3"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/budget"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
)

// maxExpressionDepth is the maximum bracket nesting depth allowed in an expression. Parsing and evaluating
// are recursive, so without a limit a deeply nested expression can overflow the stack and crash the process.
// Real expressions are written by humans and nest a handful of levels deep at most.
const maxExpressionDepth = 100

// maxEvaluationCost is the cost budget for a single expression evaluation. Cost accrues as values are
// produced - text costs its length in bytes, everything else costs 1 - so this bounds both the memory and
// the number of operations a single evaluation can consume, no matter how per-function limits are composed.
// It's set generously: real expressions cost a few hundred units at most, so this is many orders of magnitude
// of headroom whilst still bounding an attack to a few MB. It can be tightened later based on real-world usage.
const maxEvaluationCost = 10_000_000

// Evaluator evaluates templates and expressions.
type Evaluator struct{}

// NewEvaluator creates a new evaluator
func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

// Escaping is a function applied to expressions in a template after they've been evaluated
type Escaping func(string) string

// Template evaluates the passed in template
func (e *Evaluator) Template(ctx context.Context, env envs.Environment, root *types.XObject, template string, escaping Escaping) (string, []string, error) {
	var buf strings.Builder
	var allWarnings []string

	err := VisitTemplate(template, root.Properties(), true, func(tokenType XTokenType, token string) error {
		switch tokenType {
		case BODY:
			buf.WriteString(token)
		case IDENTIFIER, EXPRESSION:
			value, warnings := e.Expression(ctx, env, root, token)

			allWarnings = append(allWarnings, warnings...)

			// if we got an error, return that
			if types.IsXError(value) {
				return value.(error)
			}

			// if not, stringify value and append to the output
			asText, _ := types.ToXText(env, value)
			asString := asText.Native()

			if escaping != nil {
				asString = escaping(asString)
			}

			buf.WriteString(asString)
		}
		return nil
	})

	return buf.String(), allWarnings, err
}

// TemplateValue is equivalent to Template except in the case where the template contains
// a single identifier or expression, ie: "@contact" or "@(first(contact.urns))". In these cases we return
// the typed value from EvaluateExpression instead of stringifying the result.
func (e *Evaluator) TemplateValue(ctx context.Context, env envs.Environment, root *types.XObject, template string) (types.XValue, []string, error) {
	template = strings.TrimSpace(template)
	scanner := NewXScanner(strings.NewReader(template), root.Properties())

	// parse our first token
	tokenType, token := scanner.Scan()

	// try to scan to our next token
	nextTT, _ := scanner.Scan()

	// if we only have an identifier or an expression, evaluate it on its own
	if nextTT == EOF {
		switch tokenType {
		case IDENTIFIER, EXPRESSION:
			val, warnings := e.Expression(ctx, env, root, token)
			return val, warnings, nil
		}
	}

	// otherwise fallback to full template evaluation
	asStr, warnings, err := e.Template(ctx, env, root, template, nil)
	return types.NewXText(asStr), warnings, err
}

// Expression evalutes the passed in Excellent expression, returning the typed value it evaluates to,
// which might be an error, e.g. "2 / 3" or "contact.fields.age"
func (e *Evaluator) Expression(ctx context.Context, env envs.Environment, root *types.XObject, expression string) (types.XValue, []string) {
	parsed, err := Parse(expression, nil)
	if err != nil {
		return types.NewXError(err), nil
	}

	scope := NewScope(root, nil)

	warnings := &Warnings{}

	// a per-evaluation cost budget is added to the caller's context so that its deadline (if any) is honoured
	// alongside the budget
	ctx = budget.With(ctx, budget.New(maxEvaluationCost))

	return parsed.Evaluate(ctx, env, scope, warnings), warnings.all
}

// Parse parses an expression
func Parse(expression string, contextCallback func([]string)) (Expression, error) {
	// reject overly nested expressions before parsing to avoid a stack overflow
	if utils.NestingDepthExceeds(expression, maxExpressionDepth) {
		return nil, fmt.Errorf("expression nesting too deep")
	}

	errListener := NewErrorListener(expression)

	input := antlr.NewInputStream(expression)
	lexer := gen.NewExcellent3Lexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := gen.NewExcellent3Parser(stream)
	p.RemoveErrorListeners()
	p.AddErrorListener(errListener)
	tree := p.Parse()

	// if we ran into errors parsing, return the first one
	if len(errListener.Errors()) > 0 {
		return nil, errListener.Errors()[0]
	}

	visitor := &visitor{contextCallback: contextCallback}
	output := visitor.Visit(tree)
	return toExpression(output), nil
}

// VisitTemplate scans the given template and calls the callback for each token encountered
func VisitTemplate(template string, allowedTopLevels []string, unescapeBody bool, callback func(XTokenType, string) error) error {
	// nothing todo for an empty template
	if template == "" {
		return nil
	}

	scanner := NewXScanner(strings.NewReader(template), allowedTopLevels)
	scanner.SetUnescapeBody(unescapeBody)
	errors := NewTemplateErrors()

	for tokenType, token := scanner.Scan(); tokenType != EOF; tokenType, token = scanner.Scan() {
		if err := callback(tokenType, token); err != nil {
			var repr string
			if tokenType == IDENTIFIER {
				repr = "@" + token
			} else {
				repr = "@(" + token + ")"
			}

			errors.Add(repr, err.Error())
		}
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}

// HasExpressions returns whether the given template contains any expressions or identifiers
func HasExpressions(template string, allowedTopLevels []string) bool {
	found := false
	VisitTemplate(template, allowedTopLevels, false, func(tokenType XTokenType, token string) error {
		switch tokenType {
		case IDENTIFIER, EXPRESSION:
			found = true
			return nil
		}
		return nil
	})
	return found
}
