package diagram

import "testing"

func TestNewCLIConfigAcceptsGraphStyles(t *testing.T) {
	config, err := NewCLIConfig(false, false, false, 1, 5, 5, "LR", "rounded", "heavy")
	if err != nil {
		t.Fatalf("NewCLIConfig() error = %v", err)
	}

	if config.GraphBoxStyle != "rounded" {
		t.Fatalf("GraphBoxStyle = %q, want rounded", config.GraphBoxStyle)
	}
	if config.GraphEdgeStyle != "heavy" {
		t.Fatalf("GraphEdgeStyle = %q, want heavy", config.GraphEdgeStyle)
	}
}

func TestConfigValidateRejectsInvalidGraphStyles(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		field  string
	}{
		{
			name: "box style",
			config: &Config{
				GraphDirection:             "LR",
				GraphBoxStyle:              "bubble",
				GraphEdgeStyle:             "light",
				StyleType:                  "cli",
				SequenceSelfMessageWidth:   4,
				SequenceParticipantSpacing: 5,
			},
			field: "GraphBoxStyle",
		},
		{
			name: "edge style",
			config: &Config{
				GraphDirection:             "LR",
				GraphBoxStyle:              "square",
				GraphEdgeStyle:             "wavy",
				StyleType:                  "cli",
				SequenceSelfMessageWidth:   4,
				SequenceParticipantSpacing: 5,
			},
			field: "GraphEdgeStyle",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if err == nil {
				t.Fatal("Validate() error = nil, want error")
			}
			configErr, ok := err.(*ConfigError)
			if !ok {
				t.Fatalf("Validate() error type = %T, want *ConfigError", err)
			}
			if configErr.Field != tt.field {
				t.Fatalf("ConfigError.Field = %q, want %q", configErr.Field, tt.field)
			}
		})
	}
}
