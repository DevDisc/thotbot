package app

import (
  "log"

  "github.com/urfave/cli"
)

const (
  discordAuthTokenArg = "discord-auth"
)

// Define runtime
type RunApp func(cfg *Config) error
type CliMethods struct {
	RunApp RunApp
}

// Create CLI app
func Cli(methods *CliMethods) *cli.App {

  // Define CLI input
  app := cli.NewApp()
  app.Name = "thot-bot"
  app.HelpName = "thot-bot"
  app.Usage = "Bot for Thots discord channel"
  app.UsageText = "thot-bot --discord-auth DISCORD_AUTH"
  app.Flags = []cli.Flag{
    cli.StringFlag{
	    Name: discordAuthTokenArg,
	    Value: defaultAuthToken,
	    Usage: "Authentiction token for discord app",
	    EnvVar: "DISCORD_AUTH",
    },
  }
  app.Action = func(c *cli.Context) error{
	// Startup
	cfg := getConfig(c)
	log.Printf("Starting bot...")
	return methods.RunApp(cfg)
  }

  return app
}

func getConfig(c *cli.Context) *Config {
  discordAuth := c.String(discordAuthTokenArg)

  return NewConfig(discordAuth)
}
