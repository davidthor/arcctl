terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

locals {
  dns_name = "${var.name}.internal"
  port     = coalesce(var.port, 8080)
}

# Private DNS zone for internal service discovery
resource "google_dns_managed_zone" "internal" {
  name        = "${var.name}-internal"
  project     = var.project
  dns_name    = "${local.dns_name}."
  description = "Internal DNS zone for ${var.name}"
  visibility  = "private"

  private_visibility_config {
    networks {
      network_url = var.network
    }
  }

  labels = {
    managed-by = "cldctl"
  }
}

resource "google_dns_record_set" "main" {
  project      = var.project
  managed_zone = google_dns_managed_zone.internal.name
  name         = "${local.dns_name}."
  type         = "A"
  ttl          = 60
  rrdatas      = [var.target]
}
