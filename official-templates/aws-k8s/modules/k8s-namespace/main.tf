terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
  }
}

resource "kubernetes_namespace_v1" "this" {
  metadata {
    name = var.name

    labels = {
      "app.kubernetes.io/managed-by" = "arcctl"
      "arcctl.io/environment"        = var.name
    }
  }
}
