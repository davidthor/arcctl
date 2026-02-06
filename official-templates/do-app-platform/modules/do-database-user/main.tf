terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
}

provider "digitalocean" {
  token = var.token
}

# Look up the database cluster
data "digitalocean_database_cluster" "cluster" {
  name = var.database
}

# Create database user
resource "digitalocean_database_user" "user" {
  cluster_id = data.digitalocean_database_cluster.cluster.id
  name       = var.name
}

# Create a dedicated database for this user
resource "digitalocean_database_db" "db" {
  cluster_id = data.digitalocean_database_cluster.cluster.id
  name       = replace(var.name, "-", "_")
}

locals {
  connection_url = "${data.digitalocean_database_cluster.cluster.engine}://${digitalocean_database_user.user.name}:${digitalocean_database_user.user.password}@${data.digitalocean_database_cluster.cluster.host}:${data.digitalocean_database_cluster.cluster.port}/${digitalocean_database_db.db.name}"
}
