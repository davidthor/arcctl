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

resource "digitalocean_record" "record" {
  domain = var.domain
  type   = "A"
  name   = var.subdomain
  value  = var.target
  ttl    = 300
}
