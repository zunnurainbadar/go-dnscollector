# Logger: Loki client

Loki client to remote server

Options:
- `server-url`: (string) Loki server url
- `job-name`: (string) Job name
- `mode`: (string) output format: text, json, or flat-json
- `flush-interval`: (integer) flush batch every X seconds
- `batch-size`: (integer) batch size for log entries in bytes
- `retry-interval`: (integer) interval in second between before to retry to send batch
- `text-format`: (string) output text format, please refer to the default text format to see all available directives, use this parameter if you want a specific format
- `proxy-url`: (string) Proxy URL
- `tls-support`: (boolean) enable tls
- `tls-insecure`: (boolean) insecure skip verify
- `tls-min-version`: (string) min tls version
- `basic-auth-login`: (string) basic auth login
- `basic-auth-pwd`: (string) basic auth password
- `basic-auth-pwd-file`: (string) path to a file containing the basic auth password
- `tenant-id`: (string) tenant/organisation id. If omitted or empty, no X-Scope-OrgID header is sent.
- `relabel-configs`: (list) configuration to relabel targets. Functionality like described in https://grafana.com/docs/loki/latest/clients/promtail/configuration/#relabel_configs.
- `chan-buffer-size`: (integer) channel buffer size used on incoming dns message, number of messages before to drop it.

Default values:

```yaml
lokiclient:
  server-url: "http://localhost:3100/loki/api/v1/push"
  job-name: "dnscollector"
  mode: "text"
  flush-interval: 5
  batch-size: 1048576
  retry-interval: 10
  text-format: ""
  proxy-url: ""
  tls-insecure: false
  tls-min-version: 1.2
  basic-auth-login: ""
  basic-auth-pwd: ""
  basic-auth-pwd-file: ""
  tenant-id: ""
  relabel-configs: []
  chan-buffer-size: 65535
```

## Grafana dashboard with Loki datasource

The following [build-in](https://grafana.com/grafana/dashboards/15415) dashboard is available

<p align="center">
  <img src="dashboard_loki.png" alt="dnscollector"/>
</p>