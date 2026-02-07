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

resource "kubernetes_cron_job_v1" "main" {
  metadata {
    name      = var.name
    namespace = var.namespace

    labels = {
      managed-by = "cldctl"
    }
  }

  spec {
    schedule                      = var.schedule
    concurrency_policy            = "Forbid"
    successful_jobs_history_limit = 3
    failed_jobs_history_limit     = 1

    job_template {
      metadata {
        labels = {
          app        = var.name
          managed-by = "cldctl"
        }
      }

      spec {
        backoff_limit = 3

        template {
          metadata {
            labels = {
              app        = var.name
              managed-by = "cldctl"
            }
          }

          spec {
            container {
              name  = "cronjob"
              image = var.image

              dynamic "env" {
                for_each = local.env_vars
                content {
                  name  = env.key
                  value = env.value
                }
              }
            }

            restart_policy = "OnFailure"
          }
        }
      }
    }
  }
}
