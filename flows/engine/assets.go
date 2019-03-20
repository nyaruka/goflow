package engine

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition"

	"github.com/pkg/errors"
)

// our implementation of SessionAssets - the high-level API for asset access from the engine
type sessionAssets struct {
	source assets.AssetSource

	channels  *flows.ChannelAssets
	fields    *flows.FieldAssets
	flows     flows.FlowAssets
	groups    *flows.GroupAssets
	labels    *flows.LabelAssets
	locations *flows.LocationAssets
	resthooks *flows.ResthookAssets
}

var _ flows.SessionAssets = (*sessionAssets)(nil)

// NewSessionAssets creates a new session assets instance with the provided base URLs
func NewSessionAssets(source assets.AssetSource) (flows.SessionAssets, error) {
	channels, err := source.Channels()
	if err != nil {
		return nil, err
	}
	fields, err := source.Fields()
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

	return &sessionAssets{
		source:    source,
		channels:  flows.NewChannelAssets(channels),
		fields:    flows.NewFieldAssets(fields),
		flows:     definition.NewFlowAssets(source),
		groups:    flows.NewGroupAssets(groups),
		labels:    flows.NewLabelAssets(labels),
		locations: flows.NewLocationAssets(locations),
		resthooks: flows.NewResthookAssets(resthooks),
	}, nil
}

func (s *sessionAssets) Channels() *flows.ChannelAssets   { return s.channels }
func (s *sessionAssets) Fields() *flows.FieldAssets       { return s.fields }
func (s *sessionAssets) Flows() flows.FlowAssets          { return s.flows }
func (s *sessionAssets) Groups() *flows.GroupAssets       { return s.groups }
func (s *sessionAssets) Labels() *flows.LabelAssets       { return s.labels }
func (s *sessionAssets) Locations() *flows.LocationAssets { return s.locations }
func (s *sessionAssets) Resthooks() *flows.ResthookAssets { return s.resthooks }

// Validate ensures that the given flow exists, is correct, and all its dependent flows are also valid
func (s *sessionAssets) Validate(flowUUID assets.FlowUUID) ([]assets.Reference, error) {
	flow, err := s.Flows().Get(flowUUID)
	if err != nil {
		return nil, err
	}

	missingAssets := make([]assets.Reference, 0)
	missing := func(r assets.Reference) {
		missingAssets = append(missingAssets, r)
	}

	if err := flow.ValidateRecursively(s, missing); err != nil {
		return nil, errors.Wrapf(err, "validation failed for flow[uuid=%s]", flow.UUID())
	}
	return missingAssets, nil
}
