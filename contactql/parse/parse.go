// Package parse implements parsing of ContactQL queries using the ANTLR generated parser. It is separate
// from the contactql package so that consumers which only evaluate or inspect already parsed queries don't
// depend on the ANTLR runtime.
package parse

import (
	"fmt"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	gen "github.com/nyaruka/goflow/antlr/gen/contactql"
	"github.com/nyaruka/goflow/contactql"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils"
)

// maxQueryDepth is the maximum bracket nesting depth allowed in a query. Parsing and walking the query
// tree are recursive, so without a limit a deeply nested query can overflow the stack and crash the
// process. The limit is far above anything a real query needs but well below what overflows the stack.
const maxQueryDepth = 250

// Query parses a ContactQL query from the given input. If resolver is provided then we validate against it
// to ensure that fields and groups exist. If not provided then still validate what we can.
func Query(env envs.Environment, text string, resolver contactql.Resolver) (*contactql.ContactQuery, error) {
	// preprocess text before parsing
	text = strings.TrimSpace(text)

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
