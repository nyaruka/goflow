package utils

import (
	"slices"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"
)

// Set converts a slice to a set (a K > bool map)
func Set[K constraints.Ordered](s []K) map[K]bool {
	m := make(map[K]bool, len(s))
	for _, v := range s {
		m[v] = true
	}
	return m
}

// SortedKeys returns the keys of a set in lexical order
func SortedKeys[K constraints.Ordered, V any](m map[K]V) []K {
	keys := maps.Keys(m)
	slices.Sort(keys)
	return keys
}

// Typed is an interface of objects that are marshalled as typed envelopes
type Typed interface {
	Type() string
}

// TypedEnvelope can be mixed into envelopes that have a type field
type TypedEnvelope struct {
	Type string `json:"type" validate:"required"`
}

// ReadTypeFromJSON reads a field called `type` from the given JSON
func ReadTypeFromJSON(data []byte) (string, error) {
	t := &TypedEnvelope{}
	if err := UnmarshalAndValidate(data, t); err != nil {
		return "", err
	}
	return t.Type, nil
}
