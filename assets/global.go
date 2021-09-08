package assets

import "fmt"

// Global is a named constant.
//
//   {
//     "key": "organization_name",
//     "name": "Organization Name",
//     "value": "U-Report"
//   }
//
// @asset global
type Global interface {
	Key() string
	Name() string
	Value() string
}

// GlobalReference is a reference to a global
type GlobalReference struct {
	Key  string `json:"key" validate:"required"`
	Name string `json:"name"`
}

// NewGlobalReference creates a new global reference with the given key and name
func NewGlobalReference(key string, name string) *GlobalReference {
	return &GlobalReference{Key: key, Name: name}
}

// Type returns the name of the asset type
func (r *GlobalReference) Type() string {
	return "global"
}

// Identity returns the unique identity of the asset
func (r *GlobalReference) Identity() string {
	return string(r.Key)
}

// Variable returns whether this a variable (vs concrete) reference
func (r *GlobalReference) Variable() bool {
	return false
}

func (r *GlobalReference) String() string {
	return fmt.Sprintf("%s[key=%s,name=%s]", r.Type(), r.Identity(), r.Name)
}

var _ Reference = (*GlobalReference)(nil)
