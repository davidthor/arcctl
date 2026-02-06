terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

resource "google_secret_manager_secret" "main" {
  secret_id = var.name
  project   = var.project

  replication {
    auto {}
  }

  labels = {
    managed-by = "arcctl"
  }
}

resource "google_secret_manager_secret_version" "main" {
  secret      = google_secret_manager_secret.main.id
  secret_data = var.data
}
