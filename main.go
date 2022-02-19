package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/rrednoss/alertmanager-signl4/pkg/client"
	"github.com/rrednoss/alertmanager-signl4/pkg/config"
	"github.com/rrednoss/alertmanager-signl4/pkg/server"
)

func main() {
	ac := config.NewAppConfig("/conf/signl4.yaml") // TODO (rednoss): Rename to app.yaml because it holds more configs than Signl4 specific ones.
	sc := client.NewSignl4Client(ac)

	// adding structured logging
	log.SetFormatter(&log.JSONFormatter{})
	log.Info("starting the application with the following configuration")
	log.Info("allowInsecureTLSConfig", ac.AllowInsecureTLSConfig)
	log.Info("groupKey", ac.GroupKey)
	log.Info("statusKey", ac.StatusKey)
	log.Info("template", ac.Template)

	// If you read the code for the first time this is a good point to start.
	// The server is the entrypoint for all requests and by this for the whole application.
	server := server.NewServer(server.NewAlertHandler(ac, sc), server.NewHealthHandler())
	server.ListenAndServe()
}
