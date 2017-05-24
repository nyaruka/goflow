package excellent

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/nyaruka/goflow/excellent/gen"
	"github.com/nyaruka/goflow/utils"
)

type Token int

const (
	BODY       Token = iota // BODY: Not in expression
	IDENTIFIER              // IDENTIFIER: 'contact.age' in '@contact.age'
	EXPRESSION              // EXPRESSION: the body of an expression '1+2' in '@(1+2)'
	EOF                     // EOF: end of expression
)

const eof rune = rune(0)

func isIdentifierChar(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsNumber(ch) || ch == '.' || ch == '_'
}

// XScanner represents a lexical scanner.
type XScanner struct {
	r           *bufio.Reader
	unreadRunes []rune
	unreadCount int
}

// NewXScanner returns a new instance of Scanner.
func NewXScanner(r io.Reader) *XScanner {
	return &XScanner{r: bufio.NewReader(r), unreadRunes: make([]rune, 4)}
}

// gets the next rune or EOF if we are at the end of the string
func (s *XScanner) read() rune {
	// first see if we have any unread runes to return
	if s.unreadCount > 0 {
		ch := s.unreadRunes[s.unreadCount-1]
		s.unreadCount--
		return ch
	}

	// otherwise, read the next run
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// pops the passed in rune as the next rune to be returned
func (s *XScanner) unread(ch rune) {
	s.unreadRunes[s.unreadCount] = ch
	s.unreadCount++
}

// scanExpression consumes the current rune and all contiguous pieces until the end of the expression
// our read should be after the '('
func (s *XScanner) scanExpression() (Token, string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer

	// our parentheses depth
	parens := 1

	// Read every subsequent character until we reach the end of the expression
	for ch := s.read(); ch != eof; ch = s.read() {
		if ch == '(' {
			buf.WriteRune(ch)
			parens++
		} else if ch == ')' {
			parens--

			// we are the end of the expression
			if parens == 0 {
				break
			}
			buf.WriteRune(ch)
		} else {
			buf.WriteRune(ch)
		}
	}

	if parens == 0 {
		return EXPRESSION, buf.String()
	}

	return BODY, strings.Join([]string{"@(", buf.String()}, "")
}

// scanIdentifier consumes the current rune and all contiguous pieces until the end of the identifer
// our read should be after the '@'
func (s *XScanner) scanIdentifier() (Token, string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer

	// Read every subsequent character until we reach the end of the identifier
	for ch := s.read(); ch != eof; ch = s.read() {
		if isIdentifierChar(ch) {
			buf.WriteRune(ch)
		} else {
			s.unread(ch)
			break
		}
	}

	// If we end with a period, unread that as well
	identifier := buf.String()
	if len(identifier) > 1 && identifier[len(identifier)-1] == '.' {
		s.unread('.')
		identifier = identifier[:len(identifier)-1]
	}

	return IDENTIFIER, identifier
}

// scanBody consumes the current body until we reach the end of the file or the start of an expression
func (s *XScanner) scanBody() (Token, string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer

	// read characters until we reach the end of the file or the start of an expression or identifier
	for ch := s.read(); ch != eof; ch = s.read() {
		// could be start of an expression
		if ch == '@' {
			peek := s.read()

			// start of an expression
			if peek == '(' {
				s.unread(peek)
				s.unread('@')
				break

				// @@, means literal @
			} else if peek == '@' {
				buf.WriteRune('@')

				// this is an identifier
			} else if isIdentifierChar(peek) {
				s.unread(peek)
				s.unread('@')
				break

				// @ followed by non-letter
			} else {
				buf.WriteRune('@')
				buf.WriteRune(peek)
			}
		} else {
			buf.WriteRune(ch)
		}
	}

	return BODY, buf.String()
}

// Scan returns the next token and literal value.
func (s *XScanner) Scan() (Token, string) {
	for ch := s.read(); ch != eof; ch = s.read() {
		switch ch {
		case '@':
			peek := s.read()

			// start of an expression
			if peek == '(' {
				return s.scanExpression()

				// @@, means literal @
			} else if peek == '@' {
				s.unread('@')
				s.unread('@')
				return s.scanBody()

				// this is an identifier
			} else if isIdentifierChar(peek) {
				s.unread(peek)
				return s.scanIdentifier()

				// '@' followed by non-letter, plain body
			}

			s.unread(peek)
			s.unread('@')
			return s.scanBody()

		default:
			s.unread(ch)
			return s.scanBody()
		}
	}

	return EOF, ""
}

func EvaluateExpression(env utils.Environment, resolver utils.VariableResolver, template string) (interface{}, error) {
	errors := NewErrorListener()

	input := antlr.NewInputStream(template)
	lexer := gen.NewExcellentLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := gen.NewExcellentParser(stream)
	p.AddErrorListener(errors)
	tree := p.Parse()

	// if we ran into errors parsing, bail
	if errors.HasErrors() {
		return nil, fmt.Errorf(errors.Errors())
	}

	visitor := NewVisitor(env, resolver)
	value := visitor.Visit(tree)

	err, isErr := value.(error)

	// did our evaluation result in an error? return that
	if isErr {
		return nil, err
	}

	// all is good, return our value
	return value, nil
}

// EvaluateTemplate tries to evaluate the passed in template into an object, this only works if the template
// is a single identifier or expression, ie: "@contact" or "@(first(contact.urns))". In cases
// which are not a single identifier or expression, we return the stringified value
func EvaluateTemplate(env utils.Environment, resolver utils.VariableResolver, template string) (interface{}, error) {
	var buf bytes.Buffer
	template = strings.TrimSpace(template)
	scanner := NewXScanner(strings.NewReader(template))

	// parse our first token
	tokenType, token := scanner.Scan()

	// try to scan to our next token
	nextTT, _ := scanner.Scan()

	// if we had one, then just return our string evaluation strategy
	if nextTT != EOF {
		return EvaluateTemplateAsString(env, resolver, template)
	}

	switch tokenType {
	case IDENTIFIER:
		value := utils.ResolveVariable(env, resolver, token)

		// didn't find it, our value is empty string
		if value == nil {
			value = ""
		}

		err, isErr := value.(error)

		// we got an error, return our raw value
		if isErr {
			buf.WriteString("@")
			buf.WriteString(token)
			return buf.String(), err
		}

		// found it, return that value
		return value, nil

	case EXPRESSION:
		value, err := EvaluateExpression(env, resolver, token)
		if err != nil {
			return buf.String(), err
		}

		return value, nil
	}

	// different type of token, return the string representation
	return EvaluateTemplateAsString(env, resolver, template)
}

func EvaluateTemplateAsString(env utils.Environment, resolver utils.VariableResolver, template string) (string, error) {
	var buf bytes.Buffer
	var errors TemplateErrors
	scanner := NewXScanner(strings.NewReader(template))

	for tokenType, token := scanner.Scan(); tokenType != EOF; tokenType, token = scanner.Scan() {
		switch tokenType {
		case BODY:
			buf.WriteString(token)
		case IDENTIFIER:
			value := utils.ResolveVariable(env, resolver, token)

			// didn't find it, our value is empty string
			if value == nil {
				value = ""
			}
			_, isErr := value.(error)

			// we got an error, return our raw variable
			if isErr {
				buf.WriteString("@")
				buf.WriteString(token)
			} else {
				strValue, _ := utils.ToString(env, value)
				buf.WriteString(strValue)
			}
		case EXPRESSION:
			value, err := EvaluateExpression(env, resolver, token)
			if err != nil {
				errors = append(errors, err)
			} else {
				strValue, _ := utils.ToString(env, value)
				buf.WriteString(strValue)
			}

		}
	}

	if len(errors) > 0 {
		return buf.String(), errors
	}
	return buf.String(), nil
}

type TemplateErrors []error

func (e TemplateErrors) Error() string {
	if len(e) == 1 {
		return e[0].Error()
	}

	msg := "multiple errors:"
	for _, err := range e {
		msg += "\n" + err.Error()
	}
	return msg
}

type ErrorListener struct {
	errors bytes.Buffer
	*antlr.DefaultErrorListener
}

func NewErrorListener() *ErrorListener {
	return &ErrorListener{}
}

func (l *ErrorListener) HasErrors() bool {
	return l.errors.Len() > 0
}

func (l *ErrorListener) Errors() string {
	return l.errors.String()
}

func (l *ErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	l.errors.WriteString(fmt.Sprintln("line " + strconv.Itoa(line) + ":" + strconv.Itoa(column) + " " + msg))
}
