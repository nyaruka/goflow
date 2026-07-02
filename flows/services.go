package flows

import (
	"context"
	"github.com/nyaruka/goflow/core"
	"github.com/nyaruka/goflow/core/events"
	"net/http"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/gocommon/urns"

	"github.com/shopspring/decimal"
)

// Services groups together interfaces for several services whose implementation is provided outside of the flow engine.
type Services interface {
	Email(SessionAssets) (EmailService, error)
	Webhook(SessionAssets) (WebhookService, error)
	LLM(*LLM) (LLMService, error)
	Airtime(SessionAssets) (AirtimeService, error)
}

// EmailService provides email functionality to the engine
type EmailService interface {
	Send(addresses []string, subject, body string) error
}

// WebhookService provides webhook functionality to the engine
type WebhookService interface {
	Call(request *http.Request) (*httpx.Trace, error)
}

// LLMService provides LLM functionality to the engine
type LLMService interface {
	Response(ctx context.Context, instructions, input string, maxTokens int) (*core.LLMResponse, error)
}

// AirtimeService provides airtime functionality to the engine
type AirtimeService interface {
	// Create initiates a new airtime transfer to the given URN. For providers with a two-step lifecycle,
	// this submits the transaction in an unconfirmed state and returns its identifier in ExternalID; the
	// host then calls Confirm to actually trigger the send. For providers that initiate immediately, this
	// method does the send and Confirm is a no-op.
	//
	// transferUUID is pre-allocated by the caller and is the UUID that will be assigned to the resulting
	// airtime_created event; implementations may pass it to the provider as their own reference (e.g.
	// DT One's external_id field) so that subsequent provider callbacks can be correlated back to the
	// event/transfer by the host without needing the provider's transaction id.
	Create(ctx context.Context, transferUUID events.EventUUID, sender urns.URN, recipient urns.URN, amounts map[string]decimal.Decimal, logHTTP core.HTTPLogCallback) (*core.AirtimeTransfer, error)

	// Confirm completes initiation of a transfer previously surfaced by Create. Hosts typically call
	// Confirm after the session commits, so the airtime is only actually sent once the surrounding work
	// is durably recorded. Implementations whose Create already triggers the send should make Confirm a
	// no-op.
	//
	// The transfer argument is the value Create returned. Implementations are expected to use whichever
	// fields they need (typically ExternalID) to address the provider transaction. Confirm should only be
	// invoked for transfers that Create returned without error; hosts should not call Confirm on a
	// partially-populated transfer left over from a failed Create.
	//
	// Confirm is not required to be idempotent — implementations may return an error on a duplicate
	// confirmation. Hosts that need at-most-once delivery semantics should ensure Confirm is called at
	// most once per transfer. On error, the airtime was not sent; hosts are responsible for surfacing
	// that to their users.
	Confirm(ctx context.Context, transfer *core.AirtimeTransfer, logHTTP core.HTTPLogCallback) error
}
