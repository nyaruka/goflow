package i18n_test

import (
	"os"
	"strings"
	"testing"
	"testing/iotest"
	"time"

	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils/i18n"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComments(t *testing.T) {
	text := `#  translator
#  
#. extracted
#: src/foo.go
#: src/bar.go
#, fuzzy,go-format`

	c := i18n.POComment{}
	assert.Equal(t, "", c.String())
	assert.Equal(t, c, i18n.ParsePOComment(""))

	c = i18n.POComment{
		Translator: []string{"translator", ""},
		Extracted:  []string{"extracted"},
		References: []string{"src/foo.go", "src/bar.go"},
		Flags:      []string{"fuzzy", "go-format"},
	}
	assert.Equal(t, text, c.String())

	assert.Equal(t, c, i18n.ParsePOComment(text))
}

func TestPOCreation(t *testing.T) {
	header := i18n.NewPOHeader("Generated for testing", time.Date(2020, 3, 25, 11, 50, 30, 123456789, time.UTC), "es")
	header.Custom["Foo"] = "Bar"
	po := i18n.NewPO(header)

	po.AddEntry(&i18n.POEntry{
		MsgID:  "Yes",
		MsgStr: "",
	})
	po.AddEntry(&i18n.POEntry{
		MsgID:  "Yes",
		MsgStr: "Si",
	})

	po.AddEntry(&i18n.POEntry{
		MsgContext: "context1",
		MsgID:      "No",
		MsgStr:     "",
	})
	po.AddEntry(&i18n.POEntry{
		Comment: i18n.POComment{
			Extracted: []string{"has_text"},
		},
		MsgContext: "context1",
		MsgID:      "No",
		MsgStr:     "No",
	})

	b := &strings.Builder{}
	po.Write(b)

	assert.Equal(t, 2, len(po.Entries))
	assert.Equal(
		t, `#  Generated for testing
#  
#, fuzzy
msgid ""
msgstr ""
"POT-Creation-Date: 2020-03-25 11:50+0000\n"
"Language: es\n"
"MIME-Version: 1.0\n"
"Content-Type: text/plain; charset=UTF-8\n"
"Foo: Bar\n"

msgid "Yes"
msgstr "Si"

#. has_text
msgctxt "context1"
msgid "No"
msgstr "No"

`,
		b.String())
}

func TestReadAndWritePO(t *testing.T) {
	poFile, err := os.Open("testdata/translation_mismatches.noargs.es.po")
	require.NoError(t, err)

	defer poFile.Close()

	po, err := i18n.ReadPO(poFile)
	require.NoError(t, err)

	assert.Equal(t, "Testing\n", po.Header.InitialComment)
	assert.True(t, time.Date(2020, 3, 25, 13, 57, 0, 0, time.UTC).Equal(po.Header.POTCreationDate))
	assert.Equal(t, "es", po.Header.Language)
	assert.Equal(t, "1.0", po.Header.MIMEVersion)
	assert.Equal(t, "text/plain; charset=UTF-8", po.Header.ContentType)

	assert.Equal(t, 7, len(po.Entries))
	assert.Equal(t, []string{"Translated/43f7e69e-727d-4cfe-81b8-564e7833052b/name:0", "Translated/e42deebf-90fa-4636-81cb-d247a3d3ba75/quick_replies:1"}, po.Entries[0].Comment.References)
	assert.Equal(t, "", po.Entries[0].MsgContext)
	assert.Equal(t, "Blue", po.Entries[0].MsgID)
	assert.Equal(t, "Azul", po.Entries[0].MsgStr)

	assert.Equal(t, "d1ce3c92-7025-4607-a910-444361a6b9b3/name:0", po.Entries[2].MsgContext)
	assert.Equal(t, "Red", po.Entries[2].MsgID)
	assert.Equal(t, "Roja", po.Entries[2].MsgStr)

	// try handling an i/o error
	badReader := iotest.TimeoutReader(strings.NewReader(`# Generated`))
	_, err = i18n.ReadPO(badReader)
	assert.EqualError(t, err, "timeout")

	// we can sort the entries
	po.Sort()

	assert.Equal(t, "Blue", po.Entries[0].MsgID)
	assert.Equal(t, "Other", po.Entries[1].MsgID)
	assert.Equal(t, "Red", po.Entries[2].MsgID)
	assert.Equal(t, "", po.Entries[2].MsgContext)
	assert.Equal(t, "Red", po.Entries[3].MsgID)
	assert.Equal(t, "d1ce3c92-7025-4607-a910-444361a6b9b3/name:0", po.Entries[3].MsgContext)

	// test writing the PO file
	b := &strings.Builder{}
	po.Write(b)
	test.AssertSnapshot(t, "write_po", b.String())
}

func TestGetText(t *testing.T) {
	poFile, err := os.Open("testdata/locales/es/simple.po")
	require.NoError(t, err)

	defer poFile.Close()
	po, err := i18n.ReadPO(poFile)
	require.NoError(t, err)

	assert.Equal(t, "Rojo", po.GetText("Male", "Red"))
	assert.Equal(t, "Roja", po.GetText("Female", "Red"))
	assert.Equal(t, "Red", po.GetText("", "Red"))
	assert.Equal(t, "Azul", po.GetText("", "Blue"))
	assert.Equal(t, "Missing", po.GetText("", "Missing"))
	assert.Equal(t, "Not even an entry", po.GetText("", "Not even an entry"))
}

func TestEncodeAndDecodePOString(t *testing.T) {
	tests := []struct {
		original string
		encoded  string
	}{
		{``, `""`},
		{`FOO`, `"FOO"`},
		{
			"FOO\nBAR", `""
"FOO\n"
"BAR"`,
		},
		{
			"\nFOO\n\nBAR\n", `""
"\n"
"FOO\n"
"\n"
"BAR\n"`,
		},
		{
			"FOO\n\n\n", `""
"FOO\n"
"\n"
"\n"`,
		},
		{"FOO\tB\\AR\"", `"FOO\tB\\AR\""`},
	}

	for _, tc := range tests {
		assert.Equal(t, tc.encoded, i18n.EncodePOString(tc.original), "mismatch encoding: %s", tc.original)
		assert.Equal(t, tc.original, i18n.DecodePOString(tc.encoded), "mismatch decoding: %s", tc.encoded)
	}
}
