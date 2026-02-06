terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

provider "digitalocean" {
  token = var.do_token
}

locals {
  env_exports = var.environment != null ? join("\n", [
    for key, value in var.environment : "export ${key}='${value}'"
  ]) : ""

  env_flags = var.environment != null ? join(" ", [
    for key, value in var.environment : "-e ${key}='${value}'"
  ]) : ""

  command_str = var.command != null ? join(" ", var.command) : ""

  user_data = <<-EOT
    #!/bin/bash
    set -euo pipefail

    # Install Docker
    apt-get update -y
    apt-get install -y docker.io cron
    systemctl start docker
    systemctl enable docker

    # Login to container registry
    echo "${var.do_token}" | docker login registry.digitalocean.com -u token --password-stdin

    # Pull the image
    docker pull ${var.image}

    # Set environment variables
    ${local.env_exports}

    # Create cron wrapper script
    cat > /opt/cron-job.sh <<'SCRIPT'
    #!/bin/bash
    ${local.env_exports}
    docker run --rm \
      ${local.env_flags} \
      ${var.image} \
      ${local.command_str}
    SCRIPT
    chmod +x /opt/cron-job.sh

    # Install crontab
    echo "${var.schedule} /opt/cron-job.sh >> /var/log/arcctl-cron.log 2>&1" | crontab -

    # Ensure cron is running
    systemctl enable cron
    systemctl start cron
  EOT
}

resource "digitalocean_droplet" "droplet" {
  name     = var.name
  region   = var.region
  size     = var.size
  image    = "docker-20-04"
  ssh_keys = var.ssh_key_fingerprint != "" ? [var.ssh_key_fingerprint] : []
  vpc_uuid = var.vpc_uuid

  user_data = local.user_data

  tags = ["arcctl", "managed-by:arcctl", "cronjob"]

  lifecycle {
    create_before_destroy = true
  }
}
