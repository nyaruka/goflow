package flows

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/contactql"
	"github.com/nyaruka/goflow/utils"
)

// Contact represents a single contact
type Contact struct {
	uuid     ContactUUID
	name     string
	language utils.Language
	timezone *time.Location
	urns     URNList
	groups   *GroupList
	fields   FieldValues
	channel  Channel
}

// NewContact returns a new contact
func NewContact(uuid ContactUUID, name string, language utils.Language, timezone *time.Location) *Contact {
	return &Contact{
		uuid:     uuid,
		name:     name,
		language: language,
		timezone: timezone,
		groups:   NewGroupList([]*Group{}),
	}
}

// Clone creates a copy of this contact
func (c *Contact) Clone() *Contact {
	return &Contact{
		uuid:     c.uuid,
		name:     c.name,
		language: c.language,
		timezone: c.timezone,
		urns:     c.urns.Clone(),
		groups:   c.groups.Clone(),
		fields:   c.fields.Clone(),
		channel:  c.channel,
	}
}

// UUID returns the UUID of this contact
func (c *Contact) UUID() ContactUUID { return c.uuid }

// SetLanguage sets the language for this contact
func (c *Contact) SetLanguage(lang utils.Language) { c.language = lang }

// Language gets the language for this contact
func (c *Contact) Language() utils.Language { return c.language }

// SetTimezone sets the timezone of this contact
func (c *Contact) SetTimezone(tz *time.Location) {
	c.timezone = tz
}

// Timezone returns the timezone of this contact
func (c *Contact) Timezone() *time.Location { return c.timezone }

// SetName sets the name of this contact
func (c *Contact) SetName(name string) { c.name = name }

// Name returns the name of this contact
func (c *Contact) Name() string { return c.name }

// URNs returns the URNs of this contact
func (c *Contact) URNs() URNList { return c.urns }

// AddURN adds a new URN to this contact
func (c *Contact) AddURN(urn urns.URN) {
	// TODO normalize and check we're not adding duplicates

	c.urns = append(c.urns, urn)
}

// Groups returns the groups that this contact belongs to
func (c *Contact) Groups() *GroupList { return c.groups }

// Fields returns this contact's field values
func (c *Contact) Fields() FieldValues { return c.fields }

// Channel returns the preferred channel of this contact
func (c *Contact) Channel() Channel { return c.channel }

// SetChannel sets the preferred channel of this contact
func (c *Contact) SetChannel(channel Channel) { c.channel = channel }

// Reference returns a reference to this contact
func (c *Contact) Reference() *ContactReference { return NewContactReference(c.uuid, c.name) }

func (c *Contact) Resolve(key string) interface{} {
	switch key {
	case "name":
		return c.name
	case "first_name":
		names := utils.TokenizeString(c.name)
		if len(names) >= 1 {
			return names[0]
		}
		return ""
	case "uuid":
		return c.uuid
	case "urns":
		return c.urns
	case "language":
		return string(c.language)
	case "groups":
		return c.groups
	case "fields":
		return c.fields
	case "timezone":
		return c.timezone
	case "channel":
		return c.channel
	}

	return fmt.Errorf("no field '%s' on contact", key)
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

// UpdateDynamicGroups reevaluates membership of all dynamic groups for this contact
func (c *Contact) UpdateDynamicGroups(session Session) error {
	groups, err := session.Assets().GetGroupSet()
	if err != nil {
		return err
	}

	for _, group := range groups.All() {
		if group.IsDynamic() {
			qualifies, err := group.CheckDynamicMembership(session, c)
			if err != nil {
				return err
			}
			if qualifies {
				c.groups.Add(group)
			} else {
				c.groups.Remove(group)
			}
		}
	}

	return nil
}

// ResolveQueryKey resolves a contact query search key for this contact
func (c *Contact) ResolveQueryKey(key string) []interface{} {
	// try as a URN scheme
	if urns.IsValidScheme(key) {
		urnsWithScheme := c.urns.WithScheme(key)
		vals := make([]interface{}, len(urnsWithScheme))
		for u := range urnsWithScheme {
			vals[u] = string(urnsWithScheme[u])
		}
		return vals
	}

	// try as a contact field
	for k, value := range c.fields {
		if key == string(k) {
			return []interface{}{value.value}
		}
	}

	return nil
}

var _ contactql.Queryable = (*Contact)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type fieldValueEnvelope struct {
	Value     string    `json:"value"`
	CreatedOn time.Time `json:"created_on"`
}

type contactEnvelope struct {
	UUID     ContactUUID                     `json:"uuid" validate:"required,uuid4"`
	Name     string                          `json:"name"`
	Language utils.Language                  `json:"language"`
	Timezone string                          `json:"timezone"`
	URNs     URNList                         `json:"urns"`
	Groups   []*GroupReference               `json:"groups,omitempty" validate:"dive"`
	Fields   map[FieldKey]fieldValueEnvelope `json:"fields,omitempty"`
	Channel  *ChannelReference               `json:"channel,omitempty" validate:"omitempty,dive"`
}

// ReadContact decodes a contact from the passed in JSON
func ReadContact(session Session, data json.RawMessage) (*Contact, error) {
	var envelope contactEnvelope

	if err := utils.UnmarshalAndValidate(data, &envelope, "contact"); err != nil {
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
		c.groups = NewGroupList([]*Group{})
	} else {
		groups := make([]*Group, len(envelope.Groups))
		for g := range envelope.Groups {
			if groups[g], err = session.Assets().GetGroup(envelope.Groups[g].UUID); err != nil {
				return nil, err
			}
		}
		c.groups = NewGroupList(groups)
	}

	if envelope.Fields == nil {
		c.fields = make(FieldValues)
	} else {
		c.fields = make(FieldValues, len(envelope.Fields))
		for fieldKey, valueEnvelope := range envelope.Fields {
			field, err := session.Assets().GetField(fieldKey)
			if err != nil {
				return nil, err
			}

			value, err := field.ParseValue(session.Environment(), valueEnvelope.Value)
			if err != nil {
				return nil, err
			}

			c.fields[field.key] = NewFieldValue(field, value, valueEnvelope.CreatedOn)
		}
	}

	if envelope.Channel != nil {
		c.channel, err = session.Assets().GetChannel(envelope.Channel.UUID)
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
	if c.timezone != nil {
		ce.Timezone = c.timezone.String()
	}

	ce.Groups = make([]*GroupReference, c.groups.Count())
	for g, group := range c.groups.All() {
		ce.Groups[g] = group.Reference()
	}

	ce.Fields = make(map[FieldKey]fieldValueEnvelope, len(c.fields))
	for _, v := range c.fields {
		ce.Fields[v.field.Key()] = fieldValueEnvelope{Value: v.SerializeValue(), CreatedOn: v.createdOn}
	}

	return json.Marshal(ce)
}
