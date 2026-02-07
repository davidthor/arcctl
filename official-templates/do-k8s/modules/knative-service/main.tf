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
  name            = lower(replace(coalesce(var.name, "function"), "/[^a-z0-9-]/", "-"))

  env_list = var.environment != null ? [
    for key, value in var.environment : {
      name  = key
      value = value
    }
  ] : []
}

resource "local_file" "kubeconfig" {
  content         = var.kubeconfig
  filename        = local.kubeconfig_path
  file_permission = "0600"
}

resource "kubernetes_manifest" "knative_service" {
  manifest = {
    apiVersion = "serving.knative.dev/v1"
    kind       = "Service"
    metadata = {
      name      = local.name
      namespace = var.namespace
      labels = {
        "app.kubernetes.io/managed-by" = "cldctl"
      }
    }
    spec = {
      template = {
        metadata = {
          annotations = {
            "autoscaling.knative.dev/minScale" = tostring(var.min_scale)
            "autoscaling.knative.dev/maxScale" = tostring(var.max_scale)
          }
        }
        spec = {
          containerConcurrency = var.concurrency
          timeoutSeconds       = var.timeout
          containers = [{
            image   = var.image
            command = var.command
            env     = local.env_list
            resources = {
              requests = {
                cpu    = coalesce(var.cpu, "250m")
                memory = coalesce(var.memory, "256Mi")
              }
              limits = {
                cpu    = coalesce(var.cpu, "1000m")
                memory = coalesce(var.memory, "512Mi")
              }
            }
            ports = var.port != null ? [{
              containerPort = var.port
            }] : [{
              containerPort = 8080
            }]
          }]
        }
      }
    }
  }

  depends_on = [local_file.kubeconfig]
}
