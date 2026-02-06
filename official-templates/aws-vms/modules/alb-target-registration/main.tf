terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

# Look up the target group to get the ALB DNS name
data "aws_lb_target_group" "this" {
  arn = var.target_group_arn
}

data "aws_lb" "this" {
  arn = data.aws_lb_target_group.this.load_balancer_arns[0]
}

resource "aws_lb_target_group_attachment" "this" {
  target_group_arn = var.target_group_arn
  target_id        = var.target
  port             = var.port
}
