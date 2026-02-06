terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
  }
}

locals {
  name             = replace(try(var.name, "function"), "/[^a-zA-Z0-9-]/", "-")
  environment_vars = try(var.environment, {})
  timeout          = try(var.timeout, 300)
  memory           = try(var.memory, "256Mi")
}

resource "kubernetes_manifest" "knative_service" {
  manifest = {
    apiVersion = "serving.knative.dev/v1"
    kind       = "Service"
    metadata = {
      name      = local.name
      namespace = var.namespace
      labels = {
        "app.kubernetes.io/managed-by" = "arcctl"
      }
    }
    spec = {
      template = {
        metadata = {
          annotations = {
            "autoscaling.knative.dev/minScale" = "0"
            "autoscaling.knative.dev/maxScale" = "10"
          }
        }
        spec = {
          containerConcurrency = 0
          timeoutSeconds       = local.timeout
          containers = [{
            image = var.image
            ports = [{
              containerPort = try(var.port, 8080)
            }]
            env = [for k, v in local.environment_vars : {
              name  = k
              value = tostring(v)
            }]
            resources = {
              requests = {
                memory = local.memory
              }
            }
          }]
        }
      }
    }
  }
}
