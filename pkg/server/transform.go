package server

import (
	"bytes"
	"text/template"
)

func transform(gotpl string, input map[string]interface{}) (string, error) {
	t := template.Must(template.New("alert").Parse(gotpl))
	var b bytes.Buffer
	if err := t.Execute(&b, input); err != nil {
		return "", err
	}
	return b.String(), nil
}
