terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

data "aws_region" "current" {}

# Look up the latest Amazon Linux 2023 AMI
data "aws_ami" "al2023" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["al2023-ami-*-x86_64"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

locals {
  name             = var.name
  runtime_str      = try(tostring(var.runtime), "")
  runtime_obj      = try(var.runtime, {})
  language         = local.runtime_str != "" ? local.runtime_str : try(local.runtime_obj.language, "node:20")
  os_type          = try(local.runtime_obj.os, "linux")
  packages         = try(local.runtime_obj.packages, [])
  setup_commands   = try(local.runtime_obj.setup, [])
  environment_vars = try(var.environment, {})
  command          = try(var.command, [])

  # Parse language:version
  lang_parts = split(":", local.language)
  lang_name  = local.lang_parts[0]
  lang_ver   = length(local.lang_parts) > 1 ? local.lang_parts[1] : "latest"

  # Generate install script based on language
  install_script = local.lang_name == "node" ? <<-SCRIPT
    curl -fsSL https://rpm.nodesource.com/setup_${local.lang_ver}.x | bash -
    yum install -y nodejs
  SCRIPT
  : local.lang_name == "python" ? <<-SCRIPT
    yum install -y python${local.lang_ver} python${local.lang_ver}-pip
  SCRIPT
  : local.lang_name == "go" ? <<-SCRIPT
    wget https://go.dev/dl/go${local.lang_ver}.linux-amd64.tar.gz
    tar -C /usr/local -xzf go${local.lang_ver}.linux-amd64.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile.d/go.sh
  SCRIPT
  : ""

  env_exports = join("\n", [for k, v in local.environment_vars : "export ${k}='${v}'"])
  pkg_install = length(local.packages) > 0 ? "yum install -y ${join(" ", local.packages)}" : ""
  setup_script = join("\n", local.setup_commands)

  user_data = <<-USERDATA
    #!/bin/bash
    set -euxo pipefail

    # Install language runtime
    ${local.install_script}

    # Install system packages
    ${local.pkg_install}

    # Setup commands
    ${local.setup_script}

    # Set environment variables
    ${local.env_exports}

    # Start application
    ${length(local.command) > 0 ? join(" ", local.command) : "echo 'No command specified'"}
  USERDATA
}

resource "aws_instance" "this" {
  ami                    = data.aws_ami.al2023.id
  instance_type          = try(var.instance_type, "t3.small")
  subnet_id              = var.subnet_id
  vpc_security_group_ids = [var.security_group]

  user_data = base64encode(local.user_data)

  tags = {
    Name      = local.name
    ManagedBy = "arcctl"
  }

  lifecycle {
    create_before_destroy = true
  }
}
