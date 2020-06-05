package docs

import (
	"go/doc"
	"go/parser"
	"go/token"
	"path"
	"regexp"
	"sort"
	"strings"
)

var searchDirs = []string{
	"assets",
	"excellent/functions",
	"excellent/operators",
	"excellent/types",
	"flows",
	"flows/actions",
	"flows/definition",
	"flows/events",
	"flows/inputs",
	"flows/resumes",
	"flows/routers/cases",
	"flows/runs",
	"flows/triggers",
}

// the format of the tags which indicate a docstring is used by docgen: @name value<extra>
var tagRegex = regexp.MustCompile(`^@(?P<name>\w+)\s+(?P<value>\w+)(?P<extra>\(.*\))?(\s("(?P<title>.+)"))?$`)

// TaggedItem is any item that is documented with a @tag to indicate it will be used by docgen
type TaggedItem struct {
	typeName    string   // actual go type name
	tagName     string   // tag used to make this as a documented item
	tagValue    string   // identifier value after @tag
	tagExtra    string   // additional text after tag value in (...)
	tagTitle    string   // additional text after tag value in "" which becomes item title
	examples    []string // indented example lines
	description []string // the other lines
}

// FindAllTaggedItems finds all tagged docstrings in go files
func FindAllTaggedItems(baseDir string) (map[string][]*TaggedItem, error) {
	items := make(map[string][]*TaggedItem)

	// if tagged method is on a base class, we'll "find" it on each type that embeds that base
	// so need to ignore repeats
	seen := make(map[string]bool)

	for _, dir := range searchDirs {
		err := findTaggedItems(baseDir, dir, func(item *TaggedItem) {
			fullTag := item.tagName + ":" + item.tagValue
			if !seen[fullTag] {
				items[item.tagName] = append(items[item.tagName], item)
				seen[fullTag] = true
			}
		})
		if err != nil {
			return nil, err
		}
	}

	for _, v := range items {
		// sort items by their tag value
		sort.SliceStable(v, func(i, j int) bool { return v[i].tagValue < v[j].tagValue })
	}

	return items, nil
}

// finds all tagged docstrings in go files in the given directory
func findTaggedItems(baseDir string, searchDir string, callback func(item *TaggedItem)) error {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path.Join(baseDir, searchDir), nil, parser.ParseComments)
	if err != nil {
		return err
	}

	tryToParse := func(doc string, typeName string) {
		taggedItem := parseTaggedItem(doc, typeName)
		if taggedItem != nil {
			callback(taggedItem)
		}
	}

	for _, f := range pkgs {
		p := doc.New(f, "./", doc.AllDecls)
		for _, t := range p.Types {
			tryToParse(t.Doc, t.Name)
			for _, m := range t.Methods {
				tryToParse(m.Doc, m.Name)
			}
		}
		for _, t := range p.Funcs {
			tryToParse(t.Doc, t.Name)
		}
		for _, t := range p.Vars {
			tryToParse(t.Doc, t.Names[0])
		}
	}

	return nil
}

// tries to parse the given docstring as a tagged item
func parseTaggedItem(doc string, typeName string) *TaggedItem {
	lines := strings.Split(strings.TrimSpace(doc), "\n")
	if len(lines) == 0 {
		return nil
	}

	last := lines[len(lines)-1]
	if len(last) == 0 {
		return nil
	}

	tagParts := tagRegex.FindStringSubmatch(last)
	if len(tagParts) == 0 {
		return nil
	}

	lines = lines[:len(lines)-1]
	lines[0] = removeTypeNamePrefix(lines[0], typeName)

	examples := make([]string, 0)
	description := make([]string, 0)

	for _, l := range lines {
		trimmed := strings.TrimSpace(l)

		if strings.HasPrefix(l, "  ") { // examples are indented by at least two spaces
			trimmed = strings.Replace(trimmed, "->", "→", -1)
			examples = append(examples, trimmed)
		} else {
			description = append(description, l)
		}
	}

	title := tagParts[6]
	if title == "" {
		title = tagParts[2] + tagParts[3]
	}

	return &TaggedItem{
		typeName:    typeName,
		tagName:     tagParts[1],
		tagValue:    tagParts[2],
		tagExtra:    tagParts[3],
		tagTitle:    title,
		description: description,
		examples:    examples,
	}
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
