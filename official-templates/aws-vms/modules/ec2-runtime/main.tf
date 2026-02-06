terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

data "aws_region" "current" {}

data "aws_subnets" "private" {
  filter {
    name   = "vpc-id"
    values = [var.vpc_id]
  }
}

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
  packages         = try(local.runtime_obj.packages, [])
  setup_commands   = try(local.runtime_obj.setup, [])
  environment_vars = try(var.environment, {})
  command_list     = try(var.command, [])

  lang_parts = split(":", local.language)
  lang_name  = local.lang_parts[0]
  lang_ver   = length(local.lang_parts) > 1 ? local.lang_parts[1] : "latest"

  install_scripts = {
    node   = "curl -fsSL https://rpm.nodesource.com/setup_${local.lang_ver}.x | bash -\nyum install -y nodejs"
    python = "yum install -y python${local.lang_ver} python${local.lang_ver}-pip"
    go     = "wget https://go.dev/dl/go${local.lang_ver}.linux-amd64.tar.gz\ntar -C /usr/local -xzf go${local.lang_ver}.linux-amd64.tar.gz\necho 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile.d/go.sh"
    java   = "yum install -y java-${local.lang_ver}-amazon-corretto"
    ruby   = "yum install -y ruby${local.lang_ver}"
  }
  install_script = lookup(local.install_scripts, local.lang_name, "")

  env_exports  = join("\n", [for k, v in local.environment_vars : "export ${k}='${v}'"])
  pkg_install  = length(local.packages) > 0 ? "yum install -y ${join(" ", local.packages)}" : ""
  setup_script = join("\n", local.setup_commands)
  app_port     = try(var.port, 8080)

  user_data = <<-USERDATA
    #!/bin/bash
    set -euxo pipefail

    # Install CloudWatch agent
    yum install -y amazon-cloudwatch-agent
    cat > /opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.json <<'CWEOF'
    {
      "logs": {
        "logs_collected": {
          "files": {
            "collect_list": [{
              "file_path": "/var/log/app/*.log",
              "log_group_name": "${var.log_group}",
              "log_stream_name": "${local.name}"
            }]
          }
        }
      }
    }
    CWEOF
    /opt/aws/amazon-cloudwatch-agent/bin/amazon-cloudwatch-agent-ctl -a fetch-config -m ec2 -s -c file:/opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.json

    # Install language runtime
    ${local.install_script}

    # Install system packages
    ${local.pkg_install}

    # Setup commands
    ${local.setup_script}

    # Set environment variables
    ${local.env_exports}

    # Create log directory
    mkdir -p /var/log/app

    # Run application
    ${length(local.command_list) > 0 ? join(" ", local.command_list) : "echo 'No command specified'"} 2>&1 | tee /var/log/app/app.log &
  USERDATA
}

resource "aws_iam_role" "this" {
  name_prefix = "${local.name}-"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ec2.amazonaws.com"
      }
    }]
  })

  tags = {
    Name      = local.name
    ManagedBy = "arcctl"
  }
}

resource "aws_iam_role_policy_attachment" "ssm" {
  role       = aws_iam_role.this.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

resource "aws_iam_role_policy_attachment" "cloudwatch" {
  role       = aws_iam_role.this.name
  policy_arn = "arn:aws:iam::aws:policy/CloudWatchAgentServerPolicy"
}

resource "aws_iam_instance_profile" "this" {
  name_prefix = "${local.name}-"
  role        = aws_iam_role.this.name
}

resource "aws_instance" "this" {
  ami                    = data.aws_ami.al2023.id
  instance_type          = var.instance_type
  subnet_id              = data.aws_subnets.private.ids[0]
  vpc_security_group_ids = [var.security_group_id]
  key_name               = var.key_pair != "" ? var.key_pair : null
  iam_instance_profile   = aws_iam_instance_profile.this.name
  user_data              = base64encode(local.user_data)

  tags = {
    Name      = local.name
    ManagedBy = "arcctl"
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_lb_target_group_attachment" "this" {
  count            = var.target_group_arn != "" ? 1 : 0
  target_group_arn = var.target_group_arn
  target_id        = aws_instance.this.id
  port             = local.app_port
}
