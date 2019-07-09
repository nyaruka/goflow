package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

func init() {
	registerGenerator("template docs", generateTemplateDocs)
}

const (
	templateDir string = "cmd/docgen/templates"
)

type urlResolver func(string, string) (string, error)

var resources = []string{"styles.css"}
var templates = []struct {
	title         string
	path          string
	containsTypes []string
}{
	{"Flow Specification", "index.md", nil},
	{"Flows", "flows.md", []string{"action", "router", "wait"}},
	{"Expressions", "expressions.md", []string{"type", "operator", "function"}},
	{"Context", "context.md", []string{"context"}},
	{"Routing", "routing.md", []string{"test"}},
	{"Sessions", "sessions.md", []string{"event", "trigger", "resume"}},
	{"Assets", "assets.md", []string{"asset"}},
}

var docSets = []struct {
	tag      string
	renderer renderFunc
}{
	{"type", renderTypeDoc},
	{"operator", renderOperatorDoc},
	{"function", renderFunctionDoc},
	{"asset", renderAssetDoc},
	{"context", renderContextDoc},
	{"test", renderFunctionDoc},
	{"action", renderActionDoc},
	{"event", renderEventDoc},
	{"trigger", renderTriggerDoc},
	{"resume", renderResumeDoc},
}

type renderFunc func(output *strings.Builder, item *TaggedItem, session flows.Session) error

func generateTemplateDocs(baseDir string, outputDir string, items map[string][]*TaggedItem) error {
	if err := renderTemplateDocs(baseDir, outputDir, items); err != nil {
		return errors.Wrap(err, "error rendering templates")
	}

	// copy static resources to docs dir
	for _, res := range resources {
		src := path.Join(baseDir, templateDir, res)
		dst := path.Join(outputDir, res)
		if err := copyFile(src, dst); err != nil {
			return errors.Wrap(err, "error copying resource")
		}
		fmt.Printf(" > Copied %s > %s\n", src, dst)
	}
	return nil
}

func renderTemplateDocs(baseDir string, outputDir string, items map[string][]*TaggedItem) error {
	// render items as context for the main doc templates
	context, err := buildDocsContext(items)
	if err != nil {
		return errors.Wrap(err, "error building docs context")
	}

	// to post-process templates to resolve links between templates
	linkResolver, linkTargets := createLinkResolver(items)

	// ensure our output directory exists
	if err := os.MkdirAll(path.Join(outputDir), 0777); err != nil {
		return err
	}
	if err := os.MkdirAll(path.Join(outputDir, "md"), 0777); err != nil {
		return err
	}

	for _, template := range templates {
		templatePath := path.Join(baseDir, templateDir, template.path)
		renderedPath := path.Join(outputDir, "md", template.path)
		htmlPath := path.Join(outputDir, template.path[0:len(template.path)-3]+".html")

		if err := renderTemplate(templatePath, renderedPath, context, linkResolver, linkTargets); err != nil {
			return errors.Wrapf(err, "error rendering template %s", templatePath)
		}

		htmlTemplate := path.Join(baseDir, "cmd/docgen/templates/template.html")
		htmlContext := map[string]string{"title": template.title}

		if err := renderHTML(renderedPath, htmlPath, htmlTemplate, htmlContext); err != nil {
			return errors.Wrapf(err, "error rendering HTML from %s to %s", renderedPath, htmlPath)
		}

		fmt.Printf(" > Rendered %s > %s > %s\n", templatePath, renderedPath, htmlPath)
	}

	return nil
}

// renders a markdown template
func renderTemplate(src, dst string, context map[string]string, resolver urlResolver, linkTargets map[string]bool) error {
	// generate our complete docs
	docTpl, err := template.ParseFiles(src)
	if err != nil {
		return errors.Wrap(err, "error reading template file")
	}

	output := &strings.Builder{}
	if err := docTpl.Execute(output, context); err != nil {
		return errors.Wrap(err, "error executing template")
	}

	processed := resolveLinks(output.String(), resolver, linkTargets)

	return ioutil.WriteFile(dst, []byte(processed), 0666)
}

// converts a markdown file to HTML
func renderHTML(src, dst, htmlTemplate string, variables map[string]string) error {
	panDocArgs := []string{
		"--from=markdown",
		"--to=html",
		"-o", dst,
		"--standalone",
		"--template=" + htmlTemplate,
		"--toc",
		"--toc-depth=1",
	}

	for k, v := range variables {
		panDocArgs = append(panDocArgs, fmt.Sprintf("--variable=%s:%s", k, v))
	}

	panDocArgs = append(panDocArgs, src)
	return exec.Command("pandoc", panDocArgs...).Run()
}

func createLinkResolver(items map[string][]*TaggedItem) (urlResolver, map[string]bool) {
	linkTargets := make(map[string]bool)
	typeTemplates := make(map[string]string)

	for _, ds := range docSets {
		for _, item := range items[ds.tag] {
			linkTargets[ds.tag+":"+item.tagValue] = true
		}
	}

	for _, template := range templates {
		for _, typeTag := range template.containsTypes {
			typeTemplates[typeTag] = fmt.Sprintf("%s.html#%s:%%s", template.path[0:len(template.path)-3], typeTag)
		}
	}

	return func(tag string, val string) (string, error) {
		linkTpl := typeTemplates[tag]
		if linkTpl == "" {
			return "", errors.Errorf("no link template for type %s", tag)
		}
		return fmt.Sprintf(linkTpl, val), nil
	}, linkTargets
}

func resolveLinks(s string, resolver urlResolver, targets map[string]bool) string {
	r := regexp.MustCompile(`\[\w+:\w+\]`)
	return r.ReplaceAllStringFunc(s, func(old string) string {
		target := old[1 : len(old)-1]
		if !targets[target] {
			panic(fmt.Sprintf("found link to %s which is not a valid target", target))
		}

		groups := strings.Split(target, ":")
		url, err := resolver(groups[0], groups[1])
		if err != nil {
			panic(err.Error())
		}
		return fmt.Sprintf("[%s](%s)", groups[1], url)
	})
}

// copies a file from one path to another
func copyFile(src, dst string) error {
	contents, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dst, contents, 0666)
}

// builds the documentation generation context from the given documented items
func buildDocsContext(items map[string][]*TaggedItem) (map[string]string, error) {
	server := test.NewTestHTTPServer(49998)
	defer server.Close()

	defer utils.SetRand(utils.DefaultRand)
	defer utils.SetUUIDGenerator(utils.DefaultUUIDGenerator)
	defer utils.SetTimeSource(utils.DefaultTimeSource)

	utils.SetRand(utils.NewSeededRand(123456))
	utils.SetUUIDGenerator(test.NewSeededUUIDGenerator(123456))
	utils.SetTimeSource(test.NewFixedTimeSource(time.Date(2018, 4, 11, 18, 24, 30, 123456000, time.UTC)))

	session, _, err := test.CreateTestSession(server.URL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating example session")
	}

	context := make(map[string]string, len(docSets))

	for _, ds := range docSets {
		contextKey := fmt.Sprintf("%sDocs", ds.tag)

		if context[contextKey], err = buildTagContext(ds.tag, items[ds.tag], ds.renderer, session); err != nil {
			return nil, err
		}
	}

	return context, nil
}

// builds a docset for the given tag type
func buildTagContext(tag string, tagItems []*TaggedItem, renderer renderFunc, session flows.Session) (string, error) {
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
