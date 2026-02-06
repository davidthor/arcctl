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

resource "vercel_project_environment_variable" "this" {
  project_id = var.project_id
  key        = replace(upper(var.name), "-", "_")
  value      = var.value
  target     = [var.vercel_env]
  sensitive  = true
}
