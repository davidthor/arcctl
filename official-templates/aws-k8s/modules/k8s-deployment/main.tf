terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
  }
}

locals {
  name             = replace(try(var.name, "deployment"), "/[^a-zA-Z0-9-]/", "-")
  environment_vars = try(var.environment, {})
  replicas         = try(var.replicas, 1)
  container_port   = try(var.port, 8080)
  cpu_request      = try(var.cpu, "250m")
  memory_request   = try(var.memory, "256Mi")
}

resource "kubernetes_deployment_v1" "this" {
  metadata {
    name      = local.name
    namespace = var.namespace

    labels = {
      "app.kubernetes.io/name"       = local.name
      "app.kubernetes.io/managed-by" = "cldctl"
    }
  }

  spec {
    replicas = local.replicas

    selector {
      match_labels = {
        "app.kubernetes.io/name" = local.name
      }
    }

    template {
      metadata {
        labels = {
          "app.kubernetes.io/name"       = local.name
          "app.kubernetes.io/managed-by" = "cldctl"
        }
      }

      spec {
        container {
          name  = local.name
          image = var.image

          dynamic "port" {
            for_each = [local.container_port]
            content {
              container_port = port.value
              protocol       = "TCP"
            }
          }

          dynamic "env" {
            for_each = local.environment_vars
            content {
              name  = env.key
              value = env.value
            }
          }

          resources {
            requests = {
              cpu    = local.cpu_request
              memory = local.memory_request
            }
          }

          dynamic "command" {
            for_each = try(var.command, null) != null ? [1] : []
            content {
            }
          }
        }
      }
    }
  }
}
