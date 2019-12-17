package migrations

import (
	"encoding/json"
	"sort"

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
			return nil, errors.Wrapf(err, "unable to migrate to version %s", version.String())
		}

		migrated["spec_version"] = version.String()
	}

	// finally marshal back to JSON
	return json.Marshal(migrated)
}
