terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

locals {
  env_flags = join(" ", [for k, v in coalesce(var.environment, {}) : "-e ${k}=\"${v}\""])
  port_flag = var.port != null ? "-p ${var.port}:${var.port}" : ""

  startup_script = <<-EOT
    #!/bin/bash
    set -e

    # Install Docker
    apt-get update
    apt-get install -y docker.io
    systemctl enable docker
    systemctl start docker

    # Authenticate with Artifact Registry (if applicable)
    gcloud auth configure-docker ${split("/", var.image)[0]} --quiet 2>/dev/null || true

    # Pull the container image
    docker pull ${var.image}

    # Run the container
    docker run -d \
      --name app \
      --restart always \
      ${local.port_flag} \
      ${local.env_flags} \
      ${var.image} \
      ${var.command != null ? join(" ", var.command) : ""}

    echo "Docker deployment complete"
  EOT
}

resource "google_compute_instance" "main" {
  name         = var.name
  project      = var.project
  zone         = var.zone
  machine_type = coalesce(var.machine_type, "e2-medium")

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-12"
      size  = 50
      type  = "pd-ssd"
    }
  }

  network_interface {
    network    = var.network
    subnetwork = var.subnet

    access_config {}
  }

  metadata = {
    ssh-keys = var.ssh_key != "" ? "cldctl:${var.ssh_key}" : null
  }

  metadata_startup_script = local.startup_script

  service_account {
    scopes = ["cloud-platform"]
  }

  labels = {
    managed-by = "cldctl"
  }

  tags = coalesce(var.tags, [])

  allow_stopping_for_update = true
}
