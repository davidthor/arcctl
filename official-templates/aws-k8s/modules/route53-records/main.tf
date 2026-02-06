terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

locals {
  is_alb = var.target_type == "alb"
}

data "aws_lb" "alb" {
  count = local.is_alb ? 1 : 0
  name  = var.target
}

resource "aws_route53_record" "this" {
  zone_id = var.hosted_zone_id
  name    = var.domain
  type    = "A"

  alias {
    name                   = var.target
    zone_id                = local.is_alb && length(data.aws_lb.alb) > 0 ? data.aws_lb.alb[0].zone_id : var.hosted_zone_id
    evaluate_target_health = true
  }
}
