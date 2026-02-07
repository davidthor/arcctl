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
  environment_vars = try(var.environment, {})
  container_port   = try(var.port, 8080)

  env_flags = join(" ", [for k, v in local.environment_vars : "-e ${k}='${v}'"])
  port_flag = "-p ${local.container_port}:${local.container_port}"

  user_data = <<-USERDATA
    #!/bin/bash
    set -euxo pipefail

    # Install Docker
    yum update -y
    yum install -y docker
    systemctl enable docker
    systemctl start docker

    # Login to ECR if image is from ECR
    ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
    REGION=${data.aws_region.current.name}
    aws ecr get-login-password --region $REGION | docker login --username AWS --password-stdin $ACCOUNT_ID.dkr.ecr.$REGION.amazonaws.com || true

    # Install CloudWatch agent for log forwarding
    yum install -y amazon-cloudwatch-agent
    cat > /opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.json <<'EOF'
    {
      "logs": {
        "logs_collected": {
          "files": {
            "collect_list": [{
              "file_path": "/var/log/docker/*.log",
              "log_group_name": "${var.log_group}",
              "log_stream_name": "${local.name}"
            }]
          }
        }
      }
    }
    EOF
    /opt/aws/amazon-cloudwatch-agent/bin/amazon-cloudwatch-agent-ctl -a fetch-config -m ec2 -s -c file:/opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.json

    # Run container
    docker run -d --restart=always --name ${local.name} \
      ${local.port_flag} \
      ${local.env_flags} \
      --log-driver=awslogs \
      --log-opt awslogs-region=${data.aws_region.current.name} \
      --log-opt awslogs-group=${var.log_group} \
      --log-opt awslogs-stream=${local.name} \
      ${var.image} ${try(join(" ", var.command), "")}
  USERDATA
}

# IAM role for EC2 instances
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
    ManagedBy = "cldctl"
  }

  lifecycle {
    create_before_destroy = true
  }
}

# Register with target group for ALB routing
resource "aws_lb_target_group_attachment" "this" {
  count            = var.target_group_arn != "" ? 1 : 0
  target_group_arn = var.target_group_arn
  target_id        = aws_instance.this.id
  port             = local.container_port
}
