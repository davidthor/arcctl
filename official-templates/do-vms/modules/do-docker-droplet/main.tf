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
  env_flags = var.environment != null ? join(" ", [
    for key, value in var.environment : "-e ${key}='${value}'"
  ]) : ""

  port_flag   = var.port != null ? "-p ${var.port}:${var.port}" : ""
  command_str = var.command != null ? join(" ", var.command) : ""

  user_data = <<-EOT
    #!/bin/bash
    set -euo pipefail

    # Install Docker
    apt-get update -y
    apt-get install -y docker.io
    systemctl start docker
    systemctl enable docker

    # Login to container registry
    echo "${var.do_token}" | docker login registry.digitalocean.com -u token --password-stdin

    # Pull the image
    docker pull ${var.image}

    # Create systemd service for the container
    cat > /etc/systemd/system/cldctl-app.service <<'UNIT'
    [Unit]
    Description=cldctl managed container
    After=docker.service
    Requires=docker.service

    [Service]
    Restart=always
    RestartSec=5
    ExecStartPre=-/usr/bin/docker rm -f cldctl-app
    ExecStart=/usr/bin/docker run --name cldctl-app \
      --restart=unless-stopped \
      ${local.env_flags} \
      ${local.port_flag} \
      ${var.image} \
      ${local.command_str}
    ExecStop=/usr/bin/docker stop cldctl-app

    [Install]
    WantedBy=multi-user.target
    UNIT

    systemctl daemon-reload
    systemctl enable cldctl-app
    systemctl start cldctl-app
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

  tags = ["cldctl", "managed-by:cldctl", "docker", "cldctl-${var.name}"]

  lifecycle {
    create_before_destroy = true
  }
}
