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
  name            = lower(replace(coalesce(var.name, "route"), "/[^a-z0-9-]/", "-"))

  # Build backend refs from service/deployment/function targets
  service_name = coalesce(var.service, var.deployment, var.function, "default")
  service_port = coalesce(var.port, 80)
}

resource "local_file" "kubeconfig" {
  content         = var.kubeconfig
  filename        = local.kubeconfig_path
  file_permission = "0600"
}

resource "kubernetes_manifest" "httproute" {
  manifest = {
    apiVersion = "gateway.networking.k8s.io/v1"
    kind       = "HTTPRoute"
    metadata = {
      name      = local.name
      namespace = var.namespace
      labels = {
        "app.kubernetes.io/managed-by" = "arcctl"
      }
    }
    spec = {
      parentRefs = [{
        name      = var.gateway_name
        namespace = var.namespace
      }]
      hostnames = var.hostnames
      rules = [{
        matches = var.path != null ? [{
          path = {
            type  = "PathPrefix"
            value = var.path
          }
        }] : [{
          path = {
            type  = "PathPrefix"
            value = "/"
          }
        }]
        backendRefs = [{
          name = local.service_name
          port = local.service_port
        }]
      }]
    }
  }

  depends_on = [local_file.kubeconfig]
}
