package docs

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/random"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/test"
)

func init() {
	RegisterGenerator(&htmlDocsGenerator{})
}

const (
	templateDir string = "cmd/docgen/templates"
)

type urlResolver func(string, string) (string, error)

var Resources = []string{"styles.css"}
var Templates = []struct {
	Title         string
	Path          string
	ContainsTypes []string // used for resolving links
	TOC           bool
}{
	{"Flow Specification", "index.md", []string{"version"}, false},
	{"Flows", "flows.md", []string{"action", "router", "wait"}, true},
	{"Expressions", "expressions.md", []string{"type", "operator", "function"}, true},
	{"Context", "context.md", []string{"context"}, true},
	{"Routing", "routing.md", []string{"test"}, true},
	{"Sessions", "sessions.md", []string{"event", "trigger", "resume"}, true},
	{"Assets", "assets.md", []string{"asset"}, true},
}

// ContextFunc is a function which produces values to put the template context
type ContextFunc func(map[string][]*TaggedItem, flows.Session, flows.Session) (map[string]string, error)

var contextFuncs []ContextFunc

func registerContextFunc(f ContextFunc) {
	contextFuncs = append(contextFuncs, f)
}

type htmlDocsGenerator struct{}

func (g *htmlDocsGenerator) Name() string {
	return "html docs"
}

func (g *htmlDocsGenerator) Generate(baseDir, outputDir string, items map[string][]*TaggedItem, gettext func(string) string) error {
	if err := renderTemplateDocs(baseDir, outputDir, items); err != nil {
		return fmt.Errorf("error rendering templates: %w", err)
	}

	// copy static resources to docs dir
	for _, res := range Resources {
		src := path.Join(baseDir, templateDir, res)
		dst := path.Join(outputDir, res)
		if err := copyFile(src, dst); err != nil {
			return fmt.Errorf("error copying resource: %w", err)
		}
		fmt.Printf(" > Copied %s > %s\n", src, dst)
	}
	return nil
}

func renderTemplateDocs(baseDir string, outputDir string, items map[string][]*TaggedItem) error {
	// render items as context for the main doc templates
	context, err := buildTemplateContext(items)
	if err != nil {
		return fmt.Errorf("error building docs context: %w", err)
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

	for _, template := range Templates {
		templatePath := path.Join(baseDir, templateDir, template.Path)
		renderedPath := path.Join(outputDir, "md", template.Path)
		htmlPath := path.Join(outputDir, template.Path[0:len(template.Path)-3]+".html")

		if err := renderTemplate(templatePath, renderedPath, context, linkResolver, linkTargets); err != nil {
			return fmt.Errorf("error rendering template %s: %w", templatePath, err)
		}

		htmlTemplate := path.Join(baseDir, "cmd/docgen/templates/template.html")
		htmlContext := map[string]string{"title": template.Title}

		if err := renderHTML(renderedPath, htmlPath, htmlTemplate, template.TOC, htmlContext); err != nil {
			return fmt.Errorf("error rendering HTML from %s to %s: %w", renderedPath, htmlPath, err)
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
		return fmt.Errorf("error reading template file: %w", err)
	}

	output := &strings.Builder{}
	if err := docTpl.Execute(output, context); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	processed := resolveLinks(output.String(), resolver, linkTargets)

	return os.WriteFile(dst, []byte(processed), 0666)
}

// converts a markdown file to HTML
func renderHTML(src, dst, htmlTemplate string, toc bool, variables map[string]string) error {
	panDocArgs := []string{
		"--from=markdown",
		"--to=html",
		"-o", dst,
		"--standalone",
		"--template=" + htmlTemplate,
	}

	if toc {
		panDocArgs = append(panDocArgs, "--toc", "--toc-depth=1")
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

	for tag, itemsWithTag := range items {
		for _, item := range itemsWithTag {
			linkTargets[tag+":"+item.tagValue] = true
		}
	}

	for _, template := range Templates {
		for _, typeTag := range template.ContainsTypes {
			typeTemplates[typeTag] = fmt.Sprintf("%s.html#%s:%%s", template.Path[0:len(template.Path)-3], typeTag)
		}
	}

	return func(tag string, val string) (string, error) {
		linkTpl := typeTemplates[tag]
		if linkTpl == "" {
			return "", fmt.Errorf("no link template for type %s", tag)
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
	contents, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, contents, 0666)
}

// builds the documentation generation context from the given documented items
func buildTemplateContext(items map[string][]*TaggedItem) (map[string]string, error) {
	server := test.NewHTTPServer(49998, test.MockWebhooksHandler)
	defer server.Close()

	defer random.SetGenerator(random.DefaultGenerator)
	defer uuids.SetGenerator(uuids.DefaultGenerator)
	defer dates.SetNowFunc(time.Now)

	random.SetGenerator(random.NewSeededGenerator(123456))
	uuids.SetGenerator(uuids.NewSeededGenerator(123456, time.Now))
	dates.SetNowFunc(dates.NewFixedNow(time.Date(2018, 4, 11, 18, 24, 30, 123456000, time.UTC)))

	session, _, err := test.CreateTestSession(server.URL, envs.RedactionPolicyNone)
	if err != nil {
		return nil, fmt.Errorf("error creating example session: %w", err)
	}

	voiceSession, _, err := test.CreateTestVoiceSession(server.URL)
	if err != nil {
		return nil, fmt.Errorf("error creating example session: %w", err)
	}

	context := make(map[string]string)

	for _, f := range contextFuncs {
		newContext, err := f(items, session, voiceSession)
		if err != nil {
			return nil, err
		}

		// add to the overall context
		for k, v := range newContext {
			context[k] = v
		}
	}

	return context, nil
}
