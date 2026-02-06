terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

resource "google_compute_firewall" "allow_http" {
  name    = "${var.name}-allow-http"
  project = var.project
  network = var.network

  allow {
    protocol = "tcp"
    ports    = ["80", "443"]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = var.tags
}

resource "google_compute_firewall" "allow_ssh" {
  name    = "${var.name}-allow-ssh"
  project = var.project
  network = var.network

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = var.tags
}

resource "google_compute_firewall" "allow_internal" {
  name    = "${var.name}-allow-internal"
  project = var.project
  network = var.network

  allow {
    protocol = "tcp"
    ports    = ["0-65535"]
  }

  allow {
    protocol = "udp"
    ports    = ["0-65535"]
  }

  allow {
    protocol = "icmp"
  }

  source_tags = var.tags
  target_tags = var.tags
}

# Allow health check probes from GCP load balancers
resource "google_compute_firewall" "allow_health_check" {
  name    = "${var.name}-allow-hc"
  project = var.project
  network = var.network

  allow {
    protocol = "tcp"
  }

  source_ranges = ["130.211.0.0/22", "35.191.0.0/16"]
  target_tags   = var.tags
}
