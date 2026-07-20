// Package parse implements parsing of ContactQL queries using the ANTLR generated parser. It is separate
// from the contactql package so that consumers which only evaluate or inspect already parsed queries don't
// depend on the ANTLR runtime.
package parse

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/antlr4-go/antlr/v4"
	gen "github.com/nyaruka/goflow/antlr/gen/contactql"
	"github.com/nyaruka/goflow/contactql"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils"
)

// maxQueryDepth is the maximum bracket nesting depth allowed in a query. Parsing and walking the query
// tree are recursive, so without a limit a deeply nested query can overflow the stack and crash the
// process. Real queries are written by humans and nest a handful of levels deep at most.
const maxQueryDepth = 100

// maxQueryLength is the maximum length of query text we'll attempt to parse. The cost of parsing is
// proportional to the length of the input, and the bracket depth limit doesn't bound a query that's simply
// long rather than nested, so without this a huge query is fully parsed before being rejected for having
// too many conditions. Matches the limit applied by callers.
const maxQueryLength = 10_000

// Query parses a ContactQL query from the given input. If resolver is provided then we validate against it
// to ensure that fields and groups exist. If not provided then still validate what we can.
func Query(env envs.Environment, text string, resolver contactql.Resolver) (*contactql.ContactQuery, error) {
	// preprocess text before parsing
	text = strings.TrimSpace(text)

	// reject overly long queries before parsing. Length is counted in characters to match how callers count
	// it, and the byte length is checked first to short circuit as it's never less than the character count.
	if len(text) > maxQueryLength && utf8.RuneCountInString(text) > maxQueryLength {
		return nil, contactql.NewQueryError(contactql.ErrTooComplex, "query is too complex")
	}

	// reject overly nested queries before parsing to avoid a stack overflow
	if utils.NestingDepthExceeds(text, maxQueryDepth) {
		return nil, contactql.NewQueryError(contactql.ErrTooComplex, "query is too complex")
	}

	// if query is a valid number, rewrite as a tel = query
	if env.RedactionPolicy() != envs.RedactionPolicyURNs {
		if number := utils.ParsePhoneNumber(text, env.DefaultCountry()); number != "" {
			text = fmt.Sprintf(`tel = %s`, number)
		}
	}

	errListener := &errorListener{}
	input := antlr.NewInputStream(text)
	lexer := gen.NewContactQLLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := gen.NewContactQLParser(stream)
	p.RemoveErrorListeners()
	p.AddErrorListener(errListener)
	tree := p.Parse()

	// if we ran into errors parsing, bail
	if err := errListener.Error(); err != nil {
		return nil, err
	}

	visitor := newVisitor(env)
	rootNode := visitor.Visit(tree).(contactql.QueryNode)

	if len(visitor.errors) > 0 {
		return nil, visitor.errors[0]
	}

	return contactql.NewContactQuery(env, rootNode, resolver)
}

type errorListener struct {
	*antlr.DefaultErrorListener

	errs []*contactql.QueryError
}

func (l *errorListener) Error() error {
	if len(l.errs) > 0 {
		return l.errs[0]
	}
	return nil
}

func (l *errorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol any, line, column int, msg string, e antlr.RecognitionException) {
	l.errs = append(l.errs, contactql.NewQueryError(contactql.ErrSyntax, msg))
}
