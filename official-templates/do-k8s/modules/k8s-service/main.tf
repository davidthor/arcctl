terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
  }
}

provider "kubernetes" {
  config_path = local.kubeconfig_path
}

locals {
  kubeconfig_path = "${path.module}/.kubeconfig"
  name            = lower(replace(coalesce(var.name, "service"), "/[^a-z0-9-]/", "-"))
  target_port     = coalesce(var.port, 8080)
}

resource "local_file" "kubeconfig" {
  content         = var.kubeconfig
  filename        = local.kubeconfig_path
  file_permission = "0600"
}

resource "kubernetes_service_v1" "service" {
  metadata {
    name      = local.name
    namespace = var.namespace

    labels = {
      "app.kubernetes.io/managed-by" = "arcctl"
    }
  }

  spec {
    type = "ClusterIP"

    selector = {
      "app.kubernetes.io/name" = var.deployment
    }

    port {
      port        = local.target_port
      target_port = local.target_port
      protocol    = "TCP"
    }
  }

  depends_on = [local_file.kubeconfig]
}
