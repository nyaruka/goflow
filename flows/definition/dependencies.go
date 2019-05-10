package definition

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/pkg/errors"
)

type dependencies struct {
	Channels  []*assets.ChannelReference  `json:"channels,omitempty"`
	Contacts  []*flows.ContactReference   `json:"contacts,omitempty"`
	Fields    []*assets.FieldReference    `json:"fields,omitempty"`
	Flows     []*assets.FlowReference     `json:"flows,omitempty"`
	Groups    []*assets.GroupReference    `json:"groups,omitempty"`
	Labels    []*assets.LabelReference    `json:"labels,omitempty"`
	Templates []*assets.TemplateReference `json:"templates,omitempty"`
}

func newDependencies(refs []assets.Reference) *dependencies {
	d := &dependencies{}
	for _, r := range refs {
		switch typed := r.(type) {
		case *assets.ChannelReference:
			d.Channels = append(d.Channels, typed)
		case *flows.ContactReference:
			d.Contacts = append(d.Contacts, typed)
		case *assets.FieldReference:
			d.Fields = append(d.Fields, typed)
		case *assets.FlowReference:
			d.Flows = append(d.Flows, typed)
		case *assets.GroupReference:
			d.Groups = append(d.Groups, typed)
		case *assets.LabelReference:
			d.Labels = append(d.Labels, typed)
		case *assets.TemplateReference:
			d.Templates = append(d.Templates, typed)
		}
	}
	return d
}

// refreshes the asset dependencies and notifies the caller of missing assets via the callback
func (d *dependencies) refresh(sa flows.SessionAssets, missing assets.MissingCallback) error {
	for i, ref := range d.Channels {
		a := sa.Channels().Get(ref.UUID)
		if a == nil {
			missing(ref)
		} else {
			d.Channels[i] = a.Reference()
		}
	}
	for i, ref := range d.Fields {
		a := sa.Fields().Get(ref.Key)
		if a == nil {
			missing(ref)
		} else {
			d.Fields[i] = a.Reference()
		}
	}
	for i, ref := range d.Flows {
		a, err := sa.Flows().Get(ref.UUID)
		if err != nil {
			// flows are the one thing that aren't allowed to be missing as we wouldn't know how to
			// route from a subflow split with a missing flow
			return errors.Wrapf(err, "unable to read %s", ref)
		}
		d.Flows[i] = a.Reference()
	}
	for i, ref := range d.Groups {
		a := sa.Groups().Get(ref.UUID)
		if a == nil {
			missing(ref)
		} else {
			d.Groups[i] = a.Reference()
		}
	}
	for i, ref := range d.Labels {
		a := sa.Labels().Get(ref.UUID)
		if a == nil {
			missing(ref)
		} else {
			d.Labels[i] = a.Reference()
		}
	}
	for i, ref := range d.Templates {
		a := sa.Templates().Get(ref.UUID)
		if a == nil {
			missing(ref)
		} else {
			d.Templates[i] = a.Reference()
		}
	}

	return nil
}
