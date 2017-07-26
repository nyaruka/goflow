package flows

import (
	"encoding/json"
	"fmt"

	"time"

	"github.com/nyaruka/goflow/utils"
)

// Contact represents a single contact
type Contact struct {
	uuid     ContactUUID
	name     string
	language utils.Language
	timezone *time.Location
	urns     URNList
	groups   GroupList
	fields   *Fields
	channel  Channel
}

// SetLanguage sets the language for this contact
func (c *Contact) SetLanguage(lang utils.Language) { c.language = lang }

// Language gets the language for this contact
func (c *Contact) Language() utils.Language { return c.language }

func (c *Contact) SetTimezone(tz *time.Location) {
	c.timezone = tz
}
func (c *Contact) Timezone() *time.Location { return c.timezone }

func (c *Contact) SetName(name string) { c.name = name }
func (c *Contact) Name() string        { return c.name }

func (c *Contact) URNs() URNList     { return c.urns }
func (c *Contact) UUID() ContactUUID { return c.uuid }

func (c *Contact) Groups() GroupList { return GroupList(c.groups) }
func (c *Contact) AddGroup(uuid GroupUUID, name string) {
	c.groups = append(c.groups, &Group{uuid, name})
}
func (c *Contact) RemoveGroup(uuid GroupUUID) bool {
	for i := range c.groups {
		if c.groups[i].uuid == uuid {
			c.groups = append(c.groups[:i], c.groups[i+1:]...)
			return true
		}
	}
	return false
}

func (c *Contact) InGroup(group *Group) bool {
	for i := range c.groups {
		if c.groups[i].uuid == group.UUID() {
			return true
		}
	}
	return false
}

func (c *Contact) Fields() *Fields            { return c.fields }
func (c *Contact) Channel() Channel           { return c.channel }
func (c *Contact) SetChannel(channel Channel) { c.channel = channel }

func (c *Contact) Resolve(key string) interface{} {
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
		return GroupList(c.groups)

	case "fields":
		return c.fields

	case "timezone":
		return c.timezone

	case "urn":
		return c.urns

	case "channel":
		return c.channel
	}

	return fmt.Errorf("No field '%s' on contact", key)
}

// Default returns our default value in the context
func (c *Contact) Default() interface{} {
	return c
}

// String returns our string value in the context
func (c *Contact) String() string {
	return c.name
}

var _ utils.VariableResolver = (*Contact)(nil)

type ContactReference struct {
	UUID ContactUUID `json:"uuid"    validate:"required,uuid4"`
	Name string      `json:"name"`
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type contactEnvelope struct {
	UUID        ContactUUID    `json:"uuid" validate:"uuid4"`
	Name        string         `json:"name"`
	Language    utils.Language `json:"language"`
	Timezone    string         `json:"timezone"`
	URNs        URNList        `json:"urns"`
	Groups      GroupList      `json:"groups"`
	Fields      *Fields        `json:"fields,omitempty"`
	ChannelUUID ChannelUUID    `json:"channel_uuid" validate:"omitempty,uuid4"`
}

// ReadContact decodes a contact from the passed in JSON
func ReadContact(assets Assets, data json.RawMessage) (*Contact, error) {
	var envelope contactEnvelope
	var err error

	err = json.Unmarshal(data, &envelope)
	if err != nil {
		return nil, err
	}
	err = utils.ValidateAll(envelope)
	if err != nil {
		return nil, err
	}

	c := &Contact{}
	c.uuid = envelope.UUID
	c.name = envelope.Name
	c.language = envelope.Language

	tz, err := time.LoadLocation(envelope.Timezone)
	if err != nil {
		return nil, err
	}
	c.timezone = tz

	if envelope.URNs == nil {
		c.urns = make(URNList, 0)
	} else {
		c.urns = envelope.URNs
	}

	if envelope.Groups == nil {
		c.groups = make(GroupList, 0)
	} else {
		c.groups = envelope.Groups
	}

	if envelope.Fields == nil {
		c.fields = NewFields()
	} else {
		c.fields = envelope.Fields
	}

	if envelope.ChannelUUID != "" {
		c.channel, err = assets.GetChannel(envelope.ChannelUUID)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Contact) MarshalJSON() ([]byte, error) {
	var ce contactEnvelope

	ce.Name = c.name
	ce.UUID = c.uuid
	ce.Language = c.language
	ce.URNs = c.urns
	ce.Groups = c.groups
	ce.Fields = c.fields
	if c.timezone != nil {
		ce.Timezone = c.timezone.String()
	}

	return json.Marshal(ce)
}
