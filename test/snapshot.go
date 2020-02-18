package test

import (
	"flag"
)

// UpdateSnapshots indicates whether tests should update snapshots
var UpdateSnapshots bool

func init() {
	flag.BoolVar(&UpdateSnapshots, "update", false, "whether to update test snapshots")
}
