terraform {
  required_providers {
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

resource "kubernetes_service_v1" "main" {
  metadata {
    name      = var.name
    namespace = var.namespace

    labels = {
      app        = var.name
      managed-by = "arcctl"
    }
  }

  spec {
    type = "ClusterIP"

    selector = {
      app = var.target
    }

    port {
      port        = var.port
      target_port = var.target_port != null ? var.target_port : var.port
      protocol    = "TCP"
    }
  }
}
