terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.0"
    }
  }
}

provider "kubernetes" {
  config_path = local.kubeconfig_path
}

provider "helm" {
  kubernetes {
    config_path = local.kubeconfig_path
  }
}

locals {
  kubeconfig_path = "${path.module}/.kubeconfig"
}

resource "local_file" "kubeconfig" {
  content         = var.kubeconfig
  filename        = local.kubeconfig_path
  file_permission = "0600"
}

# Deploy Loki for log aggregation
resource "helm_release" "loki" {
  name       = "${var.name}-loki"
  repository = "https://grafana.github.io/helm-charts"
  chart      = "loki"
  namespace  = var.namespace
  version    = "5.47.0"

  set {
    name  = "loki.auth_enabled"
    value = "false"
  }

  set {
    name  = "loki.commonConfig.replication_factor"
    value = "1"
  }

  set {
    name  = "loki.storage.type"
    value = "filesystem"
  }

  set {
    name  = "singleBinary.replicas"
    value = "1"
  }

  set {
    name  = "monitoring.selfMonitoring.enabled"
    value = "false"
  }

  set {
    name  = "monitoring.lokiCanary.enabled"
    value = "false"
  }

  set {
    name  = "test.enabled"
    value = "false"
  }

  depends_on = [local_file.kubeconfig]
}

# Deploy Grafana for dashboards
resource "helm_release" "grafana" {
  name       = "${var.name}-grafana"
  repository = "https://grafana.github.io/helm-charts"
  chart      = "grafana"
  namespace  = var.namespace
  version    = "7.3.0"

  set {
    name  = "adminPassword"
    value = "arcctl-admin"
  }

  set {
    name  = "service.type"
    value = "ClusterIP"
  }

  values = [yamlencode({
    datasources = {
      "datasources.yaml" = {
        apiVersion = 1
        datasources = [
          {
            name      = "Loki"
            type      = "loki"
            url       = "http://${var.name}-loki:3100"
            access    = "proxy"
            isDefault = true
          }
        ]
      }
    }
  })]

  depends_on = [local_file.kubeconfig, helm_release.loki]
}

# Deploy OpenTelemetry Collector
resource "helm_release" "otel_collector" {
  name       = "${var.name}-collector"
  repository = "https://open-telemetry.github.io/opentelemetry-helm-charts"
  chart      = "opentelemetry-collector"
  namespace  = var.namespace
  version    = "0.82.0"

  values = [yamlencode({
    mode = "deployment"
    config = {
      receivers = {
        otlp = {
          protocols = {
            grpc = { endpoint = "0.0.0.0:4317" }
            http = { endpoint = "0.0.0.0:4318" }
          }
        }
      }
      processors = {
        batch = {
          timeout       = "5s"
          send_batch_size = 1000
        }
      }
      exporters = {
        loki = {
          endpoint = "http://${var.name}-loki:3100/loki/api/v1/push"
        }
        otlp = {
          endpoint = "http://${var.name}-loki:3100"
          tls = { insecure = true }
        }
      }
      service = {
        pipelines = {
          logs = {
            receivers  = ["otlp"]
            processors = ["batch"]
            exporters  = ["loki"]
          }
          traces = {
            receivers  = ["otlp"]
            processors = ["batch"]
            exporters  = ["otlp"]
          }
          metrics = {
            receivers  = ["otlp"]
            processors = ["batch"]
            exporters  = ["otlp"]
          }
        }
      }
    }
    service = {
      type = "ClusterIP"
    }
    ports = {
      otlp = { enabled = true }
      otlp-http = { enabled = true }
    }
  })]

  depends_on = [local_file.kubeconfig, helm_release.loki]
}
