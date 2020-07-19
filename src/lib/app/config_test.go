package app_test

import (
  "testing"

  . "github.com/DevDisc/thotbot/lib/app"
  "github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
  cfg := NewConfig("test")
  assert.Equal(t, "test", cfg.DiscordAuthToken)
}

func TestDefaultConfig(t *testing.T) {
  cfg := DefaultConfig()
  assert.Equal(t, "", cfg.DiscordAuthToken)
}
