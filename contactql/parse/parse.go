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

// Query parses a ContactQL query from the given input. If resolver is provided then we validate against it
// to ensure that fields and groups exist. If not provided then still validate what we can.
func Query(env envs.Environment, text string, resolver contactql.Resolver) (*contactql.ContactQuery, error) {
	// preprocess text before parsing
	text = strings.TrimSpace(text)

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
