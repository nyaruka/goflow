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

	"github.com/pkg/errors"
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
}{
	{"Flow Specification", "index.md", nil},
	{"Flows", "flows.md", []string{"action", "router", "wait"}},
	{"Expressions", "expressions.md", []string{"type", "operator", "function"}},
	{"Context", "context.md", []string{"context"}},
	{"Routing", "routing.md", []string{"test"}},
	{"Sessions", "sessions.md", []string{"event", "trigger", "resume"}},
	{"Assets", "assets.md", []string{"asset"}},
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
		return errors.Wrap(err, "error rendering templates")
	}

	// copy static resources to docs dir
	for _, res := range Resources {
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
	context, err := buildTemplateContext(items)
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

	for _, template := range Templates {
		templatePath := path.Join(baseDir, templateDir, template.Path)
		renderedPath := path.Join(outputDir, "md", template.Path)
		htmlPath := path.Join(outputDir, template.Path[0:len(template.Path)-3]+".html")

		if err := renderTemplate(templatePath, renderedPath, context, linkResolver, linkTargets); err != nil {
			return errors.Wrapf(err, "error rendering template %s", templatePath)
		}

		htmlTemplate := path.Join(baseDir, "cmd/docgen/templates/template.html")
		htmlContext := map[string]string{"title": template.Title}

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

	return os.WriteFile(dst, []byte(processed), 0666)
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
	contents, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, contents, 0666)
}

// builds the documentation generation context from the given documented items
func buildTemplateContext(items map[string][]*TaggedItem) (map[string]string, error) {
	server := test.NewTestHTTPServer(49998)
	defer server.Close()

	defer random.SetGenerator(random.DefaultGenerator)
	defer uuids.SetGenerator(uuids.DefaultGenerator)
	defer dates.SetNowSource(dates.DefaultNowSource)

	random.SetGenerator(random.NewSeededGenerator(123456))
	uuids.SetGenerator(uuids.NewSeededGenerator(123456))
	dates.SetNowSource(dates.NewFixedNowSource(time.Date(2018, 4, 11, 18, 24, 30, 123456000, time.UTC)))

	session, _, err := test.CreateTestSession(server.URL, envs.RedactionPolicyNone)
	if err != nil {
		return nil, errors.Wrap(err, "error creating example session")
	}

	voiceSession, _, err := test.CreateTestVoiceSession(server.URL)
	if err != nil {
		return nil, errors.Wrap(err, "error creating example session")
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
