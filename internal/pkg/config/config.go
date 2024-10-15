package config

// Main configuration struct
type Configuration struct {
	Server      ServerConfiguration
	PprofServer PprofServerConfiguration
	Client      ClientConfiguration
}

// Configuration for the HTTP server
type ServerConfiguration struct {
	Host string
	Port int
}

// Configuration for the Pprof server
type PprofServerConfiguration struct {
	Enable bool
	Host   string
	Port   int
}

// Configuration for the client
type ClientConfiguration struct {
	Host string
	Port int
}
