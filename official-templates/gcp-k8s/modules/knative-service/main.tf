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

locals {
  env_vars = { for k, v in coalesce(var.environment, {}) : k => v }
  env_list = [for k, v in local.env_vars : { name = k, value = v }]
}

resource "kubernetes_manifest" "knative_service" {
  manifest = {
    apiVersion = "serving.knative.dev/v1"
    kind       = "Service"
    metadata = {
      name      = var.name
      namespace = var.namespace
      labels = {
        managed-by = "arcctl"
      }
    }
    spec = {
      template = {
        metadata = {
          labels = {
            app        = var.name
            managed-by = "arcctl"
          }
          annotations = {
            "autoscaling.knative.dev/minScale" = "0"
            "autoscaling.knative.dev/maxScale" = tostring(coalesce(var.max_scale, 100))
          }
        }
        spec = {
          containers = [
            {
              image = var.image
              ports = var.port != null ? [{ containerPort = var.port }] : [{ containerPort = 8080 }]
              env   = local.env_list
              resources = {
                limits = {
                  cpu    = coalesce(var.cpu, "1")
                  memory = coalesce(var.memory, "512Mi")
                }
              }
            }
          ]
          timeoutSeconds = coalesce(var.timeout, 300)
        }
      }
    }
  }
}
