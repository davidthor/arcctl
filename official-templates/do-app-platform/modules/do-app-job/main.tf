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

  # Determine if this is a scheduled job (cronjob) or a run-once job
  is_scheduled = var.schedule != null && var.schedule != ""
}

resource "digitalocean_app" "job" {
  spec {
    name   = var.name
    region = var.region

    job {
      name               = var.name
      instance_count     = 1
      instance_size_slug = "basic-xxs"
      kind               = local.is_scheduled ? "PRE_DEPLOY" : "FAILED_DEPLOY"

      image {
        registry_type = "DOCR"
        registry      = split("/", var.image)[0]
        repository    = join("/", slice(split("/", split("::", var.image)[0]), 1, length(split("/", split(":", var.image)[0]))))
        tag           = try(split(":", var.image)[1], "latest")
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
