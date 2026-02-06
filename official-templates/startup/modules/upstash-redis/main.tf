terraform {
  required_providers {
    upstash = {
      source  = "upstash/upstash"
      version = "~> 1.0"
    }
  }
}

provider "upstash" {
  api_key = var.api_key
  email   = var.email
}

resource "upstash_redis_database" "this" {
  database_name = var.name
  region        = var.region
  tls           = true
}
