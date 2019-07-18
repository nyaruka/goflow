package flows

import (
	"encoding/json"
	"time"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
)

// NodeUUID is a UUID of a flow node
type NodeUUID utils.UUID

// CategoryUUID is the UUID of a node category
type CategoryUUID utils.UUID

// ExitUUID is the UUID of a node exit
type ExitUUID utils.UUID

// ActionUUID is the UUID of an action
type ActionUUID utils.UUID

// ContactID is the ID of a contact
type ContactID int64

// ContactUUID is the UUID of a contact
type ContactUUID utils.UUID

// RunUUID is the UUID of a flow run
type RunUUID utils.UUID

// StepUUID is the UUID of a run step
type StepUUID utils.UUID

// InputUUID is the UUID of an input
type InputUUID utils.UUID

// MsgID is the ID of a message
type MsgID int64

// NilMsgID is our constant for nil message ids
const NilMsgID = MsgID(0)

// MsgUUID is the UUID of a message
type MsgUUID utils.UUID

// FlowType represents the different types of flows
type FlowType string

const (
	// FlowTypeMessaging is a flow that is run over a messaging channel
	FlowTypeMessaging FlowType = "messaging"

	// FlowTypeMessagingOffline is a flow which is run over an offline messaging client like Surveyor
	FlowTypeMessagingOffline FlowType = "messaging_offline"

	// FlowTypeVoice is a flow which is run over IVR
	FlowTypeVoice FlowType = "voice"
)

// SessionStatus represents the current status of the engine session
type SessionStatus string

const (
	// SessionStatusActive represents a session that is still active
	SessionStatusActive SessionStatus = "active"

	// SessionStatusCompleted represents a session that has run to completion
	SessionStatusCompleted SessionStatus = "completed"

	// SessionStatusWaiting represents a session which is waiting for something from the caller
	SessionStatusWaiting SessionStatus = "waiting"

	// SessionStatusErrored represents a session that encountered an error
	SessionStatusErrored SessionStatus = "errored"
)

// RunStatus represents the current status of the flow run
type RunStatus string

const (
	// RunStatusActive represents a run that is still active
	RunStatusActive RunStatus = "active"

	// RunStatusCompleted represents a run that has run to completion
	RunStatusCompleted RunStatus = "completed"

	// RunStatusWaiting represents a run which is waiting for something from the caller
	RunStatusWaiting RunStatus = "waiting"

	// RunStatusErrored represents a run that encountered an error
	RunStatusErrored RunStatus = "errored"

	// RunStatusExpired represents a run that expired due to inactivity
	RunStatusExpired RunStatus = "expired"

	// RunStatusInterrupted represents a run that was interrupted by another flow
	RunStatusInterrupted RunStatus = "interrupted"
)

type FlowAssets interface {
	Get(assets.FlowUUID) (Flow, error)
}

// SessionAssets is the assets available to a session
type SessionAssets interface {
	Channels() *ChannelAssets
	Fields() *FieldAssets
	Flows() FlowAssets
	Groups() *GroupAssets
	Labels() *LabelAssets
	Locations() *LocationAssets
	Resthooks() *ResthookAssets
	Templates() *TemplateAssets
}

// Localizable is anything in the flow definition which can be localized and therefore needs a UUID
type Localizable interface {
	LocalizationUUID() utils.UUID
}

// Flow describes the ordered logic of actions and routers
type Flow interface {
	Contextable

	// spec properties
	UUID() assets.FlowUUID
	Name() string
	Revision() int
	Language() utils.Language
	Type() FlowType
	ExpireAfterMinutes() int
	Localization() Localization
	UI() json.RawMessage
	Nodes() []Node
	GetNode(uuid NodeUUID) Node
	Reference() *assets.FlowReference
	Generic() map[string]interface{}
	Clone(map[utils.UUID]utils.UUID) Flow

	Inspect() *FlowInfo
	Validate(SessionAssets, func(assets.Reference)) error
	ValidateRecursive(SessionAssets, func(assets.Reference)) error

	ExtractTemplates() []string
	ExtractDependencies() []assets.Reference
	ExtractResults() []*ResultInfo

	MarshalWithInfo() ([]byte, error)
}

// Node is a single node in a flow
type Node interface {
	Inspectable

	UUID() NodeUUID
	Actions() []Action
	Router() Router
	Exits() []Exit

	Validate(Flow, map[utils.UUID]bool) error
}

// Action is an action within a flow node
type Action interface {
	utils.Typed
	Localizable

	UUID() ActionUUID
	Execute(FlowRun, Step, ModifierCallback, EventCallback) error
	Validate() error
	AllowedFlowTypes() []FlowType
}

type Router interface {
	utils.Typed
	Inspectable

	Wait() Wait
	ResultName() string

	Validate([]Exit) error
	AllowTimeout() bool
	Route(FlowRun, Step, EventCallback) (ExitUUID, error)
	RouteTimeout(FlowRun, Step, EventCallback) (ExitUUID, error)
}

type Exit interface {
	UUID() ExitUUID
	DestinationUUID() NodeUUID
}

type Timeout interface {
	Seconds() int
	CategoryUUID() CategoryUUID
}

type Wait interface {
	utils.Typed

	Timeout() Timeout

	Begin(FlowRun, EventCallback) ActivatedWait
	End(Resume) error
}

type ActivatedWait interface {
	utils.Typed

	TimeoutSeconds() *int
}

type Hint interface {
	utils.Typed
}

// Localization provide a way to get the translations for a specific language
type Localization interface {
	AddItemTranslation(utils.Language, utils.UUID, string, []string)
	GetTranslations(utils.Language) Translations
	Languages() []utils.Language
}

// Translations provide a way to get the translation for a specific language for a uuid/key pair
type Translations interface {
	GetTextArray(utils.UUID, string) []string
	SetTextArray(utils.UUID, string, []string)
}

// Trigger represents something which can initiate a session with the flow engine
type Trigger interface {
	utils.Typed
	Contextable

	Initialize(Session, EventCallback) error
	InitializeRun(FlowRun, EventCallback) error

	Environment() utils.Environment
	Flow() *assets.FlowReference
	Contact() *Contact
	Connection() *Connection
	Params() types.XValue
	TriggeredOn() time.Time
}

// TriggerWithRun is special case of trigger that provides a parent run to the session
type TriggerWithRun interface {
	Trigger

	RunSummary() json.RawMessage
}

// Resume represents something which can resume a session with the flow engine
type Resume interface {
	utils.Typed

	Apply(FlowRun, EventCallback) error

	Environment() utils.Environment
	Contact() *Contact
	ResumedOn() time.Time
}

// Modifier is something which can modify a contact
type Modifier interface {
	utils.Typed

	Apply(utils.Environment, SessionAssets, *Contact, EventCallback)
}

// ModifierCallback is a callback invoked when a modifier has been generated
type ModifierCallback func(Modifier)

// Event describes a state change
type Event interface {
	utils.Typed

	CreatedOn() time.Time
	StepUUID() StepUUID
	SetStepUUID(StepUUID)
}

// EventCallback is a callback invoked when an event has been generated
type EventCallback func(Event)

// Input describes input from the contact and currently we only support one type of input: `msg`
type Input interface {
	utils.Typed
	Contextable

	UUID() InputUUID
	CreatedOn() time.Time
	Channel() *Channel
}

type Step interface {
	Contextable

	UUID() StepUUID
	NodeUUID() NodeUUID
	ExitUUID() ExitUUID
	ArrivedOn() time.Time

	Leave(ExitUUID)
}

type Engine interface {
	NewSession(SessionAssets, Trigger) (Session, Sprint, error)
	ReadSession(SessionAssets, json.RawMessage, assets.MissingCallback) (Session, error)

	HTTPClient() *utils.HTTPClient
	DisableWebhooks() bool
	MaxWebhookResponseBytes() int
	MaxStepsPerSprint() int
}

// Sprint is an interaction with the engine - i.e. a start or resume of a session
type Sprint interface {
	Modifiers() []Modifier
	LogModifier(Modifier)
	Events() []Event
	LogEvent(Event)
}

// Session represents the session of a flow run which may contain many runs
type Session interface {
	Assets() SessionAssets

	Type() FlowType
	SetType(FlowType)

	Environment() utils.Environment
	SetEnvironment(utils.Environment)

	Contact() *Contact
	SetContact(*Contact)

	Input() Input
	SetInput(Input)

	Status() SessionStatus
	Trigger() Trigger
	PushFlow(Flow, FlowRun, bool)
	Wait() ActivatedWait

	Resume(Resume) (Sprint, error)
	Runs() []FlowRun
	GetRun(RunUUID) (FlowRun, error)
	GetCurrentChild(FlowRun) FlowRun
	ParentRun() RunSummary

	Engine() Engine
}

// RunSummary represents the minimum information available about all runs (current or related) and is the
// representation of runs made accessible to router tests.
type RunSummary interface {
	UUID() RunUUID
	Contact() *Contact
	Flow() Flow
	Status() RunStatus
	Results() Results
}

// RunEnvironment is a run specific environment which adds location functionality required by some router tests
type RunEnvironment interface {
	utils.Environment

	FindLocations(string, utils.LocationLevel, *utils.Location) ([]*utils.Location, error)
	FindLocationsFuzzy(string, utils.LocationLevel, *utils.Location) ([]*utils.Location, error)
	LookupLocation(utils.LocationPath) (*utils.Location, error)
}

// FlowRun is a single contact's journey through a flow. It records the path they have taken,
// and the results that have been collected.
type FlowRun interface {
	Contextable
	RunSummary

	Environment() RunEnvironment
	Session() Session
	SaveResult(*Result)
	SetStatus(RunStatus)

	CreateStep(Node) Step
	Path() []Step
	PathLocation() (Step, Node, error)

	LogEvent(Step, Event)
	LogError(Step, error)
	Events() []Event

	EvaluateTemplateValue(template string) (types.XValue, error)
	EvaluateTemplate(template string) (string, error)

	GetText(utils.UUID, string, string) string
	GetTextArray(utils.UUID, string, []string) []string
	GetTranslatedTextArray(utils.UUID, string, []string, []utils.Language) []string

	Snapshot() RunSummary
	Parent() RunSummary
	ParentInSession() FlowRun
	Ancestors() []FlowRun

	CreatedOn() time.Time
	ModifiedOn() time.Time
	ExpiresOn() *time.Time
	ResetExpiration(*time.Time)
	ExitedOn() *time.Time
	Exit(RunStatus)
}

// LegacyExtraContributor is something which contributes results for constructing @legacy_extra
type LegacyExtraContributor interface {
	LegacyExtra() Results
}
