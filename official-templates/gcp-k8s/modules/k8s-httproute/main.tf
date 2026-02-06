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

resource "kubernetes_manifest" "httproute" {
  manifest = {
    apiVersion = "gateway.networking.k8s.io/v1"
    kind       = "HTTPRoute"
    metadata = {
      name      = var.name
      namespace = var.namespace
      labels = {
        managed-by = "arcctl"
      }
    }
    spec = {
      parentRefs = [
        {
          name      = var.gateway_name
          namespace = var.namespace
        }
      ]
      hostnames = var.hostnames
      rules = [
        {
          matches = [
            {
              path = {
                type  = "PathPrefix"
                value = coalesce(var.path, "/")
              }
            }
          ]
          backendRefs = [
            {
              name = var.target
              port = var.port
            }
          ]
        }
      ]
    }
  }
}
