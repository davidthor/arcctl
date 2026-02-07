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

resource "digitalocean_loadbalancer" "lb" {
  name   = var.name
  region = var.region

  # Default forwarding rule for HTTP
  forwarding_rule {
    entry_port      = 80
    entry_protocol  = "http"
    target_port     = 80
    target_protocol = "http"
  }

  # HTTPS forwarding
  forwarding_rule {
    entry_port      = 443
    entry_protocol  = "https"
    target_port     = 80
    target_protocol = "http"
    tls_passthrough = false
  }

  healthcheck {
    port     = 80
    protocol = "http"
    path     = "/healthz"
  }

  vpc_uuid             = var.vpc_id
  redirect_http_to_https = true
  enable_proxy_protocol  = false

  droplet_tag = "cldctl-${var.name}"
}
