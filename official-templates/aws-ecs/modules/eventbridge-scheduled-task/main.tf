terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

data "aws_subnets" "private" {
  filter {
    name   = "vpc-id"
    values = [try(var.vpc_id, "")]
  }
}

data "aws_ecs_cluster" "this" {
  cluster_name = var.cluster
}

locals {
  name             = replace(var.name, "/[^a-zA-Z0-9-_]/", "-")
  schedule         = try(var.schedule, "rate(1 hour)")
  environment_vars = try(var.environment, {})
  cpu              = try(var.cpu, "256")
  memory           = try(var.memory, "512")
}

# IAM role for EventBridge to invoke ECS
resource "aws_iam_role" "eventbridge" {
  name_prefix = "${local.name}-eb-"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "events.amazonaws.com"
      }
    }]
  })

  tags = {
    Name      = "${local.name}-eventbridge"
    ManagedBy = "arcctl"
  }
}

resource "aws_iam_role_policy" "eventbridge_ecs" {
  name = "${local.name}-ecs-run"
  role = aws_iam_role.eventbridge.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "ecs:RunTask",
        ]
        Resource = ["*"]
      },
      {
        Effect = "Allow"
        Action = [
          "iam:PassRole",
        ]
        Resource = ["*"]
      },
    ]
  })
}

# Task execution role
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
  cpu                      = local.cpu
  memory                   = local.memory
  execution_role_arn       = aws_iam_role.execution.arn

  container_definitions = jsonencode([{
    name      = local.name
    image     = var.image
    essential = true
    command   = try(var.command, null)

    environment = [for k, v in local.environment_vars : {
      name  = k
      value = tostring(v)
    }]

    logConfiguration = {
      logDriver = "awslogs"
      options = {
        "awslogs-group"         = try(var.log_group, "/ecs/${local.name}")
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

resource "aws_cloudwatch_event_rule" "this" {
  name                = local.name
  description         = "Scheduled task: ${local.name}"
  schedule_expression = local.schedule

  tags = {
    Name      = local.name
    ManagedBy = "arcctl"
  }
}

resource "aws_cloudwatch_event_target" "this" {
  rule      = aws_cloudwatch_event_rule.this.name
  target_id = local.name
  arn       = data.aws_ecs_cluster.this.arn
  role_arn  = aws_iam_role.eventbridge.arn

  ecs_target {
    task_count          = 1
    task_definition_arn = aws_ecs_task_definition.this.arn
    launch_type         = "FARGATE"

    network_configuration {
      subnets         = data.aws_subnets.private.ids
      security_groups = try([var.security_group_id], [])
      assign_public_ip = false
    }
  }
}
