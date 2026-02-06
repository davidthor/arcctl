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

# Create a ClusterIssuer for Let's Encrypt if TLS is enabled
resource "kubernetes_manifest" "cluster_issuer" {
  count = var.tls.enabled ? 1 : 0

  manifest = {
    apiVersion = "cert-manager.io/v1"
    kind       = "ClusterIssuer"
    metadata = {
      name = var.tls.issuer
    }
    spec = {
      acme = {
        server = "https://acme-v02.api.letsencrypt.org/directory"
        privateKeySecretRef = {
          name = "${var.tls.issuer}-account-key"
        }
        solvers = [{
          http01 = {
            ingress = {
              class = var.gateway_class
            }
          }
        }]
      }
    }
  }

  depends_on = [local_file.kubeconfig]
}

# Create the Gateway resource
resource "kubernetes_manifest" "gateway" {
  manifest = {
    apiVersion = "gateway.networking.k8s.io/v1"
    kind       = "Gateway"
    metadata = {
      name      = var.name
      namespace = var.namespace
      annotations = var.tls.enabled ? {
        "cert-manager.io/cluster-issuer" = var.tls.issuer
      } : {}
    }
    spec = {
      gatewayClassName = var.gateway_class
      listeners = [{
        name     = "http"
        port     = 80
        protocol = "HTTP"
        allowedRoutes = {
          namespaces = {
            from = "Same"
          }
        }
      }]
    }
  }

  depends_on = [local_file.kubeconfig]
}
