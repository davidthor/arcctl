terraform {
  required_providers {
    vercel = {
      source  = "vercel/vercel"
      version = "~> 2.0"
    }
  }
}

provider "vercel" {
  api_token = var.token
  team      = var.team_id != "" ? var.team_id : null
}

# Vercel internal service routing.
# Services in Vercel are accessed via their deployment URL.
# This module resolves the target deployment or function endpoint
# and provides a stable service URL for internal reference.

locals {
  # Construct the internal service URL based on the target
  service_host = "${var.name}.vercel.internal"
  service_port = var.port != null ? var.port : 443
  service_url  = "https://${local.service_host}:${local.service_port}"
}

# Store the service configuration as a project environment variable
# so other deployments can discover this service.
resource "vercel_project_environment_variable" "service_url" {
  project_id = var.project_id
  key        = "SERVICE_${replace(upper(var.name), "-", "_")}_URL"
  value      = local.service_url
  target     = ["production", "preview"]
}
