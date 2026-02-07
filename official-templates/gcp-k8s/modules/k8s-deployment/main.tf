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
  replicas = coalesce(var.replicas, 1)
}

resource "kubernetes_deployment_v1" "main" {
  metadata {
    name      = var.name
    namespace = var.namespace

    labels = {
      app        = var.name
      managed-by = "cldctl"
    }
  }

  spec {
    replicas = local.replicas

    selector {
      match_labels = {
        app = var.name
      }
    }

    template {
      metadata {
        labels = {
          app        = var.name
          managed-by = "cldctl"
        }
      }

      spec {
        container {
          name  = "main"
          image = var.image

          dynamic "port" {
            for_each = var.port != null ? [var.port] : []
            content {
              container_port = port.value
            }
          }

          dynamic "env" {
            for_each = local.env_vars
            content {
              name  = env.key
              value = env.value
            }
          }

          resources {
            requests = {
              cpu    = coalesce(var.cpu, "250m")
              memory = coalesce(var.memory, "256Mi")
            }
            limits = {
              cpu    = coalesce(var.cpu, "1")
              memory = coalesce(var.memory, "512Mi")
            }
          }
        }
      }
    }
  }
}
