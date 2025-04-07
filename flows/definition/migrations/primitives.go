package migrations

import (
	"errors"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/uuids"
)

// Flow holds a flow definition
type Flow map[string]any

// Type returns the type of this flow
func (f Flow) Type() string {
	d, _ := f["type"].(string)
	return d
}

// Localization returns the localization of this flow
func (f Flow) Localization() Localization {
	d, _ := f["localization"].(map[string]any)
	return Localization(d)
}

type ItemTranslation map[string]any

func (it ItemTranslation) Get(property string) []string {
	vs, exists := it[property].([]any)
	if !exists {
		return nil
	}

	ss := make([]string, len(vs))
	for i := range vs {
		ss[i], _ = vs[i].(string)
	}
	return ss
}

func (it ItemTranslation) Set(property string, trans []string) {
	vs := make([]any, len(trans))
	for i := range trans {
		vs[i] = trans[i]
	}
	it[property] = vs
}

func (it ItemTranslation) Delete(property string) {
	delete(it, property)
}

type LanguageTranslation map[string]any

func (lt LanguageTranslation) getItemTranslation(uuid uuids.UUID) ItemTranslation {
	it, _ := lt[string(uuid)].(map[string]any)
	return ItemTranslation(it)
}

func (lt LanguageTranslation) GetTranslation(uuid uuids.UUID, property string) []string {
	it := lt.getItemTranslation(uuid)
	if it != nil {
		return it.Get(property)
	}
	return nil
}

func (lt LanguageTranslation) SetTranslation(uuid uuids.UUID, property string, trans []string) {
	it := lt.getItemTranslation(uuid)
	if it == nil {
		m := map[string]any{}
		lt[string(uuid)] = m
		it = ItemTranslation(m)
	}
	it.Set(property, trans)
}

func (lt LanguageTranslation) DeleteTranslation(uuid uuids.UUID, property string) {
	it := lt.getItemTranslation(uuid)
	if it != nil {
		it.Delete(property)

		// if item translation is now empty, remove it completely
		if len(it) == 0 {
			delete(lt, string(uuid))
		}
	}
}

type Localization map[string]any

func (l Localization) Languages() []i18n.Language {
	langs := make([]i18n.Language, 0, len(l))
	for k := range l {
		langs = append(langs, i18n.Language(k))
	}
	return langs
}

func (l Localization) GetLanguageTranslation(lang i18n.Language) LanguageTranslation {
	lt, _ := l[string(lang)].(map[string]any)
	return LanguageTranslation(lt)
}

// Nodes returns the nodes in this flow
func (f Flow) Nodes() []Node {
	d, _ := f["nodes"].([]any)
	nodes := make([]Node, 0, len(d))
	for _, v := range d {
		n, _ := v.(map[string]any)
		if n != nil {
			nodes = append(nodes, n)
		}
	}
	return nodes
}

// Node holds a node definition
type Node map[string]any

// Actions returns the actions on this node
func (n Node) Actions() []Action {
	d, _ := n["actions"].([]any)
	actions := make([]Action, 0, len(d))
	for _, v := range d {
		a, _ := v.(map[string]any)
		if a != nil {
			actions = append(actions, a)
		}
	}
	return actions
}

// Router returns the router on this node
func (n Node) Router() Router {
	v, _ := n["router"].(map[string]any)
	if v != nil {
		return Router(v)
	}
	return nil
}

// Action holds an action definition
type Action map[string]any

// UUID returns the uuid of this action
func (a Action) UUID() uuids.UUID {
	d, _ := a["uuid"].(string)
	return uuids.UUID(d)
}

// Type returns the type of this action
func (a Action) Type() string {
	d, _ := a["type"].(string)
	return d
}

// Router holds a router definition
type Router map[string]any

// Type returns the type of this router
func (r Router) Type() string {
	d, _ := r["type"].(string)
	return d
}

// GetObjectUUID gets the UUID property of o, if o is an object, if it has "uuid" property, and if the type of that property is a string
func GetObjectUUID(o any) uuids.UUID {
	m, ok := o.(map[string]any)
	if ok {
		v, exists := m["uuid"]
		if exists {
			s, ok := v.(string)
			if ok {
				return uuids.UUID(s)
			}
		}
	}
	return ""
}

// ReadFlow reads a flow definition as a flow primitive
func ReadFlow(data []byte) (Flow, error) {
	g, err := jsonx.DecodeGeneric(data)
	if err != nil {
		return nil, err
	}

	d, _ := g.(map[string]any)
	if d == nil {
		return nil, errors.New("flow definition isn't an object")
	}

	return d, nil
}
