terraform {
  required_providers {
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.0"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

data "aws_region" "current" {}

locals {
  name = replace(var.name, "/[^a-zA-Z0-9-]/", "-")
}

resource "helm_release" "otel_collector" {
  name       = local.name
  repository = "https://open-telemetry.github.io/opentelemetry-helm-charts"
  chart      = "opentelemetry-collector"
  namespace  = var.namespace
  version    = "0.80.0"

  set {
    name  = "mode"
    value = "daemonset"
  }

  set {
    name  = "config.receivers.otlp.protocols.grpc.endpoint"
    value = "0.0.0.0:4317"
  }

  set {
    name  = "config.receivers.otlp.protocols.http.endpoint"
    value = "0.0.0.0:4318"
  }

  set {
    name  = "config.exporters.awsxray.region"
    value = data.aws_region.current.name
  }

  set {
    name  = "config.exporters.awscloudwatchlogs.log_group_name"
    value = var.log_group
  }

  set {
    name  = "config.exporters.awscloudwatchlogs.region"
    value = data.aws_region.current.name
  }

  set {
    name  = "config.service.pipelines.traces.receivers[0]"
    value = "otlp"
  }

  set {
    name  = "config.service.pipelines.traces.exporters[0]"
    value = "awsxray"
  }

  set {
    name  = "config.service.pipelines.metrics.receivers[0]"
    value = "otlp"
  }

  set {
    name  = "config.service.pipelines.metrics.exporters[0]"
    value = "awscloudwatchlogs"
  }

  set {
    name  = "ports.otlp.enabled"
    value = "true"
  }

  set {
    name  = "ports.otlp-http.enabled"
    value = "true"
  }
}
