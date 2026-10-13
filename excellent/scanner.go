package excellent

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// XTokenType is a set of types than can be scanned
type XTokenType int

const (
	// BODY - Not in expression
	BODY XTokenType = iota

	// IDENTIFIER - 'contact.age' in '@contact.age'
	IDENTIFIER

	// EXPRESSION_PARENTHESES - the body of an expression '1+2' in '@(1+2)'
	EXPRESSION_PARENTHESES

	// EXPRESSION_BRACES - the body of an expression '1+2' in '{{1+2}}'
	EXPRESSION_BRACES

	// EOF - end of expression
	EOF
)

// Source returns the template source of the given token, e.g. "@(1+2)" for an EXPRESSION_PARENTHESES token "1+2"
func Source(tokenType XTokenType, token string) string {
	switch tokenType {
	case IDENTIFIER:
		return "@" + token
	case EXPRESSION_PARENTHESES:
		return "@(" + token + ")"
	case EXPRESSION_BRACES:
		return "{{" + token + "}}"
	}
	return token
}

func isNameChar(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsNumber(ch) || ch == '_'
}

// Scanner is something which can scan tokens from input
type Scanner interface {
	Scan() (XTokenType, string)
	SetUnescapeBody(bool)
	SetBraces(bool)
}

// xscanner represents a lexical scanner.
type xscanner struct {
	input               *xinput
	identifierTopLevels []string
	unescapeBody        bool // unescape @@ and \{{ sequences in the body
	braces              bool // recognize {{...}} expressions and \{{ escapes
}

// NewXScanner returns a new instance of our excellent scanner
func NewXScanner(r io.Reader, identifierTopLevels []string) Scanner {
	return &xscanner{
		input:               newInput(bufio.NewReader(r)),
		identifierTopLevels: identifierTopLevels,
		unescapeBody:        true,
		braces:              true,
	}
}

func (s *xscanner) SetUnescapeBody(unescape bool) {
	s.unescapeBody = unescape
}

// SetBraces sets whether {{...}} expressions are recognized - disabled to scan templates written before they existed
func (s *xscanner) SetBraces(braces bool) {
	s.braces = braces
}

// scanExpression consumes the current rune and all contiguous pieces until the end of the expression
// our read should be after the '('
func (s *xscanner) scanExpression() (XTokenType, string) {
	// create a buffer and read the current character into it.
	buf := &bytes.Buffer{}

	// our parentheses depth
	parens := 1

	// read every subsequent character until we reach the end of the expression
	for ch := s.input.read(); ch != eof; ch = s.input.read() {
		if ch == '"' {
			buf.WriteRune(ch)
			s.readTextLiteral(buf)
		} else if ch == '(' {
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
		return EXPRESSION_PARENTHESES, buf.String()
	}

	return BODY, strings.Join([]string{"@(", buf.String()}, "")
}

// scanBraces consumes everything until the closing }} of a braced expression, our read should be after the '{{'
func (s *xscanner) scanBraces() (XTokenType, string) {
	buf := &bytes.Buffer{}

	for ch := s.input.read(); ch != eof; ch = s.input.read() {
		if ch == '"' {
			buf.WriteRune(ch)
			s.readTextLiteral(buf)
		} else if ch == '}' {
			peek := s.input.read()
			if peek == '}' {
				return EXPRESSION_BRACES, buf.String()
			}
			s.input.unread(peek)
			buf.WriteRune(ch)
		} else {
			buf.WriteRune(ch)
		}
	}

	// unterminated so treat as body text
	return BODY, "{{" + buf.String()
}

// reads the remainder of a " quoted text literal
func (s *xscanner) readTextLiteral(buf *bytes.Buffer) {
	escaped := false
	for ch := s.input.read(); ch != eof; ch = s.input.read() {
		buf.WriteRune(ch)

		if ch == '"' && !escaped {
			break
		} else if ch == '\\' {
			escaped = !escaped
		} else {
			escaped = false
		}
	}
}

// scanIdentifier consumes the current rune and all contiguous pieces until the end of the identifer
// our read should be after the '@'
func (s *xscanner) scanIdentifier() (XTokenType, string) {
	// Create a buffer and read the current character into it.
	buf := &strings.Builder{}
	var topLevel string

	// Read every subsequent character until we reach the end of the identifier
	for ch := s.input.read(); ch != eof; ch = s.input.read() {
		if ch == '.' && topLevel == "" {
			topLevel = buf.String()
		}

		// only include period if it's followed by a valid name char
		if ch == '.' {
			peek := s.input.read()
			if isNameChar(peek) {
				buf.WriteRune(ch)
				buf.WriteRune(peek)
			} else {
				// this period actually signifies the end of the indentifier
				s.input.unread(peek)
				s.input.unread('.')
				break
			}
		} else if isNameChar(ch) {
			buf.WriteRune(ch)
		} else {
			s.input.unread(ch)
			break
		}
	}

	identifier := buf.String()

	if topLevel == "" {
		topLevel = identifier
	}
	topLevel = strings.ToLower(topLevel)

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
func (s *xscanner) scanBody() (XTokenType, string) {
	// Create a buffer and read the current character into it.
	buf := &strings.Builder{}

	// read characters until we reach the end of the file or the start of an expression or identifier
	for ch := s.input.read(); ch != eof; ch = s.input.read() {
		// could be start of an expression
		if ch == '@' {
			peek := s.input.read()

			// start of an expression
			if peek == '(' {
				s.input.unread(peek)
				s.input.unread('@')
				break

				// @@, means literal @
			} else if peek == '@' {
				buf.WriteRune('@')

				if !s.unescapeBody {
					buf.WriteRune('@')
				}

				// this is an identifier
			} else if isNameChar(peek) {
				s.input.unread(peek)
				s.input.unread('@')
				break

				// @ at the end of the input
			} else if peek == eof {
				buf.WriteRune('@')

				// @ followed by non-letter, which is scanned normally as it could be the start of {{ or \{{
			} else {
				buf.WriteRune('@')
				s.input.unread(peek)
			}
		} else if ch == '{' && s.braces {
			peek := s.input.read()

			// start of a braces expression
			if peek == '{' {
				s.input.unread(peek)
				s.input.unread('{')
				break
			}

			s.input.unread(peek)
			buf.WriteRune(ch)
		} else if ch == '\\' && s.braces {
			// handlebars style escaping: \{{ is a literal {{ and \\{{ is a literal \ followed by an expression
			p1 := s.input.read()
			p2 := s.input.read()

			if p1 == '{' && p2 == '{' {
				if !s.unescapeBody {
					buf.WriteRune('\\')
				}
				buf.WriteString("{{")
				continue
			}

			if p1 == '\\' && p2 == '{' {
				p3 := s.input.read()
				s.input.unread(p3)

				if p3 == '{' {
					s.input.unread(p2)
					buf.WriteRune('\\')
					if !s.unescapeBody {
						buf.WriteRune('\\')
					}
					break
				}
			}

			s.input.unread(p2)
			s.input.unread(p1)
			buf.WriteRune(ch)
		} else {
			buf.WriteRune(ch)
		}
	}

	return BODY, buf.String()
}

// Scan returns the next token and literal value.
func (s *xscanner) Scan() (XTokenType, string) {
	for ch := s.input.read(); ch != eof; ch = s.input.read() {
		switch ch {
		case '@':
			peek := s.input.read()

			// start of an expression
			if peek == '(' {
				return s.scanExpression()

				// @@, means literal @
			} else if peek == '@' {
				s.input.unread('@')
				s.input.unread('@')
				return s.scanBody()

				// this is an identifier
			} else if isNameChar(peek) {
				s.input.unread(peek)
				return s.scanIdentifier()

				// '@' followed by non-letter, plain body
			}

			s.input.unread(peek)
			s.input.unread('@')
			return s.scanBody()

		case '{':
			peek := s.input.read()
			if peek == '{' && s.braces {
				return s.scanBraces()
			}

			s.input.unread(peek)
			s.input.unread('{')
			return s.scanBody()

		default:
			s.input.unread(ch)
			return s.scanBody()
		}
	}

	return EOF, ""
}
