package config

import (
	"testing"
)

func TestParseConfigFromFlags(t *testing.T) {
	config, err := NewFromFlags()

	if err != nil && config == nil {
		t.Errorf("Formatting flags failed %s", err)
	}

	if !config.Cluster && config.Output != "text" {
		t.Errorf("Config not parsed correctly")
	}
}
