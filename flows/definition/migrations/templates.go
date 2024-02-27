package migrations

import (
	_ "embed"
	"strings"
	"sync"

	"github.com/Masterminds/semver"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/flows/definition/migrations/jsonpath"
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
	txl := func(container, key, val any) any {
		localizableUUID := getObjectUUID(container)
		if localizableUUID != "" {
			prop, _ := key.(string)
			if prop != "" {
				rewriteTranslations(f, localizableUUID, prop, tx)
			}
		}

		switch vtyped := val.(type) {
		case string:
			return tx(vtyped)
		case []any:
			vnew := make([]any, len(vtyped))
			for i, v := range vtyped {
				vs, ok := v.(string)
				if ok {
					vnew[i] = tx(vs)
				} else {
					vnew[i] = v
				}
			}
			return vnew
		}
		return val
	}

	for _, n := range f.Nodes() {
		for _, a := range n.Actions() {
			for _, p := range catalog.Actions[a.Type()] {
				rewriteTemplates(a, p, txl)
			}
		}

		if n.Router() != nil {
			for _, p := range catalog.Routers[n.Router().Type()] {
				rewriteTemplates(n.Router(), p, txl)
			}
		}
	}
}

func rewriteTemplates[T ~map[string]any](o T, path string, tx func(container, key, val any) any) {
	// if path points to an array of strings, get the array instead of the individual strings
	path = strings.TrimSuffix(path, "[*]")

	jsonpath.Transform(map[string]any(o), jsonpath.ParsePath(path), tx)
}

func rewriteTranslations(f Flow, uuid, property string, tx func(string) string) {
	localization := f.Localization()
	if localization != nil {
		for _, lang := range localization.Languages() {
			langTrans := localization.GetLanguageTranslation(lang)
			if langTrans != nil {
				itemTrans := langTrans.GetItemTranslation(uuid)
				if itemTrans != nil {
					ss := itemTrans.Get(property)
					if ss != nil {
						for i, s := range ss {
							ss[i] = tx(s)
						}
						itemTrans.Set(property, ss)
					}
				}
			}
		}
	}
}

// gets the UUID property of o, if o is an object, if it has "uuid" property, and if the type of that property is a string
func getObjectUUID(o any) string {
	m, ok := o.(map[string]any)
	if ok {
		v, exists := m["uuid"]
		if exists {
			s, ok := v.(string)
			if ok {
				return s
			}
		}
	}
	return ""
}
