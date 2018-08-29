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
	{"Expressions", "expressions.md", []string{"context", "function", "test"}},
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
	context, linkTargets, err := buildDocsContext(baseDir)
	if err != nil {
		return fmt.Errorf("error building docs context: %s", err)
	}

	// to post-process templates to resolve links between templates
	linkResolver := createLinkResolver()

	// ensure our output directory exists
	if err := os.MkdirAll(path.Join(outputDir), 0777); err != nil {
		return err
	}
	if err := os.MkdirAll(path.Join(outputDir, "md"), 0777); err != nil {
		return err
	}

	fmt.Println("Rendering templates...")

	for _, template := range templates {
		templatePath := path.Join(baseDir, templateDir, template.path)
		renderedPath := path.Join(outputDir, "md", template.path)
		htmlPath := path.Join(outputDir, template.path[0:len(template.path)-3]+".html")

		if err := renderTemplate(templatePath, renderedPath, context, linkResolver, linkTargets); err != nil {
			return fmt.Errorf("error rendering template %s: %s", templatePath, err)
		}

		htmlTemplate := path.Join(baseDir, "cmd/docgen/templates/template.html")
		htmlContext := map[string]string{"title": template.title}

		if err := renderHTML(renderedPath, htmlPath, htmlTemplate, htmlContext); err != nil {
			return fmt.Errorf("error rendering HTML from %s to %s: %s", renderedPath, htmlPath, err)
		}

		fmt.Printf(" > Rendered %s > %s > %s\n", templatePath, renderedPath, htmlPath)
	}

	fmt.Println("Copying static resources...")

	// copy static resources to docs dir
	for _, res := range resources {
		src := path.Join(baseDir, templateDir, res)
		dst := path.Join(outputDir, res)
		if err := copyFile(src, dst); err != nil {
			return fmt.Errorf("error copying resource: %s", err)
		}

		fmt.Printf(" > Copied %s > %s\n", src, dst)
	}

	return nil
}

// renders a markdown template
func renderTemplate(src, dst string, context map[string]string, resolver urlResolver, linkTargets map[string]bool) error {
	// generate our complete docs
	docTpl, err := template.ParseFiles(src)
	if err != nil {
		return fmt.Errorf("error reading template file: %s", err)
	}

	output := &strings.Builder{}
	if err := docTpl.Execute(output, context); err != nil {
		return fmt.Errorf("error executing template: %s", err)
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

func createLinkResolver() urlResolver {
	typeTemplates := make(map[string]string)

	for _, template := range templates {
		for _, typeTag := range template.containsTypes {
			typeTemplates[typeTag] = fmt.Sprintf("%s.html#%s:%%s", template.path[0:len(template.path)-3], typeTag)
		}
	}

	return func(tag string, val string) (string, error) {
		linkTpl := typeTemplates[tag]
		if linkTpl == "" {
			return "", fmt.Errorf("no link template for type %s", tag)
		}
		return fmt.Sprintf(linkTpl, val), nil
	}
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
