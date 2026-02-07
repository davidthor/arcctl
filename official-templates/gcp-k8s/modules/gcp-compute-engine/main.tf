terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

locals {
  # Parse runtime string or object to get language and version
  runtime_language = try(var.runtime.language, var.runtime)
  language_parts   = split(":", local.runtime_language)
  language         = local.language_parts[0]
  lang_version     = length(local.language_parts) > 1 ? local.language_parts[1] : "latest"

  # Build startup script to install language runtime
  startup_script = <<-EOT
    #!/bin/bash
    set -e

    # Install runtime: ${local.language}:${local.lang_version}
    %{if local.language == "node" || local.language == "nodejs"}
    curl -fsSL https://deb.nodesource.com/setup_${local.lang_version}.x | bash -
    apt-get install -y nodejs
    %{endif}
    %{if local.language == "python"}
    apt-get update && apt-get install -y python${local.lang_version} python3-pip
    %{endif}
    %{if local.language == "go" || local.language == "golang"}
    wget -q https://go.dev/dl/go${local.lang_version}.linux-amd64.tar.gz
    tar -C /usr/local -xzf go${local.lang_version}.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin
    %{endif}
    %{if local.language == "java"}
    apt-get update && apt-get install -y openjdk-${local.lang_version}-jdk
    %{endif}

    # Install additional packages
    %{for pkg in try(var.runtime.packages, [])}
    apt-get install -y ${pkg}
    %{endfor}

    # Run setup commands
    %{for cmd in try(var.runtime.setup, [])}
    ${cmd}
    %{endfor}

    # Set environment variables
    %{for k, v in coalesce(var.environment, {})}
    export ${k}="${v}"
    %{endfor}

    # Start the application
    %{if var.command != null}
    exec ${join(" ", var.command)}
    %{endif}
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

    access_config {
      # Ephemeral public IP
    }
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
