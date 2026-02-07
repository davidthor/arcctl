terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "digitalocean" {
  token             = var.token
  spaces_access_id  = digitalocean_spaces_bucket.bucket.access_key_id
  spaces_secret_key = digitalocean_spaces_bucket.bucket.secret_access_key
}

locals {
  bucket_name = lower(replace(var.name, "/[^a-z0-9-]/", "-"))
}

resource "digitalocean_spaces_bucket" "bucket" {
  name   = local.bucket_name
  region = var.region

  dynamic "versioning" {
    for_each = var.versioning ? [1] : []
    content {
      enabled = true
    }
  }
}

# Configure CORS for public buckets
resource "digitalocean_spaces_bucket_cors_configuration" "cors" {
  count  = var.public ? 1 : 0
  bucket = digitalocean_spaces_bucket.bucket.id
  region = var.region

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["GET", "HEAD"]
    allowed_origins = ["*"]
    max_age_seconds = 3600
  }
}

# Generate access keys for the Spaces bucket
resource "digitalocean_spaces_bucket_object" "access_marker" {
  region       = var.region
  bucket       = digitalocean_spaces_bucket.bucket.name
  key          = ".cldctl-managed"
  content      = "managed by cldctl"
  content_type = "text/plain"
}
