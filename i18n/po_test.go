package i18n_test

import (
	"strings"
	"testing"
	"time"

	"github.com/nyaruka/goflow/i18n"

	"github.com/stretchr/testify/assert"
)

func TestComments(t *testing.T) {
	c := i18n.Comment{}
	assert.Equal(t, "", c.String())

	c = i18n.Comment{
		Translator: "translator",
		Extracted:  "extracted",
		References: []string{"src/foo.go", "src/bar.go"},
		Flags:      []string{"fuzzy"},
	}
	assert.Equal(t, `#  translator
#. extracted
#: src/foo.go
#: src/bar.go
#, fuzzy
`, c.String())
}

func TestPOs(t *testing.T) {
	po := i18n.NewPO("Generated for testing", time.Date(2020, 3, 25, 11, 50, 30, 123456789, time.UTC), "es")

	po.AddEntry(&i18n.Entry{
		MsgID:  "Yes",
		MsgStr: "",
	})
	po.AddEntry(&i18n.Entry{
		MsgID:  "Yes",
		MsgStr: "Si",
	})

	po.AddEntry(&i18n.Entry{
		MsgContext: "context1",
		MsgID:      "No",
		MsgStr:     "",
	})
	po.AddEntry(&i18n.Entry{
		Comment: i18n.Comment{
			Extracted: "has_text",
		},
		MsgContext: "context1",
		MsgID:      "No",
		MsgStr:     "No",
	})

	b := &strings.Builder{}
	po.Write(b)

	assert.Equal(t, 2, len(po.Entries))
	assert.Equal(
		t, `# Generated for testing
#
#, fuzzy
msgid ""
msgstr ""
"POT-Creation-Date: 2020-03-25 11:50+0000\n"
"Language: es\n"
"MIME-Version: 1.0\n"
"Content-Type: text/plain; charset=UTF-8\n"

msgid "Yes"
msgstr "Si"

#. has_text
msgctxt "context1"
msgid "No"
msgstr "No"

`,
		b.String())
}

func TestEncodePOString(t *testing.T) {
	assert.Equal(t, `""`, i18n.EncodePOString(""))
	assert.Equal(t, `"FOO"`, i18n.EncodePOString("FOO"))
	assert.Equal(
		t,
		`""
"FOO\n"
"BAR"`,
		i18n.EncodePOString("FOO\nBAR"),
	)
	assert.Equal(
		t,
		`""
"\n"
"FOO\n"
"\n"
"BAR\n"`,
		i18n.EncodePOString("\nFOO\n\nBAR\n"),
	)
	assert.Equal(
		t,
		`""
"FOO\n"
"\n"
"\n"`,
		i18n.EncodePOString("FOO\n\n\n"),
	)
	assert.Equal(t, `"FOO\tBAR"`, i18n.EncodePOString("FOO\tBAR"))
}
