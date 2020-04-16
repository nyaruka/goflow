package flows

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/contactql"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/dates"
	"github.com/nyaruka/goflow/utils/jsonx"
	"github.com/nyaruka/goflow/utils/uuids"
	"github.com/shopspring/decimal"

	"github.com/pkg/errors"
)

// Contact represents a person who is interacting with the flow
type Contact struct {
	uuid      ContactUUID
	id        ContactID
	name      string
	language  envs.Language
	blocked   bool
	stopped   bool
	timezone  *time.Location
	createdOn time.Time
	urns      URNList
	groups    *GroupList
	fields    FieldValues

	// transient fields
	assets SessionAssets
}

// NewContact creates a new contact with the passed in attributes
func NewContact(
	sa SessionAssets,
	uuid ContactUUID,
	id ContactID,
	name string,
	language envs.Language,
	blocked bool,
	stopped bool,
	timezone *time.Location,
	createdOn time.Time,
	urns []urns.URN,
	groups []*assets.GroupReference,
	fields map[string]*Value,
	missing assets.MissingCallback) (*Contact, error) {

	urnList, err := ReadURNList(sa, urns, missing)
	if err != nil {
		return nil, err
	}

	groupList := NewGroupList(sa, groups, missing)

	fieldValues, err := NewFieldValues(sa, fields, missing)
	if err != nil {
		return nil, err
	}

	return &Contact{
		uuid:      uuid,
		id:        id,
		name:      name,
		language:  language,
		blocked:   blocked,
		stopped:   stopped,
		timezone:  timezone,
		createdOn: createdOn,
		urns:      urnList,
		groups:    groupList,
		fields:    fieldValues,
		assets:    sa,
	}, nil
}

// NewEmptyContact creates a new empy contact with the passed in name, language and location
func NewEmptyContact(sa SessionAssets, name string, language envs.Language, timezone *time.Location) *Contact {
	return &Contact{
		uuid:      ContactUUID(uuids.New()),
		name:      name,
		language:  language,
		blocked:   false,
		stopped:   false,
		timezone:  timezone,
		createdOn: dates.Now(),
		urns:      URNList{},
		groups:    NewGroupList(sa, nil, assets.IgnoreMissing),
		fields:    make(FieldValues),
		assets:    sa,
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
		blocked:   c.blocked,
		stopped:   c.stopped,
		timezone:  c.timezone,
		createdOn: c.createdOn,
		urns:      c.urns.clone(),
		groups:    c.groups.clone(),
		fields:    c.fields.clone(),
		assets:    c.assets,
	}
}

// Equal returns true if this instance is equal to the given instance
func (c *Contact) Equal(other *Contact) bool {
	asJSON1, _ := jsonx.Marshal(c)
	asJSON2, _ := jsonx.Marshal(other)
	return string(asJSON1) == string(asJSON2)
}

// UUID returns the UUID of this contact
func (c *Contact) UUID() ContactUUID { return c.uuid }

// ID returns the numeric ID of this contact
func (c *Contact) ID() ContactID { return c.id }

// SetLanguage sets the language for this contact
func (c *Contact) SetLanguage(lang envs.Language) { c.language = lang }

// Language gets the language for this contact
func (c *Contact) Language() envs.Language { return c.language }

// Blocked returns whether the contact is blocked or not
func (c *Contact) Blocked() bool { return c.blocked }

// SetBlocked sets the blocked state of this contact
func (c *Contact) SetBlocked(blocked bool) { c.blocked = blocked }

// Stopped returns whether the contact is stopped or not
func (c *Contact) Stopped() bool { return c.stopped }

// SetStopped sets the stopped stte of this contact
func (c *Contact) SetStopped(stopped bool) { c.stopped = stopped }

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
func (c *Contact) AddURN(urn urns.URN, channel *Channel) bool {
	if c.HasURN(urn) {
		return false
	}

	c.urns = append(c.urns, NewContactURN(urn, channel))
	return true
}

// RemoveURN adds a new URN to this contact
func (c *Contact) RemoveURN(urn urns.URN) bool {
	if !c.HasURN(urn) {
		return false
	}

	newURNs := make([]*ContactURN, 0, len(c.urns)-1)
	for _, u := range c.urns {
		if u.URN().Identity() != urn.Identity() {
			newURNs = append(newURNs, u)
		}
	}

	c.urns = URNList(newURNs)
	return true
}

// HasURN checks whether the contact has the given URN
func (c *Contact) HasURN(urn urns.URN) bool {
	urn = urn.Normalize("")

	for _, u := range c.urns {
		if u.URN().Identity() == urn.Identity() {
			return true
		}
	}
	return false
}

// Fields returns this contact's field values
func (c *Contact) Fields() FieldValues { return c.fields }

// Groups returns the groups that this contact belongs to
func (c *Contact) Groups() *GroupList { return c.groups }

// Reference returns a reference to this contact
func (c *Contact) Reference() *ContactReference {
	if c == nil {
		return nil
	}
	return NewContactReference(c.uuid, c.name)
}

// Format returns a friendly string version of this contact depending on what fields are set
func (c *Contact) Format(env envs.Environment) string {
	// if contact has a name set, use that
	if c.name != "" {
		return c.name
	}

	// otherwise use either id or the highest priority URN depending on the env
	if env.RedactionPolicy() == envs.RedactionPolicyURNs {
		return strconv.Itoa(int(c.id))
	}
	if len(c.urns) > 0 {
		return c.urns[0].URN().Format()
	}

	return ""
}

// Context returns the properties available in expressions
//
//   __default__:text -> the name or URN
//   uuid:text -> the UUID of the contact
//   id:text -> the numeric ID of the contact
//   first_name:text -> the first name of the contact
//   name:text -> the name of the contact
//   language:text -> the language of the contact as 3-letter ISO code
//   created_on:datetime -> the creation date of the contact
//   urns:[]text -> the URNs belonging to the contact
//   urn:text -> the preferred URN of the contact
//   groups:[]group -> the groups the contact belongs to
//   fields:fields -> the custom field values of the contact
//   channel:channel -> the preferred channel of the contact
//
// @context contact
func (c *Contact) Context(env envs.Environment) map[string]types.XValue {
	var urn, timezone types.XValue
	if c.timezone != nil {
		timezone = types.NewXText(c.timezone.String())
	}
	preferredURN := c.PreferredURN()
	if preferredURN != nil {
		urn = preferredURN.ToXValue(env)
	}

	var firstName types.XValue
	names := utils.TokenizeString(c.name)
	if len(names) >= 1 {
		firstName = types.NewXText(names[0])
	}

	return map[string]types.XValue{
		"__default__": types.NewXText(c.Format(env)),
		"uuid":        types.NewXText(string(c.uuid)),
		"id":          types.NewXText(strconv.Itoa(int(c.id))),
		"name":        types.NewXText(c.name),
		"first_name":  firstName,
		"language":    types.NewXText(string(c.language)),
		"timezone":    timezone,
		"created_on":  types.NewXDateTime(c.createdOn),
		"urns":        c.urns.ToXValue(env),
		"urn":         urn,
		"groups":      c.groups.ToXValue(env),
		"fields":      Context(env, c.Fields()),
		"channel":     Context(env, c.PreferredChannel()),
	}
}

// Destination is a sendable channel and URN pair
type Destination struct {
	Channel *Channel
	URN     *ContactURN
}

// ResolveDestinations resolves possible URN/channel destinations
func (c *Contact) ResolveDestinations(all bool) []Destination {
	destinations := []Destination{}

	for _, u := range c.urns {
		channel := c.assets.Channels().GetForURN(u, assets.ChannelRoleSend)
		if channel != nil {
			destinations = append(destinations, Destination{URN: u, Channel: channel})
			if !all {
				break
			}
		}
	}
	return destinations
}

// PreferredURN gets the preferred URN for this contact, i.e. the URN we would use for sending
func (c *Contact) PreferredURN() *ContactURN {
	destinations := c.ResolveDestinations(false)
	if len(destinations) > 0 {
		return destinations[0].URN
	}
	return nil
}

// PreferredChannel gets the preferred channel for this contact, i.e. the channel we would use for sending
func (c *Contact) PreferredChannel() *Channel {
	destinations := c.ResolveDestinations(false)
	if len(destinations) > 0 {
		return destinations[0].Channel
	}
	return nil
}

// UpdatePreferredChannel updates the preferred channel and returns whether any change was made
func (c *Contact) UpdatePreferredChannel(channel *Channel) bool {
	oldURNs := c.urns.clone()

	// setting preferred channel to nil means clearing affinity on all URNs
	if channel == nil {
		for _, urn := range c.urns {
			urn.SetChannel(nil)
		}
	} else {
		priorityURNs := make([]*ContactURN, 0)
		otherURNs := make([]*ContactURN, 0)

		for _, urn := range c.urns {
			// tel URNs can be re-assigned, other URN schemes are considered channel specific
			if urn.URN().Scheme() == urns.TelScheme && channel.SupportsScheme(urns.TelScheme) {
				urn.SetChannel(channel)
			}

			// If URN doesn't have a channel and is a scheme supported by the channel, then we can set its
			// channel. This may result in unsendable URN/channel pairing but can't do much about that.
			if urn.Channel() == nil && channel.SupportsScheme(urn.URN().Scheme()) {
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

	return !oldURNs.Equal(c.urns)
}

// ReevaluateDynamicGroups reevaluates membership of all dynamic groups for this contact
func (c *Contact) ReevaluateDynamicGroups(env envs.Environment) ([]*Group, []*Group, []error) {
	added := make([]*Group, 0)
	removed := make([]*Group, 0)
	errs := make([]error, 0)

	if c.Blocked() || c.Stopped() {
		return added, removed, errs
	}

	for _, group := range c.assets.Groups().All() {
		if !group.IsDynamic() {
			continue
		}

		qualifies, err := group.CheckDynamicMembership(env, c)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "unable to re-evaluate membership of group '%s'", group.Name()))
		}

		if qualifies {
			if c.groups.Add(group) {
				added = append(added, group)
			}
		} else {
			if c.groups.Remove(group) {
				removed = append(removed, group)
			}
		}
	}

	return added, removed, errs
}

// QueryProperty resolves a contact query search key for this contact
func (c *Contact) QueryProperty(env envs.Environment, key string, propType contactql.PropertyType) []interface{} {
	if propType == contactql.PropertyTypeAttribute {
		switch key {
		case contactql.AttributeID:
			if c.id != 0 {
				return []interface{}{decimal.New(int64(c.id), 0)}
			}
			return nil
		case contactql.AttributeName:
			if c.name != "" {
				return []interface{}{c.name}
			}
			return nil
		case contactql.AttributeLanguage:
			if c.language != envs.NilLanguage {
				return []interface{}{string(c.language)}
			}
			return nil
		case contactql.AttributeURN:
			vals := make([]interface{}, len(c.URNs()))
			for i, urn := range c.URNs() {
				vals[i] = urn.URN().Path()
			}
			return vals
		case contactql.AttributeGroup:
			vals := make([]interface{}, c.Groups().Count())
			for i, group := range c.Groups().All() {
				vals[i] = group.Name()
			}
			return vals
		case contactql.AttributeCreatedOn:
			return []interface{}{c.createdOn}
		default:
			return nil
		}
	} else if propType == contactql.PropertyTypeScheme {
		urnsWithScheme := c.urns.WithScheme(key)
		vals := make([]interface{}, len(urnsWithScheme))
		for i := range urnsWithScheme {
			vals[i] = urnsWithScheme[i].URN().Path()
		}
		return vals
	}

	// try as a contact field
	nativeValue := c.fields[key].QueryValue()
	if nativeValue == nil {
		return nil
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

// Type returns the name of the asset type
func (r *ContactReference) Type() string {
	return "contact"
}

// Identity returns the unique identity of the asset
func (r *ContactReference) Identity() string {
	return string(r.UUID)
}

// Variable returns whether this a variable (vs concrete) reference
func (r *ContactReference) Variable() bool {
	return r.Identity() == ""
}

func (r *ContactReference) String() string {
	return fmt.Sprintf("%s[uuid=%s,name=%s]", r.Type(), r.Identity(), r.Name)
}

var _ assets.Reference = (*ContactReference)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type contactEnvelope struct {
	UUID      ContactUUID              `json:"uuid" validate:"required,uuid4"`
	ID        ContactID                `json:"id,omitempty"`
	Name      string                   `json:"name,omitempty"`
	Language  envs.Language            `json:"language,omitempty"`
	Stopped   bool                     `json:"stopped,omitempty"`
	Blocked   bool                     `json:"blocked,omitempty"`
	Timezone  string                   `json:"timezone,omitempty"`
	CreatedOn time.Time                `json:"created_on" validate:"required"`
	URNs      []urns.URN               `json:"urns,omitempty" validate:"dive,urn"`
	Groups    []*assets.GroupReference `json:"groups,omitempty" validate:"dive"`
	Fields    map[string]*Value        `json:"fields,omitempty"`
}

// ReadContact decodes a contact from the passed in JSON
func ReadContact(sa SessionAssets, data json.RawMessage, missing assets.MissingCallback) (*Contact, error) {
	var envelope contactEnvelope
	var err error

	if err := utils.UnmarshalAndValidate(data, &envelope); err != nil {
		return nil, errors.Wrap(err, "unable to read contact")
	}

	c := &Contact{
		uuid:      envelope.UUID,
		id:        envelope.ID,
		name:      envelope.Name,
		language:  envelope.Language,
		blocked:   envelope.Blocked,
		stopped:   envelope.Stopped,
		createdOn: envelope.CreatedOn,
		assets:    sa,
	}

	if envelope.Timezone != "" {
		if c.timezone, err = time.LoadLocation(envelope.Timezone); err != nil {
			return nil, err
		}
	}

	if envelope.URNs == nil {
		c.urns = make(URNList, 0)
	} else {
		if c.urns, err = ReadURNList(sa, envelope.URNs, missing); err != nil {
			return nil, errors.Wrap(err, "error reading urns")
		}
	}

	c.groups = NewGroupList(sa, envelope.Groups, missing)

	if c.fields, err = NewFieldValues(sa, envelope.Fields, missing); err != nil {
		return nil, errors.Wrap(err, "error reading fields")
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
		Blocked:   c.blocked,
		Stopped:   c.stopped,
		CreatedOn: c.createdOn,
	}

	ce.URNs = c.urns.RawURNs()
	if c.timezone != nil {
		ce.Timezone = c.timezone.String()
	}

	ce.Groups = make([]*assets.GroupReference, c.groups.Count())
	for i, group := range c.groups.All() {
		ce.Groups[i] = group.Reference()
	}

	ce.Fields = make(map[string]*Value)
	for _, v := range c.fields {
		if v != nil {
			ce.Fields[v.field.Key()] = v.Value
		}
	}

	return jsonx.Marshal(ce)
}
