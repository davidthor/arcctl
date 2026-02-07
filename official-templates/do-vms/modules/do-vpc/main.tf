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

resource "digitalocean_vpc" "vpc" {
  name        = var.name
  region      = var.region
  description = "VPC for cldctl managed infrastructure"
  ip_range    = "10.10.10.0/24"
}
