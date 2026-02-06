variable "alb_arn" {
  description = "ALB ARN"
  type        = string
}

variable "target_group_arn" {
  description = "Target group ARN to forward traffic to"
  type        = string
}

variable "domain" {
  description = "Domain name for host-header matching"
  type        = string
}
