package main

import (
	"github.com/rrednoss/alertmanager-signl4/pkg/client"
	"github.com/rrednoss/alertmanager-signl4/pkg/config"
	"github.com/rrednoss/alertmanager-signl4/pkg/server"
)

func main() {
	ac := config.NewAppConfig("/conf/signl4.yaml") // TODO (rednoss): Rename to app.yaml because it holds more configs than Signl4 specific ones.
	sc := client.NewSignl4Client(ac)

	// If you read the code for the first time this is a good point to start.
	// The server is the entrypoint for all requests and by this for the whole application.
	server := server.NewServer(server.NewAlertHandler(ac, sc), server.NewHealthHandler())
	server.ListenAndServe()
}
