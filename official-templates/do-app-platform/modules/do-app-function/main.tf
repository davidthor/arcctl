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
  env_vars = var.environment != null ? [
    for key, value in var.environment : {
      key   = key
      value = value
      type  = "GENERAL"
    }
  ] : []

  # Parse image into registry/repo:tag
  image_parts = split(":", var.image)
  image_repo  = local.image_parts[0]
  image_tag   = length(local.image_parts) > 1 ? local.image_parts[1] : "latest"
}

resource "digitalocean_app" "function" {
  spec {
    name   = var.name
    region = var.region

    # Deploy as a service with auto-scaling for function-like behavior
    service {
      name               = var.name
      instance_count     = 1
      instance_size_slug = "basic-xxs"
      http_port          = coalesce(var.port, 8080)

      image {
        registry_type = "DOCR"
        registry      = split("/", local.image_repo)[0]
        repository    = join("/", slice(split("/", local.image_repo), 1, length(split("/", local.image_repo))))
        tag           = local.image_tag
      }

      dynamic "env" {
        for_each = local.env_vars
        content {
          key   = env.value.key
          value = env.value.value
          type  = env.value.type
        }
      }
    }
  }
}
