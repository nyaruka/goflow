package flows

import (
	"time"

	"github.com/nyaruka/goflow/utils"
)

// NodeUUID is a UUID of a flow node
type NodeUUID utils.UUID

func (u NodeUUID) String() string { return string(u) }

// ExitUUID is the UUID of a node exit
type ExitUUID utils.UUID

func (u ExitUUID) String() string { return string(u) }

// FlowUUID is the UUID of a flow
type FlowUUID utils.UUID

func (u FlowUUID) String() string { return string(u) }

// ActionUUID is the UUID of an action
type ActionUUID utils.UUID

func (u ActionUUID) String() string { return string(u) }

// ContactUUID is the UUID of a contact
type ContactUUID utils.UUID

func (u ContactUUID) String() string { return string(u) }

// ChannelUUID is the UUID of a channel
type ChannelUUID utils.UUID

func (u ChannelUUID) String() string { return string(u) }

// RunUUID is the UUID of a flow run
type RunUUID utils.UUID

func (u RunUUID) String() string { return string(u) }

// StepUUID is the UUID of a run step
type StepUUID utils.UUID

func (u StepUUID) String() string { return string(u) }

// LabelUUID is the UUID of a label
type LabelUUID utils.UUID

func (u LabelUUID) String() string { return string(u) }

// GroupUUID is the UUID of a group
type GroupUUID utils.UUID

func (u GroupUUID) String() string { return string(u) }

// InputUUID is the UUID of an input
type InputUUID utils.UUID

func (u InputUUID) String() string { return string(u) }

// MsgUUID is the UUID of a message
type MsgUUID utils.UUID

func (u MsgUUID) String() string { return string(u) }

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

func (r SessionStatus) String() string { return string(r) }

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

func (r RunStatus) String() string { return string(r) }

// SessionAssets is the assets available to a session
type SessionAssets interface {
	GetChannel(ChannelUUID) (Channel, error)
	GetChannelSet() (*ChannelSet, error)

	GetField(FieldKey) (*Field, error)
	GetFieldSet() (*FieldSet, error)

	GetFlow(FlowUUID) (Flow, error)

	GetGroup(GroupUUID) (*Group, error)
	GetGroupSet() (*GroupSet, error)

	GetLabel(LabelUUID) (*Label, error)
	GetLabelSet() (*LabelSet, error)

	HasLocations() bool
	GetLocationHierarchy() (*utils.LocationHierarchy, error)
}

// Flow is a graph of nodes containing actions and routers
type Flow interface {
	UUID() FlowUUID
	Name() string
	Language() utils.Language
	ExpireAfterMinutes() int
	Localization() Localization

	Validate(SessionAssets) error
	Nodes() []Node
	GetNode(uuid NodeUUID) Node

	Reference() *FlowReference
}

// Node is a single node in a flow
type Node interface {
	UUID() NodeUUID

	Router() Router
	Actions() []Action
	Exits() []Exit
	Wait() Wait
}

// Action is an action within a flow node
type Action interface {
	UUID() ActionUUID

	Execute(FlowRun, Step, EventLog) error
	Validate(SessionAssets) error
	utils.Typed
}

type Router interface {
	PickRoute(FlowRun, []Exit, Step) (interface{}, Route, error)
	Validate([]Exit) error
	ResultName() string
	utils.Typed
}

type Route struct {
	exit  ExitUUID
	match string
}

func (r Route) Exit() ExitUUID { return r.exit }
func (r Route) Match() string  { return r.match }

var NoRoute = Route{}

func NewRoute(exit ExitUUID, match string) Route {
	return Route{exit, match}
}

type Exit interface {
	UUID() ExitUUID
	DestinationNodeUUID() NodeUUID
	Name() string
}

type Wait interface {
	utils.Typed

	Begin(FlowRun, Step)
	CanResume([]Event) bool
	HasTimedOut() bool

	Resume(FlowRun)
	ResumeByTimeOut(FlowRun)
}

// Localization provide a way to get the translations for a specific language
type Localization interface {
	AddItemTranslation(utils.Language, utils.UUID, string, []string)
	GetTranslations(utils.Language) Translations
	Languages() utils.LanguageList
}

// Translations provide a way to get the translation for a specific language for a uuid/key pair
type Translations interface {
	GetTextArray(uuid utils.UUID, key string) []string
}

type Trigger interface {
	utils.VariableResolver
	utils.Typed

	Environment() utils.Environment
	Flow() Flow
	Contact() *Contact
	Params() utils.JSONFragment
	TriggeredOn() time.Time
}

// EventOrigin is the allowed origin of an event
type EventOrigin int

const (
	// EventOriginCaller means an event can originate from the caller
	EventOriginCaller EventOrigin = 1

	// EventOriginEngine means an event can originate from the engine
	EventOriginEngine EventOrigin = 2
)

// Event describes a state change
type Event interface {
	CreatedOn() time.Time
	SetCreatedOn(time.Time)

	StepUUID() StepUUID
	SetStepUUID(StepUUID)

	FromCaller() bool
	SetFromCaller(bool)

	AllowedOrigin() EventOrigin
	Validate(SessionAssets) error

	Apply(FlowRun) error

	utils.Typed
}

// EventLog is the log of events the caller must apply after each call
type EventLog interface {
	Add(Event)
	Events() []Event
}

type Input interface {
	utils.VariableResolver
	utils.Typed

	UUID() InputUUID
	CreatedOn() time.Time
	Channel() Channel
}

type Step interface {
	UUID() StepUUID
	NodeUUID() NodeUUID
	ExitUUID() ExitUUID

	ArrivedOn() time.Time
	LeftOn() *time.Time

	Leave(ExitUUID)
}

// Session represents the session of a flow run which may contain many runs
type Session interface {
	Assets() SessionAssets

	Environment() utils.Environment
	SetEnvironment(utils.Environment)

	Contact() *Contact
	SetContact(*Contact)

	Status() SessionStatus
	Trigger() Trigger
	PushFlow(Flow, FlowRun)
	Wait() Wait
	FlowOnStack(FlowUUID) bool

	Start(Trigger, []Event) error
	Resume([]Event) error
	Runs() []FlowRun
	GetRun(RunUUID) (FlowRun, error)
	GetCurrentChild(FlowRun) FlowRun

	Events() []Event
	LogEvent(Event)
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

// FlowRun represents a run in the current session
type FlowRun interface {
	RunSummary

	Environment() utils.Environment
	Session() Session
	Context() utils.VariableResolver
	Input() Input
	Webhook() *utils.RequestResponse

	SetContact(*Contact)
	SetInput(Input)
	SetStatus(RunStatus)
	SetWebhook(*utils.RequestResponse)

	ApplyEvent(Step, Action, Event) error
	AddError(Step, Action, error)
	AddFatalError(Step, Action, error)

	CreateStep(Node) Step
	Path() []Step
	PathLocation() (Step, Node, error)

	GetText(utils.UUID, string, string) string
	GetTextArray(utils.UUID, string, []string) []string
	GetTranslatedTextArray(utils.UUID, string, []string, utils.LanguageList) []string

	Snapshot() RunSummary
	Parent() RunSummary
	SessionParent() FlowRun
	Ancestors() []FlowRun

	CreatedOn() time.Time
	ExpiresOn() *time.Time
	ResetExpiration(*time.Time)
	ExitedOn() *time.Time
	Exit(RunStatus)
}

// Channel represents a channel for sending and receiving messages
type Channel interface {
	UUID() ChannelUUID
	Name() string
	Address() string
	Schemes() []string
	SupportsScheme(string) bool
	Roles() []ChannelRole
	HasRole(ChannelRole) bool
	Reference() *ChannelReference
}
