package flows

import "github.com/buger/jsonparser"

type JSONFragment []byte

func (j JSONFragment) Default() interface{} {
	return string(j)
}

func (j JSONFragment) Resolve(key string) interface{} {
	val, _, _, err := jsonparser.Get(j, key)
	if err != nil {
		return err
	}
	return JSONFragment(val)
}
