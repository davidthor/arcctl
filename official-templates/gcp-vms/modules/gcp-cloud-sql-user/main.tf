terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
}

locals {
  # Parse database reference to extract instance and database info
  # database input is expected to be a reference resolved by the engine
  instance_name = var.database.instance_name
  database_name = var.database.database
  host          = var.database.host
  port          = var.database.port
  scheme        = var.database.scheme
}

resource "random_password" "user" {
  length  = 24
  special = false
}

resource "google_sql_user" "user" {
  name     = var.username
  project  = var.project
  instance = local.instance_name
  password = random_password.user.result
}
