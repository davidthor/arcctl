terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

provider "digitalocean" {
  token = var.do_token
}

locals {
  service_port = coalesce(var.port, 8080)
  # Use the deployment/function target's private IP for internal routing
  target_host = var.target
}

# Create an internal DNS record for service discovery between Droplets.
# Uses a DigitalOcean domain record pointing to the target Droplet's private IP.
resource "digitalocean_record" "service" {
  domain = var.domain
  type   = "A"
  name   = "svc-${var.name}"
  value  = local.target_host
  ttl    = 60
}
