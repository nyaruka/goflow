package completion

import (
	"regexp"
	"strings"
)

// matches a context property description, e.g. groups:[]group -> the groups the contact belongs to
var contextPropRegexp = regexp.MustCompile(`(\w+)\:(\[\])?(\w+)\s→\s([\w\s-]+)`)

// Type is a type that exists in the context
type Type interface {
	Name() string
	TypeRefs() []string
	EnumerateProperties(context *Context) []*Property
}

type primitiveType struct {
	name string
}

func newPrimitiveType(name string) *primitiveType {
	return &primitiveType{name: name}
}

// TypeName returns the name that is used to reference this type in a property
func (t *primitiveType) Name() string {
	return t.name
}

// TypeRefs returns any references to other types
func (t *primitiveType) TypeRefs() []string {
	return nil
}

// EnumerateProperties enumerates runtime properties
func (t *primitiveType) EnumerateProperties(context *Context) []*Property {
	return nil // primitive types never have properties
}

var primitiveTypes = []Type{
	newPrimitiveType("any"),
	newPrimitiveType("text"),
	newPrimitiveType("number"),
	newPrimitiveType("datetime"),
}

// Property is a field of a context type which can be accessed in the context with the dot operator
type Property struct {
	Key   string `json:"key"`
	Help  string `json:"help"`
	Type  string `json:"type"`
	Array bool   `json:"array,omitempty"`
}

// NewProperty creates a new property
func NewProperty(key, help string, typeRef string) *Property {
	return &Property{Key: key, Help: help, Type: typeRef, Array: false}
}

// NewArrayProperty creates a new array property
func NewArrayProperty(key, help string, typeRef string) *Property {
	return &Property{Key: key, Help: help, Type: typeRef, Array: true}
}

// ParseProperty parses a property from a docstring line
func ParseProperty(line string) *Property {
	matches := contextPropRegexp.FindStringSubmatch(line)
	if len(matches) != 5 {
		return nil
	}
	return &Property{
		Key:   matches[1],
		Help:  matches[4],
		Type:  matches[3],
		Array: len(matches[2]) > 0,
	}
}

// a type with fixed properties
type staticType struct {
	Name_      string      `json:"name"`
	Properties []*Property `json:"properties"`
}

// NewStaticType creates a new static type
func NewStaticType(name string, properties []*Property) Type {
	return &staticType{Name_: name, Properties: properties}
}

// TypeName returns the name that is used to reference this type in a property
func (t *staticType) Name() string {
	return t.Name_
}

// TypeRefs returns any references to other types
func (t *staticType) TypeRefs() []string {
	refs := make([]string, len(t.Properties))
	for i, p := range t.Properties {
		refs[i] = p.Type
	}
	return refs
}

// EnumerateProperties enumerates runtime properties
func (t *staticType) EnumerateProperties(context *Context) []*Property {
	return t.Properties
}

type dynamicType struct {
	Name_            string    `json:"name"`
	KeySource        string    `json:"key_source"`
	PropertyTemplate *Property `json:"property_template"`
}

// NewDynamicType creates a new dynamic type, i.e. properties determined at runtime
func NewDynamicType(name, keySource string, propertyTemplate *Property) Type {
	return &dynamicType{Name_: name, KeySource: keySource, PropertyTemplate: propertyTemplate}
}

// TypeName returns the name that is used to reference this type in a property
func (t *dynamicType) Name() string {
	return t.Name_
}

// TypeRefs returns any references to other types
func (t *dynamicType) TypeRefs() []string {
	return []string{t.PropertyTemplate.Type}
}

// EnumerateProperties enumerates runtime properties
func (t *dynamicType) EnumerateProperties(context *Context) []*Property {
	keyTemplate := t.PropertyTemplate.Key
	helpTemplate := t.PropertyTemplate.Help

	keys := context.KeySources[t.KeySource]
	properties := make([]*Property, len(keys))

	for i, key := range keys {
		key := strings.Replace(keyTemplate, "{key}", key, -1)
		help := strings.Replace(helpTemplate, "{key}", key, -1)
		properties[i] = NewProperty(key, help, t.PropertyTemplate.Type)
	}
	return properties
}
