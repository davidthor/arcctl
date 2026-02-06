terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

locals {
  is_cloudfront = var.target_type == "cloudfront"
  # CloudFront always uses hosted zone Z2FDTNDATAQYW2
  cloudfront_zone_id = "Z2FDTNDATAQYW2"
}

resource "aws_route53_record" "this" {
  zone_id = var.hosted_zone_id
  name    = var.domain
  type    = "A"

  alias {
    name                   = var.target
    zone_id                = local.is_cloudfront ? local.cloudfront_zone_id : var.hosted_zone_id
    evaluate_target_health = !local.is_cloudfront
  }
}
