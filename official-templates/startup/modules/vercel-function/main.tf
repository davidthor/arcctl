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

# Deploy source code as a Vercel function with framework detection.
# Vercel auto-detects frameworks (Next.js, Remix, Astro, etc.) and
# configures the build pipeline accordingly.
resource "vercel_deployment" "this" {
  project_id = var.project_id

  production = var.vercel_env == "production"

  environment = merge(var.environment, {
    CLDCTL_FUNCTION_NAME = var.name
  })
}

# Set up a project domain alias for this function
resource "vercel_project_domain" "this" {
  count = var.alias != "" ? 1 : 0

  project_id = var.project_id
  domain     = var.alias
}
