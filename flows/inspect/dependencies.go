package inspect

import (
	"fmt"
	"reflect"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
)

type Dependency struct {
	Reference_ assets.Reference `json:"-"`
	Type_      string           `json:"type"`
	Missing_   bool             `json:"missing,omitempty"`
}

func (d *Dependency) Reference() assets.Reference {
	return d.Reference_
}

func (d *Dependency) Type() string {
	return d.Type_
}

func (d *Dependency) Missing() bool {
	return d.Missing_
}

func (d Dependency) MarshalJSON() ([]byte, error) {
	type dependency Dependency // need to alias type to avoid circular calls to this method
	return jsonx.MarshalMerged(d.Reference_, dependency(d))
}

// NewDependencies inspects a list of references. If a session assets is provided,
// each dependency is checked to see if it is available or missing.
func NewDependencies(refs []flows.ExtractedReference, sa flows.SessionAssets) []flows.Dependency {
	deps := make([]flows.Dependency, 0)
	depsSeen := make(map[string]*Dependency)

	for _, er := range refs {
		key := fmt.Sprintf("%s:%s", er.Reference.Type(), er.Reference.Identity())

		// create new dependency record if we haven't seen this reference before
		if _, seen := depsSeen[key]; !seen {
			// check if this dependency is accessible
			missing := false
			if sa != nil {
				missing = !CheckReference(sa, er.Reference)
			}

			dep := &Dependency{
				Reference_: er.Reference,
				Type_:      er.Reference.Type(),
				Missing_:   missing,
			}
			deps = append(deps, dep)
			depsSeen[key] = dep
		}
	}

	return deps
}

// CheckReference determines whether this reference is accessible
func CheckReference(sa flows.SessionAssets, ref assets.Reference) bool {
	switch typed := ref.(type) {
	case *assets.ChannelReference:
		return sa.Channels().Get(typed.UUID) != nil
	case *assets.ClassifierReference:
		return sa.Classifiers().Get(typed.UUID) != nil
	case *flows.ContactReference:
		return true // have to assume contacts exist
	case *assets.FieldReference:
		return sa.Fields().Get(typed.Key) != nil
	case *assets.FlowReference:
		_, err := sa.Flows().Get(typed.UUID)
		return err == nil
	case *assets.GlobalReference:
		return sa.Globals().Get(typed.Key) != nil
	case *assets.GroupReference:
		return sa.Groups().Get(typed.UUID) != nil
	case *assets.LabelReference:
		return sa.Labels().Get(typed.UUID) != nil
	case *assets.TemplateReference:
		return sa.Templates().Get(typed.UUID) != nil
	case *assets.TicketerReference:
		return sa.Ticketers().Get(typed.UUID) != nil
	case *assets.TopicReference:
		return sa.Topics().Get(typed.UUID) != nil
	case *assets.UserReference:
		return sa.Users().Get(typed.Email) != nil
	default:
		panic(fmt.Sprintf("unknown dependency type reference: %T", ref))
	}
}

// DependencyContainer allows flow objects to declare other dependencies
type DependencyContainer interface {
	Dependencies(flows.Localization, func(envs.Language, assets.Reference))
}

// Dependencies extracts dependencies
func Dependencies(s interface{}, localization flows.Localization, include func(envs.Language, assets.Reference)) {
	dependencies(reflect.ValueOf(s), localization, include)
}

func dependencies(v reflect.Value, localization flows.Localization, include func(envs.Language, assets.Reference)) {
	walk(
		v,
		func(sv reflect.Value) {
			// anything which is a DependencyContainer can explicitly provide dependencies
			asDepCon, isDepCon := sv.Interface().(DependencyContainer)
			if isDepCon {
				asDepCon.Dependencies(localization, include)
			}
		},
		func(sv reflect.Value, fv reflect.Value, ef *EngineField) {
			// extract any asset.Reference fields automatically as dependencies
			extractAssetReferences(fv, include)
		},
	)
}

func extractAssetReferences(v reflect.Value, include func(envs.Language, assets.Reference)) {
	if v.Kind() == reflect.Slice {
		// field is a slice of asset references
		for i := 0; i < v.Len(); i++ {
			extractAssetReferences(v.Index(i), include)
		}
	} else if v.Kind() == reflect.Ptr && !v.IsNil() {
		// field is a single asset reference
		asRef, isRef := v.Interface().(assets.Reference)
		if isRef && asRef != nil && !asRef.Variable() {
			include(envs.NilLanguage, asRef)
		}
	}
}
