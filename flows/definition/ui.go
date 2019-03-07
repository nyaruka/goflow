package definition

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

//------------------------------------------------------------------------------------------
// Top level UI section
//------------------------------------------------------------------------------------------

type ui struct {
	nodes    map[flows.NodeUUID]flows.UINodeDetails
	stickies map[utils.UUID]flows.Sticky
}

// NewUI creates a new UI section
func NewUI() flows.UI {
	return &ui{
		nodes:    make(map[flows.NodeUUID]flows.UINodeDetails),
		stickies: make(map[utils.UUID]flows.Sticky),
	}
}

func (u *ui) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"nodes":    u.nodes,
		"stickies": u.stickies,
	})
}

// AddNode adds information about a node
func (u *ui) AddNode(uuid flows.NodeUUID, nodeDetails flows.UINodeDetails) {
	u.nodes[uuid] = nodeDetails
}

func (u *ui) GetNode(uuid flows.NodeUUID) flows.UINodeDetails {
	return u.nodes[uuid]
}

// AddSticky adds a new sticky note
func (u *ui) AddSticky(sticky flows.Sticky) {
	u.stickies[utils.NewUUID()] = sticky
}

//------------------------------------------------------------------------------------------
// Details for a specific node's configuration
//------------------------------------------------------------------------------------------

type position struct {
	left int
	top  int
}

func (p position) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"left": p.left,
		"top":  p.top,
	})
}

func (p position) Left() int {
	return p.left
}

func (p position) Top() int {
	return p.top
}

type uiNodeDetails struct {
	NodeType_     flows.UINodeType   `json:"type,omitempty"`
	UiNodeConfig_ flows.UINodeConfig `json:"config,omitempty"`
	Position_     flows.Position     `json:"position"`
}

func (n *uiNodeDetails) Position() flows.Position {
	return n.Position_
}

// NewUINodeDetails creates a ui configuration for a specific
func NewUINodeDetails(x, y int, nodeType flows.UINodeType, uiNodeConfig flows.UINodeConfig) flows.UINodeDetails {
	return &uiNodeDetails{
		NodeType_:     nodeType,
		UiNodeConfig_: uiNodeConfig,
		Position_: position{
			left: x,
			top:  y,
		},
	}
}
