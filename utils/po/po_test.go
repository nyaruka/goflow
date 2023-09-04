package po_test

import (
	"os"
	"strings"
	"testing"
	"testing/iotest"
	"time"

	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils/po"
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

	c := po.Comment{}
	assert.Equal(t, "", c.String())
	assert.Equal(t, c, po.ParseComment(""))

	c = po.Comment{
		Translator: []string{"translator", ""},
		Extracted:  []string{"extracted"},
		References: []string{"src/foo.go", "src/bar.go"},
		Flags:      []string{"fuzzy", "go-format"},
	}
	assert.Equal(t, text, c.String())

	assert.Equal(t, c, po.ParseComment(text))
	assert.True(t, c.HasFlag("fuzzy"))
	assert.True(t, c.HasFlag("go-format"))
	assert.False(t, c.HasFlag("python-format"))
}

func TestPOCreation(t *testing.T) {
	header := po.NewHeader("Generated for testing", time.Date(2020, 3, 25, 11, 50, 30, 123456789, time.UTC), "es")
	header.Custom["Foo"] = "Bar"
	p := po.NewPO(header)

	p.AddEntry(&po.Entry{
		MsgID:  "Yes",
		MsgStr: "",
	})
	p.AddEntry(&po.Entry{
		MsgID:  "Yes",
		MsgStr: "Si",
	})

	p.AddEntry(&po.Entry{
		MsgContext: "context1",
		MsgID:      "No",
		MsgStr:     "",
	})
	p.AddEntry(&po.Entry{
		Comment: po.Comment{
			Extracted: []string{"has_text"},
		},
		MsgContext: "context1",
		MsgID:      "No",
		MsgStr:     "No",
	})

	b := &strings.Builder{}
	p.Write(b)

	assert.Equal(t, 2, len(p.Entries))
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

	p, err := po.ReadPO(poFile)
	require.NoError(t, err)

	assert.Equal(t, "Testing\n", p.Header.InitialComment)
	assert.True(t, time.Date(2020, 3, 25, 13, 57, 0, 0, time.UTC).Equal(p.Header.POTCreationDate))
	assert.Equal(t, "es", p.Header.Language)
	assert.Equal(t, "1.0", p.Header.MIMEVersion)
	assert.Equal(t, "text/plain; charset=UTF-8", p.Header.ContentType)

	assert.Equal(t, 7, len(p.Entries))
	assert.Equal(t, []string{"Translated/43f7e69e-727d-4cfe-81b8-564e7833052b/name:0", "Translated/e42deebf-90fa-4636-81cb-d247a3d3ba75/quick_replies:1"}, p.Entries[0].Comment.References)
	assert.Equal(t, "", p.Entries[0].MsgContext)
	assert.Equal(t, "Blue", p.Entries[0].MsgID)
	assert.Equal(t, "Azul", p.Entries[0].MsgStr)

	assert.Equal(t, "d1ce3c92-7025-4607-a910-444361a6b9b3/name:0", p.Entries[2].MsgContext)
	assert.Equal(t, "Red", p.Entries[2].MsgID)
	assert.Equal(t, "Roja", p.Entries[2].MsgStr)

	// try handling an i/o error
	badReader := iotest.TimeoutReader(strings.NewReader(`# Generated`))
	_, err = po.ReadPO(badReader)
	assert.EqualError(t, err, "timeout")

	// we can sort the entries
	p.Sort()

	assert.Equal(t, "Blue", p.Entries[0].MsgID)
	assert.Equal(t, "Other", p.Entries[1].MsgID)
	assert.Equal(t, "Red", p.Entries[2].MsgID)
	assert.Equal(t, "", p.Entries[2].MsgContext)
	assert.Equal(t, "Red", p.Entries[3].MsgID)
	assert.Equal(t, "d1ce3c92-7025-4607-a910-444361a6b9b3/name:0", p.Entries[3].MsgContext)

	// test writing the PO file
	b := &strings.Builder{}
	p.Write(b)
	test.AssertSnapshot(t, "write_po", b.String())
}

func TestGetText(t *testing.T) {
	poFile, err := os.Open("testdata/locale/es/simple.po")
	require.NoError(t, err)

	defer poFile.Close()
	p, err := po.ReadPO(poFile)
	require.NoError(t, err)

	assert.Equal(t, "Rojo", p.GetText("Male", "Red"))
	assert.Equal(t, "Roja", p.GetText("Female", "Red"))
	assert.Equal(t, "Red", p.GetText("", "Red"))
	assert.Equal(t, "Azul", p.GetText("", "Blue"))
	assert.Equal(t, "Missing", p.GetText("", "Missing"))
	assert.Equal(t, "Not even an entry", p.GetText("", "Not even an entry"))
	assert.Equal(t, "Green", p.GetText("", "Green")) // entry is ignored because it's fuzzy
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
		assert.Equal(t, tc.encoded, po.EncodeString(tc.original), "mismatch encoding: %s", tc.original)
		assert.Equal(t, tc.original, po.DecodeString(tc.encoded), "mismatch decoding: %s", tc.encoded)
	}
}
