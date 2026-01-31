# SigNoz Component

This component deploys [SigNoz](https://signoz.io/), an open-source observability platform that provides logs, metrics, and distributed tracing in a single pane of glass.

## Overview

SigNoz provides:
- Distributed tracing with flame graphs and Gantt charts
- Metrics monitoring with customizable dashboards
- Log management with powerful querying
- Alerts and notifications
- OpenTelemetry-native architecture
- Application Performance Monitoring (APM)

## Architecture

This component deploys the following services:

| Service | Image | Description |
|---------|-------|-------------|
| `query-service` | signoz/signoz | Main UI and API backend |
| `otel-collector` | signoz/signoz-otel-collector | OpenTelemetry collector for data ingestion |

Plus these dependencies and databases:

| Resource | Type | Description |
|----------|------|-------------|
| `zookeeper` | Dependency | ZooKeeper for distributed coordination |
| `clickhouse` | Database | Time-series database for telemetry storage |

## Ports

| Port | Protocol | Service | Description |
|------|----------|---------|-------------|
| 8080 | HTTP | query-service | Web UI and API |
| 4317 | gRPC | otel-collector | OTLP gRPC endpoint |
| 4318 | HTTP | otel-collector | OTLP HTTP endpoint |

## System Requirements

- **RAM**: Minimum 4GB allocated to Docker
- **Storage**: 20GB+ recommended for production workloads
- ClickHouse is resource-intensive; allocate more CPU/memory for high-volume telemetry

## Dependencies

This component depends on:
- **[ZooKeeper](../zookeeper/)** - Required for ClickHouse distributed coordination

The ZooKeeper dependency is automatically included and configured.

## Example Environment Configuration

### Basic Setup

```yaml
# environment.yml
name: signoz-production
datacenter: aws-ecs

components:
  zookeeper:
    source: ./zookeeper

  signoz:
    source: ./signoz
    variables:
      retention_period_logs: "7"
      retention_period_traces: "7"
      retention_period_metrics: "30"
```

### High-Volume Configuration

```yaml
# environment.yml
name: signoz-production
datacenter: aws-ecs

components:
  zookeeper:
    source: ./zookeeper
    variables:
      heap_size: "1g"

  signoz:
    source: ./signoz
    variables:
      retention_period_logs: "14"
      retention_period_traces: "14"
      retention_period_metrics: "60"
```

## Instrumenting Your Applications

SigNoz uses OpenTelemetry for instrumentation. Send telemetry data to the OTEL collector endpoints:

### Environment Variables for Your Applications

```yaml
# In your application's component
deployments:
  my-app:
    build:
      context: ./app
    environment:
      # OTLP gRPC (recommended)
      OTEL_EXPORTER_OTLP_ENDPOINT: ${{ dependencies.signoz.services.otel-grpc.url }}
      # Or OTLP HTTP
      # OTEL_EXPORTER_OTLP_ENDPOINT: ${{ dependencies.signoz.services.otel-http.url }}
      OTEL_SERVICE_NAME: my-app
      OTEL_RESOURCE_ATTRIBUTES: "service.namespace=production,deployment.environment=prod"
```

### Language-Specific Setup

#### Node.js

```bash
npm install @opentelemetry/auto-instrumentations-node
```

```javascript
// tracing.js
const { NodeSDK } = require('@opentelemetry/sdk-node');
const { getNodeAutoInstrumentations } = require('@opentelemetry/auto-instrumentations-node');
const { OTLPTraceExporter } = require('@opentelemetry/exporter-trace-otlp-grpc');

const sdk = new NodeSDK({
  traceExporter: new OTLPTraceExporter(),
  instrumentations: [getNodeAutoInstrumentations()],
});

sdk.start();
```

#### Python

```bash
pip install opentelemetry-distro opentelemetry-exporter-otlp
opentelemetry-bootstrap -a install
```

```bash
opentelemetry-instrument python app.py
```

#### Go

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
)

// Initialize the OTLP exporter
exporter, _ := otlptracegrpc.New(ctx)
```

## Configuration Variables

### Data Retention

| Variable | Default | Description |
|----------|---------|-------------|
| `retention_period_logs` | `7` | Retention period for logs (days) |
| `retention_period_traces` | `7` | Retention period for traces (days) |
| `retention_period_metrics` | `30` | Retention period for metrics (days) |

### ClickHouse Settings

| Variable | Default | Description |
|----------|---------|-------------|
| `clickhouse_cluster` | `cluster` | ClickHouse cluster name |
| `clickhouse_replication` | `false` | Enable ClickHouse replication |

### Service Settings

| Variable | Default | Description |
|----------|---------|-------------|
| `query_service_log_level` | `info` | Log level for query service |
| `otel_collector_log_level` | `info` | Log level for OTEL collector |
| `alertmanager_enabled` | `true` | Enable built-in alert manager |

## Getting Started

After deployment:

1. Access the SigNoz UI at your configured URL (port 8080)
2. Default retention is **7 days** for logs/traces and **30 days** for metrics
3. Navigate to **Settings** > **General** to adjust retention periods
4. Instrument your applications using OpenTelemetry SDKs
5. Configure alerts in the **Alerts** section

## Creating Dashboards

SigNoz supports custom dashboards for:
- Application metrics (latency, throughput, error rates)
- Infrastructure metrics (CPU, memory, disk)
- Business metrics (custom metrics from your applications)

## Setting Up Alerts

1. Navigate to **Alerts** in the SigNoz UI
2. Create alert rules based on:
   - Metrics thresholds
   - Log patterns
   - Trace error rates
3. Configure notification channels (Slack, PagerDuty, Email, Webhooks)

## Scaling Considerations

For high-volume environments:

1. **ClickHouse**: Increase CPU/memory allocation and storage
2. **OTEL Collector**: Scale horizontally for high ingestion rates
3. **Query Service**: Add replicas for better UI/API performance
4. **ZooKeeper**: Consider a 3-node ensemble for high availability

## Troubleshooting

### No Data Appearing

1. Verify OTEL collector is receiving data:
   - Check collector logs for incoming spans/metrics
2. Verify ClickHouse connectivity:
   - Check query-service logs for database errors
3. Ensure your application's OTLP endpoint is correctly configured

### High Memory Usage

ClickHouse can consume significant memory. Consider:
- Reducing retention periods
- Increasing memory allocation
- Enabling data sampling in the OTEL collector

## Documentation

- [SigNoz Documentation](https://signoz.io/docs/)
- [Docker Installation Guide](https://signoz.io/docs/install/docker/)
- [Instrumentation Guides](https://signoz.io/docs/instrumentation/)
- [OpenTelemetry Documentation](https://opentelemetry.io/docs/)
