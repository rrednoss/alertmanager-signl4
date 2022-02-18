# alertmanager-signl4
Transforms the Alertmanager payload to a more SIGNL4 friendly one.

## Open Points:
* [ ] Write documentation on how to use and configure the app.
* [x] Fix broken unit tests. The code needs to be refactored in some way that the Server has its Client as dependency.
* [ ] Add more logging information to understand what the app is doing.
* [ ] Add context timeout so that request aren't stuck forever inside the app if something goes wrong.
* [ ] Add buffered channel to accept only a defined amount of parallel requests.
* [ ] Add TLS Support.
* [x] Add /healthz endpoint for liveness and readiness probes.
* [ ] Add /metrics endpoint to be monitored by Prometheus.
