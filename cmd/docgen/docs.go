package main

import (
	"fmt"
	"go/doc"
	"go/parser"
	"go/token"
	"path"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

var tagLineRegex = regexp.MustCompile(`@\w+\s+(?P<value>\w+)(?P<extra>.+)?`)

var docSets = []struct {
	searchDirs []string
	tag        string
	renderer   renderFunc
}{
	{[]string{"excellent/functions"}, "function", renderFunctionDoc},
	{[]string{"assets"}, "asset", renderAssetDoc},
	{[]string{"flows"}, "context", renderContextDoc},
	{[]string{"flows/routers/cases"}, "test", renderFunctionDoc},
	{[]string{"flows/actions"}, "action", renderActionDoc},
	{[]string{"flows/events"}, "event", renderEventDoc},
	{[]string{"flows/triggers"}, "trigger", renderTriggerDoc},
	{[]string{"flows/resumes"}, "resume", renderResumeDoc},
}

type documentedItem struct {
	typeName    string   // actual go type name
	tagName     string   // tag used to make this as a documented item
	tagValue    string   // identifier value after @tag
	tagExtra    string   // additional text after tag value
	examples    []string // any indented line
	description []string // any other line
}

type renderFunc func(output *strings.Builder, item *documentedItem, session flows.Session) error

// builds the documentation generation context from the given documented items
func buildDocsContext(items map[string][]*documentedItem) (map[string]string, error) {
	server := test.NewTestHTTPServer(49998)
	defer server.Close()

	defer utils.SetRand(utils.DefaultRand)
	defer utils.SetUUIDGenerator(utils.DefaultUUIDGenerator)
	defer utils.SetTimeSource(utils.DefaultTimeSource)

	utils.SetRand(utils.NewSeededRand(123456))
	utils.SetUUIDGenerator(utils.NewSeededUUID4Generator(123456))
	utils.SetTimeSource(utils.NewFixedTimeSource(time.Date(2018, 4, 11, 18, 24, 30, 123456000, time.UTC)))

	session, _, err := test.CreateTestSession(server.URL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating example session")
	}

	context := make(map[string]string, len(docSets))

	for _, ds := range docSets {
		contextKey := fmt.Sprintf("%sDocs", ds.tag)

		if context[contextKey], err = buildDocSet(ds.tag, items[ds.tag], ds.renderer, session); err != nil {
			return nil, err
		}
	}

	return context, nil
}

// builds a docset for the given tag type
func buildDocSet(tag string, tagItems []*documentedItem, renderer renderFunc, session flows.Session) (string, error) {
	// sort documented items by their tag value
	sort.SliceStable(tagItems, func(i, j int) bool { return tagItems[i].tagValue < tagItems[j].tagValue })

	buffer := &strings.Builder{}

	for _, item := range tagItems {
		if err := renderer(buffer, item, session); err != nil {
			return "", errors.Wrapf(err, "error rendering %s:%s", item.tagName, item.tagValue)
		}
	}

	return buffer.String(), nil
}

// finds all documented items for all tag types
func findAllDocumentedItems(baseDir string) (map[string][]*documentedItem, error) {
	items := make(map[string][]*documentedItem)
	for _, ds := range docSets {
		tagItems := make([]*documentedItem, 0)
		for _, searchDir := range ds.searchDirs {
			fromDir, err := findDocumentedItems(baseDir, searchDir, ds.tag)
			if err != nil {
				return nil, err
			}
			tagItems = append(tagItems, fromDir...)
		}
		items[ds.tag] = tagItems

		fmt.Printf(" > Found %d documented items with tag %s\n", len(tagItems), ds.tag)
	}
	return items, nil
}

// finds all documented items in go files in the given directory
func findDocumentedItems(baseDir string, searchDir string, tag string) ([]*documentedItem, error) {
	documentedItems := make([]*documentedItem, 0)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path.Join(baseDir, searchDir), nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	tag = "@" + tag

	for _, f := range pkgs {
		p := doc.New(f, "./", 0)
		for _, t := range p.Types {
			if strings.Contains(t.Doc, tag) {
				documentedItems = append(documentedItems, parseDocString(tag, t.Doc, t.Name))
			}
		}
		for _, t := range p.Funcs {
			if strings.Contains(t.Doc, tag) {
				documentedItems = append(documentedItems, parseDocString(tag, t.Doc, t.Name))
			}
		}
	}

	return documentedItems, nil
}

func parseDocString(tag string, docString string, typeName string) *documentedItem {
	var tagValue, tagExtra string
	examples := make([]string, 0)
	description := make([]string, 0)

	docString = removeTypeNamePrefix(docString, typeName)

	for _, l := range strings.Split(docString, "\n") {
		trimmed := strings.TrimSpace(l)

		if strings.HasPrefix(l, tag) {
			parts := tagLineRegex.FindStringSubmatch(l)
			tagValue = parts[1]
			tagExtra = parts[2]
		} else if strings.HasPrefix(l, "  ") { // examples are indented by at least two spaces
			trimmed = strings.Replace(trimmed, "->", "→", -1)
			examples = append(examples, trimmed)
		} else {
			description = append(description, l)
		}
	}

	return &documentedItem{typeName: typeName, tagName: tag[1:], tagValue: tagValue, tagExtra: tagExtra, examples: examples, description: description}
}

// Golang convention is to start all docstrings with the type name, but the actual type name can differ from how the type is
// referenced in the flow spec, so we remove it.
func removeTypeNamePrefix(docString string, typeName string) string {
	// remove type name from start of description and capitalize the next word
	if strings.HasPrefix(docString, typeName) {
		docString = strings.Replace(docString, typeName, "", 1)
		docString = strings.TrimSpace(docString)
		docString = strings.ToUpper(docString[0:1]) + docString[1:]
	}
	return docString
}
