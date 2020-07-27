package events

import (
	"encoding/json"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeSessionTriggered, func() flows.Event { return &SessionTriggeredEvent{} })
}

// TypeSessionTriggered is the type of our session triggered event
const TypeSessionTriggered string = "session_triggered"

// SessionTriggeredEvent events are created when an action wants to start other people in a flow.
//
//   {
//     "type": "session_triggered",
//     "created_on": "2006-01-02T15:04:05Z",
//     "flow": {"uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a", "name": "Registration"},
//     "groups": [
//       {"uuid": "8f8e2cae-3c8d-4dce-9c4b-19514437e427", "name": "New contacts"}
//     ],
//     "run_summary": {
//       "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
//       "flow": {"uuid": "93c554a1-b90d-4892-b029-a2a87dec9b87", "name": "Other Flow"},
//       "contact": {
//         "uuid": "c59b0033-e748-4240-9d4c-e85eb6800151",
//         "name": "Bob",
//         "fields": {"state": {"value": "Azuay", "created_on": "2000-01-01T00:00:00.000000000-00:00"}}
//       },
//       "results": {
//         "age": {
//           "name": "Age",
//           "value": "33",
//           "node_uuid": "cd2be8c4-59bc-453c-8777-dec9a80043b8",
//           "created_on": "2000-01-01T00:00:00.000000000-00:00"
//         }
//       }
//     },
//     "history": {
//       "parent_uuid": "55105da5-abb5-4690-b1f6-ec2e5762a561",
//       "ancestors": 3,
//       "ancestors_since_input": 1
//     }
//   }
//
// @event session_triggered
type SessionTriggeredEvent struct {
	baseEvent

	Flow          *assets.FlowReference     `json:"flow" validate:"required"`
	Groups        []*assets.GroupReference  `json:"groups,omitempty" validate:"dive"`
	Contacts      []*flows.ContactReference `json:"contacts,omitempty" validate:"dive"`
	ContactQuery  string                    `json:"contact_query,omitempty"`
	CreateContact bool                      `json:"create_contact,omitempty"`
	URNs          []urns.URN                `json:"urns,omitempty" validate:"dive,urn"`
	RunSummary    json.RawMessage           `json:"run_summary"`
	History       *flows.SessionHistory     `json:"history"`
}

// NewSessionTriggered returns a new session triggered event
func NewSessionTriggered(flow *assets.FlowReference, groups []*assets.GroupReference, contacts []*flows.ContactReference, contactQuery string, createContact bool, urns []urns.URN, runSummary json.RawMessage, history *flows.SessionHistory) *SessionTriggeredEvent {
	return &SessionTriggeredEvent{
		baseEvent:     newBaseEvent(TypeSessionTriggered),
		Flow:          flow,
		Groups:        groups,
		Contacts:      contacts,
		ContactQuery:  contactQuery,
		CreateContact: createContact,
		URNs:          urns,
		RunSummary:    runSummary,
		History:       history,
	}
}
