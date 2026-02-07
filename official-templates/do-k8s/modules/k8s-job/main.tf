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

resource "kubernetes_job_v1" "job" {
  metadata {
    name      = var.name
    namespace = var.namespace

    labels = {
      "app.kubernetes.io/managed-by" = "cldctl"
      "cldctl/component"             = var.name
    }
  }

  spec {
    backoff_limit              = var.backoff_limit
    ttl_seconds_after_finished = var.ttl_seconds_after_finished

    template {
      metadata {
        labels = {
          "app.kubernetes.io/managed-by" = "cldctl"
          "cldctl/component"             = var.name
        }
      }

      spec {
        container {
          name    = "task"
          image   = var.image
          command = var.command

          dynamic "env" {
            for_each = local.env_vars
            content {
              name  = env.value.name
              value = env.value.value
            }
          }
        }

        restart_policy = "Never"
      }
    }
  }

  wait_for_completion = true

  timeouts {
    create = "10m"
    update = "10m"
  }

  depends_on = [local_file.kubeconfig]
}
