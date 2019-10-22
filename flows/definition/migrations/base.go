package migrations

import (
	"encoding/json"
	"errors"
	"sort"

	"github.com/Masterminds/semver"
	"github.com/nyaruka/goflow/utils"
)

type MigrationFunc func(data []byte) ([]byte, error)

var registered = map[*semver.Version]MigrationFunc{}

// registers a new type of action
func registerMigration(version *semver.Version, fn MigrationFunc) {
	registered[version] = fn
}

// gets all migration functions after the given version
func fromVersion(from *semver.Version) []MigrationFunc {
	// get all newer versions
	versions := make([]*semver.Version, 0)
	for v := range registered {
		if v.GreaterThan(from) {
			versions = append(versions, v)
		}
	}

	// sorted by earliest first
	sort.SliceStable(versions, func(i, j int) bool { return versions[i].LessThan(versions[j]) })

	// get the migrations
	migrations := make([]MigrationFunc, len(versions))
	for i, version := range versions {
		migrations[i] = registered[version]
	}
	return migrations
}

// MigrateDefinition migrates the given flow definition to the latest version
func MigrateDefinition(data []byte, from *semver.Version) ([]byte, error) {
	migrations := fromVersion(from)

	migrated := data
	var err error

	for _, migration := range migrations {
		migrated, err = migration(migrated)
		if err != nil {
			return nil, err
		}
	}

	return migrated, nil
}

type Flow struct {
	Def map[string]interface{}
}

func (f Flow) Nodes() []Node {
	d, _ := f.Def["nodes"].([]interface{})
	nodes := make([]Node, len(d))
	for i := range d {
		nodes[i] = Node{d[i].(map[string]interface{})}
	}
	return nodes
}

type Node struct {
	Def map[string]interface{}
}

func (n Node) Actions() []Action {
	d, _ := n.Def["actions"].([]interface{})
	actions := make([]Action, len(d))
	for i := range d {
		actions[i] = Action{d[i].(map[string]interface{})}
	}
	return actions
}

func (n Node) Router() *Router {
	d, _ := n.Def["router"].(map[string]interface{})
	if d == nil {
		return nil
	}
	return &Router{d}
}

type Action struct {
	Def map[string]interface{}
}

func (a Action) Type() string {
	d, _ := a.Def["type"].(string)
	return d
}

type Router struct {
	Def map[string]interface{}
}

func (r Router) Type() string {
	d, _ := r.Def["type"].(string)
	return d
}

// AsParsed is a wrapper for a migration func which migrates a flow as a generic map
func AsParsed(fn func(Flow) Flow) MigrationFunc {
	return func(data []byte) ([]byte, error) {
		g, err := utils.JSONDecodeGeneric(data)
		if err != nil {
			return nil, err
		}

		definition, isMap := g.(map[string]interface{})
		if !isMap {
			return nil, errors.New("can't migrate definition which isn't a flow")
		}

		migrated := fn(Flow{definition})

		return json.Marshal(migrated.Def)
	}
}
