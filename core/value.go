package core

import (
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
)

// Value represents a value in each of the field types
type Value struct {
	Text     *types.XText      `json:"text" validate:"required"`
	Datetime *types.XDateTime  `json:"datetime,omitempty"`
	Number   *types.XNumber    `json:"number,omitempty"`
	State    envs.LocationPath `json:"state,omitempty"`
	District envs.LocationPath `json:"district,omitempty"`
	Ward     envs.LocationPath `json:"ward,omitempty"`
}

// NewValue creates an empty value
func NewValue(text *types.XText, datetime *types.XDateTime, number *types.XNumber, state envs.LocationPath, district envs.LocationPath, ward envs.LocationPath) *Value {
	return &Value{
		Text:     text,
		Datetime: datetime,
		Number:   number,
		State:    state,
		District: district,
		Ward:     ward,
	}
}

// Equals determines whether two values are equal
func (v *Value) Equals(o *Value) bool {
	if v == nil && o == nil {
		return true
	}
	if (v == nil && o != nil) || (v != nil && o == nil) {
		return false
	}

	dateEqual := (v.Datetime == nil && o.Datetime == nil) || (v.Datetime != nil && o.Datetime != nil && v.Datetime.Equals(o.Datetime))
	numEqual := (v.Number == nil && o.Number == nil) || (v.Number != nil && o.Number != nil && v.Number.Equals(o.Number))

	return v.Text.Equals(o.Text) && dateEqual && numEqual && v.State == o.State && v.District == o.District && v.Ward == o.Ward
}
