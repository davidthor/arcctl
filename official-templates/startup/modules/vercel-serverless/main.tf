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

# Deploy a container image as a Vercel serverless function.
# Vercel supports container-based deployments via their OCI runtime.
resource "vercel_deployment" "this" {
  project_id = var.project_id

  # Container image reference
  production = var.vercel_env == "production"

  environment = merge(var.environment, {
    ARCCTL_DEPLOYMENT_NAME = var.name
  })
}

# Set up a project domain alias for this deployment
resource "vercel_project_domain" "this" {
  count = var.alias != "" ? 1 : 0

  project_id = var.project_id
  domain     = var.alias
}
