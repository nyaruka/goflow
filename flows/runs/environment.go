package runs

import (
	"time"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
)

// overrides some environment values to use values from the contact
type runEnvironment struct {
	envs.Environment

	run *flowRun
}

// creates a run environment based on the given run
func newRunEnvironment(base envs.Environment, run *flowRun) *runEnvironment {
	return &runEnvironment{
		flows.NewEnvironment(base, run.Session().Assets().Locations()),
		run,
	}
}

func (e *runEnvironment) Timezone() *time.Location {
	contact := e.run.Contact()

	// if run has a contact with a timezone, that overrides the enviroment's timezone
	if contact != nil && contact.Timezone() != nil {
		return contact.Timezone()
	}
	return e.run.Session().Environment().Timezone()
}

func (e *runEnvironment) DefaultCountry() envs.Country {
	contact := e.run.Contact()

	// if run has a contact with a preferred channel with a country, that overrides the environment's country
	if contact != nil {
		cc := contact.Country()
		if cc != envs.NilCountry {
			return cc
		}
	}
	return e.run.Session().Environment().DefaultCountry()
}
