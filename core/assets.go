package core

// Assets provides access to the asset types used by the domain model. The flow engine extends this
// with access to flow definitions as flows.SessionAssets.
type Assets interface {
	Campaigns() *CampaignAssets
	Channels() *ChannelAssets
	Fields() *FieldAssets
	Globals() *GlobalAssets
	Groups() *GroupAssets
	Labels() *LabelAssets
	LLMs() *LLMAssets
	Locations() *LocationAssets
	OptIns() *OptInAssets
	Resthooks() *ResthookAssets
	Templates() *TemplateAssets
	Topics() *TopicAssets
	Users() *UserAssets
}
