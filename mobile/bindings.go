package mobile

// To build an Android Archive:
//
// go get golang.org/x/mobile/cmd/gomobile
// gomobile bind -target android -javapkg=com.nyaruka.goflow -o mobile/goflow.aar github.com/nyaruka/goflow/mobile

import (
	"encoding/json"
	"time"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/definition/migrations"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/resumes"
	"github.com/nyaruka/goflow/flows/routers/waits"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/utils"

	"github.com/Masterminds/semver"
)

// CurrentSpecVersion returns the current flow spec version
func CurrentSpecVersion() string {
	return definition.CurrentSpecVersion.String()
}

// IsVersionSupported returns whether the given spec version is supported
func IsVersionSupported(version string) bool {
	v, err := semver.NewVersion(version)
	if err != nil {
		return false
	}

	return definition.IsVersionSupported(v)
}

// Environment defines the environment for expression evaluation etc
type Environment struct {
	target envs.Environment
}

// NewEnvironment creates a new environment.
func NewEnvironment(dateFormat string, timeFormat string, timezone string, allowedLanguages *StringSlice, defaultCountry string, redactionPolicy string) (*Environment, error) {
	tz, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, err
	}

	langs := make([]envs.Language, allowedLanguages.Length())
	for i := 0; i < allowedLanguages.Length(); i++ {
		langs[i] = envs.Language(allowedLanguages.Get(i))
	}

	return &Environment{
		target: envs.NewBuilder().
			WithDateFormat(envs.DateFormat(dateFormat)).
			WithTimeFormat(envs.TimeFormat(timeFormat)).
			WithTimezone(tz).
			WithAllowedLanguages(langs).
			WithDefaultCountry(envs.Country(defaultCountry)).
			WithRedactionPolicy(envs.RedactionPolicy(redactionPolicy)).
			Build(),
	}, nil
}

// AssetsSource is a static asset source
type AssetsSource struct {
	target *static.StaticSource
}

// NewAssetsSource creates a new static asset source
func NewAssetsSource(src string) (*AssetsSource, error) {
	s, err := static.NewSource(json.RawMessage(src))
	if err != nil {
		return nil, err
	}
	return &AssetsSource{target: s}, nil
}

// SessionAssets provides optimized access to assets
type SessionAssets struct {
	target flows.SessionAssets
}

// NewSessionAssets creates a new session assets
func NewSessionAssets(environment *Environment, source *AssetsSource) (*SessionAssets, error) {
	s, err := engine.NewSessionAssets(environment.target, source.target, &migrations.Config{BaseMediaURL: ""})
	if err != nil {
		return nil, err
	}
	return &SessionAssets{target: s}, nil
}

// Contact represents a person who is interacting with a flow
type Contact struct {
	target *flows.Contact
}

// NewEmptyContact creates a new contact
func NewEmptyContact(sa *SessionAssets) *Contact {
	return &Contact{
		target: flows.NewEmptyContact(sa.target, "", envs.NilLanguage, nil),
	}
}

// MsgIn is an incoming message
type MsgIn struct {
	target *flows.MsgIn
}

// NewMsgIn creates a new incoming message
func NewMsgIn(uuid string, text string, attachments *StringSlice) *MsgIn {
	var convertedAttachments []utils.Attachment
	if attachments != nil {
		convertedAttachments = make([]utils.Attachment, attachments.Length())
		for i := 0; i < attachments.Length(); i++ {
			convertedAttachments[i] = utils.Attachment(attachments.Get(i))
		}
	}

	return &MsgIn{
		target: flows.NewMsgIn(flows.MsgUUID(uuid), urns.NilURN, nil, text, convertedAttachments),
	}
}

func (m *MsgIn) Text() string {
	return m.target.Text()
}

func (m *MsgIn) Attachments() *StringSlice {
	attachments := NewStringSlice(len(m.target.Attachments()))
	for _, attachment := range m.target.Attachments() {
		attachments.Add(string(attachment))
	}
	return attachments
}

// FlowReference is a reference to a flow
type FlowReference struct {
	uuid string
	name string
}

// NewFlowReference creates a new flow reference
func NewFlowReference(uuid string, name string) *FlowReference {
	return &FlowReference{uuid: uuid, name: name}
}

// Trigger represents something which can initiate a session
type Trigger struct {
	target flows.Trigger
}

// NewManualTrigger creates a new manual trigger
func NewManualTrigger(environment *Environment, contact *Contact, flow *FlowReference) *Trigger {
	flowRef := assets.NewFlowReference(assets.FlowUUID(flow.uuid), flow.name)
	return &Trigger{
		target: triggers.NewBuilder(environment.target, flowRef, contact.target).Manual().Build(),
	}
}

// Resume represents something which can resume a session
type Resume struct {
	target flows.Resume
}

// NewMsgResume creates a new message resume
func NewMsgResume(environment *Environment, contact *Contact, msg *MsgIn) *Resume {
	var e envs.Environment
	if environment != nil {
		e = environment.target
	}
	var c *flows.Contact
	if contact != nil {
		c = contact.target
	}

	return &Resume{
		target: resumes.NewMsg(e, c, msg.target),
	}
}

type Event struct {
	type_   string
	payload string
}

func (e *Event) Type() string {
	return e.type_
}

func (e *Event) Payload() string {
	return e.payload
}

type Modifier struct {
	type_   string
	payload string
}

func (m *Modifier) Type() string {
	return m.type_
}

func (m *Modifier) Payload() string {
	return m.payload
}

// Sprint is an interaction with the engine - i.e. a start or resume of a session
type Sprint struct {
	target flows.Sprint
}

// Modifiers returns the modifiers created during this sprint
func (s *Sprint) Modifiers() *ModifierSlice {
	mods := NewModifierSlice(len(s.target.Modifiers()))
	for _, mod := range s.target.Modifiers() {
		marshaled, _ := jsonx.Marshal(mod)
		mods.Add(&Modifier{type_: mod.Type(), payload: string(marshaled)})
	}
	return mods
}

// Events returns the events created during this sprint
func (s *Sprint) Events() *EventSlice {
	events := NewEventSlice(len(s.target.Events()))
	for _, event := range s.target.Events() {
		marshaled, _ := jsonx.Marshal(event)
		events.Add(&Event{type_: event.Type(), payload: string(marshaled)})
	}
	return events
}

// Session represents a session with the flow engine
type Session struct {
	target flows.Session
}

// Status returns the status of this session
func (s *Session) Status() string {
	return string(s.target.Status())
}

// Assets returns the assets associated with this session
func (s *Session) Assets() *SessionAssets {
	return &SessionAssets{target: s.target.Assets()}
}

// Resume resumes this session
func (s *Session) Resume(resume *Resume) (*Sprint, error) {
	sprint, err := s.target.Resume(resume.target)
	if err != nil {
		return nil, err
	}
	return &Sprint{target: sprint}, nil
}

// GetWait gets the current wait of this session.. can't call this Wait() because Object in Java already has a wait() method
func (s *Session) GetWait() *Wait {
	if s.target.Wait() != nil {
		return &Wait{target: s.target.Wait()}
	}
	return nil
}

// ToJSON serializes this session as JSON
func (s *Session) ToJSON() (string, error) {
	data, err := jsonx.Marshal(s.target)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type Hint struct {
	target flows.Hint
}

func (h *Hint) Type() string {
	return string(h.target.Type())
}

type Wait struct {
	target flows.ActivatedWait
}

func (w *Wait) Type() string {
	return string(w.target.Type())
}

func (w *Wait) Hint() *Hint {
	asMsgWait, isMsgWait := w.target.(*waits.ActivatedMsgWait)
	if isMsgWait && asMsgWait.Hint() != nil {
		return &Hint{target: asMsgWait.Hint()}
	}
	return nil
}

type Engine struct {
	target flows.Engine
}

func NewEngine() *Engine {
	return &Engine{
		target: engine.NewBuilder().Build(),
	}
}

// NewSession creates a new session
func (e *Engine) NewSession(sa *SessionAssets, trigger *Trigger) (*SessionAndSprint, error) {
	session, sprint, err := e.target.NewSession(sa.target, trigger.target)
	if err != nil {
		return nil, err
	}

	return &SessionAndSprint{
		session: &Session{target: session},
		sprint:  &Sprint{target: sprint},
	}, nil
}

// ReadSession reads an existing session from JSON
func (e *Engine) ReadSession(a *SessionAssets, data string) (*Session, error) {
	s, err := e.target.ReadSession(a.target, []byte(data), assets.IgnoreMissing)
	if err != nil {
		return nil, err
	}
	return &Session{target: s}, nil
}

// SessionAndSprint holds a session and a sprint.. because a Java method can't return two values
type SessionAndSprint struct {
	session *Session
	sprint  *Sprint
}

func (ss *SessionAndSprint) Session() *Session {
	return ss.session
}

func (ss *SessionAndSprint) Sprint() *Sprint {
	return ss.sprint
}
