package app_test

import (
  "testing"

  "github.com/DevDisc/thotbot/lib/app"

  "github.com/stretchr/testify/assert"
)

func TestCli(t *testing.T) {
  // Mock app
  runApp := func(*Config) error { return nil }

  cli := Cli(&CliMethods{
    RunApp: runApp,
  })

  assert.Equal(t, "thot-bot", cli.Name)
}
