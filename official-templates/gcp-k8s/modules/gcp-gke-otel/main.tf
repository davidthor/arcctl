terraform {
  required_providers {
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
  }
}

provider "kubernetes" {
  host                   = var.kubeconfig.host
  cluster_ca_certificate = base64decode(var.kubeconfig.cluster_ca_certificate)
  token                  = var.kubeconfig.token
}

provider "helm" {
  kubernetes {
    host                   = var.kubeconfig.host
    cluster_ca_certificate = base64decode(var.kubeconfig.cluster_ca_certificate)
    token                  = var.kubeconfig.token
  }
}

# Deploy OTel Collector DaemonSet via Helm
resource "helm_release" "otel_collector" {
  name       = var.name
  namespace  = var.namespace
  repository = "https://open-telemetry.github.io/opentelemetry-helm-charts"
  chart      = "opentelemetry-collector"
  version    = "0.73.1"

  set {
    name  = "mode"
    value = "daemonset"
  }

  set {
    name  = "config.exporters.googlecloud.project"
    value = var.project
  }

  set {
    name  = "config.exporters.googlecloud.log.default_log_name"
    value = "opentelemetry.io/collector-exported-log"
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
    name  = "config.service.pipelines.traces.exporters[0]"
    value = "googlecloud"
  }

  set {
    name  = "config.service.pipelines.metrics.exporters[0]"
    value = "googlecloud"
  }

  set {
    name  = "config.service.pipelines.logs.exporters[0]"
    value = "googlecloud"
  }

  set {
    name  = "image.repository"
    value = "otel/opentelemetry-collector-contrib"
  }

  set {
    name  = "resources.limits.cpu"
    value = "500m"
  }

  set {
    name  = "resources.limits.memory"
    value = "512Mi"
  }
}

# Expose the collector as a ClusterIP service
resource "kubernetes_service_v1" "otel_collector" {
  metadata {
    name      = "${var.name}-collector"
    namespace = var.namespace

    labels = {
      managed-by = "cldctl"
    }
  }

  spec {
    type = "ClusterIP"

    selector = {
      "app.kubernetes.io/name"     = "opentelemetry-collector"
      "app.kubernetes.io/instance" = var.name
    }

    port {
      name        = "otlp-grpc"
      port        = 4317
      target_port = 4317
      protocol    = "TCP"
    }

    port {
      name        = "otlp-http"
      port        = 4318
      target_port = 4318
      protocol    = "TCP"
    }
  }

  depends_on = [helm_release.otel_collector]
}
