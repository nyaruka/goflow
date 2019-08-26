package engine

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type services struct {
	airtime flows.AirtimeService
}

func (s *services) Airtime() flows.AirtimeService {
	return s.airtime
}

type nilAirtimeService struct{}

// Transfer in this case is a failure
func (s *nilAirtimeService) Transfer(session flows.Session, from urns.URN, to urns.URN, amounts map[string]decimal.Decimal) (*flows.AirtimeTransfer, error) {
	return nil, errors.New("no airtime service available")
}
