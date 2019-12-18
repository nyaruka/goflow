package migrations

// Flow holds a flow definition
type Flow map[string]interface{}

// Nodes returns the nodes in this flow
func (f Flow) Nodes() []Node {
	d, _ := f["nodes"].([]interface{})
	nodes := make([]Node, len(d))
	for i := range d {
		nodes[i] = Node(d[i].(map[string]interface{}))
	}
	return nodes
}

// Node holds a node definition
type Node map[string]interface{}

// Actions returns the actions on this node
func (n Node) Actions() []Action {
	d, _ := n["actions"].([]interface{})
	actions := make([]Action, len(d))
	for i := range d {
		actions[i] = Action(d[i].(map[string]interface{}))
	}
	return actions
}

// Router returns the router on this node
func (n Node) Router() Router {
	d, _ := n["router"].(map[string]interface{})
	if d == nil {
		return nil
	}
	return Router(d)
}

// Action holds an action definition
type Action map[string]interface{}

// Type returns the type of this action
func (a Action) Type() string {
	d, _ := a["type"].(string)
	return d
}

// Router holds a router definition
type Router map[string]interface{}

// Type returns the type of this router
func (r Router) Type() string {
	d, _ := r["type"].(string)
	return d
}
