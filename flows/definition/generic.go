package definition

import (
	"encoding/json"
	"strings"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// ReadFlowFromGeneric tries to read a flow in the current spec from the given generic map
func ReadFlowFromGeneric(data map[string]interface{}) (flows.Flow, error) {
	marshaled, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return ReadFlow(marshaled)
}

// MustReadFlowFromGeneric tries to read a flow from the given generic map, panics if it can't
func MustReadFlowFromGeneric(data map[string]interface{}) flows.Flow {
	f, err := ReadFlowFromGeneric(data)
	if err != nil {
		panic(err.Error())
	}
	return f
}

// remap all UUIDs in the flow
func remapUUIDs(data map[string]interface{}, depMapping map[utils.UUID]utils.UUID) {
	// copy in the dependency mappings into a master mapping of all UUIDs
	mapping := make(map[utils.UUID]utils.UUID)
	for k, v := range depMapping {
		mapping[k] = v
	}

	replaceUUID := func(u utils.UUID) utils.UUID {
		if u == utils.UUID("") {
			return utils.UUID("")
		}
		mapped, exists := mapping[u]
		if !exists {
			mapped = utils.NewUUID()
			mapping[u] = mapped
		}
		return mapped
	}

	walkObjects(data, func(obj map[string]interface{}) {
		for k, v := range obj {
			if k == "uuid" || strings.HasSuffix(k, "_uuid") {
				asString, isString := v.(string)
				if isString {
					obj[k] = replaceUUID(utils.UUID(asString))
				}
			} else if utils.IsUUIDv4(k) {
				newKey := string(replaceUUID(utils.UUID(k)))
				obj[newKey] = v
				delete(obj, k)
			}
		}
	})
}

// walks the given generic JSON invoking the given callback for each object found
func walkObjects(j interface{}, callback func(map[string]interface{})) {
	switch typed := j.(type) {
	case map[string]interface{}:
		callback(typed)

		for _, v := range typed {
			walkObjects(v, callback)
		}
	case []interface{}:
		for _, v := range typed {
			walkObjects(v, callback)
		}
	}
}
