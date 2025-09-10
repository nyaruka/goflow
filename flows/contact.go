package flows

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/contactql"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/obfuscate"
	"github.com/shopspring/decimal"
)

func init() {
	utils.RegisterValidatorAlias("contact_status", "eq=active|eq=blocked|eq=stopped|eq=archived", func(validator.FieldError) string {
		return "is not a valid contact status"
	})
}

// ContactID is the ID of a contact
type ContactID int64

// ContactUUID is the UUID of a contact
type ContactUUID uuids.UUID

// NewContactUUID generates a new UUID for a contact
func NewContactUUID() ContactUUID { return ContactUUID(uuids.NewV4()) }

// ContactStatus is status in which a contact is in
type ContactStatus string

const (
	// ContactStatusActive is the contact status of active
	ContactStatusActive ContactStatus = "active"

	// ContactStatusBlocked is the contact status of blocked
	ContactStatusBlocked ContactStatus = "blocked"

	// ContactStatusStopped is the contact status of stopped
	ContactStatusStopped ContactStatus = "stopped"

	// ContactStatusArchived is the contact status of archived
	ContactStatusArchived ContactStatus = "archived"

	// MaxContactURNs is maximum number of URNs a contact can have
	MaxContactURNs = 50
)

// schemes of URNs that aren't tied to a specific channel
var portableURNSchemes = map[string]bool{urns.Phone.Prefix: true, urns.WhatsApp.Prefix: true}

// Contact represents a person who is interacting with the flow
type Contact struct {
	uuid       ContactUUID
	id         ContactID
	name       string
	language   i18n.Language
	status     ContactStatus
	timezone   *time.Location
	createdOn  time.Time
	lastSeenOn *time.Time
	urns       URNList
	groups     *GroupList
	fields     FieldValues
	tickets    *TicketList

	// transient fields
	assets SessionAssets
}

// NewContact creates a new contact with the passed in attributes
func NewContact(
	sa SessionAssets,
	uuid ContactUUID,
	id ContactID,
	name string,
	language i18n.Language,
	status ContactStatus,
	timezone *time.Location,
	createdOn time.Time,
	lastSeenOn *time.Time,
	urns []urns.URN,
	groups []*assets.GroupReference,
	fields map[string]*Value,
	tickets []*Ticket,
	missing assets.MissingCallback) (*Contact, error) {

	urnList, err := ReadURNList(sa, urns, missing)
	if err != nil {
		return nil, err
	}

	groupList := NewGroupList(sa, groups, missing)
	fieldValues := NewFieldValues(sa, fields, missing)
	ticketsList := NewTicketList(tickets)

	return &Contact{
		uuid:       uuid,
		id:         id,
		name:       name,
		language:   language,
		status:     status,
		timezone:   timezone,
		createdOn:  createdOn,
		lastSeenOn: lastSeenOn,
		urns:       urnList,
		groups:     groupList,
		fields:     fieldValues,
		tickets:    ticketsList,
		assets:     sa,
	}, nil
}

// NewEmptyContact creates a new empy contact with the passed in name, language and location
func NewEmptyContact(sa SessionAssets, name string, language i18n.Language, timezone *time.Location) *Contact {
	return &Contact{
		uuid:       NewContactUUID(),
		name:       name,
		language:   language,
		status:     ContactStatusActive,
		timezone:   timezone,
		createdOn:  dates.Now(),
		lastSeenOn: nil,
		urns:       URNList{},
		groups:     NewGroupList(sa, nil, assets.IgnoreMissing),
		fields:     make(FieldValues),
		tickets:    NewTicketList(nil),
		assets:     sa,
	}
}

// Clone clones this contact when we need to snapshot as part of parent run summary.
func (c *Contact) Clone() *Contact {
	return &Contact{
		uuid:       c.uuid,
		id:         c.id,
		name:       c.name,
		language:   c.language,
		status:     c.status,
		timezone:   c.timezone,
		createdOn:  c.createdOn,
		lastSeenOn: c.lastSeenOn,
		urns:       c.urns.clone(),
		groups:     c.groups.clone(),
		fields:     c.fields.clone(),
		tickets:    NewTicketList(nil), // tickets not included
		assets:     c.assets,
	}
}

// UUID returns the UUID of this contact
func (c *Contact) UUID() ContactUUID { return c.uuid }

// ID returns the numeric ID of this contact
func (c *Contact) ID() ContactID { return c.id }

// SetLanguage sets the language for this contact
func (c *Contact) SetLanguage(lang i18n.Language) { c.language = lang }

// Language gets the language for this contact
func (c *Contact) Language() i18n.Language { return c.language }

// Country gets the country for this contact..
//
// TODO: currently this is derived from their preferred channel or any tel URNs but probably should become an explicit
// field at some point
func (c *Contact) Country() i18n.Country {
	ch := c.PreferredChannel()
	if ch != nil && ch.Country() != i18n.NilCountry {
		return ch.Country()
	}

	for _, u := range c.urns {
		if u.urn.Scheme() == urns.Phone.Prefix {
			c := i18n.DeriveCountryFromTel(u.urn.Path())
			if c != i18n.NilCountry {
				return c
			}
		}
	}

	return i18n.NilCountry
}

// Locale gets the locale for this contact, using the environment country if contact doesn't have one
func (c *Contact) Locale(env envs.Environment) i18n.Locale {
	country := c.Country()
	if country == i18n.NilCountry {
		country = env.DefaultCountry()
	}
	return i18n.NewLocale(c.language, country)
}

// Status returns the contact status
func (c *Contact) Status() ContactStatus { return c.status }

// SetStatus sets the status of this contact (blocked, stopped or active)
func (c *Contact) SetStatus(status ContactStatus) { c.status = status }

// SetTimezone sets the timezone of this contact
func (c *Contact) SetTimezone(tz *time.Location) { c.timezone = tz }

// Timezone returns the timezone of this contact
func (c *Contact) Timezone() *time.Location { return c.timezone }

// CreatedOn returns the created on time of this contact
func (c *Contact) CreatedOn() time.Time { return c.createdOn }

// LastSeenOn returns the last seen on time of this contact
func (c *Contact) LastSeenOn() *time.Time { return c.lastSeenOn }

// SetLastSeenOn sets the last seen on time of this contact
func (c *Contact) SetLastSeenOn(t time.Time) { c.lastSeenOn = &t }

// SetName sets the name of this contact
func (c *Contact) SetName(name string) { c.name = name }

// Name returns the name of this contact
func (c *Contact) Name() string { return c.name }

// URNs returns the URNs of this contact
func (c *Contact) URNs() URNList { return c.urns }

// AddURN adds a new URN to this contact
func (c *Contact) AddURN(urn urns.URN) bool {
	if c.hasURN(urn) {
		return false
	}

	c.urns = append(c.urns, NewContactURN(urn, nil))
	return true
}

// RemoveURN adds a new URN to this contact
func (c *Contact) RemoveURN(urn urns.URN) bool {
	if !c.hasURN(urn) {
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

// SetURNs sets the URNs of this contact
func (c *Contact) SetURNs(urn []urns.URN) bool {
	isSame := func() bool {
		if len(c.urns) != len(urn) {
			return false
		}
		for i, u := range c.urns {
			if u.URN().Identity() != urn[i].Identity() {
				return false
			}
		}
		return true
	}
	if isSame() {
		return false
	}

	c.urns = URNList{}
	for _, u := range urn {
		c.urns = append(c.urns, NewContactURN(u.Normalize(), nil))
	}

	return true
}

func (c *Contact) hasURN(urn urns.URN) bool {
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

// Tickets returns the tickets for this contact
func (c *Contact) Tickets() *TicketList { return c.tickets }

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
//	__default__:text -> the name or URN
//	uuid:text -> the UUID of the contact
//	ref:text -> the reference identifier of the contact
//	first_name:text -> the first name of the contact
//	name:text -> the name of the contact
//	language:text -> the language of the contact as 3-letter ISO code
//	status:text -> the status of the contact
//	created_on:datetime -> the creation date of the contact
//	last_seen_on:any -> the last seen date of the contact
//	urns:[]text -> the URNs belonging to the contact
//	urn:text -> the preferred URN of the contact
//	groups:[]group -> the groups the contact belongs to
//	fields:fields -> the custom field values of the contact
//	channel:channel -> the preferred channel of the contact
//
// @context contact
func (c *Contact) Context(env envs.Environment) map[string]types.XValue {
	var firstName, urn, timezone, lastSeenOn types.XValue

	id := types.NewXText(strconv.Itoa(int(c.id)))
	id.SetDeprecated("id: use ref instead")

	ref, _ := obfuscate.EncodeID(int64(c.id), env.ObfuscationKey())

	if c.timezone != nil {
		timezone = types.NewXText(c.timezone.String())
	}

	preferredURN := c.PreferredURN()
	if preferredURN != nil {
		urn = preferredURN.ToXValue(env)
	}

	names := utils.TokenizeString(c.name)
	if len(names) >= 1 {
		firstName = types.NewXText(names[0])
	}

	if c.lastSeenOn != nil {
		lastSeenOn = types.NewXDateTime(*c.lastSeenOn)
	}

	return map[string]types.XValue{
		"__default__":  types.NewXText(c.Format(env)),
		"uuid":         types.NewXText(string(c.uuid)),
		"id":           id,
		"ref":          types.NewXText(ref),
		"name":         types.NewXText(c.name),
		"first_name":   firstName,
		"language":     types.NewXText(string(c.language)),
		"timezone":     timezone,
		"status":       types.NewXText(string(c.status)),
		"created_on":   types.NewXDateTime(c.createdOn),
		"last_seen_on": lastSeenOn,
		"urns":         c.urns.ToXValue(env),
		"urn":          urn,
		"groups":       c.groups.ToXValue(env),
		"fields":       Context(env, c.Fields()),
		"channel":      Context(env, c.PreferredChannel()),
		"tickets":      c.tickets.ToXValue(env),
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
		if !channel.HasRole(assets.ChannelRoleSend) {
			return false
		}

		priorityURNs := make([]*ContactURN, 0)
		otherURNs := make([]*ContactURN, 0)

		for _, urn := range c.urns {
			// portable URNs can be re-assigned when supported by channel
			if portableURNSchemes[urn.URN().Scheme()] && channel.SupportsScheme(urn.URN().Scheme()) {
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

// ReevaluateQueryBasedGroups reevaluates membership of all query based groups for this contact
func (c *Contact) ReevaluateQueryBasedGroups(env envs.Environment) ([]*Group, []*Group) {
	added := make([]*Group, 0)
	removed := make([]*Group, 0)

	for _, group := range c.assets.Groups().All() {
		if !group.UsesQuery() {
			continue
		}

		qualifies := group.CheckQueryBasedMembership(env, c)

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

	return added, removed
}

// QueryProperty resolves a contact query search key for this contact
//
// Note that this method excludes id, ref, group and flow search attributes as those are disallowed
// query based groups.
func (c *Contact) QueryProperty(env envs.Environment, key string, propType contactql.PropertyType) []any {
	if propType == contactql.PropertyTypeAttribute {
		switch key {
		case contactql.AttributeUUID:
			return []any{string(c.uuid)}
		case contactql.AttributeName:
			if c.name != "" {
				return []any{c.name}
			}
			return nil
		case contactql.AttributeLanguage:
			if c.language != i18n.NilLanguage {
				return []any{string(c.language)}
			}
			return nil
		case contactql.AttributeURN:
			vals := make([]any, len(c.URNs()))
			for i, urn := range c.URNs() {
				vals[i] = urn.URN().Path()
			}
			return vals
		case contactql.AttributeTickets:
			return []any{decimal.NewFromInt(int64(c.tickets.OpenCount()))}
		case contactql.AttributeCreatedOn:
			return []any{c.createdOn}
		case contactql.AttributeLastSeenOn:
			if c.lastSeenOn != nil {
				return []any{*c.lastSeenOn}
			}
			return nil
		default:
			return nil
		}
	} else if propType == contactql.PropertyTypeURN {
		urnsWithScheme := c.urns.WithScheme(key)
		vals := make([]any, len(urnsWithScheme))
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

	return []any{nativeValue}
}

var _ contactql.Queryable = (*Contact)(nil)

// ContactReference is used to reference a contact
type ContactReference struct {
	UUID ContactUUID `json:"uuid" validate:"required,uuid"`
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

type ContactEnvelope struct {
	UUID       ContactUUID              `json:"uuid"                validate:"required,uuid"`
	ID         ContactID                `json:"id,omitempty"`
	Name       string                   `json:"name,omitempty"`
	Language   i18n.Language            `json:"language,omitempty"`
	Status     ContactStatus            `json:"status,omitempty"    validate:"required,contact_status"`
	Timezone   string                   `json:"timezone,omitempty"`
	CreatedOn  time.Time                `json:"created_on"          validate:"required"`
	LastSeenOn *time.Time               `json:"last_seen_on,omitempty"`
	URNs       []urns.URN               `json:"urns,omitempty"      validate:"dive,urn"`
	Groups     []*assets.GroupReference `json:"groups,omitempty"    validate:"dive"`
	Fields     map[string]*Value        `json:"fields,omitempty"`
	Tickets    []*TicketEnvelope        `json:"tickets,omitempty"`
}

func (e *ContactEnvelope) Unmarshal(sa SessionAssets, missing assets.MissingCallback) (*Contact, error) {
	c := &Contact{
		uuid:       e.UUID,
		id:         e.ID,
		name:       e.Name,
		language:   e.Language,
		status:     e.Status,
		createdOn:  e.CreatedOn,
		lastSeenOn: e.LastSeenOn,
		assets:     sa,
	}

	if e.Timezone != "" {
		var err error
		if c.timezone, err = time.LoadLocation(e.Timezone); err != nil {
			return nil, err
		}
	}

	if e.URNs == nil {
		c.urns = make(URNList, 0)
	} else {
		var err error
		if c.urns, err = ReadURNList(sa, e.URNs, missing); err != nil {
			return nil, fmt.Errorf("error reading urns: %w", err)
		}
	}

	c.groups = NewGroupList(sa, e.Groups, missing)
	c.fields = NewFieldValues(sa, e.Fields, missing)

	tickets := make([]*Ticket, len(e.Tickets))
	for i, t := range e.Tickets {
		tickets[i] = t.Unmarshal(sa, missing)
	}
	c.tickets = NewTicketList(tickets)

	return c, nil
}

func (c *Contact) Marshal() *ContactEnvelope {
	e := &ContactEnvelope{
		Name:       c.name,
		UUID:       c.uuid,
		ID:         c.id,
		Status:     c.status,
		Language:   c.language,
		CreatedOn:  c.createdOn,
		LastSeenOn: c.lastSeenOn,
		URNs:       c.urns.RawURNs(),
		Groups:     c.groups.references(),
	}

	if c.timezone != nil {
		e.Timezone = c.timezone.String()
	}

	e.Fields = make(map[string]*Value)
	for _, v := range c.fields {
		if v != nil {
			e.Fields[v.field.Key()] = v.Value
		}
	}

	e.Tickets = c.Tickets().Marshal()

	return e
}

// ReadContact decodes a contact from the passed in JSON
func ReadContact(sa SessionAssets, data []byte, missing assets.MissingCallback) (*Contact, error) {
	e := &ContactEnvelope{}

	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, fmt.Errorf("unable to read contact: %w", err)
	}

	return e.Unmarshal(sa, missing)
}

// MarshalJSON marshals this contact into JSON
func (c *Contact) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(c.Marshal())
}
