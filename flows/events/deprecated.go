package events

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeClassifierCalled, func() flows.Event { return &ClassifierCalled{} })
	registerType(TypeContactRefreshed, func() flows.Event { return &ContactRefreshed{} })
	registerType(TypeEmailCreated, func() flows.Event { return &EmailCreated{} })
	registerType(TypeEnvironmentRefreshed, func() flows.Event { return &EnvironmentRefreshed{} })
	registerType(TypeRunExpired, func() flows.Event { return &RunExpired{} })
}

// TypeClassifierCalled is our type for the classification event
const TypeClassifierCalled string = "classifier_called"

// ClassifierCalled events have been replaced by service_called.
type ClassifierCalled struct {
	BaseEvent

	Classifier *assets.ClassifierReference `json:"classifier" validate:"required"`
	HTTPLogs   []*flows.HTTPLog            `json:"http_logs"`
}

// TypeContactRefreshed is the type of our contact refreshed event
const TypeContactRefreshed string = "contact_refreshed"

// ContactRefreshed events are generated when the resume has a contact with differences to the current session contact.
type ContactRefreshed struct {
	BaseEvent

	Contact json.RawMessage `json:"contact"`
}

// TypeEnvironmentRefreshed is the type of our environment changed event
const TypeEnvironmentRefreshed string = "environment_refreshed"

// EnvironmentRefreshed events are sent by the caller to tell the engine to update the session environment.
type EnvironmentRefreshed struct {
	BaseEvent

	Environment json.RawMessage `json:"environment"`
}

// TypeEmailCreated is our type for the email event
const TypeEmailCreated string = "email_created"

// EmailCreated is no longer used but old sessions might include these
type EmailCreated struct {
	BaseEvent

	Addresses []string `json:"addresses" validate:"required,min=1"`
	Subject   string   `json:"subject" validate:"required"`
	Body      string   `json:"body"`
}

// TypeRunExpired is the type of our flow expired event
const TypeRunExpired string = "run_expired"

// RunExpired events are sent by the caller to tell the engine that a run has expired.
type RunExpired struct {
	BaseEvent

	RunUUID flows.RunUUID `json:"run_uuid"    validate:"required,uuid"`
}

// NewRunExpired creates a new run expired event
func NewRunExpired(run flows.Run) *RunExpired {
	return &RunExpired{
		BaseEvent: NewBaseEvent(TypeRunExpired),
		RunUUID:   run.UUID(),
	}
}

var _ flows.Event = (*RunExpired)(nil)
