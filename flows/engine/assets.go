package engine

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/definition/migrations"
)

// our implementation of SessionAssets - the high-level API for asset access from the engine
type sessionAssets struct {
	source assets.Source

	channels    *flows.ChannelAssets
	classifiers *flows.ClassifierAssets
	fields      *flows.FieldAssets
	flows       flows.FlowAssets
	globals     *flows.GlobalAssets
	groups      *flows.GroupAssets
	labels      *flows.LabelAssets
	locations   *flows.LocationAssets
	resthooks   *flows.ResthookAssets
	templates   *flows.TemplateAssets
	ticketers   *flows.TicketerAssets
	users       *flows.UserAssets
}

var _ flows.SessionAssets = (*sessionAssets)(nil)

// NewSessionAssets creates a new session assets instance with the provided base URLs
func NewSessionAssets(env envs.Environment, source assets.Source, migrationConfig *migrations.Config) (flows.SessionAssets, error) {
	channels, err := source.Channels()
	if err != nil {
		return nil, err
	}
	classifiers, err := source.Classifiers()
	if err != nil {
		return nil, err
	}
	fields, err := source.Fields()
	if err != nil {
		return nil, err
	}
	globals, err := source.Globals()
	if err != nil {
		return nil, err
	}
	groups, err := source.Groups()
	if err != nil {
		return nil, err
	}
	labels, err := source.Labels()
	if err != nil {
		return nil, err
	}
	locations, err := source.Locations()
	if err != nil {
		return nil, err
	}
	resthooks, err := source.Resthooks()
	if err != nil {
		return nil, err
	}
	templates, err := source.Templates()
	if err != nil {
		return nil, err
	}
	ticketers, err := source.Ticketers()
	if err != nil {
		return nil, err
	}
	users, err := source.Users()
	if err != nil {
		return nil, err
	}

	fieldAssets := flows.NewFieldAssets(fields)
	groupAssets, _ := flows.NewGroupAssets(env, fieldAssets, groups)

	return &sessionAssets{
		source:      source,
		channels:    flows.NewChannelAssets(channels),
		classifiers: flows.NewClassifierAssets(classifiers),
		fields:      fieldAssets,
		flows:       definition.NewFlowAssets(source, migrationConfig),
		globals:     flows.NewGlobalAssets(globals),
		groups:      groupAssets,
		labels:      flows.NewLabelAssets(labels),
		locations:   flows.NewLocationAssets(locations),
		resthooks:   flows.NewResthookAssets(resthooks),
		templates:   flows.NewTemplateAssets(templates),
		ticketers:   flows.NewTicketerAssets(ticketers),
		users:       flows.NewUserAssets(users),
	}, nil
}

func (s *sessionAssets) Source() assets.Source                { return s.source }
func (s *sessionAssets) Channels() *flows.ChannelAssets       { return s.channels }
func (s *sessionAssets) Classifiers() *flows.ClassifierAssets { return s.classifiers }
func (s *sessionAssets) Fields() *flows.FieldAssets           { return s.fields }
func (s *sessionAssets) Flows() flows.FlowAssets              { return s.flows }
func (s *sessionAssets) Globals() *flows.GlobalAssets         { return s.globals }
func (s *sessionAssets) Groups() *flows.GroupAssets           { return s.groups }
func (s *sessionAssets) Labels() *flows.LabelAssets           { return s.labels }
func (s *sessionAssets) Locations() *flows.LocationAssets     { return s.locations }
func (s *sessionAssets) Resthooks() *flows.ResthookAssets     { return s.resthooks }
func (s *sessionAssets) Templates() *flows.TemplateAssets     { return s.templates }
func (s *sessionAssets) Ticketers() *flows.TicketerAssets     { return s.ticketers }
func (s *sessionAssets) Users() *flows.UserAssets             { return s.users }

func (s *sessionAssets) ResolveField(key string) assets.Field {
	f := s.Fields().Get(key)
	if f == nil {
		return nil
	}
	return f
}
func (s *sessionAssets) ResolveGroup(name string) assets.Group {
	g := s.Groups().FindByName(name)
	if g == nil {
		return nil
	}
	return g
}
