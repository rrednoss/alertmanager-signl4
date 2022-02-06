package server

import (
	"testing"
)

func TestReadConfig(t *testing.T) {
	tests := []struct {
		name               string
		path               string
		expectedErrMsg     string
		expectedStatusKey  string
		expectedTeamSecret string
		expectedTemplate   string
	}{
		{
			name:               "should read config",
			path:               "../../test/config.yaml",
			expectedErrMsg:     "",
			expectedStatusKey:  "status",
			expectedTeamSecret: "random",
			expectedTemplate:   "{ }",
		},
		{
			name:               "should throw error when file is missing",
			path:               "../../test/non-existing.yaml",
			expectedErrMsg:     "open ../../test/non-existing.yaml: no such file or directory",
			expectedStatusKey:  "",
			expectedTeamSecret: "",
			expectedTemplate:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := ReadConfig(tt.path)
			if err != nil && err.Error() != tt.expectedErrMsg {
				t.Errorf("got %s, want %s", err.Error(), tt.expectedErrMsg)
			}
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
