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

  # Build env vars for user_data script
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
  }

  packages_install = length(local.runtime_packages) > 0 ? "apt-get install -y ${join(" ", local.runtime_packages)}" : ""
  setup_commands   = length(local.runtime_setup) > 0 ? join("\n", local.runtime_setup) : ""

  user_data = <<-EOT
    #!/bin/bash
    set -euo pipefail

    # Update system
    apt-get update -y
    apt-get upgrade -y

    # Install language runtime
    ${lookup(local.install_commands, local.language, "")}

    # Install system packages
    ${local.packages_install}

    # Run setup commands
    ${local.setup_commands}

    # Set environment variables
    ${local.env_exports}

    # Start application
    cd /opt/app
    ${var.command != null ? join(" ", var.command) : ""}
  EOT

  command_str = var.command != null ? join(" ", var.command) : ""
}

resource "digitalocean_droplet" "droplet" {
  name     = var.name
  region   = var.region
  size     = coalesce(var.size, "s-1vcpu-1gb")
  image    = "ubuntu-22-04-x64"
  ssh_keys = var.ssh_key_fingerprint != "" ? [var.ssh_key_fingerprint] : []

  user_data = local.user_data

  tags = ["arcctl", "managed-by:arcctl", "runtime"]

  lifecycle {
    create_before_destroy = true
  }
}
