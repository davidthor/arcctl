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

resource "kubernetes_manifest" "gateway" {
  manifest = {
    apiVersion = "gateway.networking.k8s.io/v1"
    kind       = "Gateway"
    metadata = {
      name      = var.name
      namespace = var.namespace
      annotations = var.tls.enabled ? {
        "networking.gke.io/certmap" = var.name
      } : {}
    }
    spec = {
      gatewayClassName = var.gateway_class
      listeners = concat(
        [
          {
            name     = "http"
            protocol = "HTTP"
            port     = 80
          }
        ],
        var.tls.enabled ? [
          {
            name     = "https"
            protocol = "HTTPS"
            port     = 443
            tls = {
              mode = "Terminate"
            }
          }
        ] : []
      )
    }
  }
}
