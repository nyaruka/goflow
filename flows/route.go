package flows

// Route describes leaving a node
type Route struct {
	exit  ExitUUID
	match string
}

// Exit returns the UUID of the chosen exit
func (r Route) Exit() ExitUUID { return r.exit }

// Match returns the match which led to this route being chosen
func (r Route) Match() string { return r.match }

// NoRoute is used when a router can't find a route
var NoRoute = Route{}

// NewRoute creates a new route
func NewRoute(exit ExitUUID, match string) Route {
	return Route{exit, match}
}
