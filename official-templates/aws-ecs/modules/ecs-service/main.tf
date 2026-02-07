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
    values = [var.vpc_id]
  }
}

locals {
  name       = replace("${var.cluster}-${lookup(var, "name", "service")}", "/[^a-zA-Z0-9-]/", "-")
  cpu        = coalesce(lookup(var, "cpu", null), "256")
  memory     = coalesce(lookup(var, "memory", null), "512")
  replicas   = coalesce(lookup(var, "replicas", null), 1)

  environment_vars = lookup(var, "environment", {})
  container_port   = lookup(var, "port", 8080)
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
    ManagedBy = "cldctl"
  }
}

resource "aws_iam_role_policy_attachment" "execution" {
  role       = aws_iam_role.execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# ECS Task Role
resource "aws_iam_role" "task" {
  name_prefix = "${local.name}-task-"

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
    Name      = "${local.name}-task"
    ManagedBy = "cldctl"
  }
}

resource "aws_ecs_task_definition" "this" {
  family                   = local.name
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = local.cpu
  memory                   = local.memory
  execution_role_arn       = aws_iam_role.execution.arn
  task_role_arn            = aws_iam_role.task.arn

  container_definitions = jsonencode([{
    name      = local.name
    image     = var.image
    essential = true
    command   = lookup(var, "command", null)

    portMappings = [{
      containerPort = local.container_port
      hostPort      = local.container_port
      protocol      = "tcp"
    }]

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
    ManagedBy = "cldctl"
  }
}

resource "aws_ecs_service" "this" {
  name            = local.name
  cluster         = var.cluster
  task_definition = aws_ecs_task_definition.this.arn
  desired_count   = local.replicas
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = data.aws_subnets.private.ids
    security_groups  = [var.security_group_id]
    assign_public_ip = false
  }

  deployment_circuit_breaker {
    enable   = true
    rollback = true
  }

  tags = {
    Name      = local.name
    ManagedBy = "cldctl"
  }

  lifecycle {
    ignore_changes = [desired_count]
  }
}
