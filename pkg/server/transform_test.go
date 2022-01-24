package server

import (
	"testing"
)

func TestTransformAlert(t *testing.T) {
	tests := []struct {
		name           string
		gotpl          string
		input          map[string]interface{}
		expectedOutput string
	}{
		{
			name:  "should transform json",
			gotpl: `{"z": "{{ index . "x" }}"}`,
			input: map[string]interface{}{
				"x": "Lorem Ipsum",
				"y": "Ipsum Lorem",
			},
			expectedOutput: `{"z": "Lorem Ipsum"}`,
		},
		{
			name:  "should transform to broken json with invalid gotpl",
			gotpl: `{"z": "{{ index . "abc" }}"}`,
			input: map[string]interface{}{
				"abc": "Lorem Ipsum",
				"def": "Ipsum Lorem",
			},
			expectedOutput: `{"z": "<no value>"}`,
		},
		{
			name:  "should transform to broken json with invalid input",
			gotpl: `{"z": "{{ index . "x" }}"}`,
			input: map[string]interface{}{
				"abc": "Lorem Ipsum",
				"def": "Ipsum Lorem",
			},
			expectedOutput: `{"z": "<no value>"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := transform(tt.gotpl, tt.input)
			if err != nil {
				t.Errorf(err.Error())
			}

			if got != tt.expectedOutput {
				t.Errorf("got %s, expected %s", got, tt.expectedOutput)
			}
		})
	}
}
