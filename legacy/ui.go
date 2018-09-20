package legacy

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/shopspring/decimal"
)

// the different node types supported by the editor
const (
	UINodeTyepActionSet                 flows.UINodeType = "execute_actions"
	UINodeTypeWaitForResponse           flows.UINodeType = "wait_for_response"
	UINodeTypeSplitByAirtime            flows.UINodeType = "split_by_airtime"
	UINodeTypeSplitBySubflow            flows.UINodeType = "split_by_subflow"
	UINodeTypeSplitByWebhook            flows.UINodeType = "split_by_webhook"
	UINodeTypeSplitByResthook           flows.UINodeType = "split_by_resthook"
	UINodeTypeSplitByGroups             flows.UINodeType = "split_by_groups"
	UINodeTypeSplitByExpression         flows.UINodeType = "split_by_expression"
	UINodeTypeSplitByContactField       flows.UINodeType = "split_by_contact_field"
	UINodeTypeSplitByRunResult          flows.UINodeType = "split_by_run_result"
	UINodeTypeSplitByRunResultDelimited flows.UINodeType = "split_by_run_result_delimited"
	UINodeTypeSplitByRandom             flows.UINodeType = "split_by_random"
)

// Note is a legacy sticky note
type Note struct {
	X     decimal.Decimal `json:"x"`
	Y     decimal.Decimal `json:"y"`
	Title string          `json:"title"`
	Body  string          `json:"body"`
}

// Migrate migrates this note to a new sticky note
func (n *Note) Migrate() flows.Sticky {
	return flows.Sticky{
		"position": map[string]interface{}{"left": n.X.IntPart(), "top": n.Y.IntPart()},
		"title":    n.Title,
		"body":     n.Body,
		"color":    "yellow",
	}
}
