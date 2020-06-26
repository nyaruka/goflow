package envs

import "github.com/nyaruka/goflow/utils"

// LocationResolver is used to resolve locations from names or hierarchical paths
type LocationResolver interface {
	FindLocations(string, utils.LocationLevel, *utils.Location) []*utils.Location
	FindLocationsFuzzy(string, utils.LocationLevel, *utils.Location) []*utils.Location
	LookupLocation(utils.LocationPath) *utils.Location
}
