package mobile

// To build an Android Archive:
//
// gomobile bind -target android -o mobile/goflow.aar github.com/nyaruka/goflow/mobile

import (
	"encoding/json"
	"time"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/resumes"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/utils"
)

// Environment defines the environment for expression evaluation etc
type Environment struct {
	target utils.Environment
}

// NewEnvironment creates a new environment
func NewEnvironment(dateFormat string, timeFormat string, timezone string, defaultLanguage string, allowedLanguages []string) (*Environment, error) {
	tz, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, err
	}

	langs := make([]utils.Language, len(allowedLanguages))
	for l := range allowedLanguages {
		langs[l] = utils.Language(allowedLanguages[l])
	}

	return &Environment{
		target: utils.NewEnvironment(
			utils.DateFormat(dateFormat),
			utils.TimeFormat(timeFormat),
			tz,
			utils.Language(defaultLanguage),
			langs,
			utils.DefaultNumberFormat,
			utils.RedactionPolicyNone,
		),
	}, nil
}

// AssetsSource is a static asset source
type AssetsSource struct {
	target *static.StaticSource
}

// NewAssetsSource creates a new static asset source
func NewAssetsSource(src string) (*AssetsSource, error) {
	s, err := static.NewStaticSource(json.RawMessage(src))
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
func NewSessionAssets(source *AssetsSource) (*SessionAssets, error) {
	s, err := engine.NewSessionAssets(source.target)
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
func NewEmptyContact() *Contact {
	return &Contact{
		target: flows.NewEmptyContact("", utils.NilLanguage, nil),
	}
}

// Trigger represents something which can initiate a session
type Trigger struct {
	target flows.Trigger
}

// NewManualTrigger creates a new manual trigger
func NewManualTrigger(environment *Environment, contact *Contact, flowUUID string, flowName string) *Trigger {
	flow := assets.NewFlowReference(assets.FlowUUID(flowUUID), flowName)
	return &Trigger{
		target: triggers.NewManualTrigger(environment.target, contact.target, flow, nil, utils.Now()),
	}
}

// Resume represents something which can resume a session
type Resume struct {
	target flows.Resume
}

// NewMsgResume creates a new message resume
func NewMsgResume(text string) *Resume {
	msg := flows.NewMsgIn(flows.MsgUUID(utils.NewUUID()), 0, urns.NilURN, nil, text, nil)
	return &Resume{
		target: resumes.NewMsgResume(nil, nil, msg),
	}
}

// Session represents a session with the flow engine
type Session struct {
	target flows.Session
}

// NewSession creates a new session
func NewSession(a *SessionAssets, httpUserAgent string) *Session {
	httpClient := utils.NewHTTPClient(httpUserAgent)
	s := engine.NewSession(a.target, engine.NewDefaultConfig(), httpClient)
	return &Session{target: s}
}

// Start starts this session using the given trigger
func (s *Session) Start(trigger *Trigger) error {
	return s.target.Start(trigger.target)
}

// Resume resumes this session
func (s *Session) Resume(resume *Resume) error {
	return s.target.Resume(resume.target)
}
