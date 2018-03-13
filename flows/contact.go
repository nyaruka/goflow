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
	if c == nil {
		return nil
	}

	return &Contact{
		uuid:     c.uuid,
		name:     c.name,
		language: c.language,
		timezone: c.timezone,
		urns:     c.urns.Clone(),
		groups:   c.groups.Clone(),
		fields:   c.fields.Clone(),
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
func (c *Contact) AddURN(urn urns.URN) error {
	if c.HasURN(urn) {
		return fmt.Errorf("contact already has the URN %s", urn)
	}
	c.urns = append(c.urns, &ContactURN{URN: urn.Normalize("")})
	return nil
}

// HasURN checks whether the contact has the given URN
func (c *Contact) HasURN(urn urns.URN) bool {
	urn = urn.Normalize("")

	for u := range c.urns {
		if c.urns[u].URN == urn {
			return true
		}
	}
	return false
}

// Groups returns the groups that this contact belongs to
func (c *Contact) Groups() *GroupList { return c.groups }

// Fields returns this contact's field values
func (c *Contact) Fields() FieldValues { return c.fields }

// Reference returns a reference to this contact
func (c *Contact) Reference() *ContactReference { return NewContactReference(c.uuid, c.name) }

func (c *Contact) Resolve(key string) interface{} {
	switch key {
	case "uuid":
		return c.uuid
	case "name":
		return c.name
	case "first_name":
		names := utils.TokenizeString(c.name)
		if len(names) >= 1 {
			return names[0]
		}
		return ""
	case "language":
		return string(c.language)
	case "timezone":
		return c.timezone
	case "urns":
		return c.urns
	case "groups":
		return c.groups
	case "fields":
		return c.fields
	case "channel":
		if len(c.urns) > 0 {
			return c.urns[0].Channel()
		}
		return nil
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

// UpdatePreferredChannel updates the preferred channel
func (c *Contact) UpdatePreferredChannel(channel Channel) {
	priorityURNs := make([]*ContactURN, 0)
	otherURNs := make([]*ContactURN, 0)

	for _, urn := range c.urns {
		// tel URNs can be re-assigned, other URN schemes are considered channel specific
		if urn.URN.Scheme() == urns.TelScheme && channel.SupportsScheme(urns.TelScheme) {
			urn.SetChannel(channel)
		}

		// move any URNs with this channel to the front of the list
		if urn.Channel() == channel {
			priorityURNs = append(priorityURNs, urn)
		} else {
			otherURNs = append(otherURNs, urn)
		}
	}

	c.urns = append(priorityURNs, otherURNs...)
}

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
			vals[u] = string(urnsWithScheme[u].URN)
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
	URNs     []urns.URN                      `json:"urns" validate:"dive,urn"`
	Groups   []*GroupReference               `json:"groups,omitempty" validate:"dive"`
	Fields   map[FieldKey]fieldValueEnvelope `json:"fields,omitempty"`
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
		if c.urns, err = ReadURNList(session, envelope.URNs); err != nil {
			return nil, err
		}
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

	return c, nil
}

func (c *Contact) MarshalJSON() ([]byte, error) {
	var ce contactEnvelope

	ce.Name = c.name
	ce.UUID = c.uuid
	ce.Language = c.language
	ce.URNs = c.urns.RawURNs(true)
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
