terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

locals {
  # For Cloud Run services, internal routing uses the Cloud Run service URL directly
  service_url = var.target
  port        = coalesce(var.port, 443)
  host        = replace(replace(local.service_url, "https://", ""), "/", "")
}
