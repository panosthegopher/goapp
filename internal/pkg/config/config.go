package config

// Main configuration struct
type Configuration struct {
	Server ServerConfiguration
}

// Configuration for the HTTP server
type ServerConfiguration struct {
	Host string
	Port int
}
