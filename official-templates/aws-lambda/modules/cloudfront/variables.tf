variable "api_gateway_url" {
  description = "API Gateway invoke URL"
  type        = string
}

variable "domain" {
  description = "Custom domain name"
  type        = string
  default     = ""
}

variable "certificate_arn" {
  description = "ACM certificate ARN (must be in us-east-1 for CloudFront)"
  type        = string
  default     = ""
}
