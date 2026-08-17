package events

func init() {
	registerType(TypeWarning, func() Event { return &Warning{} })
}

// TypeWarning is the type of our warning events
const TypeWarning string = "warning"

const (
	WarningCodeDeprecatedContext   = "context:deprecated"
	WarningCodeWebhookResponseSize = "webhook:response_size"
)

// Warning events are created for things like accessing deprecated context values. Some warnings have
// a `code` which identifies the type of warning, and `extra` values with more details.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "warning",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "text": "@legacy_extra is deprecated",
//	  "code": "context:deprecated"
//	}
//
// @event warning
type Warning struct {
	BaseEvent

	Text  string            `json:"text" validate:"required"`
	Code  string            `json:"code,omitempty"`
	Extra map[string]string `json:"extra,omitempty"`
}

// NewWarning returns a new warning event
func NewWarning(text, code string, extra ...string) *Warning {
	return &Warning{
		BaseEvent: NewBaseEvent(TypeWarning),
		Text:      text,
		Code:      code,
		Extra:     extraFromPairs(extra),
	}
}

// converts a flat list of key/value pairs into a map for an event's extra
func extraFromPairs(pairs []string) map[string]string {
	if len(pairs)%2 != 0 {
		panic("extra fields must be key/value pairs")
	}
	if len(pairs) == 0 {
		return nil
	}
	extra := make(map[string]string, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		extra[pairs[i]] = pairs[i+1]
	}
	return extra
}
