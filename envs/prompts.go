package envs

import "text/template"

// PromptResolver is a function that resolves a prompt name to a template.
type PromptResolver func(string) *template.Template

// EmptyPromptResolver is an LLM prompt resolver that always returns nil.
func EmptyPromptResolver(string) *template.Template { return nil }

// NewPromptResolver creates a new LLM prompt resolver that uses the provided map of templates.
func NewPromptResolver(static map[string]*template.Template) PromptResolver {
	return func(name string) *template.Template {
		return static[name]
	}
}
