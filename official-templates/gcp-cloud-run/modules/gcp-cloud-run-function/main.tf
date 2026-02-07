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

# Cloud Run service configured for serverless function behavior (scale-to-zero)
resource "google_cloud_run_v2_service" "main" {
  name     = var.name
  project  = var.project
  location = var.region
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {
    containers {
      image = var.image

      dynamic "ports" {
        for_each = var.port != null ? [var.port] : [8080]
        content {
          container_port = ports.value
        }
      }

      resources {
        limits = {
          cpu    = coalesce(var.cpu, "1")
          memory = coalesce(var.memory, "512Mi")
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

    # Scale to zero for serverless behavior
    scaling {
      min_instance_count = 0
      max_instance_count = coalesce(var.max_scale, 100)
    }

    timeout = "${coalesce(var.timeout, 300)}s"

    vpc_access {
      connector = var.vpc_connector
      egress    = "ALL_TRAFFIC"
    }

    labels = {
      managed-by = "cldctl"
    }
  }
}

# Allow unauthenticated invocation
resource "google_cloud_run_v2_service_iam_member" "invoker" {
  project  = google_cloud_run_v2_service.main.project
  location = google_cloud_run_v2_service.main.location
  name     = google_cloud_run_v2_service.main.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
