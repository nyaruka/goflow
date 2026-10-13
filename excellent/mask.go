package excellent

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

// PlaceholderSyntax is how a masked template represents its expressions. Different consumers of masked text, e.g.
// translation services, preserve different kinds of markup, so this is chosen by the caller.
type PlaceholderSyntax struct {
	render  func(int) string
	pattern *regexp.Regexp
}

// NewPlaceholderSyntax creates a placeholder syntax from a function which renders the placeholder with the given id,
// and a pattern which matches any placeholder, capturing its id as its first group. The pattern can be looser than
// the rendered form, e.g. to tolerate whitespace that a translation service might add.
func NewPlaceholderSyntax(render func(int) string, pattern string) *PlaceholderSyntax {
	re := regexp.MustCompile(pattern)
	if re.NumSubexp() < 1 {
		panic(fmt.Sprintf("placeholder pattern %q has no capture group for the id", pattern))
	}
	return &PlaceholderSyntax{render: render, pattern: re}
}

// BracePlaceholders renders placeholders as `{{1}}`
var BracePlaceholders = NewPlaceholderSyntax(
	func(n int) string { return fmt.Sprintf("{{%d}}", n) },
	`\{\{\s*(\d+)\s*\}\}`,
)

// TagPlaceholders renders placeholders as self-closing XML tags like `<x id="1"/>`
var TagPlaceholders = NewPlaceholderSyntax(
	func(n int) string { return fmt.Sprintf(`<x id="%d"/>`, n) },
	`<x\s+id\s*=\s*"(\d+)"\s*/>`,
)

// MaskedTemplate is a template whose expressions have been replaced by numbered placeholders, so that its text can be
// transformed, e.g. translated, without the transformation having to recognize or preserve expressions itself.
type MaskedTemplate struct {
	// Text is the template's unescaped text with placeholders, e.g. `Hi {{1}}, email bob@gmail.com`
	Text string

	syntax           *PlaceholderSyntax
	allowedTopLevels []string
	exprs            []maskedExpression // indexed by placeholder id - 1
	hasText          bool
}

type maskedExpression struct {
	src        string // original source, e.g. "@contact.name" or "@(1 + 2)"
	identifier string // set if this is an @identifier, e.g. "contact.name"
}

// MaskTemplate replaces the expressions in the given template with placeholders of the given syntax. Identifiers whose
// top-level isn't allowed, e.g. the @gmail in an email address, are left as text. It errors if the template's text
// already contains something which looks like a placeholder, as that couldn't be told apart from a real one.
func MaskTemplate(src string, allowedTopLevels []string, syntax *PlaceholderSyntax) (*MaskedTemplate, error) {
	m := &MaskedTemplate{syntax: syntax, allowedTopLevels: allowedTopLevels}
	var buf strings.Builder

	// scanned without unescaping so that a body starting with @( can only be an unterminated expression
	VisitTemplate(src, allowedTopLevels, false, func(tt XTokenType, token string) error {
		var expr maskedExpression
		switch tt {
		case IDENTIFIER:
			expr = maskedExpression{src: "@" + token, identifier: token}
		case EXPRESSION:
			expr = maskedExpression{src: "@(" + token + ")"}
		default:
			if !strings.HasPrefix(token, "@(") {
				text := strings.ReplaceAll(token, "@@", "@")
				m.hasText = m.hasText || strings.IndexFunc(text, func(r rune) bool { return unicode.IsLetter(r) || unicode.IsNumber(r) }) >= 0
				buf.WriteString(text)
				return nil
			}
			expr = maskedExpression{src: token} // kept verbatim so that it fails the same way when evaluated
		}
		m.exprs = append(m.exprs, expr)
		buf.WriteString(syntax.render(len(m.exprs)))
		return nil
	})

	m.Text = buf.String()

	// the only things that look like placeholders must be our placeholders, in order
	matches := syntax.pattern.FindAllStringSubmatch(m.Text, -1)
	ok := len(matches) == len(m.exprs)
	for i := 0; ok && i < len(matches); i++ {
		ok = matches[i][0] == syntax.render(i+1)
	}
	if !ok {
		return nil, errors.New("template text contains placeholder-like content")
	}

	return m, nil
}

// HasText returns whether there is anything other than punctuation and whitespace outside of the placeholders, i.e.
// whether there's anything worth transforming.
func (m *MaskedTemplate) HasText() bool { return m.hasText }

// Unmask restores the expressions in a transformed version of this template's text, returning a template. It errors
// if the transformed text doesn't contain each placeholder exactly once. Text between placeholders is escaped so that
// it can't introduce new expressions.
func (m *MaskedTemplate) Unmask(transformed string) (string, error) {
	matches := m.syntax.pattern.FindAllStringSubmatchIndex(transformed, -1)
	if len(matches) != len(m.exprs) {
		return "", fmt.Errorf("expected %d placeholders, found %d", len(m.exprs), len(matches))
	}

	var buf strings.Builder
	seen := make([]bool, len(m.exprs))
	var prev *maskedExpression
	last := 0

	writeText := func(text string, followedByExpr bool) {
		// an identifier followed by text which would extend it is wrapped so that it stays the same
		if prev != nil && prev.identifier != "" && extendsIdentifier(text) {
			buf.WriteString("@(" + prev.identifier + ")")
		} else if prev != nil {
			buf.WriteString(prev.src)
		}
		buf.WriteString(m.escapeText(text, followedByExpr))
	}

	for _, match := range matches {
		id, err := strconv.Atoi(transformed[match[2]:match[3]])
		if err != nil || id < 1 || id > len(m.exprs) {
			return "", fmt.Errorf("unknown placeholder %s", transformed[match[0]:match[1]])
		}
		if seen[id-1] {
			return "", fmt.Errorf("duplicate placeholder %s", transformed[match[0]:match[1]])
		}
		seen[id-1] = true

		writeText(transformed[last:match[0]], true)
		prev = &m.exprs[id-1]
		last = match[1]
	}
	writeText(transformed[last:], false)

	return buf.String(), nil
}

// escapes each @ in text which would otherwise start an expression or an escape sequence, leaving others, e.g. in an
// email address, as they are. Every expression starts with @, so a trailing @ needs escaping if an expression follows.
func (m *MaskedTemplate) escapeText(s string, followedByExpr bool) string {
	var buf strings.Builder
	for i, r := range s {
		buf.WriteRune(r)
		if r != '@' {
			continue
		}
		rest := s[i+1:]
		if (rest == "" && followedByExpr) || strings.HasPrefix(rest, "@") || strings.HasPrefix(rest, "(") || strings.HasSuffix(s[:i], "@") {
			buf.WriteRune('@')
			continue
		}
		end := strings.IndexFunc(rest, func(r rune) bool { return !isNameChar(r) })
		if end < 0 {
			end = len(rest)
		}
		if end > 0 && (m.allowedTopLevels == nil || slices.Contains(m.allowedTopLevels, strings.ToLower(rest[:end]))) {
			buf.WriteRune('@')
		}
	}
	return buf.String()
}

// whether the given text, if it directly followed an identifier, would be scanned as part of it
func extendsIdentifier(s string) bool {
	s = strings.TrimPrefix(s, ".")
	for _, r := range s {
		return isNameChar(r)
	}
	return false
}
