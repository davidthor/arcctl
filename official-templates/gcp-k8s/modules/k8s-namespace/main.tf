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

resource "kubernetes_namespace_v1" "main" {
  metadata {
    name = var.name

    labels = {
      managed-by = "cldctl"
    }
  }
}
