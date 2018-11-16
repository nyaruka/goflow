package flows

import (
	"github.com/hashicorp/go-version"
)

const (
	minSpecVersion string = "12.0"
	maxSpecVersion string = "13.0"
)

// IsVersionSupported returns whether the given flow spec version is supported
func IsVersionSupported(ver string) bool {
	v, err := version.NewVersion(ver)
	if err != nil {
		return false
	}

	vMin, _ := version.NewVersion(minSpecVersion)
	vMax, _ := version.NewVersion(maxSpecVersion)

	return (v.Equal(vMin) || v.GreaterThan(vMin)) && v.LessThan(vMax)
}
