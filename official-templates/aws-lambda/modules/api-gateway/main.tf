terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

resource "aws_apigatewayv2_api" "this" {
  name          = var.name
  protocol_type = "HTTP"
  description   = "HTTP API Gateway for ${var.name}"

  cors_configuration {
    allow_origins = ["*"]
    allow_methods = ["GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"]
    allow_headers = ["*"]
    max_age       = 300
  }

  tags = {
    Name      = var.name
    ManagedBy = "arcctl"
  }
}

# Custom domain mapping
resource "aws_apigatewayv2_domain_name" "this" {
  count       = var.certificate_arn != "" ? 1 : 0
  domain_name = var.name

  domain_name_configuration {
    certificate_arn = var.certificate_arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }

  tags = {
    Name      = var.name
    ManagedBy = "arcctl"
  }
}
