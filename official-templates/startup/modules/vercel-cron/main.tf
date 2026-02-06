terraform {
  required_providers {
    vercel = {
      source  = "vercel/vercel"
      version = "~> 2.0"
    }
    null = {
      source  = "hashicorp/null"
      version = "~> 3.0"
    }
  }
}

provider "vercel" {
  api_token = var.token
  team      = var.team_id != "" ? var.team_id : null
}

# Vercel Cron Jobs are configured via vercel.json in the project.
# This module uses the Vercel API to create a cron job configuration
# that triggers a serverless function on the specified schedule.
resource "null_resource" "cron_config" {
  triggers = {
    name     = var.name
    schedule = var.schedule
  }

  # Create the cron job via Vercel API
  provisioner "local-exec" {
    command = <<-EOT
      curl -s -X POST "https://api.vercel.com/v1/projects/${var.project_id}/crons" \
        -H "Authorization: Bearer ${var.token}" \
        -H "Content-Type: application/json" \
        ${var.team_id != "" ? "-H \"x-vercel-team: ${var.team_id}\"" : ""} \
        -d '{
          "name": "${var.name}",
          "schedule": "${var.schedule}",
          "path": "/api/cron/${var.name}"
        }' > /tmp/vercel-cron-${var.name}.json
    EOT
  }

  # Clean up on destroy
  provisioner "local-exec" {
    when    = destroy
    command = "rm -f /tmp/vercel-cron-${self.triggers.name}.json"
  }
}

# Store environment variables for the cron job
resource "vercel_project_environment_variable" "cron_env" {
  for_each = var.environment

  project_id = var.project_id
  key        = each.key
  value      = each.value
  target     = ["production", "preview"]
}
