package utils

import (
	"golang.org/x/exp/constraints"
)

// Set converts a slice to a set (a K > bool map)
func Set[K constraints.Ordered](s []K) map[K]bool {
	m := make(map[K]bool, len(s))
	for _, v := range s {
		m[v] = true
	}
	return m
}

// Until encoding/json/v2 there's no easy way to ensure nil slices are marshalled as empty arrays
// see https://github.com/golang/go/discussions/63397
func EnsureNonNil[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
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
