terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

# Service account for HMAC key generation
resource "google_service_account" "hmac" {
  account_id   = "${substr(replace(var.name, "/[^a-z0-9-]/", "-"), 0, 28)}-sa"
  display_name = "HMAC key service account for ${var.name}"
  project      = var.project
}

resource "google_project_iam_member" "hmac_storage" {
  project = var.project
  role    = "roles/storage.objectAdmin"
  member  = "serviceAccount:${google_service_account.hmac.email}"
}

resource "google_storage_bucket" "main" {
  name          = var.name
  project       = var.project
  location      = var.region
  force_destroy = true

  uniform_bucket_level_access = true

  versioning {
    enabled = var.versioning
  }

  labels = {
    managed-by = "arcctl"
  }

  dynamic "cors" {
    for_each = var.public ? [1] : []
    content {
      origin          = ["*"]
      method          = ["GET", "HEAD"]
      response_header = ["Content-Type"]
      max_age_seconds = 3600
    }
  }
}

# Make bucket public if requested
resource "google_storage_bucket_iam_member" "public_read" {
  count  = var.public ? 1 : 0
  bucket = google_storage_bucket.main.name
  role   = "roles/storage.objectViewer"
  member = "allUsers"
}

# HMAC key for S3-compatible access
resource "google_storage_hmac_key" "main" {
  service_account_email = google_service_account.hmac.email
  project               = var.project
}
