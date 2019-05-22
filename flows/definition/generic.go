package definition

import (
	"encoding/json"
	"strings"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// GenericFlow is a flow stored as a generic hierarchy of maps and slices
type GenericFlow struct {
	data utils.GenericJSON
}

// convert a flow to a generic flow
func newGenericFlow(flow flows.Flow) GenericFlow {
	marshaled, err := json.Marshal(flow)
	if err != nil {
		panic(err.Error())
	}

	g, err := utils.ReadGenericJSON(marshaled)
	if err != nil {
		panic(err.Error())
	}

	return GenericFlow{g}
}

// Read tries to read this generic flow as a flow in the current spec
func (g GenericFlow) Read() (flows.Flow, error) {
	marshaled, err := json.Marshal(g.data)
	if err != nil {
		return nil, err
	}

	return ReadFlow(marshaled)
}

// MustRead tries to read this generic flow, panicking if there's a problem
func (g GenericFlow) MustRead() flows.Flow {
	f, err := g.Read()
	if err != nil {
		panic(err.Error())
	}
	return f
}

// remap all UUIDs in the flow
func (g GenericFlow) remap(mapping map[utils.UUID]utils.UUID) {
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

	g.data.WalkObjects(func(obj map[string]interface{}) {
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
