package flows

import (
	"strconv"
)

// RunContextTopLevels are the allowed top-level variables for expression evaluations
var RunContextTopLevels = []string{
	"child",
	"contact",
	"fields",
	"globals",
	"input",
	"legacy_extra",
	"node",
	"parent",
	"locals",
	"results",
	"resume",
	"run",
	"ticket",
	"trigger",
	"urns",
	"webhook",
}

// ContactQueryEscaping is the escaping function used for expressions in contact queries
func ContactQueryEscaping(s string) string {
	return strconv.Quote(s)
}
