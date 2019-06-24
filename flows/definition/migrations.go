package definition

import (
	"encoding/json"

	"github.com/nyaruka/goflow/utils"

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
)

type migratorFunc func(data []byte) ([]byte, error)

func migrateAsGeneric(fn func(map[string]interface{}) map[string]interface{}) migratorFunc {
	return func(data []byte) ([]byte, error) {
		g, err := utils.JSONDecodeGeneric(data)
		if err != nil {
			return nil, err
		}

		definition, isMap := g.(map[string]interface{})
		if !isMap {
			return nil, errors.New("can't migrate definition which isn't a flow")
		}

		definition = fn(definition)

		return json.Marshal(definition)
	}
}

var migrations = []struct {
	version  *semver.Version
	migrator migratorFunc
}{
	{semver.MustParse("13.1"), migrate13_1},
}

func migrateDefinition(data []byte, fromVersion *semver.Version) ([]byte, error) {
	migrated := data
	var err error

	for _, m := range migrations {
		if m.version.GreaterThan(fromVersion) {
			migrated, err = m.migrator(migrated)
			if err != nil {
				return nil, err
			}
		}
	}

	return migrated, nil
}

func migrate13_1(data []byte) ([]byte, error) {
	return data, nil
}
