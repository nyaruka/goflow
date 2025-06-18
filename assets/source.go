package assets

// Source is a source of assets
type Source interface {
	Campaigns() ([]Campaign, error)
	Channels() ([]Channel, error)
	Classifiers() ([]Classifier, error)
	Fields() ([]Field, error)
	FlowByUUID(FlowUUID) (Flow, error)
	FlowByName(string) (Flow, error)
	Globals() ([]Global, error)
	Groups() ([]Group, error)
	Labels() ([]Label, error)
	LLMs() ([]LLM, error)
	Locations() ([]LocationHierarchy, error)
	OptIns() ([]OptIn, error)
	Resthooks() ([]Resthook, error)
	Templates() ([]Template, error)
	Topics() ([]Topic, error)
	Users() ([]User, error)
}
