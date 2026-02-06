variable "api_id" {
  description = "API Gateway ID"
  type        = string
}

variable "stage" {
  description = "API Gateway stage name"
  type        = string
}

variable "domain" {
  description = "Domain name for the route"
  type        = string
}

variable "certificate_arn" {
  description = "ACM certificate ARN"
  type        = string
  default     = ""
}

variable "hosted_zone_id" {
  description = "Route53 hosted zone ID"
  type        = string
  default     = ""
}
