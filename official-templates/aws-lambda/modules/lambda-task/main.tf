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
  name             = replace(var.name, "/[^a-zA-Z0-9-_]/", "-")
  environment_vars = var.environment != null ? var.environment : {}
}

# Lambda execution role
resource "aws_iam_role" "lambda" {
  name_prefix = "${local.name}-task-"

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
    Name      = "${local.name}-task"
    ManagedBy = "cldctl"
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
  timeout       = var.timeout
  memory_size   = var.memory

  image_config {
    command = var.command
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
    ManagedBy = "cldctl"
  }
}

# Invoke the Lambda function synchronously
resource "null_resource" "invoke" {
  triggers = {
    function = aws_lambda_function.this.arn
    always   = timestamp()
  }

  provisioner "local-exec" {
    command = <<-EOT
      aws lambda invoke \
        --function-name ${aws_lambda_function.this.function_name} \
        --invocation-type RequestResponse \
        --region ${data.aws_region.current.name} \
        /tmp/${local.name}-response.json
    EOT
  }
}
