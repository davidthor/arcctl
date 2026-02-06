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

resource "vercel_project" "this" {
  name = var.name
}
