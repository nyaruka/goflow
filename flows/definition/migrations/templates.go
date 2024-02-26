package migrations

import (
	_ "embed"
	"fmt"
	"sync"

	"github.com/Masterminds/semver"
	"github.com/nyaruka/gocommon/jsonx"
)

//go:embed specdata/templates.json
var templateCatalogsJSON []byte

var templateCatalogs map[string]*TemplateCatalog
var catalogInit sync.Once

type TemplateCatalog struct {
	Actions map[string][]string `json:"actions"`
	Routers map[string][]string `json:"routers"`
}

func init() {
	catalogInit.Do(func() {
		jsonx.MustUnmarshal(templateCatalogsJSON, &templateCatalogs)
	})
}

func GetTemplateCatalog(v *semver.Version) *TemplateCatalog {
	c := templateCatalogs[v.String()]
	if c == nil {
		panic("no template catalog for version " + v.String())
	}
	return c
}

func RewriteTemplates(f Flow, catalog *TemplateCatalog, tx func(string) string) {
	for _, n := range f.Nodes() {
		for _, a := range n.Actions() {
			rewriteObjTemplates(a, catalog.Actions[a.Type()], tx)
		}

		if n.Router() != nil {
			rewriteObjTemplates(n.Router(), catalog.Routers[n.Router().Type()], tx)
		}
	}
}

func rewriteObjTemplates(obj map[string]any, paths []string, tx func(string) string) {
	onObject := func(path string, m map[string]any) {
		for k, v := range m {
			p := fmt.Sprintf("%s.%s", path, k)

			if pathMatches(p, paths) {
				m[k] = tx(v.(string))
			}
		}
	}
	onArray := func(path string, a []any) {
		for i, v := range a {
			p := fmt.Sprintf("%s[%d]", path, i)

			if pathMatches(p, paths) {
				a[i] = tx(v.(string))
			}
		}
	}

	walk(obj, onObject, onArray, "")
}

func pathMatches(p string, paths []string) bool {
	return false
}
