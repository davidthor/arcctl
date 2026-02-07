terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

locals {
  redis_version = "REDIS_${replace(var.engine_version, ".", "_")}"
}

resource "google_redis_instance" "main" {
  name               = var.name
  project            = var.project
  region             = var.region
  tier               = var.tier
  memory_size_gb     = 1
  redis_version      = local.redis_version
  authorized_network = var.network

  labels = {
    environment = var.name
    managed-by  = "cldctl"
  }
}
