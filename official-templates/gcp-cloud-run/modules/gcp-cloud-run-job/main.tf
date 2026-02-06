terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

locals {
  env_vars = { for k, v in coalesce(var.environment, {}) : k => v }
}

resource "google_cloud_run_v2_job" "main" {
  name     = var.name
  project  = var.project
  location = var.region

  template {
    template {
      containers {
        image = var.image

        dynamic "command" {
          for_each = var.command != null ? [1] : []
          content {
          }
        }

        resources {
          limits = {
            cpu    = "1"
            memory = "512Mi"
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

      vpc_access {
        connector = var.vpc_connector
        egress    = "ALL_TRAFFIC"
      }

      timeout     = "600s"
      max_retries = 1
    }

    task_count = 1
  }

  labels = {
    managed-by = "arcctl"
  }
}
