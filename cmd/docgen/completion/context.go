package completion

// Context is the runtime information required to generate completions
type Context struct {
	KeySources map[string][]string
}

// NewContext creates a new completion context
func NewContext(keySources map[string][]string) *Context {
	return &Context{KeySources: keySources}
}
