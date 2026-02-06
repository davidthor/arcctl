terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

resource "google_compute_network" "main" {
  name                    = var.name
  project                 = var.project
  auto_create_subnetworks = false

  labels = {
    managed-by = "arcctl"
  }
}

resource "google_compute_subnetwork" "main" {
  name          = "${var.name}-subnet"
  project       = var.project
  region        = var.region
  network       = google_compute_network.main.id
  ip_cidr_range = "10.0.0.0/20"

  private_ip_google_access = true
}

# Reserve an IP range for VPC peering (Cloud SQL, Memorystore, etc.)
resource "google_compute_global_address" "private_services" {
  name          = "${var.name}-private-services"
  project       = var.project
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.main.id
}

# Establish private services connection for Cloud SQL / Memorystore
resource "google_service_networking_connection" "private" {
  network                 = google_compute_network.main.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_services.name]
}
