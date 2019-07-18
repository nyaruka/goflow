package inspect

import (
	"reflect"

	"github.com/nyaruka/goflow/flows"
)

// ResultGenerator allows flow objects to declare that they can generate a result
type ResultGenerator interface {
	ResultInfo(flows.Node) *flows.ResultInfo
}

// Results extracts result infos
func Results(node flows.Node, s interface{}, include func(*flows.ResultInfo)) {
	results(node, reflect.ValueOf(s), include)
}

func results(node flows.Node, v reflect.Value, include func(*flows.ResultInfo)) {
	walk(v, func(s reflect.Value) {
		asResultGen, isResultGen := s.Interface().(ResultGenerator)
		if isResultGen {
			ri := asResultGen.ResultInfo(node)
			if ri != nil {
				include(ri)
			}
		}
	}, nil)
}
