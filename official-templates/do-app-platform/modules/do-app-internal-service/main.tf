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
  service_port = coalesce(var.port, 8080)
  # Internal service hostname follows App Platform naming convention
  internal_host = "${var.name}.internal"
}

# App Platform internal services are configured as part of the app spec.
# We use a data source to look up the app and compute the internal routing.
# In practice, the internal service is part of the main app deployment.

resource "digitalocean_app" "internal_service" {
  spec {
    name   = var.name
    region = "nyc"

    service {
      name               = var.name
      instance_count     = 1
      instance_size_slug = "basic-xxs"
      http_port          = local.service_port

      # Internal services route to the deployment's container
      internal_ports = [local.service_port]

      image {
        registry_type = "DOCR"
        registry      = "library"
        repository    = "nginx"
        tag           = "latest"
      }
    }
  }
}
