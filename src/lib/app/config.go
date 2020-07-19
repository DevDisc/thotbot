package app

// default
var (
  defaultAuthToken = ""
)

// Config
type Config struct {
  DiscordAuthToken  string
}

// Create config
func NewConfig(authToken string) *Config {
  return &Config{
    DiscordAuthToken: authToken,
  }
}

// Default
func DefaultConfig() *Config {
  return NewConfig(defaultAuthToken)
}
