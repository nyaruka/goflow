package inspect

import (
	"reflect"

	"github.com/nyaruka/goflow/flows"
)

// ResultContainer allows flow objects to declare that they can generate a result
type ResultContainer interface {
	Results(func(*flows.ResultInfo))
}

// Results extracts result infos
func Results(s interface{}, include func(*flows.ResultInfo)) {
	results(reflect.ValueOf(s), include)
}

func results(v reflect.Value, include func(*flows.ResultInfo)) {
	walk(v, func(s reflect.Value) {
		asResultGen, isResultGen := s.Interface().(ResultContainer)
		if isResultGen {
			asResultGen.Results(include)
		}
	}, nil)
}
