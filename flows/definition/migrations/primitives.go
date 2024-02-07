package migrations

// Flow holds a flow definition
type Flow map[string]any

// Localization returns the localization of this flow
func (f Flow) Localization() Localization {
	d, _ := f["localization"].(map[string]any)
	return Localization(d)
}

type Localization map[string]any

// Nodes returns the nodes in this flow
func (f Flow) Nodes() []Node {
	d, _ := f["nodes"].([]any)
	nodes := make([]Node, len(d))
	for i := range d {
		if d[i] != nil {
			nodes[i] = Node(d[i].(map[string]any))
		}
	}
	return nodes
}

// Node holds a node definition
type Node map[string]any

// Actions returns the actions on this node
func (n Node) Actions() []Action {
	d, _ := n["actions"].([]any)
	actions := make([]Action, len(d))
	for i := range d {
		if d[i] != nil {
			actions[i] = Action(d[i].(map[string]any))
		}
	}
	return actions
}

// Router returns the router on this node
func (n Node) Router() Router {
	d, _ := n["router"].(map[string]any)
	return Router(d)
}

// Action holds an action definition
type Action map[string]any

// Type returns the type of this action
func (a Action) Type() string {
	d, _ := a["type"].(string)
	return d
}

// Router holds a router definition
type Router map[string]any

// Type returns the type of this router
func (r Router) Type() string {
	d, _ := r["type"].(string)
	return d
}
