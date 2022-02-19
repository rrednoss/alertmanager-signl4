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

## Future Improvements
* [ ] Add TLS Support.
* [ ] Add Tracing Support.
* [ ] Limit the number of simultaneous performing requests.
* [ ] Limit the number of waiting/queued requests.
* [ ] Limit how long a single requests can run.

### Reflections
* Adding logging and metrics information in some methods makes the code sometimes quite messy. There might be a more elegant or uniform way of doing this.
* This application uses the logrus package for structured logging. It is in maintenance-mode. It might be worse to log for an alternative solution like Zap or Apex.
