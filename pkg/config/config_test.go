package config

import (
	"testing"
)

func TestNewAppConfig(t *testing.T) {
	tests := []struct {
		name               string
		expectedStatusKey  string
		expectedTeamSecret string
		expectedTemplate   string
	}{
		{
			name:               "should create a new valid AppConfig",
			expectedStatusKey:  "status",
			expectedTeamSecret: "random",
			expectedTemplate:   "{ }",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewAppConfig("../../test/config.yaml")
			if config.StatusKey != tt.expectedStatusKey {
				t.Errorf("got %s, want %s", config.StatusKey, tt.expectedStatusKey)
			}
			if config.TeamSecret != tt.expectedTeamSecret {
				t.Errorf("got %s, want %s", config.StatusKey, tt.expectedTeamSecret)
			}
			if config.Template != tt.expectedTemplate {
				t.Errorf("got %s, want %s", config.StatusKey, tt.expectedTemplate)
			}
		})
	}
}
