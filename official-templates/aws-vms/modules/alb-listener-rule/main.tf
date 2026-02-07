terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

data "aws_lb" "this" {
  arn = var.alb_arn
}

data "aws_lb_listener" "https" {
  load_balancer_arn = var.alb_arn
  port              = 443
}

resource "aws_lb_listener_rule" "this" {
  listener_arn = data.aws_lb_listener.https.arn
  priority     = null # Auto-assign priority

  action {
    type             = "forward"
    target_group_arn = var.target_group_arn
  }

  condition {
    host_header {
      values = [var.domain]
    }
  }

  tags = {
    Name      = var.domain
    ManagedBy = "cldctl"
  }
}
