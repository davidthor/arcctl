terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

# Service account for the OTel collector to write to Cloud Monitoring/Logging/Trace
resource "google_service_account" "otel" {
  account_id   = "${substr(replace(var.name, "/[^a-z0-9-]/", "-"), 0, 28)}-sa"
  display_name = "OTel Collector SA for ${var.name}"
  project      = var.project
}

resource "google_project_iam_member" "monitoring_writer" {
  project = var.project
  role    = "roles/monitoring.metricWriter"
  member  = "serviceAccount:${google_service_account.otel.email}"
}

resource "google_project_iam_member" "logging_writer" {
  project = var.project
  role    = "roles/logging.logWriter"
  member  = "serviceAccount:${google_service_account.otel.email}"
}

resource "google_project_iam_member" "trace_writer" {
  project = var.project
  role    = "roles/cloudtrace.agent"
  member  = "serviceAccount:${google_service_account.otel.email}"
}

# Deploy OTel Collector as a Cloud Run service
resource "google_cloud_run_v2_service" "otel_collector" {
  name     = var.name
  project  = var.project
  location = var.region
  ingress  = "INGRESS_TRAFFIC_INTERNAL_ONLY"

  template {
    service_account = google_service_account.otel.email

    containers {
      image = "otel/opentelemetry-collector-contrib:latest"

      ports {
        container_port = 4317
        name           = "h2c"
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
      }

      env {
        name  = "GOOGLE_CLOUD_PROJECT"
        value = var.project
      }
    }

    scaling {
      min_instance_count = 1
      max_instance_count = 3
    }

    labels = {
      managed-by = "arcctl"
    }
  }
}

# Allow internal services to send telemetry
resource "google_cloud_run_v2_service_iam_member" "invoker" {
  project  = google_cloud_run_v2_service.otel_collector.project
  location = google_cloud_run_v2_service.otel_collector.location
  name     = google_cloud_run_v2_service.otel_collector.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
