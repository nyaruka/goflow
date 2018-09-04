package legacy

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"

	"github.com/shopspring/decimal"
)

// UINodeType tells the editor how to render a particular node
type UINodeType string

// UINodeConfig contains config unique to its type
type UINodeConfig map[string]string

// the different node types supported by the editor
const (
	UINodeTypeWaitForResponse           UINodeType = "wait_for_response"
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

// Note is a legacy sticky note
type Note struct {
	X     decimal.Decimal `json:"x"`
	Y     decimal.Decimal `json:"y"`
	Title string          `json:"title"`
	Body  string          `json:"body"`
}

// Sticky is a migrated note
type Sticky map[string]interface{}

// Migrate migrates this note to a new sticky note
func (n *Note) Migrate() Sticky {
	return Sticky{
		"position": map[string]interface{}{"left": n.X.IntPart(), "top": n.Y.IntPart()},
		"title":    n.Title,
		"body":     n.Body,
		"color":    "yellow",
	}
}

// UI is a optional section in a flow definition with editor specific information
type UI map[string]interface{}

// NewUI creates a new UI section
func NewUI() UI {
	return UI{
		"nodes":    make(map[flows.NodeUUID]interface{}),
		"stickies": make(map[utils.UUID]Sticky),
	}
}

// AddNode adds information about a node
func (u UI) AddNode(uuid flows.NodeUUID, x, y int, nodeConf uiConfig) {
	node := make(map[string]interface{})
	node["position"] = map[string]int{
		"left": x,
		"top":  y,
	}

	if nodeConf.nodeType != "" {
		node["type"] = nodeConf.nodeType
	}

	if nodeConf.uiNodeConfig != nil {
		node["config"] = nodeConf.uiNodeConfig
	}

	u["nodes"].(map[flows.NodeUUID]interface{})[uuid] = node
}

// AddSticky adds a new sticky note
func (u UI) AddSticky(sticky Sticky) {
	u["stickies"].(map[utils.UUID]Sticky)[utils.NewUUID()] = sticky
}
