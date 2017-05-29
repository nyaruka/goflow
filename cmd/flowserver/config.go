package main

// Server is our configuration file for FlowServer
// This is our multiconfig configuration object. Right now can just specify a static
// directory (instead of serving from statik) or a port, neither is required
type Server struct {
	Port   int    `default:"8080"`
	Static string `default:""`
}
