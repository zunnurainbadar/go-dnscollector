
# Logger: InfluxDB client

InfluxDB client to remote InfluxDB server

Options:

- `server-url`: (string) InfluxDB server url
- `auth-token`: (string) authentication token
- `bucket`: (string) bucket name
- `organization`: (string) organization name
- `tls-support`: (boolean) enable tls
- `tls-insecure`: (boolean) insecure skip verify
- `tls-min-version`: (string) min tls version
- `ca-file`: (string) provide CA file to verify the server certificate
- `cert-file`: (string) provide client certificate file for mTLS
- `key-file`: (string) provide client private key file for mTLS
- `chan-buffer-size`: (integer) channel buffer size used on incoming dns message, number of messages before to drop it.

Default values:

```yaml
influxdb:
  server-url: "http://localhost:8086"
  auth-token: ""
  bucket: "db_dns"
  organization: "dnscollector"
  tls-support: false
  tls-insecure: false
  tls-min-version: 1.2
  ca-file: ""
  cert-file: ""
  key-file: ""
  chan-buffer-size: 65535
```