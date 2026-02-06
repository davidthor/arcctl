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
  name             = replace(var.name, "/[^a-zA-Z0-9-_]/", "-")
  schedule         = try(var.schedule, "rate(1 hour)")
  environment_vars = try(var.environment, {})
  timeout          = try(var.timeout, 900)
  memory_size      = try(var.memory, 512)
}

# Lambda execution role
resource "aws_iam_role" "lambda" {
  name_prefix = "${local.name}-cron-"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })

  tags = {
    Name      = "${local.name}-cronjob"
    ManagedBy = "arcctl"
  }
}

resource "aws_iam_role_policy_attachment" "lambda_basic" {
  role       = aws_iam_role.lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "lambda_vpc" {
  role       = aws_iam_role.lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_lambda_function" "this" {
  function_name = local.name
  role          = aws_iam_role.lambda.arn
  package_type  = "Image"
  image_uri     = var.image
  timeout       = local.timeout
  memory_size   = local.memory_size

  image_config {
    command = try(var.command, null)
  }

  vpc_config {
    subnet_ids         = data.aws_subnets.private.ids
    security_group_ids = [var.security_group_id]
  }

  environment {
    variables = local.environment_vars
  }

  tags = {
    Name      = local.name
    ManagedBy = "arcctl"
  }
}

# EventBridge rule for scheduling
resource "aws_cloudwatch_event_rule" "this" {
  name                = local.name
  description         = "Scheduled Lambda: ${local.name}"
  schedule_expression = local.schedule

  tags = {
    Name      = local.name
    ManagedBy = "arcctl"
  }
}

resource "aws_cloudwatch_event_target" "this" {
  rule      = aws_cloudwatch_event_rule.this.name
  target_id = local.name
  arn       = aws_lambda_function.this.arn
}

resource "aws_lambda_permission" "eventbridge" {
  statement_id  = "AllowEventBridgeInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.this.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.this.arn
}
