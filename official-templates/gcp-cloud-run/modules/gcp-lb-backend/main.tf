terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

# Serverless Network Endpoint Group for Cloud Run
resource "google_compute_region_network_endpoint_group" "main" {
  name    = "${var.name}-neg"
  project = var.project
  region  = var.region

  network_endpoint_type = "SERVERLESS"

  cloud_run {
    # Extract the Cloud Run service name from the target reference
    service = var.target
  }
}

# Backend service pointing to the serverless NEG
resource "google_compute_backend_service" "main" {
  name    = "${var.name}-backend"
  project = var.project

  load_balancing_scheme = "EXTERNAL_MANAGED"
  protocol              = "HTTPS"

  backend {
    group = google_compute_region_network_endpoint_group.main.id
  }

  log_config {
    enable = true
  }
}
