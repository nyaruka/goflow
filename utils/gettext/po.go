package gettext

import (
	"fmt"
	"io"
	"strings"
	"time"
)

const poDatetimeformat = "2006-01-02 15:04-0700"

// Header contains metadata about a PO file
type Header struct {
	POTCreationDate time.Time // POT-Creation-Date: YEAR-MO-DA HO:MI+ZONE
	Language        string    // Language: e.g. en-US
	MIMEVersion     string    // MIME-Version: 1.0
	ContentType     string    // Content-Type: text/plain; charset=UTF-8
}

func newHeader(creationDate time.Time, lang string) Header {
	return Header{
		POTCreationDate: creationDate,
		Language:        lang,
		MIMEVersion:     "1.0",
		ContentType:     "text/plain; charset=UTF-8",
	}
}

func (h *Header) asEntry() *Entry {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("POT-Creation-Date: %s\n", h.POTCreationDate.Format(poDatetimeformat)))
	b.WriteString(fmt.Sprintf("Language: %s\n", h.Language))
	b.WriteString(fmt.Sprintf("MIME-Version: %s\n", h.MIMEVersion))
	b.WriteString(fmt.Sprintf("Content-Type: %s\n", h.ContentType))

	return &Entry{
		Comment: Comment{
			Flags: []string{"fuzzy"},
		},
		MsgID:  "",
		MsgStr: b.String(),
	}
}

type Comment struct {
	Translator string   // #  translator-comments
	Extracted  string   // #. extracted-comments
	References []string // #: references
	Flags      []string // #, e.g. fuzzy,python-format
}

func (c *Comment) String() string {
	b := strings.Builder{}
	if c.Translator != "" {
		b.WriteString(fmt.Sprintf("#  %s\n", c.Translator))
	}
	if c.Extracted != "" {
		b.WriteString(fmt.Sprintf("#. %s\n", c.Extracted))
	}
	if len(c.References) > 0 {
		b.WriteString(fmt.Sprintf("#: %s\n", strings.Join(c.References, ",")))
	}
	if len(c.Flags) > 0 {
		b.WriteString(fmt.Sprintf("#, %s\n", strings.Join(c.Flags, ",")))
	}
	return b.String()
}

type Entry struct {
	Comment           // Comment
	MsgContext string // msgctxt context
	MsgID      string // msgid untranslated-string
	MsgStr     string // msgstr translated-string
}

func (e *Entry) String() string {
	b := &strings.Builder{}
	comment := e.Comment.String()
	if comment != "" {
		b.WriteString(comment)
	}
	if e.MsgContext != "" {
		fmt.Fprintf(b, "msgctxt %s\n", EncodePOString(e.MsgContext))
	}
	fmt.Fprintf(b, "msgid %s\n", EncodePOString(e.MsgID))
	fmt.Fprintf(b, "msgstr %s\n", EncodePOString(e.MsgStr))

	return b.String()
}

type PO struct {
	Comment string
	Header  Header
	Entries []*Entry

	contexts map[string]map[string]*Entry
}

func NewPO(comment string, creationDate time.Time, lang string) *PO {
	return &PO{
		Comment:  comment,
		Header:   newHeader(creationDate, lang),
		Entries:  make([]*Entry, 0),
		contexts: make(map[string]map[string]*Entry),
	}
}

func (p *PO) AddEntry(e *Entry) {
	context, exists := p.contexts[e.MsgContext]
	if !exists {
		context = make(map[string]*Entry)
		p.contexts[e.MsgContext] = context
	}

	existing := context[e.MsgID]
	if existing != nil {
		*existing = *e
	} else {
		context[e.MsgID] = e
		p.Entries = append(p.Entries, e)
	}
}

func (p *PO) Write(w io.Writer) {
	io.WriteString(w, fmt.Sprintf("# %s\n", p.Comment))
	io.WriteString(w, "#\n")
	io.WriteString(w, p.Header.asEntry().String())
	io.WriteString(w, "\n")

	for _, entry := range p.Entries {
		io.WriteString(w, entry.String())
		io.WriteString(w, "\n")
	}
}

// EncodePOString encodes the string values that appear after msgid, mgstr etc
func EncodePOString(text string) string {
	if text == "" {
		return `""`
	}

	runes := []rune(text)

	b := strings.Builder{}
	lineCount := 0
	insideLine := false
	for _, r := range runes {
		if !insideLine {
			lineCount++
			if lineCount > 1 {
				b.WriteRune('\n')
			}
			b.WriteRune('"')
			insideLine = true
		}

		switch r {
		case '\\':
			b.WriteString(`\\`)
		case '"':
			b.WriteString(`\"`)
		case '\n':
			b.WriteString(`\n`)

			// finish this line
			b.WriteRune('"')
			insideLine = false

		case '\t':
			b.WriteString(`\t`)
		default:
			b.WriteRune(r)
		}
	}

	if insideLine {
		b.WriteRune('"')
	}

	// multiline strings always start with "" on its own line
	if lineCount > 1 {
		return "\"\"\n" + b.String()
	}

	return b.String()
}
