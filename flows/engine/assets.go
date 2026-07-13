package engine

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/contactql"
	"github.com/nyaruka/goflow/contactql/parse"
	"github.com/nyaruka/goflow/core"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/definition/migrations"
)

// our implementation of SessionAssets - the high-level API for asset access from the engine
type sessionAssets struct {
	source assets.Source

	campaigns *core.CampaignAssets
	channels  *core.ChannelAssets
	fields    *core.FieldAssets
	flows     flows.FlowAssets
	globals   *core.GlobalAssets
	groups    *core.GroupAssets
	labels    *core.LabelAssets
	llms      *core.LLMAssets
	locations *core.LocationAssets
	optIns    *core.OptInAssets
	resthooks *core.ResthookAssets
	templates *core.TemplateAssets
	topics    *core.TopicAssets
	users     *core.UserAssets
}

var _ flows.SessionAssets = (*sessionAssets)(nil)

// NewSessionAssets creates a new session assets instance with the provided base URLs
func NewSessionAssets(env envs.Environment, source assets.Source, migrationConfig *migrations.Config) (flows.SessionAssets, error) {
	campaigns, err := source.Campaigns()
	if err != nil {
		return nil, err
	}
	channels, err := source.Channels()
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
	llms, err := source.LLMs()
	if err != nil {
		return nil, err
	}
	locations, err := source.Locations()
	if err != nil {
		return nil, err
	}
	optIns, err := source.OptIns()
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
	topics, err := source.Topics()
	if err != nil {
		return nil, err
	}
	users, err := source.Users()
	if err != nil {
		return nil, err
	}

	fieldAssets := core.NewFieldAssets(fields)

	// parse queries of any query based groups, skipping those which fail to parse
	parsedGroups := make([]*core.Group, 0, len(groups))
	for _, asset := range groups {
		var query *contactql.ContactQuery
		if asset.Query() != "" {
			var err error
			if query, err = parse.Query(env, asset.Query(), fieldAssets); err != nil {
				continue
			}
		}
		parsedGroups = append(parsedGroups, core.NewGroup(asset, query))
	}
	groupAssets := core.NewGroupAssets(parsedGroups)

	return &sessionAssets{
		source:    source,
		campaigns: core.NewCampaignAssets(campaigns),
		channels:  core.NewChannelAssets(channels),
		fields:    fieldAssets,
		flows:     definition.NewFlowAssets(source, migrationConfig),
		globals:   core.NewGlobalAssets(globals),
		groups:    groupAssets,
		labels:    core.NewLabelAssets(labels),
		llms:      core.NewLLMAssets(llms),
		locations: core.NewLocationAssets(locations),
		optIns:    core.NewOptInAssets(optIns),
		resthooks: core.NewResthookAssets(resthooks),
		templates: core.NewTemplateAssets(templates),
		topics:    core.NewTopicAssets(topics),
		users:     core.NewUserAssets(users),
	}, nil
}

func (s *sessionAssets) Source() assets.Source           { return s.source }
func (s *sessionAssets) Campaigns() *core.CampaignAssets { return s.campaigns }
func (s *sessionAssets) Channels() *core.ChannelAssets   { return s.channels }
func (s *sessionAssets) Fields() *core.FieldAssets       { return s.fields }
func (s *sessionAssets) Flows() flows.FlowAssets         { return s.flows }
func (s *sessionAssets) Globals() *core.GlobalAssets     { return s.globals }
func (s *sessionAssets) Groups() *core.GroupAssets       { return s.groups }
func (s *sessionAssets) Labels() *core.LabelAssets       { return s.labels }
func (s *sessionAssets) LLMs() *core.LLMAssets           { return s.llms }
func (s *sessionAssets) Locations() *core.LocationAssets { return s.locations }
func (s *sessionAssets) OptIns() *core.OptInAssets       { return s.optIns }
func (s *sessionAssets) Resthooks() *core.ResthookAssets { return s.resthooks }
func (s *sessionAssets) Templates() *core.TemplateAssets { return s.templates }
func (s *sessionAssets) Topics() *core.TopicAssets       { return s.topics }
func (s *sessionAssets) Users() *core.UserAssets         { return s.users }

// Resolver methods used by contactql

func (s *sessionAssets) ResolveField(key string) assets.Field {
	f := s.Fields().Get(key)
	if f == nil {
		return nil
	}
	return f.Asset()
}

func (s *sessionAssets) ResolveGroup(name string) assets.Group {
	g := s.Groups().FindByName(name)
	if g == nil {
		return nil
	}
	return g.Asset()
}

func (s *sessionAssets) ResolveFlow(name string) assets.Flow {
	f, _ := s.Flows().FindByName(name)
	if f == nil {
		return nil
	}
	return f.Asset()
}
