terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

data "aws_apigatewayv2_api" "this" {
  api_id = var.api_id
}

locals {
  endpoint_host = replace(data.aws_apigatewayv2_api.this.api_endpoint, "https://", "")
}

resource "aws_apigatewayv2_integration" "this" {
  api_id             = var.api_id
  integration_type   = "HTTP_PROXY"
  integration_uri    = "http://${var.target}:${var.port}"
  integration_method = "ANY"

  request_parameters = {
    "overwrite:path" = "/$request.path"
  }
}
