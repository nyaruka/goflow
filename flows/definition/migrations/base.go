package migrations

import (
	"encoding/json"
	"errors"
	"sort"

	"github.com/nyaruka/goflow/utils"

	"github.com/Masterminds/semver"
)

// MigrationFunc is a function that can migrate a flow definition from one version to another
type MigrationFunc func(Flow) (Flow, error)

var registered = map[*semver.Version]MigrationFunc{}

// registers a new type of action
func registerMigration(version *semver.Version, fn MigrationFunc) {
	registered[version] = fn
}

// Registered gets all registered migrations
func Registered() map[*semver.Version]MigrationFunc {
	return registered
}

// MigrateToLatest migrates the given flow definition to the latest version
func MigrateToLatest(data []byte, from *semver.Version) ([]byte, error) {
	return MigrateToVersion(data, from, nil)
}

// MigrateToVersion migrates the given flow definition to the given version
func MigrateToVersion(data []byte, from *semver.Version, to *semver.Version) ([]byte, error) {
	// get all newer versions
	versions := make([]*semver.Version, 0)
	for v := range registered {
		if v.GreaterThan(from) && (to == nil || v.Compare(to) <= 0) {
			versions = append(versions, v)
		}
	}

	// sorted by earliest first
	sort.SliceStable(versions, func(i, j int) bool { return versions[i].LessThan(versions[j]) })

	g, err := utils.JSONDecodeGeneric(data)
	if err != nil {
		return nil, err
	}

	d, _ := g.(map[string]interface{})
	if d == nil {
		return nil, errors.New("can't migrate definition which isn't a flow")
	}

	migrated := Flow(d)

	for _, version := range versions {
		migrated, err = registered[version](migrated)
		if err != nil {
			return nil, err
		}

		migrated["spec_version"] = version.String()
	}

	// finally marshal back to JSON
	return json.Marshal(migrated)
}

//------------------------------------------------------------------------------------------
// Generic definition primitives
//------------------------------------------------------------------------------------------

type Flow map[string]interface{}

func (f Flow) Nodes() []Node {
	d, _ := f["nodes"].([]interface{})
	nodes := make([]Node, len(d))
	for i := range d {
		nodes[i] = Node(d[i].(map[string]interface{}))
	}
	return nodes
}

type Node map[string]interface{}

func (n Node) Actions() []Action {
	d, _ := n["actions"].([]interface{})
	actions := make([]Action, len(d))
	for i := range d {
		actions[i] = Action(d[i].(map[string]interface{}))
	}
	return actions
}

func (n Node) Router() Router {
	d, _ := n["router"].(map[string]interface{})
	if d == nil {
		return nil
	}
	return Router(d)
}

type Action map[string]interface{}

func (a Action) Type() string {
	d, _ := a["type"].(string)
	return d
}

type Router map[string]interface{}

func (r Router) Type() string {
	d, _ := r["type"].(string)
	return d
}
