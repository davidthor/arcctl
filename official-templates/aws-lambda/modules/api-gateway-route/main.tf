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

# Create API Gateway custom domain mapping for the environment
resource "aws_apigatewayv2_domain_name" "this" {
  count       = var.certificate_arn != "" ? 1 : 0
  domain_name = var.domain

  domain_name_configuration {
    certificate_arn = var.certificate_arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }

  tags = {
    Name      = var.domain
    ManagedBy = "cldctl"
  }
}

resource "aws_apigatewayv2_api_mapping" "this" {
  count           = var.certificate_arn != "" ? 1 : 0
  api_id          = var.api_id
  domain_name     = aws_apigatewayv2_domain_name.this[0].id
  stage           = var.stage
  api_mapping_key = ""
}

# DNS record pointing to API Gateway custom domain
resource "aws_route53_record" "this" {
  count   = var.certificate_arn != "" && var.hosted_zone_id != "" ? 1 : 0
  zone_id = var.hosted_zone_id
  name    = var.domain
  type    = "A"

  alias {
    name                   = aws_apigatewayv2_domain_name.this[0].domain_name_configuration[0].target_domain_name
    zone_id                = aws_apigatewayv2_domain_name.this[0].domain_name_configuration[0].hosted_zone_id
    evaluate_target_health = false
  }
}
