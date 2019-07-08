package context

import "regexp"

// matches a context property description, e.g. groups:[]group -> the groups the contact belongs to
var contextPropRegexp = regexp.MustCompile(`(\w+)\:(\[\])?(\w+)\s→\s([\w\s]+)`)

// Type is a type that exists in the context
type Type interface {
	TypeName() string
	TypeRefs() []string
}

type primitiveType struct {
	name string
}

func newPrimitiveType(name string) *primitiveType {
	return &primitiveType{name: name}
}

func (t *primitiveType) TypeName() string {
	return t.name
}

// TypeRefs returns any references to other types
func (t *primitiveType) TypeRefs() []string {
	return nil
}

var primitiveTypes = []Type{
	newPrimitiveType("any"),
	newPrimitiveType("text"),
	newPrimitiveType("number"),
	newPrimitiveType("datetime"),
}

// Property is a field of a context type which can be accessed in the context with the dot operator
type Property struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	TypeRef     string `json:"type_ref"`
	Array       bool   `json:"array,omitempty"`
}

// NewProperty creates a new property
func NewProperty(name, description string, typeRef string) *Property {
	return &Property{Name: name, Description: description, TypeRef: typeRef, Array: false}
}

// NewArrayProperty creates a new array property
func NewArrayProperty(name, description string, typeRef string) *Property {
	return &Property{Name: name, Description: description, TypeRef: typeRef, Array: true}
}

func ParseProperty(line string) *Property {
	matches := contextPropRegexp.FindStringSubmatch(line)
	if len(matches) != 5 {
		return nil
	}
	return &Property{
		Name:        matches[1],
		Description: matches[4],
		TypeRef:     matches[3],
		Array:       len(matches[2]) > 0,
	}
}

type staticType struct {
	Name       string      `json:"type"`
	Properties []*Property `json:"properties"`
}

func NewStaticType(name string, properties []*Property) Type {
	return &staticType{Name: name, Properties: properties}
}

func (t *staticType) TypeName() string {
	return t.Name
}

// TypeRefs returns any references to other types
func (t *staticType) TypeRefs() []string {
	refs := make([]string, len(t.Properties))
	for i, p := range t.Properties {
		refs[i] = p.TypeRef
	}
	return refs
}

type dynamicType struct {
	Name             string    `json:"type"`
	Source           string    `json:"source"`
	PropertyTemplate *Property `json:"property_template"`
}

func NewDynamicType(name, source string, propertyTemplate *Property) Type {
	return &dynamicType{Name: name, Source: source, PropertyTemplate: propertyTemplate}
}

func (t *dynamicType) TypeName() string {
	return t.Name
}

// TypeRefs returns any references to other types
func (t *dynamicType) TypeRefs() []string {
	return []string{t.PropertyTemplate.TypeRef}
}
