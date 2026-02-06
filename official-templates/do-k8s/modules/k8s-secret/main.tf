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
}

resource "local_file" "kubeconfig" {
  content         = var.kubeconfig
  filename        = local.kubeconfig_path
  file_permission = "0600"
}

resource "kubernetes_secret_v1" "secret" {
  metadata {
    name      = var.name
    namespace = var.namespace

    labels = {
      "app.kubernetes.io/managed-by" = "arcctl"
    }
  }

  data = var.data
  type = "Opaque"

  depends_on = [local_file.kubeconfig]
}
