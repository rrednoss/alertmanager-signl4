package server

import (
	"bytes"
	"text/template"

	log "github.com/sirupsen/logrus"
)

func transform(gotpl string, input map[string]interface{}) (string, error) {
	t := template.Must(template.New("alert").Parse(gotpl))
	var b bytes.Buffer
	if err := t.Execute(&b, input); err != nil {
		log.Error("couldn't transform alert")
		return "", err
	}
	log.Info("transformed alert")
	return b.String(), nil
}
