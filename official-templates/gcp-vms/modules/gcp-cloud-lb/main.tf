terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

# Reserve a global static IP for the load balancer
resource "google_compute_global_address" "main" {
  name    = "${var.name}-ip"
  project = var.project
}

# Managed SSL certificate
resource "google_compute_managed_ssl_certificate" "main" {
  name    = "${var.name}-cert"
  project = var.project

  managed {
    domains = ["${var.domain}", "*.${var.domain}"]
  }
}

# Default URL map (will be updated by backend modules)
resource "google_compute_url_map" "main" {
  name    = var.name
  project = var.project

  default_url_redirect {
    strip_query    = false
    https_redirect = true
  }
}

# HTTPS proxy
resource "google_compute_target_https_proxy" "main" {
  name    = "${var.name}-https"
  project = var.project
  url_map = google_compute_url_map.main.id

  ssl_certificates = [google_compute_managed_ssl_certificate.main.id]
}

# Global forwarding rule
resource "google_compute_global_forwarding_rule" "https" {
  name       = "${var.name}-https"
  project    = var.project
  target     = google_compute_target_https_proxy.main.id
  port_range = "443"
  ip_address = google_compute_global_address.main.address
}

# HTTP-to-HTTPS redirect
resource "google_compute_url_map" "http_redirect" {
  name    = "${var.name}-http-redirect"
  project = var.project

  default_url_redirect {
    strip_query            = false
    https_redirect         = true
    redirect_response_code = "MOVED_PERMANENTLY_DEFAULT"
  }
}

resource "google_compute_target_http_proxy" "redirect" {
  name    = "${var.name}-http-redirect"
  project = var.project
  url_map = google_compute_url_map.http_redirect.id
}

resource "google_compute_global_forwarding_rule" "http_redirect" {
  name       = "${var.name}-http-redirect"
  project    = var.project
  target     = google_compute_target_http_proxy.redirect.id
  port_range = "80"
  ip_address = google_compute_global_address.main.address
}
