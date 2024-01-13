package config

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	t.Setenv("JWT_SECRET", "secret")

	_ = NewConfig()
}
