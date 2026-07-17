package excellent

import (
	"slices"
	"strings"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
)

// Template is a parsed representation of a template string, e.g. "Hi @contact.name!". It is immutable after parsing
// and thus safe for concurrent use - parsing a template once and evaluating it many times is equivalent to passing
// the original string to Evaluator.Template each time.
type Template struct {
	segments []templateSegment
}

// ParseTemplate parses the given template source. It never returns an error - a segment which is a malformed
// expression is retained as-is and generates the same error when the template is evaluated.
func ParseTemplate(src string) *Template {
	t := &Template{}

	// scan the source twice - once without unescaping so that text segments retain their original source including
	// @@ sequences, and once with unescaping to get the values that text segments evaluate to.. the two scans always
	// produce the same sequence of tokens, differing only in the text of BODY tokens
	rawScanner := NewXScanner(strings.NewReader(src), nil)
	rawScanner.SetUnescapeBody(false)
	valScanner := NewXScanner(strings.NewReader(src), nil)

	for {
		tokenType, rawToken := rawScanner.Scan()
		_, valToken := valScanner.Scan()

		switch tokenType {
		case EOF:
			return t
		case BODY:
			t.segments = append(t.segments, &textSegment{src: rawToken, val: valToken})
		case IDENTIFIER:
			expression, err := Parse(rawToken, nil)
			topLevel, _, _ := strings.Cut(rawToken, ".")
			t.segments = append(t.segments, &expressionSegment{
				src:        "@" + rawToken,
				topLevel:   strings.ToLower(topLevel),
				expression: expression,
				err:        err,
			})
		case EXPRESSION:
			expression, err := Parse(rawToken, nil)
			t.segments = append(t.segments, &expressionSegment{
				src:        "@(" + rawToken + ")",
				expression: expression,
				err:        err,
			})
		}
	}
}

// String returns the original source of this template, reconstructed from the retained source of each segment.
func (t *Template) String() string {
	var buf strings.Builder
	for _, s := range t.segments {
		buf.WriteString(s.source())
	}
	return buf.String()
}

// Evaluate evaluates this template against the given context, producing the same output, warnings and error as
// passing the original source to Evaluator.Template.
func (t *Template) Evaluate(env envs.Environment, ctx *types.XObject, escaping Escaping) (string, []string, error) {
	var buf strings.Builder
	var warnings []string
	errors := NewTemplateErrors()
	topLevels := ctx.Properties()

	for _, segment := range t.segments {
		switch s := segment.(type) {
		case *textSegment:
			buf.WriteString(s.val)
		case *expressionSegment:
			// an identifier whose top-level isn't a property of the context is treated as literal text,
			// e.g. the @gmail in an email address
			if s.topLevel != "" && !slices.Contains(topLevels, s.topLevel) {
				buf.WriteString(s.src)
				continue
			}

			value, segWarnings := s.evaluate(env, ctx)
			warnings = append(warnings, segWarnings...)

			// if we got an error, add it to our list and move on
			if types.IsXError(value) {
				errors.Add(s.src, value.(error).Error())
				continue
			}

			// if not, stringify value and append to the output
			asText, _ := types.ToXText(env, value)
			asString := asText.Native()

			if escaping != nil {
				asString = escaping(asString)
			}

			buf.WriteString(asString)
		}
	}

	if errors.HasErrors() {
		return buf.String(), warnings, errors
	}
	return buf.String(), warnings, nil
}

// EvaluateValue evaluates this template and returns a typed value, producing the same output, warnings and error
// as passing the original source to Evaluator.TemplateValue - except that it doesn't trim the source, so callers
// wanting that behavior should parse the trimmed source.
func (t *Template) EvaluateValue(env envs.Environment, ctx *types.XObject) (types.XValue, []string, error) {
	// if the template is a single expression, return the typed value it evaluates to
	if len(t.segments) == 1 {
		if s, ok := t.segments[0].(*expressionSegment); ok {
			if s.topLevel == "" || slices.Contains(ctx.Properties(), s.topLevel) {
				value, warnings := s.evaluate(env, ctx)
				return value, warnings, nil
			}
		}
	}

	// otherwise fallback to full template evaluation
	asStr, warnings, err := t.Evaluate(env, ctx, nil)
	return types.NewXText(asStr), warnings, err
}

// a segment of a parsed template - either literal text or an expression
type templateSegment interface {
	source() string
}

// a segment of literal text
type textSegment struct {
	src string // original source, e.g. "hi @@there"
	val string // unescaped value, e.g. "hi @there"
}

func (s *textSegment) source() string { return s.src }

// a segment which is an expression, in either @identifier or @(...) form
type expressionSegment struct {
	src        string     // original source, e.g. "@contact.name" or "@(1 + 2)"
	topLevel   string     // lowercased top-level of an @identifier segment, empty for @(...) segments
	expression Expression // parsed expression, nil if source is malformed
	err        error      // error from parsing if source is malformed
}

func (s *expressionSegment) source() string { return s.src }

func (s *expressionSegment) evaluate(env envs.Environment, ctx *types.XObject) (types.XValue, []string) {
	if s.err != nil {
		return types.NewXError(s.err), nil
	}

	warnings := &Warnings{}
	value := s.expression.Evaluate(env, NewScope(ctx, nil), warnings)
	return value, warnings.all
}
