terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

# Service account for Cloud Scheduler to invoke Cloud Run
resource "google_service_account" "scheduler" {
  account_id   = "${substr(replace(var.name, "/[^a-z0-9-]/", "-"), 0, 28)}-sch"
  display_name = "Cloud Scheduler SA for ${var.name}"
  project      = var.project
}

resource "google_project_iam_member" "scheduler_invoker" {
  project = var.project
  role    = "roles/run.invoker"
  member  = "serviceAccount:${google_service_account.scheduler.email}"
}

resource "google_cloud_scheduler_job" "main" {
  name     = var.name
  project  = var.project
  region   = var.region
  schedule = var.schedule

  http_target {
    uri         = var.image != null ? "" : ""
    http_method = "POST"

    oidc_token {
      service_account_email = google_service_account.scheduler.email
    }
  }

  retry_config {
    retry_count = 1
  }
}
