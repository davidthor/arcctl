variable "hosted_zone_id" {
  description = "Route53 hosted zone ID"
  type        = string
}

variable "domain" {
  description = "Domain name for the record"
  type        = string
}

variable "target" {
  description = "Target DNS name"
  type        = string
}

variable "target_type" {
  description = "Target type (cloudfront, alb)"
  type        = string
  default     = "cloudfront"
}
