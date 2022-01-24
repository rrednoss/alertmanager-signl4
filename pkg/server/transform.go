package server

import (
	"html/template"
	"os"
)

func transformAlert(alert string, input map[string]interface{}) {
	t := template.Must(template.New("alert").Parse(alert))
	err := t.Execute(os.Stdout, input)
	if err != nil {
		panic(err)
	}
}
