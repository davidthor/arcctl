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
  name            = lower(replace(coalesce(var.name, "cronjob"), "/[^a-z0-9-]/", "-"))

  env_vars = var.environment != null ? [
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

resource "kubernetes_cron_job_v1" "cronjob" {
  metadata {
    name      = local.name
    namespace = var.namespace

    labels = {
      "app.kubernetes.io/managed-by" = "arcctl"
    }
  }

  spec {
    schedule                      = var.schedule
    concurrency_policy            = "Forbid"
    successful_jobs_history_limit = 3
    failed_jobs_history_limit     = 3

    job_template {
      metadata {
        labels = {
          "app.kubernetes.io/managed-by" = "arcctl"
          "arcctl/component"             = local.name
        }
      }

      spec {
        backoff_limit = 3

        template {
          metadata {
            labels = {
              "app.kubernetes.io/managed-by" = "arcctl"
              "arcctl/component"             = local.name
            }
          }

          spec {
            container {
              name    = local.name
              image   = var.image
              command = var.command

              dynamic "env" {
                for_each = local.env_vars
                content {
                  name  = env.value.name
                  value = env.value.value
                }
              }

              resources {
                requests = {
                  cpu    = coalesce(var.cpu, "250m")
                  memory = coalesce(var.memory, "256Mi")
                }
                limits = {
                  cpu    = coalesce(var.cpu, "500m")
                  memory = coalesce(var.memory, "512Mi")
                }
              }
            }

            restart_policy = "OnFailure"
          }
        }
      }
    }
  }

  depends_on = [local_file.kubeconfig]
}
