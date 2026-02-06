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

  app_port    = coalesce(var.port, 8080)
  command_str = var.command != null ? join(" ", var.command) : ""

  user_data = <<-EOT
    #!/bin/bash
    set -euo pipefail

    # Install Docker and Caddy
    apt-get update -y
    apt-get install -y docker.io debian-keyring debian-archive-keyring apt-transport-https curl
    systemctl start docker
    systemctl enable docker

    # Install Caddy
    curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
    curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | tee /etc/apt/sources.list.d/caddy-stable.list
    apt-get update -y
    apt-get install -y caddy

    # Login to container registry
    echo "${var.do_token}" | docker login registry.digitalocean.com -u token --password-stdin

    # Pull the image
    docker pull ${var.image}

    # Configure Caddy as reverse proxy
    cat > /etc/caddy/Caddyfile <<'CADDY'
    :8080 {
      reverse_proxy localhost:${local.app_port}
      health_uri /healthz
    }
    CADDY

    # Create systemd service for the function container
    cat > /etc/systemd/system/arcctl-function.service <<'UNIT'
    [Unit]
    Description=arcctl managed function
    After=docker.service
    Requires=docker.service

    [Service]
    Restart=always
    RestartSec=5
    ExecStartPre=-/usr/bin/docker rm -f arcctl-function
    ExecStart=/usr/bin/docker run --name arcctl-function \
      --restart=unless-stopped \
      -p ${local.app_port}:${local.app_port} \
      ${local.env_flags} \
      ${var.image} \
      ${local.command_str}
    ExecStop=/usr/bin/docker stop arcctl-function

    [Install]
    WantedBy=multi-user.target
    UNIT

    systemctl daemon-reload
    systemctl enable arcctl-function
    systemctl start arcctl-function
    systemctl restart caddy
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

  tags = ["arcctl", "managed-by:arcctl", "function", "arcctl-${var.name}"]

  lifecycle {
    create_before_destroy = true
  }
}
