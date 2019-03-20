package routers

import (
	"github.com/nyaruka/goflow/flows"
)

type route struct {
	input        *string
	match        string
	categoryUUID flows.CategoryUUID
	extra        map[string]string
}

func newRoute(input *string, match string, categoryUUID flows.CategoryUUID, extra map[string]string) *route {
	return &route{input, match, categoryUUID, extra}
}
