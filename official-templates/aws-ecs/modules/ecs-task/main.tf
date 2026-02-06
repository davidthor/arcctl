terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    null = {
      source  = "hashicorp/null"
      version = "~> 3.0"
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

locals {
  name             = replace(var.name, "/[^a-zA-Z0-9-]/", "-")
  environment_vars = var.environment != null ? var.environment : {}
}

# ECS Task Execution Role
resource "aws_iam_role" "execution" {
  name_prefix = "${local.name}-exec-"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ecs-tasks.amazonaws.com"
      }
    }]
  })

  tags = {
    Name      = "${local.name}-execution"
    ManagedBy = "arcctl"
  }
}

resource "aws_iam_role_policy_attachment" "execution" {
  role       = aws_iam_role.execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_ecs_task_definition" "this" {
  family                   = local.name
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = var.cpu
  memory                   = var.memory
  execution_role_arn       = aws_iam_role.execution.arn

  container_definitions = jsonencode([{
    name      = local.name
    image     = var.image
    essential = true
    command   = var.command

    environment = [for k, v in local.environment_vars : {
      name  = k
      value = tostring(v)
    }]

    logConfiguration = {
      logDriver = "awslogs"
      options = {
        "awslogs-group"         = var.log_group
        "awslogs-region"        = data.aws_region.current.name
        "awslogs-stream-prefix" = local.name
      }
    }
  }])

  tags = {
    Name      = local.name
    ManagedBy = "arcctl"
  }
}

# Run the ECS task
resource "null_resource" "run_task" {
  triggers = {
    task_def = aws_ecs_task_definition.this.arn
    always   = timestamp()
  }

  provisioner "local-exec" {
    command = <<-EOT
      aws ecs run-task \
        --cluster ${var.cluster} \
        --task-definition ${aws_ecs_task_definition.this.arn} \
        --launch-type FARGATE \
        --capacity-provider-strategy capacityProvider=${var.capacity_provider},weight=1 \
        --network-configuration "awsvpcConfiguration={subnets=[${join(",", data.aws_subnets.private.ids)}],securityGroups=[${var.security_group_id}],assignPublicIp=DISABLED}" \
        --region ${data.aws_region.current.name}
    EOT
  }
}
