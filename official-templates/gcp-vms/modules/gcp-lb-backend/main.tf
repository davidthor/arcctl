terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

locals {
  port = coalesce(var.port, 8080)
}

# Health check for the backend VMs
resource "google_compute_health_check" "main" {
  name    = "${var.name}-hc"
  project = var.project

  http_health_check {
    port = local.port
  }

  check_interval_sec  = 10
  timeout_sec         = 5
  healthy_threshold   = 2
  unhealthy_threshold = 3
}

# Unmanaged instance group for the backend VMs
resource "google_compute_instance_group" "main" {
  name    = "${var.name}-ig"
  project = var.project
  zone    = "${var.region}-a"

  named_port {
    name = "http"
    port = local.port
  }
}

# Backend service attached to the load balancer
resource "google_compute_backend_service" "main" {
  name    = "${var.name}-backend"
  project = var.project

  load_balancing_scheme = "EXTERNAL_MANAGED"
  protocol              = "HTTP"
  port_name             = "http"
  health_checks         = [google_compute_health_check.main.id]
  timeout_sec           = 30

  backend {
    group           = google_compute_instance_group.main.id
    balancing_mode  = "UTILIZATION"
    max_utilization = 0.8
  }

  log_config {
    enable = true
  }
}
