package main

// generate full docs with:
//
// go install github.com/nyaruka/goflow/cmd/docgen; docgen

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"text/template"

	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

const (
	templateDir string = "cmd/docgen/templates"
	outputDir          = "docs"
)

type urlResolver func(string, string) (string, error)

var resources = []string{"styles.css"}
var templates = []struct {
	title         string
	path          string
	containsTypes []string
}{
	{"Flow Specification", "index.md", nil},
	{"Expressions", "expressions.md", []string{"type", "operator", "function"}},
	{"Context", "context.md", []string{"context", "test"}},
	{"Flows", "flows.md", []string{"action", "router", "wait"}},
	{"Sessions", "sessions.md", []string{"event", "trigger"}},
	{"Assets", "assets.md", []string{"asset"}},
}

func main() {
	if err := GenerateDocs(".", outputDir); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

// GenerateDocs generates out HTML documentation
func GenerateDocs(baseDir string, outputDir string) error {
	fmt.Println("Processing sources...")

	// extract all documented items from the source
	items, err := findAllDocumentedItems(baseDir)
	if err != nil {
		return errors.Wrap(err, "error extracting documented items")
	}

	fmt.Println("Rendering templates...")

	if err := renderDocs(baseDir, outputDir, items); err != nil {
		return errors.Wrap(err, "error rendering templates")
	}

	fmt.Println("Copying static resources...")

	// copy static resources to docs dir
	for _, res := range resources {
		src := path.Join(baseDir, templateDir, res)
		dst := path.Join(outputDir, res)
		if err := copyFile(src, dst); err != nil {
			return errors.Wrap(err, "error copying resource")
		}
		fmt.Printf(" > Copied %s > %s\n", src, dst)
	}

	fmt.Println("Generating function listing...")

	return generateFunctionListing(outputDir, items["function"])
}

func renderDocs(baseDir string, outputDir string, items map[string][]*documentedItem) error {
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

func createLinkResolver(items map[string][]*documentedItem) (urlResolver, map[string]bool) {
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

type functionExample struct {
	Template string `json:"template"`
	Output   string `json:"output"`
}

type functionListing struct {
	Signature string             `json:"signature"`
	Summary   string             `json:"summary"`
	Detail    string             `json:"detail"`
	Examples  []*functionExample `json:"examples"`
}

func generateFunctionListing(outputDir string, funcItems []*documentedItem) error {
	listings := make([]*functionListing, len(funcItems))
	for f, funcItem := range funcItems {
		summary := funcItem.description[0]
		detail := strings.TrimSpace(strings.Join(funcItem.description[1:len(funcItem.description)-1], "\n"))

		examples := make([]*functionExample, len(funcItem.examples))
		for e := range funcItem.examples {
			parts := strings.Split(funcItem.examples[e], "→")
			examples[e] = &functionExample{Template: strings.TrimSpace(parts[0]), Output: strings.TrimSpace(parts[1])}
		}

		listings[f] = &functionListing{
			Signature: funcItem.tagValue + funcItem.tagExtra,
			Summary:   summary,
			Detail:    detail,
			Examples:  examples,
		}
	}

	data, err := utils.JSONMarshalPretty(listings)
	if err != nil {
		return err
	}

	listingPath := path.Join(outputDir, "functions.json")

	return ioutil.WriteFile(listingPath, []byte(data), 0666)
}
