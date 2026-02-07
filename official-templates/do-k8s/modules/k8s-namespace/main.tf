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

resource "kubernetes_namespace_v1" "namespace" {
  metadata {
    name = var.name

    labels = {
      "app.kubernetes.io/managed-by" = "cldctl"
      "cldctl/environment"           = var.name
    }
  }

  depends_on = [local_file.kubeconfig]
}
