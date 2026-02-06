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

  # Sanitize name for Kubernetes
  name = lower(replace(var.name, "/[^a-z0-9-]/", "-"))

  env_vars = var.environment != null ? [
    for key, value in var.environment : {
      name  = key
      value = value
    }
  ] : []

  # Parse CPU and memory with defaults
  cpu_request    = coalesce(var.cpu, "250m")
  memory_request = coalesce(var.memory, "256Mi")
}

resource "local_file" "kubeconfig" {
  content         = var.kubeconfig
  filename        = local.kubeconfig_path
  file_permission = "0600"
}

resource "kubernetes_deployment_v1" "deployment" {
  metadata {
    name      = local.name
    namespace = var.namespace

    labels = {
      "app.kubernetes.io/name"       = local.name
      "app.kubernetes.io/managed-by" = "arcctl"
    }
  }

  spec {
    replicas = var.replicas

    selector {
      match_labels = {
        "app.kubernetes.io/name" = local.name
      }
    }

    template {
      metadata {
        labels = {
          "app.kubernetes.io/name"       = local.name
          "app.kubernetes.io/managed-by" = "arcctl"
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
              cpu    = local.cpu_request
              memory = local.memory_request
            }
            limits = {
              cpu    = local.cpu_request
              memory = local.memory_request
            }
          }

          dynamic "port" {
            for_each = var.port != null ? [var.port] : []
            content {
              container_port = port.value
              protocol       = "TCP"
            }
          }

          dynamic "liveness_probe" {
            for_each = var.port != null ? [1] : []
            content {
              http_get {
                path = var.health_check_path
                port = var.port
              }
              initial_delay_seconds = 30
              period_seconds        = 10
              timeout_seconds       = 5
              failure_threshold     = 3
            }
          }

          dynamic "readiness_probe" {
            for_each = var.port != null ? [1] : []
            content {
              http_get {
                path = var.health_check_path
                port = var.port
              }
              initial_delay_seconds = 5
              period_seconds        = 5
              timeout_seconds       = 3
              failure_threshold     = 3
            }
          }
        }
      }
    }
  }

  depends_on = [local_file.kubeconfig]
}
