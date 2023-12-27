package assets

// Source is a source of assets
type Source interface {
	Channels() ([]Channel, error)
	Classifiers() ([]Classifier, error)
	ExternalServices() ([]ExternalService, error)
	Fields() ([]Field, error)
	Flow(FlowUUID) (Flow, error)
	Globals() ([]Global, error)
	Groups() ([]Group, error)
	Labels() ([]Label, error)
	Locations() ([]LocationHierarchy, error)
	Resthooks() ([]Resthook, error)
	Templates() ([]Template, error)
	Ticketers() ([]Ticketer, error)
	Topics() ([]Topic, error)
	Users() ([]User, error)
	MsgCatalogs() ([]MsgCatalog, error)
	OrgContexts() ([]OrgContext, error)
}
