terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
  }
}

locals {
  name = try(var.name, "service")
  port = try(var.port, 80)
}

resource "kubernetes_service_v1" "this" {
  metadata {
    name      = local.name
    namespace = var.namespace

    labels = {
      "app.kubernetes.io/name"       = local.name
      "app.kubernetes.io/managed-by" = "arcctl"
    }
  }

  spec {
    type = "ClusterIP"

    selector = {
      "app.kubernetes.io/name" = var.target
    }

    port {
      port        = local.port
      target_port = var.target_port != null ? var.target_port : local.port
      protocol    = "TCP"
    }
  }
}
