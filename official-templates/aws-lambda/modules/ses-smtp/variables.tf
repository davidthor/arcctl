variable "region" {
  description = "AWS region for SES endpoint"
  type        = string
}

variable "identity_arn" {
  description = "SES verified identity ARN"
  type        = string
  default     = ""
}
