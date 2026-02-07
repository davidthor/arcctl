terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
  }
}

locals {
  name             = replace(var.name, "/[^a-zA-Z0-9-]/", "-")
  environment_vars = var.environment != null ? var.environment : {}
}

resource "kubernetes_job_v1" "this" {
  metadata {
    name      = local.name
    namespace = var.namespace

    labels = {
      "app.kubernetes.io/name"       = local.name
      "app.kubernetes.io/managed-by" = "cldctl"
    }
  }

  spec {
    backoff_limit                = var.backoff_limit
    ttl_seconds_after_finished   = var.ttl_seconds_after_finished

    template {
      metadata {
        labels = {
          "app.kubernetes.io/name"       = local.name
          "app.kubernetes.io/managed-by" = "cldctl"
        }
      }

      spec {
        restart_policy = "Never"

        container {
          name    = local.name
          image   = var.image
          command = var.command

          dynamic "env" {
            for_each = local.environment_vars
            content {
              name  = env.key
              value = env.value
            }
          }
        }
      }
    }
  }

  wait_for_completion = true

  timeouts {
    create = "10m"
  }
}
