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
  environment_vars = try(var.environment, {})
  timeout          = try(var.timeout, 900)
  memory_size      = try(var.memory, 512)
}

# Lambda execution role
resource "aws_iam_role" "lambda" {
  name_prefix = "${local.name}-"

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
    Name      = local.name
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
  timeout       = local.timeout
  memory_size   = local.memory_size

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

# Lambda function URL for direct invocation
resource "aws_lambda_function_url" "this" {
  function_name      = aws_lambda_function.this.function_name
  authorization_type = "NONE"
}

# API Gateway integration if api_id is provided
resource "aws_apigatewayv2_integration" "this" {
  count = var.api_id != null ? 1 : 0

  api_id                 = var.api_id
  integration_type       = "AWS_PROXY"
  integration_uri        = aws_lambda_function.this.invoke_arn
  integration_method     = "POST"
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "this" {
  count = var.api_id != null ? 1 : 0

  api_id    = var.api_id
  route_key = "$default"
  target    = "integrations/${aws_apigatewayv2_integration.this[0].id}"
}

resource "aws_lambda_permission" "apigw" {
  count = var.api_id != null ? 1 : 0

  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.this.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${var.api_id}/*/*"
}
