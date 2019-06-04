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

	objectCallback := func(obj map[string]interface{}) {
		props := objectProperties(obj)

		for _, p := range props {
			v := obj[p]

			if p == "uuid" || strings.HasSuffix(p, "_uuid") {
				asString, isString := v.(string)
				if isString {
					obj[p] = replaceUUID(utils.UUID(asString))
				}
			} else if utils.IsUUIDv4(p) {
				newProperty := string(replaceUUID(utils.UUID(p)))
				obj[newProperty] = v
				delete(obj, p)
			}
		}
	}

	arrayCallback := func(arr []interface{}) {
		for i, v := range arr {
			asString, isString := v.(string)
			if isString && utils.IsUUIDv4(asString) {
				arr[i] = replaceUUID(utils.UUID(asString))
			}
		}
	}

	walk(data, objectCallback, arrayCallback)
}

// extract the property names from a generic JSON object
func objectProperties(obj map[string]interface{}) []string {
	props := make([]string, 0, len(obj))
	for k := range obj {
		props = append(props, k)
	}
	return props
}

// walks the given generic JSON invoking the given callbacks for each thing found
func walk(j interface{}, objectCallback func(map[string]interface{}), arrayCallback func([]interface{})) {
	switch typed := j.(type) {
	case map[string]interface{}:
		objectCallback(typed)

		for _, v := range typed {
			walk(v, objectCallback, arrayCallback)
		}
	case []interface{}:
		arrayCallback(typed)

		for _, v := range typed {
			walk(v, objectCallback, arrayCallback)
		}
	}
}
