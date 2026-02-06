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
  # Map engine + version to Cloud SQL database_version format
  database_version = var.engine == "POSTGRES" ? "POSTGRES_${var.engine_version}" : "MYSQL_${replace(var.engine_version, ".", "_")}"
  port             = var.engine == "POSTGRES" ? 5432 : 3306
  scheme           = var.engine == "POSTGRES" ? "postgresql" : "mysql"
}

resource "random_password" "admin" {
  length  = 24
  special = false
}

resource "google_sql_database_instance" "main" {
  name                = var.name
  project             = var.project
  region              = var.region
  database_version    = local.database_version
  deletion_protection = false

  settings {
    tier              = var.tier
    availability_type = "ZONAL"
    disk_autoresize   = true
    disk_size         = 10
    disk_type         = "PD_SSD"

    ip_configuration {
      ipv4_enabled                                  = false
      private_network                               = var.network
      enable_private_path_for_google_cloud_services = true
    }

    backup_configuration {
      enabled = true
    }

    user_labels = {
      environment = var.name
      managed-by  = "arcctl"
    }
  }

  depends_on = [var.network]
}

resource "google_sql_database" "main" {
  name     = var.database_name
  project  = var.project
  instance = google_sql_database_instance.main.name
}

resource "google_sql_user" "admin" {
  name     = "admin"
  project  = var.project
  instance = google_sql_database_instance.main.name
  password = random_password.admin.result
}
