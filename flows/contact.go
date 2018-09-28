package flows

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
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
//  * `timezone` the timezone name of the contact
//  * `created_on` the datetime when the contact was created
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
//   @contact.timezone -> America/Guayaquil
//   @contact.created_on -> 2018-06-20T11:40:30.123456Z
//   @contact.urns -> ["tel:+12065551212?channel=57f1078f-88aa-46f4-a59a-948a5739c03d","twitterid:54784326227#nyaruka","mailto:foo@bar.com"]
//   @contact.urns.0 -> tel:+12065551212?channel=57f1078f-88aa-46f4-a59a-948a5739c03d
//   @contact.urns.tel -> ["tel:+12065551212?channel=57f1078f-88aa-46f4-a59a-948a5739c03d"]
//   @contact.urns.mailto.0 -> mailto:foo@bar.com
//   @contact.urn -> (206) 555-1212
//   @contact.groups -> ["Testers","Males"]
//   @contact.fields -> {"activation_token":"AACC55","age":23,"gender":"Male","join_date":"2017-12-02T00:00:00-02:00"}
//   @contact.fields.activation_token -> AACC55
//   @contact.fields.gender -> Male
//
// @context contact
type Contact struct {
	uuid      ContactUUID
	id        ContactID
	name      string
	language  utils.Language
	timezone  *time.Location
	createdOn time.Time
	urns      URNList
	groups    *GroupList
	fields    FieldValues
}

// NewContact creates a new contact with the passed in attributes
func NewContact(
	uuid ContactUUID, id ContactID, name string, language utils.Language, timezone *time.Location, createdOn time.Time,
	urns URNList, groups *GroupList, fields FieldValues) *Contact {
	return &Contact{
		uuid:      uuid,
		id:        id,
		name:      name,
		language:  language,
		timezone:  timezone,
		createdOn: createdOn,
		urns:      urns,
		groups:    groups,
		fields:    fields,
	}
}

// NewContactFromAssets creates a new contact using assets
func NewContactFromAssets(
	a SessionAssets,
	uuid ContactUUID,
	id ContactID,
	name string,
	language utils.Language,
	timezone *time.Location,
	createdOn time.Time,
	urns []urns.URN,
	groups []assets.Group,
	fields map[string]*Value) (*Contact, error) {

	urnList, err := ReadURNList(a, urns)
	if err != nil {
		return nil, err
	}

	groupList, err := NewGroupListFromAssets(a, groups)
	if err != nil {
		return nil, err
	}

	fieldValues, err := NewFieldValues(a, fields, true)
	if err != nil {
		return nil, err
	}

	return &Contact{
		uuid:      uuid,
		id:        id,
		name:      name,
		language:  language,
		timezone:  timezone,
		createdOn: createdOn,
		urns:      urnList,
		groups:    groupList,
		fields:    fieldValues,
	}, nil
}

// NewEmptyContact creates a new empy contact with the passed in name, language and location
func NewEmptyContact(name string, language utils.Language, timezone *time.Location) *Contact {
	return &Contact{
		uuid:      ContactUUID(utils.NewUUID()),
		name:      name,
		language:  language,
		timezone:  timezone,
		createdOn: utils.Now(),
		urns:      URNList{},
		groups:    NewGroupList([]*Group{}),
		fields:    make(FieldValues),
	}
}

// Clone creates a copy of this contact
func (c *Contact) Clone() *Contact {
	if c == nil {
		return nil
	}

	return &Contact{
		uuid:      c.uuid,
		id:        c.id,
		name:      c.name,
		language:  c.language,
		timezone:  c.timezone,
		createdOn: c.createdOn,
		urns:      c.urns.clone(),
		groups:    c.groups.clone(),
		fields:    c.fields.clone(),
	}
}

// UUID returns the UUID of this contact
func (c *Contact) UUID() ContactUUID { return c.uuid }

// ID returns the numeric ID of this contact
func (c *Contact) ID() ContactID { return c.id }

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

// SetCreatedOn sets the created on time of this contact
func (c *Contact) SetCreatedOn(createdOn time.Time) {
	c.createdOn = createdOn
}

// CreatedOn returns the created on time of this contact
func (c *Contact) CreatedOn() time.Time { return c.createdOn }

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
		if u.URN.Identity() == urn.Identity() {
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

// Format returns a friendly string version of this contact depending on what fields are set
func (c *Contact) Format(env utils.Environment) string {
	// if contact has a name set, use that
	if c.name != "" {
		return c.name
	}

	// otherwise use either id or the higest priority URN depending on the env
	if env.RedactionPolicy() == utils.RedactionPolicyURNs {
		return strconv.Itoa(int(c.id))
	}
	if len(c.urns) > 0 {
		return c.urns[0].Format()
	}

	return ""
}

// Resolve resolves the given key when this contact is referenced in an expression
func (c *Contact) Resolve(env utils.Environment, key string) types.XValue {
	switch key {
	case "uuid":
		return types.NewXText(string(c.uuid))
	case "id":
		return types.NewXNumberFromInt(int(c.id))
	case "name":
		return types.NewXText(c.name)
	case "first_name":
		names := utils.TokenizeString(c.name)
		if len(names) >= 1 {
			return types.NewXText(names[0])
		}
		return nil
	case "language":
		return types.NewXText(string(c.language))
	case "timezone":
		if c.timezone != nil {
			return types.NewXText(c.timezone.String())
		}
		return nil
	case "created_on":
		return types.NewXDateTime(c.createdOn)
	case "urns":
		return c.urns
	case "urn":
		if len(c.urns) > 0 {
			return types.NewXText(c.urns[0].Format())
		}
		return nil
	case "groups":
		return c.groups
	case "fields":
		return c.fields
	case "channel":
		return c.PreferredChannel()
	}

	return types.NewXResolveError(c, key)
}

// Describe returns a representation of this type for error messages
func (c *Contact) Describe() string { return "contact" }

// Reduce is called when this object needs to be reduced to a primitive
func (c *Contact) Reduce(env utils.Environment) types.XPrimitive {
	return types.NewXText(c.Format(env))
}

// ToXJSON is called when this type is passed to @(json(...))
func (c *Contact) ToXJSON(env utils.Environment) types.XText {
	return types.ResolveKeys(env, c, "uuid", "name", "language", "timezone", "created_on", "urns", "groups", "fields", "channel").ToXJSON(env)
}

var _ types.XValue = (*Contact)(nil)
var _ types.XResolvable = (*Contact)(nil)

// SetFieldValue updates the given contact field value for this contact
func (c *Contact) SetFieldValue(env utils.Environment, fields *FieldAssets, key string, rawValue string) (*Value, error) {
	runEnv := env.(RunEnvironment)

	// lookup the actual field object for this key
	field, err := fields.Get(key)
	if err != nil {
		return nil, err
	}

	return c.fields.setValue(runEnv, field, rawValue, fields), nil
}

// PreferredChannel gets the preferred channel for this contact, i.e. the preferred channel of their highest priority URN
func (c *Contact) PreferredChannel() *Channel {
	if len(c.urns) > 0 {
		return c.urns[0].Channel()
	}
	return nil
}

// UpdatePreferredChannel updates the preferred channel
func (c *Contact) UpdatePreferredChannel(channel *Channel) {
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

// ReevaluateDynamicGroups reevaluates membership of all dynamic groups for this contact
func (c *Contact) ReevaluateDynamicGroups(session Session) ([]*Group, []*Group, []error) {
	added := make([]*Group, 0)
	removed := make([]*Group, 0)
	errors := make([]error, 0)

	for _, group := range session.Assets().Groups().All() {
		if !group.IsDynamic() {
			continue
		}

		qualifies, err := group.CheckDynamicMembership(session.Environment(), c)
		if err != nil {
			errors = append(errors, err)
		} else if qualifies {
			if c.groups.Add(group) {
				added = append(added, group)
			}
		} else {
			if c.groups.Remove(group) {
				removed = append(removed, group)
			}
		}
	}

	return added, removed, errors
}

// ResolveQueryKey resolves a contact query search key for this contact
func (c *Contact) ResolveQueryKey(env utils.Environment, key string) []interface{} {
	if key == "language" {
		if c.language != utils.NilLanguage {
			return []interface{}{string(c.language)}
		}
		return nil
	} else if key == "created_on" {
		return []interface{}{c.createdOn}
	}

	// try as a URN scheme
	if urns.IsValidScheme(key) {
		if env.RedactionPolicy() != utils.RedactionPolicyURNs {
			urnsWithScheme := c.urns.WithScheme(key)
			vals := make([]interface{}, len(urnsWithScheme))
			for u := range urnsWithScheme {
				vals[u] = string(urnsWithScheme[u].URN)
			}
			return vals
		}
		return nil
	}

	// try as a contact field
	var nativeValue interface{}

	switch typed := c.fields[key].TypedValue().(type) {
	case nil:
		return nil
	case LocationPath:
		nativeValue = typed.Name()
	case types.XText:
		nativeValue = typed.Native()
	case types.XNumber:
		nativeValue = typed.Native()
	case types.XDateTime:
		nativeValue = typed.Native()
	}

	return []interface{}{nativeValue}
}

var _ contactql.Queryable = (*Contact)(nil)

// ContactReference is used to reference a contact
type ContactReference struct {
	UUID ContactUUID `json:"uuid" validate:"required,uuid4"`
	Name string      `json:"name"`
}

// NewContactReference creates a new contact reference with the given UUID and name
func NewContactReference(uuid ContactUUID, name string) *ContactReference {
	return &ContactReference{UUID: uuid, Name: name}
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type contactEnvelope struct {
	UUID      ContactUUID              `json:"uuid" validate:"required,uuid4"`
	ID        ContactID                `json:"id"`
	Name      string                   `json:"name"`
	Language  utils.Language           `json:"language"`
	Timezone  string                   `json:"timezone"`
	CreatedOn time.Time                `json:"created_on"`
	URNs      []urns.URN               `json:"urns" validate:"dive,urn"`
	Groups    []*assets.GroupReference `json:"groups,omitempty" validate:"dive"`
	Fields    map[string]*Value        `json:"fields,omitempty"`
}

// ReadContact decodes a contact from the passed in JSON
func ReadContact(assets SessionAssets, data json.RawMessage, strict bool) (*Contact, error) {
	var envelope contactEnvelope
	var err error

	if err := utils.UnmarshalAndValidate(data, &envelope); err != nil {
		return nil, fmt.Errorf("unable to read contact: %s", err)
	}

	c := &Contact{
		uuid:      envelope.UUID,
		id:        envelope.ID,
		name:      envelope.Name,
		language:  envelope.Language,
		createdOn: envelope.CreatedOn,
	}

	if envelope.Timezone != "" {
		if c.timezone, err = time.LoadLocation(envelope.Timezone); err != nil {
			return nil, err
		}
	}

	if envelope.URNs == nil {
		c.urns = make(URNList, 0)
	} else {
		if c.urns, err = ReadURNList(assets, envelope.URNs); err != nil {
			return nil, fmt.Errorf("error reading urns: %s", err)
		}
	}

	if envelope.Groups == nil {
		c.groups = NewGroupList([]*Group{})
	} else {
		groups := make([]*Group, 0, len(envelope.Groups))
		for g := range envelope.Groups {
			group, err := assets.Groups().Get(envelope.Groups[g].UUID)
			if err != nil && strict {
				return nil, fmt.Errorf("error reading groups: %s", err)
			}
			groups = append(groups, group)
		}
		c.groups = NewGroupList(groups)
	}

	if c.fields, err = NewFieldValues(assets, envelope.Fields, strict); err != nil {
		return nil, fmt.Errorf("error reading fields: %s", err)
	}

	return c, nil
}

// MarshalJSON marshals this contact into JSON
func (c *Contact) MarshalJSON() ([]byte, error) {
	ce := &contactEnvelope{
		Name:      c.name,
		UUID:      c.uuid,
		ID:        c.id,
		Language:  c.language,
		CreatedOn: c.createdOn,
	}

	ce.URNs = c.urns.RawURNs(true)
	if c.timezone != nil {
		ce.Timezone = c.timezone.String()
	}

	ce.Groups = make([]*assets.GroupReference, c.groups.Count())
	for g, group := range c.groups.All() {
		ce.Groups[g] = group.Reference()
	}

	ce.Fields = make(map[string]*Value)
	for _, v := range c.fields {
		if v != nil {
			ce.Fields[v.field.Key()] = v.Value
		}
	}

	return json.Marshal(ce)
}
