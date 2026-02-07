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
}

resource "kubernetes_job_v1" "main" {
  metadata {
    name      = var.name
    namespace = var.namespace

    labels = {
      managed-by = "cldctl"
    }
  }

  spec {
    backoff_limit              = var.backoff_limit
    ttl_seconds_after_finished = var.ttl_seconds_after_finished

    template {
      metadata {
        labels = {
          app        = var.name
          managed-by = "cldctl"
        }
      }

      spec {
        container {
          name  = "task"
          image = var.image

          dynamic "command" {
            for_each = var.command != null ? [1] : []
            content {
            }
          }

          dynamic "env" {
            for_each = local.env_vars
            content {
              name  = env.key
              value = env.value
            }
          }
        }

        restart_policy = "Never"
      }
    }
  }

  wait_for_completion = false
}
