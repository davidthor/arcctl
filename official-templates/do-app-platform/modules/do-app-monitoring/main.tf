terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

provider "digitalocean" {
  token = var.token
}

# DigitalOcean App Platform provides built-in monitoring and logging.
# This module configures an OTel-compatible endpoint that forwards
# to DigitalOcean's monitoring infrastructure.

# Deploy a lightweight OTel collector as an App Platform worker
# to bridge between OTel and DO monitoring
resource "digitalocean_app" "monitoring" {
  spec {
    name   = var.name
    region = var.region

    service {
      name               = "${var.name}-collector"
      instance_count     = 1
      instance_size_slug = "basic-xxs"
      http_port          = 4318

      image {
        registry_type = "DOCKER_HUB"
        registry      = "otel"
        repository    = "opentelemetry-collector-contrib"
        tag           = "0.96.0"
      }

      env {
        key   = "OTEL_CONFIG"
        value = <<-EOT
          receivers:
            otlp:
              protocols:
                http:
                  endpoint: 0.0.0.0:4318
          processors:
            batch:
              timeout: 5s
          exporters:
            logging:
              loglevel: info
          service:
            pipelines:
              logs:
                receivers: [otlp]
                processors: [batch]
                exporters: [logging]
              traces:
                receivers: [otlp]
                processors: [batch]
                exporters: [logging]
              metrics:
                receivers: [otlp]
                processors: [batch]
                exporters: [logging]
        EOT
        type  = "GENERAL"
      }
    }
  }
}
