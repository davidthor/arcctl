variable "name" {
  description = "API Gateway name"
  type        = string
}

variable "region" {
  description = "AWS region"
  type        = string
}

variable "certificate_arn" {
  description = "ACM certificate ARN for custom domain"
  type        = string
  default     = ""
}
