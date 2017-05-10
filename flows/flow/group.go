package flow

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/flows"
)

type group struct {
	uuid flows.GroupUUID
	name string
}

type groupList []*group

func (g *group) Resolve(key string) interface{} {
	switch key {

	case "name":
		return g.name

	case "uuid":
		return g.uuid
	}

	return fmt.Errorf("No field '%s' on group", key)
}

func (g *group) Default() interface{} {
	return g.name
}

func (g *group) Name() string          { return g.name }
func (g *group) UUID() flows.GroupUUID { return g.uuid }

func (l groupList) FindGroup(uuid flows.GroupUUID) flows.Group {
	for i := range l {
		if l[i].uuid == uuid {
			return l[i]
		}
	}
	return nil
}

func (l groupList) Resolve(key string) interface{} {
	// TODO: decide what to do here for @contact.groups.[]
	// Do we want to allow any kind of filtering?
	return l
}

func (l groupList) Default() interface{} {
	return l
}

func (l groupList) String() string {
	names := make([]string, len(l))
	for i := range l {
		names[i] = l[i].name
	}
	return strings.Join(names, ", ")
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type groupEnvelope struct {
	UUID flows.GroupUUID `json:"uuid"`
	Name string          `json:"name"`
}

func (g *group) UnmarshalJSON(data []byte) error {
	var ge groupEnvelope
	var err error

	err = json.Unmarshal(data, &ge)
	g.uuid = ge.UUID
	g.name = ge.Name

	return err
}

func (g *group) MarshalJSON() ([]byte, error) {
	var ge groupEnvelope

	ge.Name = g.name
	ge.UUID = g.uuid

	return json.Marshal(ge)
}
