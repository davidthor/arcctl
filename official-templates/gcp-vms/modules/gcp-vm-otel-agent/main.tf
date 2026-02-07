terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

# The GCP Ops Agent is the recommended way to collect telemetry from VMs.
# This module creates a configuration that can be applied to VMs via
# the OS Config agent or metadata startup scripts.

# Service account for writing metrics/logs/traces
resource "google_service_account" "otel" {
  account_id   = "${substr(replace(var.name, "/[^a-z0-9-]/", "-"), 0, 28)}-sa"
  display_name = "OTel Agent SA for ${var.name}"
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

# Store the OTLP collector config in Secret Manager for VMs to pull
resource "google_secret_manager_secret" "otel_config" {
  secret_id = "${var.name}-otel-config"
  project   = var.project

  replication {
    auto {}
  }

  labels = {
    managed-by = "cldctl"
  }
}

resource "google_secret_manager_secret_version" "otel_config" {
  secret = google_secret_manager_secret.otel_config.id
  secret_data = yamlencode({
    receivers = {
      otlp = {
        protocols = {
          grpc = { endpoint = "0.0.0.0:4317" }
          http = { endpoint = "0.0.0.0:4318" }
        }
      }
    }
    exporters = {
      googlecloud = {
        project = var.project
      }
    }
    service = {
      pipelines = {
        traces  = { receivers = ["otlp"], exporters = ["googlecloud"] }
        metrics = { receivers = ["otlp"], exporters = ["googlecloud"] }
        logs    = { receivers = ["otlp"], exporters = ["googlecloud"] }
      }
    }
  })
}
