package flows

import (
	"github.com/nyaruka/gocommon/urns"

	"github.com/shopspring/decimal"
)

// Services groups together interfaces for several services whose implementation existsi outside of the flow engine.
type Services interface {
	Airtime() AirtimeService
}

// AirtimeTransferStatus is a status of a airtime transfer
type AirtimeTransferStatus string

// possible values for airtime transfer statuses
const (
	AirtimeTransferStatusSuccess AirtimeTransferStatus = "success"
	AirtimeTransferStatusFailed  AirtimeTransferStatus = "failed"
)

// AirtimeTransfer is the result of an attempted airtime transfer
type AirtimeTransfer struct {
	Currency      string
	DesiredAmount decimal.Decimal
	ActualAmount  decimal.Decimal
	Status        AirtimeTransferStatus
}

// AirtimeService is the interface for an airtime transfer service
type AirtimeService interface {
	// Transfer transfers airtime to the given URN
	Transfer(session Session, from urns.URN, to urns.URN, amounts map[string]decimal.Decimal) (*AirtimeTransfer, error)
}
