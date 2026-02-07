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
  # Parse runtime configuration
  runtime_language = try(var.runtime.language, var.runtime)
  runtime_parts    = split(":", local.runtime_language)
  language         = local.runtime_parts[0]
  language_version = length(local.runtime_parts) > 1 ? local.runtime_parts[1] : "latest"

  runtime_os       = try(var.runtime.os, "linux")
  runtime_packages = try(var.runtime.packages, [])
  runtime_setup    = try(var.runtime.setup, [])

  env_exports = var.environment != null ? join("\n", [
    for key, value in var.environment : "export ${key}='${value}'"
  ]) : ""

  # Language-specific install commands
  install_commands = {
    node = <<-EOT
      curl -fsSL https://deb.nodesource.com/setup_${local.language_version}.x | bash -
      apt-get install -y nodejs
    EOT
    python = <<-EOT
      apt-get install -y software-properties-common
      add-apt-repository -y ppa:deadsnakes/ppa
      apt-get update
      apt-get install -y python${local.language_version} python${local.language_version}-venv python3-pip
      ln -sf /usr/bin/python${local.language_version} /usr/bin/python
    EOT
    go = <<-EOT
      wget -q https://go.dev/dl/go${local.language_version}.linux-amd64.tar.gz
      tar -C /usr/local -xzf go${local.language_version}.linux-amd64.tar.gz
      export PATH=$PATH:/usr/local/go/bin
      echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
    EOT
    ruby = <<-EOT
      apt-get install -y ruby${local.language_version} ruby${local.language_version}-dev build-essential
    EOT
    java = <<-EOT
      apt-get install -y openjdk-${local.language_version}-jdk
    EOT
  }

  packages_install = length(local.runtime_packages) > 0 ? "apt-get install -y ${join(" ", local.runtime_packages)}" : ""
  setup_commands   = length(local.runtime_setup) > 0 ? join("\n", local.runtime_setup) : ""
  command_str      = var.command != null ? join(" ", var.command) : ""

  user_data = <<-EOT
    #!/bin/bash
    set -euo pipefail

    # Update system
    apt-get update -y
    apt-get upgrade -y

    # Install language runtime
    ${lookup(local.install_commands, local.language, "echo 'Unknown language: ${local.language}'")}

    # Install system packages
    ${local.packages_install}

    # Set environment variables
    ${local.env_exports}

    # Create app directory
    mkdir -p /opt/app
    cd /opt/app

    # Run setup commands
    ${local.setup_commands}

    # Create systemd service for the application
    cat > /etc/systemd/system/cldctl-app.service <<'UNIT'
    [Unit]
    Description=cldctl managed application
    After=network.target

    [Service]
    Type=simple
    Restart=always
    RestartSec=5
    WorkingDirectory=/opt/app
    ${local.env_exports}
    ExecStart=/bin/bash -c '${local.command_str}'

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
  image    = "ubuntu-22-04-x64"
  ssh_keys = var.ssh_key_fingerprint != "" ? [var.ssh_key_fingerprint] : []
  vpc_uuid = var.vpc_uuid

  user_data = local.user_data

  tags = ["cldctl", "managed-by:cldctl", "runtime", "cldctl-${var.name}"]

  lifecycle {
    create_before_destroy = true
  }
}
