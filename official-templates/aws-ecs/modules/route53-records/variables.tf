variable "hosted_zone_id" {
  description = "Route53 hosted zone ID"
  type        = string
}

variable "domain" {
  description = "Domain name for the record"
  type        = string
}

variable "alb_dns_name" {
  description = "ALB DNS name for alias record"
  type        = string
}

variable "alb_zone_id" {
  description = "ALB hosted zone ID for alias record"
  type        = string
}
