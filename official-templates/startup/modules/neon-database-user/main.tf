terraform {
  required_providers {
    neon = {
      source  = "kislerdm/neon"
      version = "~> 0.6"
    }
    http = {
      source  = "hashicorp/http"
      version = "~> 3.0"
    }
  }
}

provider "neon" {
  api_key = var.api_key
}

# Look up the branch endpoints to get the host for the connection URL
data "http" "branch_endpoints" {
  url = "https://console.neon.tech/api/v2/projects/${var.project_id}/branches/${var.branch}/endpoints"

  request_headers = {
    Authorization = "Bearer ${var.api_key}"
    Accept        = "application/json"
  }
}

locals {
  endpoints = jsondecode(data.http.branch_endpoints.response_body)
  host      = local.endpoints.endpoints[0].host
}

# Create a new role on the specified branch
resource "neon_role" "this" {
  project_id = var.project_id
  branch_id  = var.branch
  name       = var.name
}
