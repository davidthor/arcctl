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
  hostname = "${var.subdomain}.${var.domain}"
}

# Create a DNS record for this route pointing to the load balancer
resource "digitalocean_record" "route" {
  domain = var.domain
  type   = "A"
  name   = var.subdomain
  value  = data.digitalocean_loadbalancer.lb.ip
  ttl    = 300
}

# Reference the existing load balancer
data "digitalocean_loadbalancer" "lb" {
  id = var.load_balancer_id
}

# Create a certificate for the domain
resource "digitalocean_certificate" "cert" {
  name    = "${var.name}-cert"
  type    = "lets_encrypt"
  domains = [local.hostname]
}
