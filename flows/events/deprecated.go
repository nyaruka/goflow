package events

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeClassifierCalled, func() flows.Event { return &ClassifierCalledEvent{} })
	registerType(TypeContactRefreshed, func() flows.Event { return &ContactRefreshedEvent{} })
	registerType(TypeEnvironmentRefreshed, func() flows.Event { return &EnvironmentRefreshedEvent{} })
}

// TypeClassifierCalled is our type for the classification event
const TypeClassifierCalled string = "classifier_called"

// ClassifierCalledEvent events have been replaced by service_called.
type ClassifierCalledEvent struct {
	BaseEvent

	Classifier *assets.ClassifierReference `json:"classifier" validate:"required"`
	HTTPLogs   []*flows.HTTPLog            `json:"http_logs"`
}

// TypeContactRefreshed is the type of our contact refreshed event
const TypeContactRefreshed string = "contact_refreshed"

// ContactRefreshedEvent events are generated when the resume has a contact with differences to the current session contact.
//
//	{
//	  "type": "contact_refreshed",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "contact": {
//	    "uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a",
//	    "name": "Bob",
//	    "urns": ["tel:+11231234567"]
//	  }
//	}
//
// @event contact_refreshed
type ContactRefreshedEvent struct {
	BaseEvent

	Contact json.RawMessage `json:"contact"`
}

// TypeEnvironmentRefreshed is the type of our environment changed event
const TypeEnvironmentRefreshed string = "environment_refreshed"

// EnvironmentRefreshedEvent events are sent by the caller to tell the engine to update the session environment.
//
//	{
//	  "type": "environment_refreshed",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "environment": {
//	    "date_format": "YYYY-MM-DD",
//	    "time_format": "hh:mm",
//	    "timezone": "Africa/Kigali",
//	    "allowed_languages": ["eng", "fra"]
//	  }
//	}
//
// @event environment_refreshed
type EnvironmentRefreshedEvent struct {
	BaseEvent

	Environment json.RawMessage `json:"environment"`
}
