package legacy

import (
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

//------------------------------------------------------------------------------------------
// Top level UI section
//------------------------------------------------------------------------------------------

// UI is the _ui section of the flow definition used by the editor
type UI struct {
	Nodes    map[flows.NodeUUID]*UINodeDetails `json:"nodes"`
	Stickies map[utils.UUID]Sticky             `json:"stickies"`
}

// NewUI creates a new UI section
func NewUI() *UI {
	return &UI{
		Nodes:    make(map[flows.NodeUUID]*UINodeDetails),
		Stickies: make(map[utils.UUID]Sticky),
	}
}

// AddNode adds information about a node
func (u *UI) AddNode(uuid flows.NodeUUID, nodeDetails *UINodeDetails) {
	u.Nodes[uuid] = nodeDetails
}

// AddSticky adds a new sticky note
func (u *UI) AddSticky(sticky Sticky) {
	u.Stickies[utils.NewUUID()] = sticky
}

// Position is a position of a node in the editor canvas
type Position struct {
	Left int `json:"left"`
	Top  int `json:"top"`
}

// UINodeDetails are the node specific UI details
type UINodeDetails struct {
	Type     UINodeType   `json:"type,omitempty"`
	Config   UINodeConfig `json:"config,omitempty"`
	Position Position     `json:"position"`
}

// NewUINodeDetails creates a ui configuration for a specific
func NewUINodeDetails(x, y int, nodeType UINodeType, uiNodeConfig UINodeConfig) *UINodeDetails {
	return &UINodeDetails{
		Type:   nodeType,
		Config: uiNodeConfig,
		Position: Position{
			Left: x,
			Top:  y,
		},
	}
}

// Sticky is a user note
type Sticky struct {
	Position Position `json:"position"`
	Title    string   `json:"title"`
	Body     string   `json:"body"`
	Color    string   `json:"color"`
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
		Position: Position{Left: int(n.X.IntPart()), Top: int(n.Y.IntPart())},
		Title:    n.Title,
		Body:     n.Body,
		Color:    "yellow",
	}
}
