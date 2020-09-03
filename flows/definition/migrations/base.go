package migrations

import (
	"sort"
	"strings"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows/definition/legacy"
	"github.com/nyaruka/goflow/utils"

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
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

// Header13 is the set of fields common to all 13+ flow spec versions
type Header13 struct {
	UUID        assets.FlowUUID `json:"uuid" validate:"required,uuid4"`
	Name        string          `json:"name"`
	SpecVersion *semver.Version `json:"spec_version" validate:"required"`
}

// Config configures how flow migrations are handled
type Config struct {
	BaseMediaURL string
}

// MigrateToLatest migrates the given flow definition to the latest version
func MigrateToLatest(data []byte, config *Config) ([]byte, error) {
	return MigrateToVersion(data, nil, config)
}

// MigrateToVersion migrates the given flow definition to the given version
func MigrateToVersion(data []byte, to *semver.Version, config *Config) ([]byte, error) {
	// try to read new style header (uuid, name, spec_version)
	header := &Header13{}
	err := utils.UnmarshalAndValidate(data, header)

	if err != nil {
		// could this be a legacy definition?
		if legacy.IsPossibleDefinition(data) {
			if config == nil {
				return nil, errors.New("unable to migrate what appears to be a legacy definition without a migration config")
			}

			// try to migrate it forwards to 13.0.0
			var err error
			data, err = legacy.MigrateDefinition(data, config.BaseMediaURL)
			if err != nil {
				return nil, errors.Wrap(err, "error migrating what appears to be a legacy definition")
			}
		}

		// try reading header again
		err = utils.UnmarshalAndValidate(data, header)
	}

	if err != nil {
		return nil, errors.Wrap(err, "unable to read flow header")
	}

	return migrate(data, header.SpecVersion, to)
}

func migrate(data []byte, from *semver.Version, to *semver.Version) ([]byte, error) {
	// get all newer versions than this version
	versions := make([]*semver.Version, 0)
	for v := range registered {
		if v.GreaterThan(from) && (to == nil || v.Compare(to) <= 0) {
			versions = append(versions, v)
		}
	}

	// we're already at least as new as this version of the engine
	if len(versions) == 0 {
		return data, nil
	}

	// sorted by earliest first
	sort.SliceStable(versions, func(i, j int) bool { return versions[i].LessThan(versions[j]) })

	migrated, err := readFlow(data)
	if err != nil {
		return nil, err
	}

	for _, version := range versions {
		migrated, err = registered[version](migrated)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to migrate to version %s", version.String())
		}

		migrated["spec_version"] = version.String()
	}

	// finally marshal back to JSON
	return jsonx.Marshal(migrated)
}

// Clone clones the given flow definition by replacing all UUIDs using the provided mapping and
// generating new random UUIDs if they aren't in the mapping
func Clone(data []byte, depMapping map[uuids.UUID]uuids.UUID) ([]byte, error) {
	clone, err := readFlow(data)
	if err != nil {
		return nil, err
	}

	remapUUIDs(clone, depMapping)

	// finally marshal back to JSON
	return jsonx.Marshal(clone)
}

// reads a flow definition as a flow primitive
func readFlow(data []byte) (Flow, error) {
	g, err := jsonx.DecodeGeneric(data)
	if err != nil {
		return nil, err
	}

	d, _ := g.(map[string]interface{})
	if d == nil {
		return nil, errors.New("flow definition isn't an object")
	}

	return d, nil
}

// remap all UUIDs in the flow
func remapUUIDs(data map[string]interface{}, depMapping map[uuids.UUID]uuids.UUID) {
	// copy in the dependency mappings into a master mapping of all UUIDs
	mapping := make(map[uuids.UUID]uuids.UUID)
	for k, v := range depMapping {
		mapping[k] = v
	}

	replaceUUID := func(u uuids.UUID) uuids.UUID {
		if u == uuids.UUID("") {
			return uuids.UUID("")
		}
		mapped, exists := mapping[u]
		if !exists {
			mapped = uuids.New()
			mapping[u] = mapped
		}
		return mapped
	}

	objectCallback := func(obj map[string]interface{}) {
		props := objectProperties(obj)

		for _, p := range props {
			v := obj[p]

			if p == "uuid" || strings.HasSuffix(p, "_uuid") {
				asString, isString := v.(string)
				if isString {
					obj[p] = replaceUUID(uuids.UUID(asString))
				}
			} else if uuids.IsV4(p) {
				newProperty := string(replaceUUID(uuids.UUID(p)))
				obj[newProperty] = v
				delete(obj, p)
			}
		}
	}

	arrayCallback := func(arr []interface{}) {
		for i, v := range arr {
			asString, isString := v.(string)
			if isString && uuids.IsV4(asString) {
				arr[i] = replaceUUID(uuids.UUID(asString))
			}
		}
	}

	walk(data, objectCallback, arrayCallback)
}

// extract the property names from a generic JSON object, sorted A-Z
func objectProperties(obj map[string]interface{}) []string {
	props := make([]string, 0, len(obj))
	for k := range obj {
		props = append(props, k)
	}
	sort.Strings(props)
	return props
}

// walks the given generic JSON invoking the given callbacks for each thing found
func walk(j interface{}, objectCallback func(map[string]interface{}), arrayCallback func([]interface{})) {
	switch typed := j.(type) {
	case map[string]interface{}:
		objectCallback(typed)

		for _, p := range objectProperties(typed) {
			walk(typed[p], objectCallback, arrayCallback)
		}
	case []interface{}:
		arrayCallback(typed)

		for _, v := range typed {
			walk(v, objectCallback, arrayCallback)
		}
	}
}
