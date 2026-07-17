package flows

import (
	"github.com/nyaruka/goflow/excellent"
)

// Cache is a container for caches of parsed things whose lifetime should be tied to the assets they are parsed
// from rather than to a single session - e.g. templates which are shared by all sessions using the same assets.
type Cache struct {
	templates *excellent.TemplateCache
}

// NewCache creates a new empty cache
func NewCache() *Cache {
	return &Cache{templates: excellent.NewTemplateCache()}
}

// Templates returns the cache of parsed templates
func (c *Cache) Templates() *excellent.TemplateCache { return c.templates }
