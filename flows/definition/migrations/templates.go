package migrations

import (
	_ "embed"
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
	/*for _, n := range f.Nodes() {
		for _, a := range n.Actions() {
			for _, p := range catalog.Actions[a.Type()] {
				rewriteObjTemplates(a, parsePath(p), tx)
			}
		}

		if n.Router() != nil {
			for _, p := range catalog.Routers[n.Router().Type()] {
				rewriteObjTemplates(n.Router(), parsePath(p), tx)
			}
		}
	}*/
}

func rewriteValues(obj any, path []string, tx func(string) string) {

}
