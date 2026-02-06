terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

locals {
  env_exports = join("\n", [for k, v in coalesce(var.environment, {}) : "export ${k}=\"${v}\""])
  cmd_string  = var.command != null ? join(" ", var.command) : "echo 'No command specified'"

  startup_script = <<-EOT
    #!/bin/bash
    set -e

    # Install Docker
    apt-get update
    apt-get install -y docker.io
    systemctl start docker

    # Set environment variables
    ${local.env_exports}

    # Pull and run the task container
    docker pull ${var.image}
    docker run --rm \
      ${join(" ", [for k, v in coalesce(var.environment, {}) : "-e ${k}=\"${v}\""])} \
      ${var.image} ${local.cmd_string}

    # Signal completion and shut down
    echo "Task completed" > /tmp/task_status
    shutdown -h now
  EOT
}

resource "google_compute_instance" "main" {
  name         = var.name
  project      = var.project
  zone         = var.zone
  machine_type = "e2-medium"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-12"
      size  = 20
      type  = "pd-standard"
    }
  }

  network_interface {
    network    = var.network
    subnetwork = var.subnet

    access_config {}
  }

  metadata = {
    ssh-keys = var.ssh_key != "" ? "arcctl:${var.ssh_key}" : null
  }

  metadata_startup_script = local.startup_script

  service_account {
    scopes = ["cloud-platform"]
  }

  labels = {
    managed-by = "arcctl"
    task       = "true"
  }

  scheduling {
    preemptible       = true
    automatic_restart = false
  }

  allow_stopping_for_update = true
}
