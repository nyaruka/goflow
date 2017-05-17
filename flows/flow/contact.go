package flow

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/flows"
)

type contact struct {
	uuid     flows.ContactUUID
	name     string
	language flows.Language
	urns     flows.URNList
	groups   []*group
	fields   fields
}

func (c *contact) SetLanguage(lang flows.Language) { c.language = lang }
func (c *contact) Language() flows.Language        { return c.language }

func (c *contact) SetName(name string) { c.name = name }
func (c *contact) Name() string        { return c.name }

func (c *contact) URNs() flows.URNList     { return c.urns }
func (c *contact) UUID() flows.ContactUUID { return c.uuid }

func (c *contact) Groups() flows.GroupList { return groupList(c.groups) }
func (c *contact) AddGroup(uuid flows.GroupUUID, name string) {
	c.groups = append(c.groups, &group{uuid, name})
}
func (c *contact) RemoveGroup(uuid flows.GroupUUID) bool {
	for i := range c.groups {
		if c.groups[i].uuid == uuid {
			c.groups = append(c.groups[:i], c.groups[i+1])
			return true
		}
	}
	return false
}

func (c *contact) Fields() flows.Fields { return fields(c.fields) }

func (c *contact) Resolve(key string) interface{} {
	switch key {

	case "name":
		return c.name

	case "uuid":
		return c.uuid

	case "urns":
		return c.urns

	case "language":
		return c.language

	case "groups":
		return groupList(c.groups)

	case "fields":
		return c.fields

	case "urn":
		return c.urns
	}

	return fmt.Errorf("No field '%s' on contact", key)
}

func (c *contact) Default() interface{} {
	return c
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

// ReadContact decodes a contact from the passed in JSON
func ReadContact(data json.RawMessage) (flows.Contact, error) {
	contact := &contact{}
	err := json.Unmarshal(data, contact)
	if err == nil {
		// err = run.Validate()
	}
	return contact, err
}

type contactEnvelope struct {
	UUID     flows.ContactUUID `json:"uuid"`
	Name     string            `json:"name"`
	Language flows.Language    `json:"language"`
	URNs     flows.URNList     `json:"urns"`
	Groups   groupList         `json:"groups"`
	Fields   fields            `json:"fields"`
}

func (c *contact) UnmarshalJSON(data []byte) error {
	var ce contactEnvelope
	var err error

	err = json.Unmarshal(data, &ce)
	if err != nil {
		return err
	}

	c.name = ce.Name
	c.uuid = ce.UUID
	c.language = ce.Language

	if ce.URNs == nil {
		c.urns = make(flows.URNList, 0)
	} else {
		c.urns = ce.URNs
	}

	if ce.Groups == nil {
		c.groups = make(groupList, 0)
	} else {
		c.groups = ce.Groups
	}

	if ce.Fields == nil {
		c.fields = newFields()
	} else {
		c.fields = ce.Fields
	}

	return err
}

func (c *contact) MarshalJSON() ([]byte, error) {
	var ce contactEnvelope

	ce.Name = c.name
	ce.UUID = c.uuid
	ce.Language = c.language
	ce.URNs = c.urns
	ce.Groups = c.groups
	ce.Fields = c.fields

	return json.Marshal(ce)
}
