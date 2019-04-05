package legacy

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"

	"github.com/shopspring/decimal"
)

// UINodeType tells the editor how to render a particular node
type UINodeType string

// UINodeConfig contains config unique to its type
type UINodeConfig map[string]interface{}

// the different node types supported by the editor
const (
	UINodeTypeActionSet                 UINodeType = "execute_actions"
	UINodeTypeWaitForResponse           UINodeType = "wait_for_response"
	UINodeTypeSplitByAirtime            UINodeType = "split_by_airtime"
	UINodeTypeSplitBySubflow            UINodeType = "split_by_subflow"
	UINodeTypeSplitByWebhook            UINodeType = "split_by_webhook"
	UINodeTypeSplitByResthook           UINodeType = "split_by_resthook"
	UINodeTypeSplitByGroups             UINodeType = "split_by_groups"
	UINodeTypeSplitByExpression         UINodeType = "split_by_expression"
	UINodeTypeSplitByContactField       UINodeType = "split_by_contact_field"
	UINodeTypeSplitByRunResult          UINodeType = "split_by_run_result"
	UINodeTypeSplitByRunResultDelimited UINodeType = "split_by_run_result_delimited"
	UINodeTypeSplitByRandom             UINodeType = "split_by_random"
)

// UI is a optional section in a flow definition with editor specific information
type UI interface {
	AddNode(uuid flows.NodeUUID, details UINodeDetails)
	AddSticky(sticky Sticky)

	GetNode(uuid flows.NodeUUID) UINodeDetails
}

// UINodeDetails is the top level ui details for a node
type UINodeDetails interface {
	Position() Position
}

// Position holds coordinates for a node
type Position interface {
	Left() int
	Top() int
}

//------------------------------------------------------------------------------------------
// Top level UI section
//------------------------------------------------------------------------------------------

type ui struct {
	nodes    map[flows.NodeUUID]UINodeDetails
	stickies map[utils.UUID]Sticky
}

// NewUI creates a new UI section
func NewUI() UI {
	return &ui{
		nodes:    make(map[flows.NodeUUID]UINodeDetails),
		stickies: make(map[utils.UUID]Sticky),
	}
}

func (u *ui) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"nodes":    u.nodes,
		"stickies": u.stickies,
	})
}

// AddNode adds information about a node
func (u *ui) AddNode(uuid flows.NodeUUID, nodeDetails UINodeDetails) {
	u.nodes[uuid] = nodeDetails
}

func (u *ui) GetNode(uuid flows.NodeUUID) UINodeDetails {
	return u.nodes[uuid]
}

// AddSticky adds a new sticky note
func (u *ui) AddSticky(sticky Sticky) {
	u.stickies[utils.NewUUID()] = sticky
}

// Sticky is a migrated note
type Sticky map[string]interface{}

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
	NodeType_     UINodeType   `json:"type,omitempty"`
	UiNodeConfig_ UINodeConfig `json:"config,omitempty"`
	Position_     Position     `json:"position"`
}

func (n *uiNodeDetails) Position() Position {
	return n.Position_
}

// NewUINodeDetails creates a ui configuration for a specific
func NewUINodeDetails(x, y int, nodeType UINodeType, uiNodeConfig UINodeConfig) UINodeDetails {
	return &uiNodeDetails{
		NodeType_:     nodeType,
		UiNodeConfig_: uiNodeConfig,
		Position_: position{
			left: x,
			top:  y,
		},
	}
}

// Note is a legacy sticky note
type Note struct {
	X     decimal.Decimal `json:"x"`
	Y     decimal.Decimal `json:"y"`
	Title string          `json:"title"`
	Body  string          `json:"body"`
}

// Migrate migrates this note to a new sticky note
func (n *Note) Migrate() Sticky {
	return Sticky{
		"position": map[string]interface{}{"left": n.X.IntPart(), "top": n.Y.IntPart()},
		"title":    n.Title,
		"body":     n.Body,
		"color":    "yellow",
	}
}
