# domain-exporter

A simple Prometheus exporter that queries top-level domain (TLD) authoritative nameservers and exports DNS return codes (RCODE) as metrics.

It can be used to monitor domains directly at the registry level, bypassing DNS caching to catch critical registrar or registry issues that standard monitoring might miss. A real-world example of this would be the application of a problematic EPP domain status code, such as `serverHold`, which would result in an `NXDOMAIN` response.

## Configuration

Create a YAML configuration file containing the list of domains you wish to monitor:

```yaml
domains:
  - example.com
  - example.org
  - expired-domain.com
```

## Usage

Build the application:

```bash
go build -o domain-exporter ./cmd/domain-exporter
```

Start the exporter:

```bash
./domain-exporter --config_path=/path/to/config.yaml --http_listen_address=0.0.0.0:9553
```

Verify that the exporter is running by cURLing the `/metrics` endpoint:

```bash
curl http://localhost:9553/metrics
```

## Metrics

```txt
# HELP domain_rcode_status DNS RCODE of the domain query (0=NOERROR, 2=SERVFAIL, 3=NXDOMAIN, etc.). -1 indicates a network/query error.
# TYPE domain_rcode_status gauge
domain_rcode_status{domain="example.com"} 0
domain_rcode_status{domain="example.org"} 0
domain_rcode_status{domain="expired-domain.com"} 3
```
