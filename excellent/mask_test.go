package excellent_test

import (
	"fmt"
	"testing"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var maskTopLevels = []string{"contact", "fields", "results"}

func TestMaskTemplate(t *testing.T) {
	tcs := []struct {
		src     string
		masked  string
		hasText bool
	}{
		{``, ``, false},
		{`Hello`, `Hello`, true},
		{`Hi @contact.name!`, `Hi {{1}}!`, true},
		{`Hi @CONTACT.name`, `Hi {{1}}`, true},
		{`@contact.name @(fields.age)`, `{{1}} {{2}}`, false},
		{`@contact.name!?`, `{{1}}!?`, false},
		{`@contact.name 42`, `{{1}} 42`, true},
		{`Total @(fields.total * 2) due`, `Total {{1}} due`, true},
		{`@(upper(left(contact.name, 3)))`, `{{1}}`, false},
		{`Say @("a ) b" & contact.name) now`, `Say {{1}} now`, true},
		{`Oops @(upper(contact.name`, `Oops {{1}}`, true},
		{`Email bob@gmail.com`, `Email bob@gmail.com`, true},
		{`Pay @@contact`, `Pay @contact`, true},
		{`Pay @@@contact.name`, `Pay @{{1}}`, true},
		{`@@(1 + 2)`, `@(1 + 2)`, true},
		{`@contact.name@fields.age`, `{{1}}{{2}}`, false},
		{`See @contact.name.`, `See {{1}}.`, true},
		{`Unknown @foo.bar`, `Unknown @foo.bar`, true},
		{`Trailing @`, `Trailing @`, true},
	}

	for _, tc := range tcs {
		m, err := excellent.MaskTemplate(tc.src, maskTopLevels, excellent.BracePlaceholders)
		require.NoError(t, err, "unexpected error masking %q", tc.src)
		assert.Equal(t, tc.masked, m.Text, "text mismatch masking %q", tc.src)
		assert.Equal(t, tc.hasText, m.HasText(), "has text mismatch masking %q", tc.src)

		// unmasking the untransformed text restores the original template
		unmasked, err := m.Unmask(m.Text)
		assert.NoError(t, err, "unexpected error unmasking %q", tc.src)
		assert.Equal(t, tc.src, unmasked, "round trip mismatch for %q", tc.src)
	}

	// every kind of template segment survives a round trip, i.e. evaluates the same afterwards
	env := envs.NewBuilder().Build()
	ctx := makeTemplateContext()

	for _, src := range templateTestCases {
		m, err := excellent.MaskTemplate(src, ctx.Properties(), excellent.TagPlaceholders)
		require.NoError(t, err, "unexpected error masking %q", src)
		unmasked, err := m.Unmask(m.Text)
		require.NoError(t, err, "unexpected error unmasking %q", src)

		expectedVal, expectedWarns, expectedErr := excellent.ParseTemplate(src).Evaluate(env, ctx, nil)
		actualVal, actualWarns, actualErr := excellent.ParseTemplate(unmasked).Evaluate(env, ctx, nil)
		assert.Equal(t, expectedVal, actualVal, "output mismatch for %q", src)
		assert.Equal(t, expectedWarns, actualWarns, "warnings mismatch for %q", src)
		assert.Equal(t, fmt.Sprint(expectedErr), fmt.Sprint(actualErr), "error mismatch for %q", src)
	}

	// nil top-levels means any identifier is an expression
	m, err := excellent.MaskTemplate(`Email bob@gmail.com`, nil, excellent.BracePlaceholders)
	assert.NoError(t, err)
	assert.Equal(t, `Email bob{{1}}`, m.Text)

	// text which looks like a placeholder can't be masked
	for _, src := range []string{`Hi {{1}}`, `Hi {{ 3 }} @contact.name`, `Hi @contact.name {{1}}`, `{{2}}@contact.name`} {
		_, err := excellent.MaskTemplate(src, maskTopLevels, excellent.BracePlaceholders)
		assert.EqualError(t, err, "template text contains placeholder-like content", "expected error masking %q", src)
	}

	// but is fine for a different syntax
	m, err = excellent.MaskTemplate(`Hi {{1}} @contact.name`, maskTopLevels, excellent.TagPlaceholders)
	assert.NoError(t, err)
	assert.Equal(t, `Hi {{1}} <x id="1"/>`, m.Text)
}

func TestUnmaskTemplate(t *testing.T) {
	tcs := []struct {
		src         string
		syntax      *excellent.PlaceholderSyntax
		masked      string
		transformed string
		unmasked    string
		err         string
	}{
		{
			src:         `Hi @contact.name, your total is @(fields.total * 2) - email bob@gmail.com or pay @@contact`,
			syntax:      excellent.BracePlaceholders,
			masked:      `Hi {{1}}, your total is {{2}} - email bob@gmail.com or pay @contact`,
			transformed: `Hola {{1}}, tu total es {{2}} - escribe a bob@gmail.com o paga @contact`,
			unmasked:    `Hola @contact.name, tu total es @(fields.total * 2) - escribe a bob@gmail.com o paga @@contact`,
		},
		{
			src:         `Hi @contact.name, your total is @(fields.total * 2) - email bob@gmail.com or pay @@contact`,
			syntax:      excellent.TagPlaceholders,
			masked:      `Hi <x id="1"/>, your total is <x id="2"/> - email bob@gmail.com or pay @contact`,
			transformed: `Hola <x id="1" />, tu total es <x  id="2"/> - escribe a bob@gmail.com o paga @contact`,
			unmasked:    `Hola @contact.name, tu total es @(fields.total * 2) - escribe a bob@gmail.com o paga @@contact`,
		},
		{ // reordered placeholders
			src:         `@contact.first_name @contact.last_name`,
			syntax:      excellent.TagPlaceholders,
			masked:      `<x id="1"/> <x id="2"/>`,
			transformed: `<x id="2"/>, <x id="1"/>`,
			unmasked:    `@contact.last_name, @contact.first_name`,
		},
		{ // nested parens and string literals kept intact
			src:         `Say @(concat("a)", upper(contact.name))) now`,
			syntax:      excellent.BracePlaceholders,
			masked:      `Say {{1}} now`,
			transformed: `Di {{1}} ahora`,
			unmasked:    `Di @(concat("a)", upper(contact.name))) ahora`,
		},
		{ // unterminated expression restored verbatim
			src:         `Oops @(upper(contact.name`,
			syntax:      excellent.BracePlaceholders,
			masked:      `Oops {{1}}`,
			transformed: `Uy {{1}}`,
			unmasked:    `Uy @(upper(contact.name`,
		},
		{ // introduced @ which would start expressions are escaped, others aren't
			src:         `Hello`,
			syntax:      excellent.BracePlaceholders,
			masked:      `Hello`,
			transformed: `@contact @Contact.name @(1 + 2) @@ @foo me@x.com @ @`,
			unmasked:    `@@contact @@Contact.name @@(1 + 2) @@@@ @foo me@x.com @ @`,
		},
		{ // a @ directly before a placeholder is escaped
			src:         `Pay @@@contact.name`,
			syntax:      excellent.BracePlaceholders,
			masked:      `Pay @{{1}}`,
			transformed: `Paga @{{1}}`,
			unmasked:    `Paga @@@contact.name`,
		},
		{ // an identifier followed by text that would extend it is wrapped
			src:         `@contact.name and @fields.age!`,
			syntax:      excellent.BracePlaceholders,
			masked:      `{{1}} and {{2}}!`,
			transformed: `{{1}}y {{2}}.x`,
			unmasked:    `@(contact.name)y @(fields.age).x`,
		},
		{ // but not if what follows can't extend it
			src:         `@contact.name and @fields.age!`,
			syntax:      excellent.BracePlaceholders,
			masked:      `{{1}} and {{2}}!`,
			transformed: `{{1}}{{2}}. x`,
			unmasked:    `@contact.name@fields.age. x`,
		},
		{ // missing placeholder
			src:         `Hi @contact.name and @fields.age`,
			syntax:      excellent.BracePlaceholders,
			masked:      `Hi {{1}} and {{2}}`,
			transformed: `Hola {{1}}`,
			err:         "expected 2 placeholders, found 1",
		},
		{ // duplicated placeholder
			src:         `Hi @contact.name and @fields.age`,
			syntax:      excellent.BracePlaceholders,
			masked:      `Hi {{1}} and {{2}}`,
			transformed: `Hola {{1}} y {{1}}`,
			err:         "duplicate placeholder {{1}}",
		},
		{ // unknown placeholder
			src:         `Hi @contact.name`,
			syntax:      excellent.TagPlaceholders,
			masked:      `Hi <x id="1"/>`,
			transformed: `Hola <x id="2"/>`,
			err:         `unknown placeholder <x id="2"/>`,
		},
		{ // zero isn't a valid id
			src:         `Hi @contact.name`,
			syntax:      excellent.BracePlaceholders,
			masked:      `Hi {{1}}`,
			transformed: `Hola {{0}}`,
			err:         `unknown placeholder {{0}}`,
		},
		{ // placeholder added to a template without expressions
			src:         `Hello`,
			syntax:      excellent.BracePlaceholders,
			masked:      `Hello`,
			transformed: `Hola {{1}}`,
			err:         "expected 0 placeholders, found 1",
		},
	}

	for _, tc := range tcs {
		m, err := excellent.MaskTemplate(tc.src, maskTopLevels, tc.syntax)
		require.NoError(t, err, "unexpected error masking %q", tc.src)
		assert.Equal(t, tc.masked, m.Text, "masked text mismatch for %q", tc.src)

		unmasked, err := m.Unmask(tc.transformed)
		if tc.err != "" {
			assert.EqualError(t, err, tc.err, "error mismatch unmasking %q", tc.transformed)
		} else {
			assert.NoError(t, err, "unexpected error unmasking %q", tc.transformed)
			assert.Equal(t, tc.unmasked, unmasked, "unmasked mismatch for %q", tc.transformed)
		}
	}
}

func TestCustomPlaceholderSyntax(t *testing.T) {
	syntax := excellent.NewPlaceholderSyntax(
		func(n int) string { return fmt.Sprintf(`<span translate="no" data-id="%d"></span>`, n) },
		`<span translate="no" data-id="(\d+)"></span>`,
	)

	m, err := excellent.MaskTemplate(`Hi @contact.name`, maskTopLevels, syntax)
	assert.NoError(t, err)
	assert.Equal(t, `Hi <span translate="no" data-id="1"></span>`, m.Text)

	unmasked, err := m.Unmask(`Hola <span translate="no" data-id="1"></span>`)
	assert.NoError(t, err)
	assert.Equal(t, `Hola @contact.name`, unmasked)

	assert.Panics(t, func() { excellent.NewPlaceholderSyntax(func(n int) string { return "" }, `\{\d+\}`) })
}
