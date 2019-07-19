package inspect

import (
	"reflect"
)

func walk(v reflect.Value, visitStruct func(reflect.Value), visitField func(reflect.Value, reflect.Value, *EngineField)) {
	// get the real underlying value
	rv := derefValue(v)

	if rv.Kind() == reflect.Slice {
		for i := 0; i < rv.Len(); i++ {
			walk(rv.Index(i), visitStruct, visitField)
		}
	} else if rv.Kind() == reflect.Struct {
		if visitStruct != nil {
			visitStruct(v)
		}

		fields := extractEngineFields(v.Type(), rv.Type())

		for _, ef := range fields {
			fv := ef.Getter(rv)

			if visitField != nil {
				visitField(v, fv, ef)
			}

			walk(fv, visitStruct, visitField)
		}
	}
}

// gets the actual value if we've been given an interface or pointer
func derefValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

func walkTypes(t reflect.Type, path string, visitField func(string, *EngineField)) {
	// get the real underlying type
	rt := derefType(t)

	if rt.Kind() == reflect.Slice {
		walkTypes(rt.Elem(), path+"[*]", visitField)
	} else if rt.Kind() == reflect.Struct {
		fields := extractEngineFields(t, rt)

		for _, ef := range fields {
			fp := path + "." + ef.JSONName
			if visitField != nil {
				visitField(fp, ef)
			}

			walkTypes(ef.Type, fp, visitField)
		}
	}
}

// gets the actual type if we've been given an interface or pointer type
func derefType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Interface || t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}
