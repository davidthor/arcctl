terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

provider "digitalocean" {
  token = var.do_token
}

locals {
  user_data = <<-EOT
    #!/bin/bash
    set -euo pipefail

    # Install Docker and Docker Compose
    apt-get update -y
    apt-get install -y docker.io docker-compose
    systemctl start docker
    systemctl enable docker

    # Create config directory
    mkdir -p /opt/otel

    # Create OTel Collector config
    cat > /opt/otel/otel-collector-config.yaml <<'CONFIG'
    receivers:
      otlp:
        protocols:
          grpc:
            endpoint: 0.0.0.0:4317
          http:
            endpoint: 0.0.0.0:4318
    processors:
      batch:
        timeout: 5s
        send_batch_size: 1000
    exporters:
      loki:
        endpoint: http://loki:3100/loki/api/v1/push
      otlp:
        endpoint: http://loki:3100
        tls:
          insecure: true
    service:
      pipelines:
        logs:
          receivers: [otlp]
          processors: [batch]
          exporters: [loki]
        traces:
          receivers: [otlp]
          processors: [batch]
          exporters: [otlp]
        metrics:
          receivers: [otlp]
          processors: [batch]
          exporters: [otlp]
    CONFIG

    # Create Loki config
    cat > /opt/otel/loki-config.yaml <<'CONFIG'
    auth_enabled: false
    server:
      http_listen_port: 3100
    common:
      path_prefix: /loki
      storage:
        filesystem:
          chunks_directory: /loki/chunks
          rules_directory: /loki/rules
      replication_factor: 1
      ring:
        kvstore:
          store: inmemory
    schema_config:
      configs:
        - from: 2020-10-24
          store: tsdb
          object_store: filesystem
          schema: v13
          index:
            prefix: index_
            period: 24h
    CONFIG

    # Create Grafana provisioning
    mkdir -p /opt/otel/grafana/provisioning/datasources
    cat > /opt/otel/grafana/provisioning/datasources/datasources.yaml <<'CONFIG'
    apiVersion: 1
    datasources:
      - name: Loki
        type: loki
        url: http://loki:3100
        access: proxy
        isDefault: true
    CONFIG

    # Create docker-compose file
    cat > /opt/otel/docker-compose.yaml <<'COMPOSE'
    version: "3.8"
    services:
      otel-collector:
        image: otel/opentelemetry-collector-contrib:0.96.0
        restart: always
        ports:
          - "4317:4317"
          - "4318:4318"
        volumes:
          - ./otel-collector-config.yaml:/etc/otelcol-contrib/config.yaml
        depends_on:
          - loki

      loki:
        image: grafana/loki:2.9.0
        restart: always
        ports:
          - "3100:3100"
        volumes:
          - ./loki-config.yaml:/etc/loki/local-config.yaml
          - loki-data:/loki
        command: -config.file=/etc/loki/local-config.yaml

      grafana:
        image: grafana/grafana:10.3.0
        restart: always
        ports:
          - "3000:3000"
        environment:
          - GF_SECURITY_ADMIN_PASSWORD=cldctl-admin
          - GF_AUTH_ANONYMOUS_ENABLED=true
          - GF_AUTH_ANONYMOUS_ORG_ROLE=Viewer
        volumes:
          - ./grafana/provisioning:/etc/grafana/provisioning
          - grafana-data:/var/lib/grafana
        depends_on:
          - loki

    volumes:
      loki-data:
      grafana-data:
    COMPOSE

    # Start the observability stack
    cd /opt/otel
    docker-compose up -d
  EOT
}

resource "digitalocean_droplet" "otel" {
  name     = var.name
  region   = var.region
  size     = var.size
  image    = "docker-20-04"
  ssh_keys = var.ssh_key_fingerprint != "" ? [var.ssh_key_fingerprint] : []
  vpc_uuid = var.vpc_uuid

  user_data = local.user_data

  tags = ["cldctl", "managed-by:cldctl", "observability"]

  lifecycle {
    create_before_destroy = true
  }
}
