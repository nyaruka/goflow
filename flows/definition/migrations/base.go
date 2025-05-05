package migrations

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows/definition/legacy"
	"github.com/nyaruka/goflow/utils"
)

// MigrationFunc is a function that can migrate a flow definition from one version to another
type MigrationFunc func(Flow, *Config) (Flow, error)

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
	UUID        assets.FlowUUID `json:"uuid"         validate:"required,uuid"`
	Name        string          `json:"name"         validate:"required,max=64"`
	SpecVersion *semver.Version `json:"spec_version" validate:"required"`
}

// Config configures how flow migrations are handled
type Config struct {
	BaseMediaURL string
}

var DefaultConfig = &Config{}

// MigrateToLatest migrates the given flow definition to the latest version
func MigrateToLatest(data []byte, cfg *Config) ([]byte, error) {
	return MigrateToVersion(data, nil, cfg)
}

// MigrateToVersion migrates the given flow definition to the given version
func MigrateToVersion(data []byte, to *semver.Version, cfg *Config) ([]byte, error) {
	// try to read new style header (uuid, name, spec_version)
	header := &Header13{}
	err := utils.UnmarshalAndValidate(data, header)

	if err != nil {
		// could this be a legacy definition?
		if legacy.IsPossibleDefinition(data) {
			// try to migrate it forwards to 13.0.0
			var err error
			data, err = legacy.MigrateDefinition(data, cfg.BaseMediaURL)
			if err != nil {
				return nil, fmt.Errorf("error migrating what appears to be a legacy definition: %w", err)
			}
		}

		// try reading header again
		err = utils.UnmarshalAndValidate(data, header)
	}

	if err != nil {
		return nil, fmt.Errorf("unable to read flow header: %w", err)
	}

	return migrate(data, header.SpecVersion, to, cfg)
}

func migrate(data []byte, from *semver.Version, to *semver.Version, cfg *Config) ([]byte, error) {
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

	for _, version := range versions {
		// we read the flow each time to ensure what we pass to the migration function uses the types it expects
		flow, err := ReadFlow(data)
		if err != nil {
			return nil, err
		}

		flow, err = registered[version](flow, cfg)
		if err != nil {
			return nil, fmt.Errorf("unable to migrate to version %s: %w", version.String(), err)
		}

		flow["spec_version"] = version.String()

		data = jsonx.MustMarshal(flow)
	}

	return data, nil
}

// Clone clones the given flow definition by replacing all UUIDs using the provided mapping and
// generating new random UUIDs if they aren't in the mapping
func Clone(data []byte, depMapping map[uuids.UUID]uuids.UUID) ([]byte, error) {
	clone, err := ReadFlow(data)
	if err != nil {
		return nil, err
	}

	remapUUIDs(clone, depMapping)

	// finally marshal back to JSON
	return jsonx.Marshal(clone)
}

// remap all UUIDs in the flow
func remapUUIDs(data map[string]any, depMapping map[uuids.UUID]uuids.UUID) {
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
			mapped = uuids.NewV4()
			mapping[u] = mapped
		}
		return mapped
	}

	objectCallback := func(path string, obj map[string]any) {
		props := objectProperties(obj)

		for _, p := range props {
			v := obj[p]

			if p == "uuid" || strings.HasSuffix(p, "_uuid") {
				asString, isString := v.(string)
				if isString {
					obj[p] = replaceUUID(uuids.UUID(asString))
				}
			} else if uuids.Is(p) {
				newProperty := string(replaceUUID(uuids.UUID(p)))
				obj[newProperty] = v
				delete(obj, p)
			}
		}
	}

	arrayCallback := func(path string, arr []any) {
		for i, v := range arr {
			asString, isString := v.(string)
			if isString && uuids.Is(asString) {
				arr[i] = replaceUUID(uuids.UUID(asString))
			}
		}
	}

	walk(data, objectCallback, arrayCallback, "")
}

// extract the property names from a generic JSON object, sorted A-Z
func objectProperties(obj map[string]any) []string {
	props := make([]string, 0, len(obj))
	for k := range obj {
		props = append(props, k)
	}
	sort.Strings(props)
	return props
}

// walks the given generic JSON invoking the given callbacks for each thing found
func walk(j any, objectCallback func(string, map[string]any), arrayCallback func(string, []any), path string) {
	switch typed := j.(type) {
	case map[string]any:
		objectCallback(path, typed)

		for _, p := range objectProperties(typed) {
			walk(typed[p], objectCallback, arrayCallback, fmt.Sprintf("%s.%s", path, p))
		}
	case []any:
		arrayCallback(path, typed)

		for i, v := range typed {
			walk(v, objectCallback, arrayCallback, fmt.Sprintf("%s[%d]", path, i))
		}
	}
}
