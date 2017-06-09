package flows

import (
	"time"

	"encoding/json"

	"github.com/nyaruka/goflow/utils"
)

type UUID string
type NodeUUID UUID
type ExitUUID UUID
type FlowUUID UUID
type ActionUUID UUID
type ContactUUID UUID
type FieldUUID UUID
type ChannelUUID UUID
type RunUUID UUID
type StepUUID UUID
type LabelUUID UUID
type GroupUUID UUID

type Flow interface {
	Name() string
	Language() utils.Language
	UUID() FlowUUID
	Translations() FlowTranslations

	Nodes() []Node
	GetNode(uuid NodeUUID) Node

	Validate() error

	CreateRun(env FlowEnvironment, contact *Contact, parent FlowRun) FlowRun
}

// RunStatus represents the current status of the flow run
type RunStatus string

const (
	// StatusActive represents an active flow run that is awaiting input
	StatusActive RunStatus = "A"

	// StatusCompleted represents a flow run that has run to completion
	StatusCompleted RunStatus = "C"

	// StatusErrored represents a flow run that encountered an error
	StatusErrored RunStatus = "E"

	// StatusExpired represents a flow run that expired due to inactivity
	StatusExpired RunStatus = "X"

	// StatusInterrupted represents a flow run that was interrupted by another flow
	StatusInterrupted RunStatus = "I"
)

func (r RunStatus) String() string { return string(r) }

type FlowEnvironment interface {
	GetFlow(FlowUUID) (Flow, error)
	GetRun(RunUUID) (FlowRun, error)
	GetContact(ContactUUID) (*Contact, error)
	utils.Environment
}

type Node interface {
	UUID() NodeUUID

	Router() Router
	Actions() []Action
	Exits() []Exit
	Wait() Wait
}

type Action interface {
	Execute(FlowRun, Step) error
	Validate() error
	utils.Typed
}

type Router interface {
	PickRoute(FlowRun, []Exit, Step) (Route, error)
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
	Begin(FlowRun, Step) error
	GetEndEvent(FlowRun, Step) (Event, error)
	End(FlowRun, Step, Event) error
	utils.Typed
	utils.VariableResolver
}

// FlowTranslations provide a way to get the Translations for a flow for a specific language
type FlowTranslations interface {
	GetTranslations(utils.Language) Translations
}

// Translations provide a way to get the translation for a specific language for a uuid/key pair
type Translations interface {
	GetText(uuid UUID, key string, backdown string) string
}

type Context interface {
	utils.VariableResolver
	Contact() *Contact
	Run() FlowRun
}

type Event interface {
	CreatedOn() *time.Time
	SetCreatedOn(time.Time)

	Step() StepUUID
	SetStep(StepUUID)

	utils.Typed
}

type Input interface {
	Event
	utils.VariableResolver
}

type Step interface {
	utils.VariableResolver

	UUID() StepUUID
	NodeUUID() NodeUUID
	ExitUUID() ExitUUID

	ArrivedOn() time.Time
	LeftOn() *time.Time

	Leave(ExitUUID)

	Events() []Event
}

// Session represents the session of a flow run which may contain many runs
type Session interface {
	Runs() []FlowRun
	AddRun(FlowRun)

	ActiveRun() FlowRun

	Events() []Event
	AddEvent(Event)
	ClearEvents()
}

// FlowRun represents a single run on a flow by a single contact
type FlowRun interface {
	UUID() RunUUID
	FlowUUID() FlowUUID
	Flow() Flow

	Hydrate(FlowEnvironment) error

	ContactUUID() ContactUUID
	Contact() *Contact

	ChannelUUID() ChannelUUID
	Channel() *Channel
	SetChannel(*Channel)

	Context() Context
	Results() *Results
	Environment() FlowEnvironment

	SetExtra(json.RawMessage)
	Extra() utils.JSONFragment

	Session() Session
	SetSession(Session)
	ResetSession()

	Status() RunStatus
	Exit(RunStatus)
	IsComplete() bool

	Wait() Wait
	SetWait(Wait)

	Input() Input
	SetInput(Input)

	AddEvent(Step, Event)
	AddError(Step, error)

	CreateStep(Node) Step
	Path() []Step

	SetLanguage(utils.Language)
	SetFlowTranslations(FlowTranslations)
	GetText(uuid UUID, key string, backdown string) string

	Webhook() utils.RequestResponse
	SetWebhook(utils.RequestResponse)

	Child() FlowRunReference
	Parent() FlowRunReference

	CreatedOn() time.Time
	ModifiedOn() time.Time
	ExpiresOn() *time.Time
	TimesOutOn() *time.Time
	ExitedOn() *time.Time
}

// FlowRunReference represents a flow run reference within a flow
type FlowRunReference interface {
	UUID() RunUUID
	FlowUUID() FlowUUID
	ContactUUID() ContactUUID
	ChannelUUID() ChannelUUID

	Results() *Results
	Status() RunStatus

	CreatedOn() time.Time
	ModifiedOn() time.Time
	ExpiresOn() *time.Time
	TimesOutOn() *time.Time
	ExitedOn() *time.Time
}

// ChannelType represents the type of a Channel
type ChannelType string

func (ct ChannelType) String() string { return string(ct) }

// MsgDirection is the direction of a Msg (either in or out)
type MsgDirection string

const (
	// MsgOut represents an outgoing message
	MsgOut MsgDirection = "O"

	// MsgIn represents an incoming message
	MsgIn MsgDirection = "I"
)
