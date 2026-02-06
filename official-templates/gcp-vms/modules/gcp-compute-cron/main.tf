terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

locals {
  env_flags  = join(" ", [for k, v in coalesce(var.environment, {}) : "-e ${k}=\"${v}\""])
  cmd_string = var.command != null ? join(" ", var.command) : "echo 'No command'"

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

    # Set up cron job
    cat > /etc/cron.d/arcctl-cronjob <<'CRON'
    ${var.schedule} root docker run --rm ${local.env_flags} ${var.image} ${local.cmd_string} >> /var/log/arcctl-cron.log 2>&1
    CRON
    chmod 644 /etc/cron.d/arcctl-cronjob

    echo "Cron deployment complete"
  EOT
}

resource "google_compute_instance" "main" {
  name         = var.name
  project      = var.project
  zone         = var.zone
  machine_type = "e2-small"

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
  }

  tags = coalesce(var.tags, [])

  allow_stopping_for_update = true
}
