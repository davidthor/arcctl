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
  # Map component types to DigitalOcean database engine names
  engine_map = {
    postgres = "pg"
    mysql    = "mysql"
    redis    = "redis"
    mongodb  = "mongodb"
  }

  engine = lookup(local.engine_map, var.type, var.type)

  # Default versions per engine
  default_versions = {
    pg      = "16"
    mysql   = "8"
    redis   = "7"
    mongodb = "7"
  }

  version = coalesce(var.version, lookup(local.default_versions, local.engine, null))

  # Build connection URL based on engine type
  connection_url = local.engine == "redis" ? (
    "rediss://${digitalocean_database_cluster.db.user}:${digitalocean_database_cluster.db.password}@${digitalocean_database_cluster.db.host}:${digitalocean_database_cluster.db.port}"
  ) : (
    "${var.type}://${digitalocean_database_cluster.db.user}:${digitalocean_database_cluster.db.password}@${digitalocean_database_cluster.db.private_host}:${digitalocean_database_cluster.db.port}/${digitalocean_database_cluster.db.database}"
  )
}

resource "digitalocean_database_cluster" "db" {
  name       = var.name
  engine     = local.engine
  version    = local.version
  size       = var.size
  region     = var.region
  node_count = 1

  # Place database in the VPC for private networking with Droplets
  private_network_uuid = var.vpc_uuid

  tags = ["arcctl", "managed-by:arcctl"]
}
