variable "name" {
  description = "Redis database name"
  type        = string
}

variable "api_key" {
  description = "Upstash API key"
  type        = string
  sensitive   = true
}

variable "email" {
  description = "Upstash account email"
  type        = string
}

variable "region" {
  description = "Upstash region (e.g., us-east-1, eu-west-1)"
  type        = string
  default     = "us-east-1"
}
