terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

resource "google_vpc_access_connector" "main" {
  name    = var.name
  project = var.project
  region  = var.region
  network = var.network

  ip_cidr_range = "10.8.0.0/28"
  min_instances = 2
  max_instances = 10

  machine_type = "e2-micro"
}
