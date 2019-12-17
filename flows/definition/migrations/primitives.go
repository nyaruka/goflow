package migrations

type Flow map[string]interface{}

func (f Flow) Nodes() []Node {
	d, _ := f["nodes"].([]interface{})
	nodes := make([]Node, len(d))
	for i := range d {
		nodes[i] = Node(d[i].(map[string]interface{}))
	}
	return nodes
}

type Node map[string]interface{}

func (n Node) Actions() []Action {
	d, _ := n["actions"].([]interface{})
	actions := make([]Action, len(d))
	for i := range d {
		actions[i] = Action(d[i].(map[string]interface{}))
	}
	return actions
}

func (n Node) Router() Router {
	d, _ := n["router"].(map[string]interface{})
	if d == nil {
		return nil
	}
	return Router(d)
}

type Action map[string]interface{}

func (a Action) Type() string {
	d, _ := a["type"].(string)
	return d
}

type Router map[string]interface{}

func (r Router) Type() string {
	d, _ := r["type"].(string)
	return d
}
