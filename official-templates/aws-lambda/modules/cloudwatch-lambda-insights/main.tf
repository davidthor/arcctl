terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

data "aws_region" "current" {}

# CloudWatch Lambda Insights provides enhanced monitoring for Lambda functions.
# The OTel collector sidecar is not needed since Lambda Insights integrates
# directly with CloudWatch. We provide an OTLP-compatible endpoint via
# the Lambda ADOT (AWS Distro for OpenTelemetry) layer.

locals {
  name = var.name

  # ADOT Lambda layer ARN for the current region
  adot_layer_arn = "arn:aws:lambda:${data.aws_region.current.name}:901920570463:layer:aws-otel-collector-amd64-ver-0-98-0:1"
}

# Create a dedicated log group for OTel data
resource "aws_cloudwatch_log_group" "otel" {
  name              = "/aws/lambda/${local.name}-otel"
  retention_in_days = 30

  tags = {
    Name      = "${local.name}-otel"
    ManagedBy = "arcctl"
  }
}
