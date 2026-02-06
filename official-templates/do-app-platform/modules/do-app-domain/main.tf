terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

provider "digitalocean" {
  token = var.token
}

locals {
  # Build the hostname from the domain and service name
  hostname = var.domain != "" ? "${var.name}.${var.domain}" : ""
  url      = local.hostname != "" ? "https://${local.hostname}" : ""

  # If no custom domain, use the app's default URL
  effective_host = local.hostname != "" ? local.hostname : "${var.name}.ondigitalocean.app"
  effective_url  = local.hostname != "" ? "https://${local.hostname}" : "https://${var.name}.ondigitalocean.app"
}

# Configure domain for the App Platform app
# Note: In practice, domain configuration is part of the app spec.
# This module manages the DNS record pointing to the app.
resource "digitalocean_record" "domain" {
  count  = var.domain != "" ? 1 : 0
  domain = var.domain
  type   = "CNAME"
  name   = var.name
  value  = "${var.name}.ondigitalocean.app."
  ttl    = 300
}
