variable "name" {
  description = "Target group name"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "alb_arn" {
  description = "ALB ARN"
  type        = string
}
