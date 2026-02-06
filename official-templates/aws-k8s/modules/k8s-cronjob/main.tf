terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
  }
}

locals {
  name             = replace(try(var.name, "cronjob"), "/[^a-zA-Z0-9-]/", "-")
  environment_vars = try(var.environment, {})
  schedule         = try(var.schedule, "*/5 * * * *")
}

resource "kubernetes_cron_job_v1" "this" {
  metadata {
    name      = local.name
    namespace = var.namespace

    labels = {
      "app.kubernetes.io/name"       = local.name
      "app.kubernetes.io/managed-by" = "arcctl"
    }
  }

  spec {
    schedule                      = local.schedule
    concurrency_policy            = "Forbid"
    successful_jobs_history_limit = 3
    failed_jobs_history_limit     = 1

    job_template {
      metadata {
        labels = {
          "app.kubernetes.io/name"       = local.name
          "app.kubernetes.io/managed-by" = "arcctl"
        }
      }

      spec {
        backoff_limit = 3

        template {
          metadata {
            labels = {
              "app.kubernetes.io/name" = local.name
            }
          }

          spec {
            restart_policy = "Never"

            container {
              name    = local.name
              image   = var.image
              command = try(var.command, null)

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
    }
  }
}
