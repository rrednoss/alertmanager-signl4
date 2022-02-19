# alertmanager-signl4

The *alertmanager-signl4* application is an adapter that sits between the Alertmanager and the Signl4 app.
This means that you have to point the Alertmanager [webhook_url](https://www.prometheus.io/docs/alerting/latest/configuration/#webhook_config) to this application.

It provides two basic features:
* URL redirection for firing and resolving alerts and
* payload transformation to show only the neccessary details inside Signl4

See the [values.yaml](chart/alertmanager-signl4/values.yaml) descption for more details on how to use and configure the application.

## Installation
```bash
$ helm repo add alertmanager-signl4 <TODO: Add Github URL!>

$ helm upgrade --install --namespace alertmanager-signl4 -f values.yaml alertmanager-signl4 alertmanager-signl4/alertmanager-signl4
```

## Observability
One aspect of cloud-native applications is observability.
It can be defined as a measure of how well the internal state of the application can be observed from its output.
For this reason, this application provides structured logging and Prometheus-enabled metrics that can be received via the `/metrics` endpoint.
The following custom metrics are provided:

* `http_total_requests`
* `http_response_status_total`
* `http_response_time_seconds`
* `alerts_send_total`

## Open Points
* [x] Write documentation on how to use and configure the app.
* [x] Fix broken unit tests. The code needs to be refactored in some way that the Server has its Client as dependency.
* [x] Add structured logging to understand what the app is doing.
* [ ] Add context timeout so that request aren't stuck forever inside the app if something goes wrong.
* [ ] Add buffered channel to accept only a defined amount of parallel requests.
* [ ] Add TLS Support.
* [x] Add /healthz endpoint for liveness and readiness probes.
* [x] Add /metrics endpoint to be monitored by Prometheus.
