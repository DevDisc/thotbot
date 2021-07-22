package app

// default
var (
	defaultAuthToken = ""
	defaultPortPath  = "./"
)

// Config
type Config struct {
	DiscordAuthToken string
	PortPath         string
}

// Create config
func NewConfig(authToken string, portPath string) *Config {
	return &Config{
		DiscordAuthToken: authToken,
		PortPath:         portPath,
	}
}

// Default
func DefaultConfig() *Config {
	return NewConfig(defaultAuthToken, defaultPortPath)
}
