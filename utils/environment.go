package utils

import (
	"encoding/json"
	"time"
)

// Environment defines the Environment that the Excellent function is running in, this includes
// the timezone the user is in as well as the preferred date and time formats.
type Environment interface {
	DateFormat() DateFormat
	SetDateFormat(DateFormat)

	TimeFormat() TimeFormat
	SetTimeFormat(TimeFormat)

	Timezone() *time.Location
	SetTimezone(*time.Location)
}

// NewDefaultEnvironment creates a new Environment with our usual defaults in the UTC timezone
func NewDefaultEnvironment() Environment {
	return &environment{YYYY_MM_DD, HH_MM, time.UTC}
}

// NewEnvironment creates a new Environment with the passed in date and time formats and timezone
func NewEnvironment(dateFormat DateFormat, timeFormat TimeFormat, timezone *time.Location) Environment {
	return &environment{dateFormat, timeFormat, timezone}
}

type environment struct {
	dateFormat DateFormat
	timeFormat TimeFormat
	timezone   *time.Location
}

func (e *environment) DateFormat() DateFormat              { return e.dateFormat }
func (e *environment) SetDateFormat(dateFormat DateFormat) { e.dateFormat = dateFormat }

func (e *environment) TimeFormat() TimeFormat              { return e.timeFormat }
func (e *environment) SetTimeFormat(timeFormat TimeFormat) { e.timeFormat = timeFormat }

func (e *environment) Timezone() *time.Location            { return e.timezone }
func (e *environment) SetTimezone(timezone *time.Location) { e.timezone = timezone }

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type envEnvelope struct {
	DateFormat DateFormat     `json:"date_format"`
	TimeFormat TimeFormat     `json:"time_format"`
	Timezone   *time.Location `json:"timezone"`
}

func (e *environment) UnmarshalJSON(data []byte) error {
	var envelope envEnvelope
	var err error

	err = json.Unmarshal(data, &envelope)
	if err != nil {
		return err
	}

	e.dateFormat = envelope.DateFormat
	e.timeFormat = envelope.TimeFormat
	e.timezone = envelope.Timezone
	return nil
}

func (e *environment) MarshalJSON() ([]byte, error) {
	ee := envEnvelope{e.dateFormat, e.timeFormat, e.timezone}
	return json.Marshal(ee)
}
