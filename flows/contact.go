package flows

import (
	"encoding/json"
	"time"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/contactql"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
)

// Contact represents a person who is interacting with the flow. It renders as the person's name
// (or perferred URN if name isn't set) in a template, and has the following properties which can be accessed:
//
//  * `uuid` the UUID of the contact
//  * `name` the full name of the contact
//  * `first_name` the first name of the contact
//  * `language` the [ISO-639-3](http://www-01.sil.org/iso639-3/) language code of the contact
//  * `urns` all [URNs](#context:urn) the contact has set
//  * `urns.[scheme]` all the [URNs](#context:urn) the contact has set for the particular URN scheme
//  * `urn` shorthand for `@(format_urn(c.urns.0))`, i.e. the contact's preferred [URN](#context:urn) in friendly formatting
//  * `groups` all the [groups](#context:group) that the contact belongs to
//  * `fields` all the custom contact fields the contact has set
//  * `fields.[snaked_field_name]` the value of the specific field
//  * `channel` shorthand for `contact.urns.0.channel`, i.e. the [channel](#context:channel) of the contact's preferred URN
//
// Examples:
//
//   @contact -> Ryan Lewis
//   @contact.name -> Ryan Lewis
//   @contact.first_name -> Ryan
//   @contact.language -> eng
//   @contact.urns -> ["tel:+12065551212","twitterid:54784326227#nyaruka","mailto:foo@bar.com"]
//   @contact.urns.0 -> tel:+12065551212
//   @contact.urns.tel -> ["tel:+12065551212"]
//   @contact.urns.mailto.0 -> mailto:foo@bar.com
//   @contact.urn -> (206) 555-1212
//   @contact.groups -> ["Testers","Males"]
//   @contact.fields -> {"activation_token":"AACC55","gender":"Male"}
//   @contact.fields.activation_token -> AACC55
//   @contact.fields.gender -> Male
//
// @context contact
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
		urns:     c.urns.clone(),
		groups:   c.groups.clone(),
		fields:   c.fields.clone(),
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
func (c *Contact) AddURN(urn urns.URN) bool {
	if c.HasURN(urn) {
		return false
	}
	c.urns = append(c.urns, &ContactURN{URN: urn.Normalize("")})
	return true
}

// HasURN checks whether the contact has the given URN
func (c *Contact) HasURN(urn urns.URN) bool {
	urn = urn.Normalize("")

	for _, u := range c.urns {
		if u.URN == urn {
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

// Resolve resolves the given key when this contact is referenced in an expression
func (c *Contact) Resolve(key string) types.XValue {
	switch key {
	case "uuid":
		return types.NewXString(string(c.uuid))
	case "name":
		return types.NewXString(c.name)
	case "first_name":
		names := utils.TokenizeString(c.name)
		if len(names) >= 1 {
			return types.NewXString(names[0])
		}
		return nil
	case "language":
		return types.NewXString(string(c.language))
	case "timezone":
		if c.timezone != nil {
			return types.NewXString(c.timezone.String())
		}
		return nil
	case "urns":
		return c.urns
	case "urn":
		if len(c.urns) > 0 {
			return types.NewXString(c.urns[0].Format())
		}
		return nil
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

	return types.NewXResolveError(c, key)
}

// Reduce is called when this object needs to be reduced to a primitive
func (c *Contact) Reduce() types.XPrimitive {
	return types.NewXString(c.name)
}

func (c *Contact) ToXJSON() types.XString { return types.NewXString("TODO") }

var _ types.XValue = (*Contact)(nil)
var _ types.XResolvable = (*Contact)(nil)

// SetFieldValue updates the given contact field value for this contact
func (c *Contact) SetFieldValue(env utils.Environment, field *Field, rawValue string) {
	c.fields.setValue(env, field, types.NewXString(rawValue))
}

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
			fieldValue := value.TypedValue()
			var nativeValue interface{}

			switch typed := fieldValue.(type) {
			case nil:
				return nil
			case *Location:
				nativeValue = typed.Name()
			case types.XString:
				nativeValue = typed.Native()
			case types.XNumber:
				nativeValue = typed.Native()
			case types.XDate:
				nativeValue = typed.Native()
			}

			return []interface{}{nativeValue}
		}
	}

	return nil
}

var _ contactql.Queryable = (*Contact)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type fieldValueEnvelope struct {
	Text     types.XString  `json:"text,omitempty"`
	Datetime *types.XDate   `json:"datetime,omitempty"`
	Decimal  *types.XNumber `json:"decimal,omitempty"`
	State    string         `json:"state,omitempty"`
	District string         `json:"district,omitempty"`
	Ward     string         `json:"ward,omitempty"`
}

type contactEnvelope struct {
	UUID     ContactUUID                      `json:"uuid" validate:"required,uuid4"`
	Name     string                           `json:"name"`
	Language utils.Language                   `json:"language"`
	Timezone string                           `json:"timezone"`
	URNs     []urns.URN                       `json:"urns" validate:"dive,urn"`
	Groups   []*GroupReference                `json:"groups,omitempty" validate:"dive"`
	Fields   map[FieldKey]*fieldValueEnvelope `json:"fields,omitempty"`
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

	fieldSet, err := session.Assets().GetFieldSet()
	if err != nil {
		return nil, err
	}

	c.fields = make(FieldValues, len(fieldSet.All()))

	for _, field := range fieldSet.All() {
		value := &FieldValue{field: field}

		if envelope.Fields != nil {
			valueEnvelope := envelope.Fields[field.key]
			if valueEnvelope != nil {
				value.text = valueEnvelope.Text
				value.decimal = valueEnvelope.Decimal
				value.datetime = valueEnvelope.Datetime

				// TODO parse locations
			}
		}

		c.fields[field.key] = value
	}

	return c, nil
}

// MarshalJSON marshals this contact into JSON
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

	ce.Fields = make(map[FieldKey]*fieldValueEnvelope)
	for _, v := range c.fields {
		if !v.IsEmpty() {
			ce.Fields[v.field.Key()] = &fieldValueEnvelope{
				Text:     v.text,
				Decimal:  v.decimal,
				Datetime: v.datetime,
				//State:    v.state,
				//District: v.district,
				//Ward:     v.ward,
			}
		}
	}

	return json.Marshal(ce)
}
