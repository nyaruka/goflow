package events

func init() {
	registerType(TypeWarning, func() Event { return &Warning{} })
}

// TypeWarning is the type of our warning events
const TypeWarning string = "warning"

const (
	WarningCodeGraylistedURL = "url:graylisted"
)

// Warning events are created for things like accessing deprecated context values. Some warnings have
// a `code` which identifies the type of warning.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "warning",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "text": "deprecated context value accessed: legacy_extra"
//	}
//
// @event warning
type Warning struct {
	BaseEvent

	Text string `json:"text" validate:"required"`
	Code string `json:"code,omitempty"`
}

// NewWarning returns a new warning event
func NewWarning(text, code string) *Warning {
	return &Warning{
		BaseEvent: NewBaseEvent(TypeWarning),
		Text:      text,
		Code:      code,
	}
}
