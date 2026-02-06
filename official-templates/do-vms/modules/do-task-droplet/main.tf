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

  command_str = var.command != null ? join(" ", var.command) : "echo 'No command specified'"

  user_data = <<-EOT
    #!/bin/bash
    set -euo pipefail

    # Install Docker
    apt-get update -y
    apt-get install -y docker.io
    systemctl start docker
    systemctl enable docker

    # Set environment variables
    ${local.env_exports}

    # Log into registry and pull image
    echo "${var.do_token}" | docker login registry.digitalocean.com -u token --password-stdin

    # Run the task container
    docker run --rm \
      ${var.environment != null ? join(" ", [for k, v in var.environment : "-e ${k}='${v}'"]) : ""} \
      ${var.image} \
      ${local.command_str}

    # Self-terminate after task completion
    # The droplet stays for log inspection, cleaned up by TTL
    echo "Task completed at $(date)" > /tmp/task-status
  EOT
}

resource "digitalocean_droplet" "task" {
  name     = var.name
  region   = var.region
  size     = var.size
  image    = "docker-20-04"
  ssh_keys = var.ssh_key_fingerprint != "" ? [var.ssh_key_fingerprint] : []
  vpc_uuid = var.vpc_uuid

  user_data = local.user_data

  tags = ["arcctl", "managed-by:arcctl", "task"]

  lifecycle {
    create_before_destroy = true
  }
}
