package main

// generate full docs with:
//
// go install github.com/nyaruka/goflow/cmd/docgen; docgen

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"text/template"
)

const (
	templateDir string = "cmd/docgen/templates"
	outputDir          = "docs"
)

var resources = []string{"styles.css"}
var templates = []struct {
	title string
	path  string
}{
	{"Flow Specification", "index.md"},
	{"Flows", "flows.md"},
	{"Sessions", "sessions.md"},
}

func main() {
	if err := generateDocs(".", outputDir); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func generateDocs(baseDir string, outputDir string) error {
	context, err := buildDocsContext(baseDir)
	if err != nil {
		return fmt.Errorf("error building docs context: %s", err)
	}

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

		if err := renderTemplate(templatePath, renderedPath, context); err != nil {
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
func renderTemplate(src, dst string, context map[string]string) error {
	// generate our complete docs
	docTpl, err := template.ParseFiles(src)
	if err != nil {
		return fmt.Errorf("error reading template file: %s", err)
	}

	output := &bytes.Buffer{}
	if err := docTpl.Execute(output, context); err != nil {
		return fmt.Errorf("error executing template: %s", err)
	}

	return ioutil.WriteFile(dst, output.Bytes(), 0666)
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

// copies a file from one path to another
func copyFile(src, dst string) error {
	contents, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dst, contents, 0666)
}
