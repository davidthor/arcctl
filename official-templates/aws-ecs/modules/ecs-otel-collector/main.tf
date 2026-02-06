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

locals {
  name = replace(var.name, "/[^a-zA-Z0-9-]/", "-")
}

# Execution role for OTel collector
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

# Task role with CloudWatch and X-Ray permissions
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
    ManagedBy = "arcctl"
  }
}

resource "aws_iam_role_policy" "otel" {
  name = "${local.name}-otel"
  role = aws_iam_role.task.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "logs:PutLogEvents",
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:DescribeLogStreams",
          "logs:DescribeLogGroups",
          "cloudwatch:PutMetricData",
          "xray:PutTraceSegments",
          "xray:PutTelemetryRecords",
          "xray:GetSamplingRules",
          "xray:GetSamplingTargets",
        ]
        Resource = "*"
      },
    ]
  })
}

resource "aws_ecs_task_definition" "this" {
  family                   = local.name
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = "512"
  memory                   = "1024"
  execution_role_arn       = aws_iam_role.execution.arn
  task_role_arn            = aws_iam_role.task.arn

  container_definitions = jsonencode([{
    name      = "otel-collector"
    image     = "amazon/aws-otel-collector:latest"
    essential = true

    portMappings = [
      {
        containerPort = 4317
        hostPort      = 4317
        protocol      = "tcp"
      },
      {
        containerPort = 4318
        hostPort      = 4318
        protocol      = "tcp"
      },
    ]

    environment = [
      {
        name  = "AOT_CONFIG_CONTENT"
        value = yamlencode({
          receivers = {
            otlp = {
              protocols = {
                grpc = { endpoint = "0.0.0.0:4317" }
                http = { endpoint = "0.0.0.0:4318" }
              }
            }
          }
          exporters = {
            awsxray = {}
            awsemf = {
              log_group_name = var.log_group
              region         = data.aws_region.current.name
            }
          }
          service = {
            pipelines = {
              traces = {
                receivers = ["otlp"]
                exporters = ["awsxray"]
              }
              metrics = {
                receivers = ["otlp"]
                exporters = ["awsemf"]
              }
            }
          }
        })
      },
    ]

    logConfiguration = {
      logDriver = "awslogs"
      options = {
        "awslogs-group"         = var.log_group
        "awslogs-region"        = data.aws_region.current.name
        "awslogs-stream-prefix" = "otel-collector"
      }
    }
  }])

  tags = {
    Name      = local.name
    ManagedBy = "arcctl"
  }
}

resource "aws_ecs_service" "this" {
  name            = local.name
  cluster         = var.cluster
  task_definition = aws_ecs_task_definition.this.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = data.aws_subnets.private.ids
    security_groups  = [var.security_group_id]
    assign_public_ip = false
  }

  tags = {
    Name      = local.name
    ManagedBy = "arcctl"
  }
}

# Service discovery for OTel collector
resource "aws_service_discovery_private_dns_namespace" "otel" {
  name = "${local.name}.otel.local"
  vpc  = var.vpc_id

  tags = {
    Name      = "${local.name}-otel"
    ManagedBy = "arcctl"
  }
}

resource "aws_service_discovery_service" "otel" {
  name = "collector"

  dns_config {
    namespace_id = aws_service_discovery_private_dns_namespace.otel.id
    dns_records {
      ttl  = 10
      type = "A"
    }
    routing_policy = "MULTIVALUE"
  }

  health_check_custom_config {
    failure_threshold = 1
  }

  tags = {
    Name      = "${local.name}-collector"
    ManagedBy = "arcctl"
  }
}
