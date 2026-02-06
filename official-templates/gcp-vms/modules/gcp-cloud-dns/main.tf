terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

resource "google_dns_record_set" "main" {
  project      = var.project
  managed_zone = var.zone_name
  name         = "${var.subdomain}.${var.zone_name}."
  type         = "A"
  ttl          = 300
  rrdatas      = [var.target]
}
