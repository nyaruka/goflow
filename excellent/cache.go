package excellent

import (
	"sync"
)

// TemplateCache is a thread-safe content-addressed cache of parsed templates. Because a parsed template is a pure
// function of its source and immutable, entries never go stale and the cache needs no invalidation - its lifetime
// should instead be bound to that of the content the templates come from, e.g. a set of assets.
type TemplateCache struct {
	templates sync.Map // string -> *Template
}

// NewTemplateCache creates a new empty template cache
func NewTemplateCache() *TemplateCache {
	return &TemplateCache{}
}

// Get returns the parsed form of the given template source, parsing and caching it if not already cached.
func (c *TemplateCache) Get(src string) *Template {
	if t, ok := c.templates.Load(src); ok {
		return t.(*Template)
	}

	// concurrent misses may both parse but LoadOrStore ensures they get the same instance
	t, _ := c.templates.LoadOrStore(src, ParseTemplate(src))
	return t.(*Template)
}
