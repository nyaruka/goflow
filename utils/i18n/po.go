package i18n

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

const poDatetimeformat = "2006-01-02 15:04-0700"

// POHeader contains metadata about a PO file
type POHeader struct {
	InitialComment  string
	POTCreationDate time.Time         // POT-Creation-Date: YEAR-MO-DA HO:MI+ZONE
	Language        string            // Language: e.g. en-US
	MIMEVersion     string            // MIME-Version: 1.0
	ContentType     string            // Content-Type: text/plain; charset=UTF-8
	Custom          map[string]string // other custom values
}

// NewPOHeader creates a new PO header with the given values
func NewPOHeader(initialComment string, creationDate time.Time, lang string) *POHeader {
	return &POHeader{
		InitialComment:  initialComment + "\n",
		POTCreationDate: creationDate,
		Language:        lang,
		MIMEVersion:     "1.0",
		ContentType:     "text/plain; charset=UTF-8",
		Custom:          make(map[string]string),
	}
}

// headers are deserialized as regular entries and converted here
func newPOHeaderFromEntry(e *POEntry) *POHeader {
	h := &POHeader{
		InitialComment: strings.Join(e.Comment.Translator, "\n"),
		Custom:         make(map[string]string),
	}

	for _, line := range strings.Split(e.MsgStr, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := parts[0]
			value := strings.TrimSpace(parts[1])
			switch key {
			case "POT-Creation-Date":
				h.POTCreationDate, _ = time.Parse(poDatetimeformat, value)
			case "Language":
				h.Language = value
			case "MIME-Version":
				h.MIMEVersion = value
			case "Content-Type":
				h.ContentType = value
			default:
				h.Custom[key] = value
			}
		}
	}

	return h
}

// convert header to an entry for serialization
func (h *POHeader) asEntry() *POEntry {
	b := &strings.Builder{}
	fmt.Fprintf(b, "POT-Creation-Date: %s\n", h.POTCreationDate.Format(poDatetimeformat))
	fmt.Fprintf(b, "Language: %s\n", h.Language)
	fmt.Fprintf(b, "MIME-Version: %s\n", h.MIMEVersion)
	fmt.Fprintf(b, "Content-Type: %s\n", h.ContentType)

	customKeys := make([]string, 0, len(h.Custom))
	for key := range h.Custom {
		customKeys = append(customKeys, key)
	}
	sort.Strings(customKeys)
	for _, key := range customKeys {
		fmt.Fprintf(b, "%s: %s\n", key, h.Custom[key])
	}

	return &POEntry{
		Comment: POComment{
			Translator: strings.Split(h.InitialComment, "\n"),
			Flags:      []string{"fuzzy"},
		},
		MsgID:  "",
		MsgStr: b.String(),
	}
}

// POComment is a comment for an entry
type POComment struct {
	Translator []string // #  translator-comments
	Extracted  []string // #. extracted-comments
	References []string // #: references
	Flags      []string // #, e.g. fuzzy,python-format
}

// ParsePOComment parses a PO file comment from the given string
func ParsePOComment(s string) POComment {
	c := POComment{}
	if s == "" {
		return c
	}

	for _, line := range strings.Split(s, "\n") {
		if line == "#" {
			c.Translator = append(c.Translator, line[1:])
			continue
		}

		trimmed := strings.TrimSpace(line[2:])

		if strings.HasPrefix(line, "# ") {
			c.Translator = append(c.Translator, trimmed)
		} else if strings.HasPrefix(line, "#.") {
			c.Extracted = append(c.Extracted, trimmed)
		} else if strings.HasPrefix(line, "#:") {
			for _, val := range strings.Split(trimmed, ",") {
				val = strings.TrimSpace(val)
				c.References = append(c.References, val)
			}
		} else if strings.HasPrefix(line, "#,") {
			for _, val := range strings.Split(trimmed, ",") {
				val = strings.TrimSpace(val)
				c.Flags = append(c.Flags, val)
			}
		}
	}
	return c
}

// HasFlag returns true if this comment contains the given flag
func (c *POComment) HasFlag(flag string) bool {
	for _, f := range c.Flags {
		if f == flag {
			return true
		}
	}
	return false
}

func (c *POComment) String() string {
	lines := make([]string, 0)

	for _, line := range c.Translator {
		lines = append(lines, fmt.Sprintf("#  %s", line))
	}
	for _, line := range c.Extracted {
		lines = append(lines, fmt.Sprintf("#. %s", line))
	}
	for _, ref := range c.References {
		lines = append(lines, fmt.Sprintf("#: %s", ref))
	}
	if len(c.Flags) > 0 {
		lines = append(lines, fmt.Sprintf("#, %s", strings.Join(c.Flags, ",")))
	}
	return strings.Join(lines, "\n")
}

// POEntry is an entry in a PO catalog
type POEntry struct {
	Comment    POComment // Comment
	MsgContext string    // msgctxt context
	MsgID      string    // msgid untranslated-string
	MsgStr     string    // msgstr translated-string
}

func (e *POEntry) Write(w io.Writer) {
	comment := e.Comment.String()
	if comment != "" {
		fmt.Fprintf(w, "%s\n", comment)
	}
	if e.MsgContext != "" {
		fmt.Fprintf(w, "msgctxt %s\n", EncodePOString(e.MsgContext))
	}
	fmt.Fprintf(w, "msgid %s\n", EncodePOString(e.MsgID))
	fmt.Fprintf(w, "msgstr %s\n", EncodePOString(e.MsgStr))
	fmt.Fprintln(w)
}

// PO is a PO file of translation entries
type PO struct {
	Header  *POHeader
	Entries []*POEntry

	contexts map[string]map[string]*POEntry
}

// NewPO creates a new PO catalog
func NewPO(h *POHeader) *PO {
	return &PO{
		Header:   h,
		Entries:  make([]*POEntry, 0),
		contexts: make(map[string]map[string]*POEntry),
	}
}

// AddEntry adds the given entry to this PO
func (p *PO) AddEntry(e *POEntry) {
	context, exists := p.contexts[e.MsgContext]
	if !exists {
		context = make(map[string]*POEntry)
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

// Sort sorts entries by ID and context
func (p *PO) Sort() {
	sort.SliceStable(p.Entries, func(i, j int) bool {
		e1 := p.Entries[i]
		e2 := p.Entries[j]
		cmp := strings.Compare(e1.MsgID, e2.MsgID)
		if cmp == 0 {
			return strings.Compare(e1.MsgContext, e2.MsgContext) < 0
		}
		return cmp < 0
	})
}

// GetText gets the translations of text with the given context (optional)
func (p *PO) GetText(context, text string) string {
	c, exists := p.contexts[context]
	if !exists {
		return text
	}
	entry := c[text]
	if entry == nil || entry.MsgStr == "" || entry.Comment.HasFlag("fuzzy") {
		return text
	}
	return entry.MsgStr
}

// Write writes this PO to the given writer
func (p *PO) Write(w io.Writer) {
	if p.Header != nil {
		p.Header.asEntry().Write(w)
	}
	for _, entry := range p.Entries {
		entry.Write(w)
	}
}

// ReadPO reads a PO file from the given reader
func ReadPO(r io.Reader) (*PO, error) {
	br := bufio.NewReader(r)
	nextLine := func() (string, error) {
		return br.ReadString('\n')
	}

	po := NewPO(nil)

	for {
		entry, err := readPOEntry(nextLine)
		if err != nil {
			return nil, err
		}
		if entry != nil {
			if entry.MsgID == "" {
				po.Header = newPOHeaderFromEntry(entry)
			} else {
				po.AddEntry(entry)
			}
		} else {
			break
		}
	}

	return po, nil
}

// reads a single entry
func readPOEntry(nextLine func() (string, error)) (*POEntry, error) {
	// read lines until we hit EOF or empty line
	lines := make([]string, 0)
	for {
		line, err := nextLine()
		line = strings.TrimSpace(line)
		if err == io.EOF || line == "" {
			break
		}
		if err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}

	// there wasn't another entry to read
	if len(lines) == 0 {
		return nil, nil
	}

	comment := ""
	values := make(map[string]string)
	currentKey := ""
	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			if comment != "" {
				comment += "\n"
			}
			comment += line
		} else if strings.HasPrefix(line, "msg") {
			parts := strings.Fields(line)
			currentKey = parts[0]
			rest := strings.TrimSpace(line[len(currentKey):])
			values[currentKey] = values[currentKey] + rest
		} else if strings.HasPrefix(line, `"`) {
			values[currentKey] += "\n" + line
		}
	}

	return &POEntry{
		Comment:    ParsePOComment(comment),
		MsgContext: DecodePOString(values["msgctxt"]),
		MsgID:      DecodePOString(values["msgid"]),
		MsgStr:     DecodePOString(values["msgstr"]),
	}, nil
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

// DecodePOString decodes the string values that appear after msgid, mgstr etc
func DecodePOString(s string) string {
	if s == "" {
		return ""
	}

	lines := strings.Split(s, "\n")
	b := &strings.Builder{}

	for _, line := range lines {
		line = line[1 : len(line)-1] // strip quotes

		unescaping := false
		for _, r := range []rune(line) {
			if unescaping {
				switch r {
				case '\\':
					b.WriteRune('\\')
				case '"':
					b.WriteRune('"')
				case 'n':
					b.WriteRune('\n')
				case 't':
					b.WriteRune('\t')
				}
				unescaping = false
			} else if r == '\\' {
				unescaping = true
			} else {
				b.WriteRune(r)
			}
		}
	}

	return b.String()
}
