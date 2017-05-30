package utils

import (
	"strconv"

	"github.com/buger/jsonparser"
)

type JSONFragment []byte

func (j JSONFragment) Default() interface{} {
	return j
}

func (j JSONFragment) Resolve(key string) interface{} {
	_, err := strconv.Atoi(key)

	// this is a numerical index, convert to jsonparser format
	if err == nil {
		jIdx := "[" + key + "]"
		val, _, _, err := jsonparser.Get(j, jIdx)
		if err == nil {
			return JSONFragment(val)
		}
	}

	val, _, _, err := jsonparser.Get(j, key)
	if err != nil {
		return err
	}
	return JSONFragment(val)
}

func (j JSONFragment) String() string {
	return string(j)
}
