package envs

import "text/template"

// LLMPromptResolver is a function that resolves a prompt name to a template.
type LLMPromptResolver func(string) *template.Template

// EmptyLLMPromptResolver is an LLM prompt resolver that always returns nil.
func EmptyLLMPromptResolver(string) *template.Template { return nil }

// NewLLMPromptResolver creates a new LLM prompt resolver that uses the provided map of templates.
func NewLLMPromptResolver(static map[string]*template.Template) LLMPromptResolver {
	return func(name string) *template.Template {
		return static[name]
	}
}
