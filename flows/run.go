package flows

import (
	"time"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/types"
)

// RunSummary represents the minimum information available about all runs (current or related) and is the
// representation of runs made accessible to router tests.
type RunSummary interface {
	UUID() RunUUID
	Contact() *Contact
	Flow() Flow
	Status() RunStatus
	Results() Results
}

// Run is a single contact's journey through a flow. It records the path they have taken,
// and the results that have been collected.
type Run interface {
	Contextable
	RunSummary
	FlowReference() *assets.FlowReference

	Session() Session
	Locals() *Locals
	SetResult(*Result) (*Result, bool)
	Webhook() *WebhookCall
	SetWebhook(*WebhookCall)

	CreateStep(Node) Step
	Path() []Step
	PathLocation() (Step, Node, error)
	HadInput() bool

	EvaluateTemplateValue(string, EventLogger) (types.XValue, bool)
	EvaluateTemplateText(string, excellent.Escaping, bool, EventLogger) (string, bool)
	EvaluateTemplate(string, EventLogger) (string, bool)
	RootContext(envs.Environment) map[string]types.XValue

	GetText(uuids.UUID, string, string) (string, i18n.Language)
	GetTextArray(uuids.UUID, string, []string, []i18n.Language) ([]string, i18n.Language)

	Snapshot() RunSummary
	Parent() RunSummary
	Ancestors() []Run

	CreatedOn() time.Time
	ModifiedOn() time.Time
	ExitedOn() *time.Time
	Exit(RunStatus)
}
