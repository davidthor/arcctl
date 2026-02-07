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
  name             = replace(var.name, "/[^a-zA-Z0-9-]/", "-")
  environment_vars = var.environment != null ? var.environment : {}

  env_flags = join(" ", [for k, v in local.environment_vars : "-e ${k}='${v}'"])

  user_data = <<-USERDATA
    #!/bin/bash
    set -euxo pipefail

    # Install Docker
    yum update -y
    yum install -y docker
    systemctl enable docker
    systemctl start docker

    # ECR login
    ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
    REGION=${data.aws_region.current.name}
    aws ecr get-login-password --region $REGION | docker login --username AWS --password-stdin $ACCOUNT_ID.dkr.ecr.$REGION.amazonaws.com || true

    # Run the task container (one-shot)
    docker run --rm \
      ${local.env_flags} \
      ${var.image} ${var.command != null ? join(" ", var.command) : ""}

    # Self-terminate after task completion
    INSTANCE_ID=$(curl -s http://169.254.169.254/latest/meta-data/instance-id)
    aws ec2 terminate-instances --instance-ids $INSTANCE_ID --region $REGION
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
    ManagedBy = "cldctl"
  }
}

resource "aws_iam_role_policy_attachment" "ecr_read" {
  role       = aws_iam_role.this.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
}

resource "aws_iam_role_policy_attachment" "ssm" {
  role       = aws_iam_role.this.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

resource "aws_iam_role_policy" "self_terminate" {
  name = "${local.name}-self-terminate"
  role = aws_iam_role.this.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect   = "Allow"
      Action   = "ec2:TerminateInstances"
      Resource = "*"
      Condition = {
        StringEquals = {
          "ec2:ResourceTag/ManagedBy" = "cldctl"
        }
      }
    }]
  })
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

  instance_initiated_shutdown_behavior = "terminate"

  tags = {
    Name      = local.name
    ManagedBy = "cldctl"
  }
}
