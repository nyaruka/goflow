package excellent

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type xToken int

const (
	// BODY - Not in expression
	BODY xToken = iota

	// IDENTIFIER - 'contact.age' in '@contact.age'
	IDENTIFIER

	// EXPRESSION - the body of an expression '1+2' in '@(1+2)'
	EXPRESSION

	// EOF - end of expression
	EOF
)

const eof rune = rune(0)

func isIdentifierChar(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsNumber(ch) || ch == '.' || ch == '_'
}

// xscanner represents a lexical scanner.
type xscanner struct {
	reader              *bufio.Reader
	unreadRunes         []rune
	unreadCount         int
	identifierTopLevels []string
}

// NewXScanner returns a new instance of our excellent scanner
func NewXScanner(r io.Reader, identifierTopLevels []string) *xscanner {
	return &xscanner{
		reader:              bufio.NewReader(r),
		unreadRunes:         make([]rune, 4),
		identifierTopLevels: identifierTopLevels,
	}
}

// gets the next rune or EOF if we are at the end of the string
func (s *xscanner) read() rune {
	// first see if we have any unread runes to return
	if s.unreadCount > 0 {
		ch := s.unreadRunes[s.unreadCount-1]
		s.unreadCount--
		return ch
	}

	// otherwise, read the next run
	ch, _, err := s.reader.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// pops the passed in rune as the next rune to be returned
func (s *xscanner) unread(ch rune) {
	s.unreadRunes[s.unreadCount] = ch
	s.unreadCount++
}

// scanExpression consumes the current rune and all contiguous pieces until the end of the expression
// our read should be after the '('
func (s *xscanner) scanExpression() (xToken, string) {
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
func (s *xscanner) scanIdentifier() (xToken, string) {
	// Create a buffer and read the current character into it.
	var buf strings.Builder
	var topLevel string

	// Read every subsequent character until we reach the end of the identifier
	for ch := s.read(); ch != eof; ch = s.read() {
		if ch == '.' && topLevel == "" {
			topLevel = buf.String()
		}

		if isIdentifierChar(ch) {
			buf.WriteRune(ch)
		} else {
			s.unread(ch)
			break
		}
	}

	identifier := buf.String()

	if topLevel == "" {
		topLevel = identifier
	}

	// if we end with a period, unread that as well
	if len(identifier) > 1 && identifier[len(identifier)-1] == '.' {
		s.unread('.')
		identifier = identifier[:len(identifier)-1]
	}

	// only return as an identifier if the toplevel scope is valid
	if s.identifierTopLevels != nil {
		for _, validTopLevel := range s.identifierTopLevels {
			if topLevel == validTopLevel {
				return IDENTIFIER, identifier
			}
		}
	} else {
		return IDENTIFIER, identifier
	}

	// this was something that looked like an identifier but wasn't an allowed top-level variable, e.g. email address
	return BODY, fmt.Sprintf("@%s", identifier)
}

// scanBody consumes the current body until we reach the end of the file or the start of an expression
func (s *xscanner) scanBody() (xToken, string) {
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
func (s *xscanner) Scan() (xToken, string) {
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
