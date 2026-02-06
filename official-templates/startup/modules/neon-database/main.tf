terraform {
  required_providers {
    neon = {
      source  = "kislerdm/neon"
      version = "~> 0.6"
    }
  }
}

provider "neon" {
  api_key = var.api_key
}

# Create a branch for non-production environments.
# Production uses the main branch directly (parent_branch == null).
resource "neon_branch" "this" {
  count = var.parent_branch != null ? 1 : 0

  project_id = var.project_id
  parent_id  = var.parent_branch
  name       = var.branch_name
}

locals {
  branch_id = var.parent_branch != null ? neon_branch.this[0].id : null
}

# Read-write endpoint for the branch
resource "neon_endpoint" "this" {
  project_id = var.project_id
  branch_id  = local.branch_id
  type       = "read_write"
}

# Database role (owner of the database)
resource "neon_role" "this" {
  project_id = var.project_id
  branch_id  = local.branch_id
  name       = var.name
}

# Database instance
resource "neon_database" "this" {
  project_id = var.project_id
  branch_id  = local.branch_id
  name       = var.name
  owner_name = neon_role.this.name
}
