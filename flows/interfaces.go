package flows

import (
	"time"

	"github.com/nyaruka/goflow/utils"
)

type UUID string
type NodeUUID UUID
type ExitUUID UUID
type FlowUUID UUID
type ActionUUID UUID
type ContactUUID UUID
type GroupUUID UUID
type FieldUUID UUID
type ChannelUUID UUID
type RunUUID UUID

type Language string

type Flow interface {
	Name() string
	Language() Language
	UUID() FlowUUID
	Translations() FlowTranslations

	Nodes() []Node
	GetNode(uuid NodeUUID) Node

	Validate() error

	CreateRun(env FlowEnvironment, contact Contact, parent FlowRun) FlowRun
}

// RunStatus represents the current status of the flow run
type RunStatus string

const (
	// RunActive represents an active flow run that is awaiting input
	RunActive RunStatus = "A"

	// RunCompleted represents a flow run that has run to completion
	RunCompleted RunStatus = "C"

	// RunExpired represents a flow run that expired due to inactivity
	RunExpired RunStatus = "E"

	// RunInterrupted represents a flow run that was interrupted by another flow
	RunInterrupted RunStatus = "I"
)

func (r RunStatus) String() string { return string(r) }

type FlowEnvironment interface {
	GetFlow(FlowUUID) (Flow, error)
	GetRun(RunUUID) (FlowRun, error)
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
	Name() string
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
	Destination() NodeUUID
	Name() string
}

type Wait interface {
	Begin(FlowRun, Step) error
	ShouldEnd(FlowRun, Step) (Event, error)
	End(FlowRun, Step, Event) error
	utils.Typed
	utils.VariableResolver
}

// FlowTranslations provide a way to get the Translations for a flow for a specific language
type FlowTranslations interface {
	GetTranslations(Language) Translations
}

// Translations provide a way to get the translation for a specific language for a uuid/key pair
type Translations interface {
	GetText(uuid UUID, key string, backdown string) string
}

type Context interface {
	utils.VariableResolver
	Contact() Contact
	Run() FlowRun
}

type Input interface {
	utils.VariableResolver
	utils.Typed
}

type Event interface {
	CreatedOn() *time.Time
	SetCreatedOn(time.Time)

	Run() RunUUID
	SetRun(RunUUID)

	utils.Typed
}

type Step interface {
	utils.VariableResolver

	Node() NodeUUID
	Exit() ExitUUID

	ArrivedOn() time.Time
	LeftOn() *time.Time

	Leave(ExitUUID)

	Events() []Event
}

type Contact interface {
	utils.VariableResolver

	UUID() ContactUUID

	Language() Language
	SetLanguage(Language)

	Groups() GroupList
	AddGroup(uuid GroupUUID, name string)
	RemoveGroup(uuid GroupUUID) bool

	Fields() Fields

	URNs() URNList
}

type Group interface {
	utils.VariableResolver

	UUID() GroupUUID
	Name() string
}

type GroupList interface {
	FindGroup(GroupUUID) Group
	utils.VariableResolver
}

type Results interface {
	Save(node NodeUUID, name string, value string, category string, createdOn time.Time)
	utils.VariableResolver
}

type Result interface {
	utils.VariableResolver
}

type Fields interface {
	Save(uuid FieldUUID, name string, value string, createdOn time.Time)
	utils.VariableResolver
}

type Field interface {
	utils.VariableResolver
}

// RunOutput represents the output of a Run in its last execution cycle
type RunOutput interface {
	Runs() []FlowRun
	AddRun(FlowRun)

	ActiveRun() FlowRun

	Events() []Event
	AddEvent(Event)
}

// FlowRun represents a single run on a flow by a single contact
type FlowRun interface {
	UUID() RunUUID
	FlowUUID() FlowUUID
	Flow() Flow
	SetFlow(Flow) error

	ContactUUID() ContactUUID
	Contact() Contact
	SetContact(Contact) error

	ChannelUUID() ChannelUUID
	Channel() Channel
	SetChannel(Channel)

	Context() Context
	Results() Results
	Environment() FlowEnvironment

	Output() RunOutput
	ResetOutput()

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

	SetLanguage(Language)
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

	Results() Results
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

type Channel interface {
	UUID() ChannelUUID
	Name() string
	Type() ChannelType
}
