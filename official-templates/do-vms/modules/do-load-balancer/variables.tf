variable "name" {
  description = "Load balancer name"
  type        = string
}

variable "region" {
  description = "DigitalOcean region"
  type        = string
}

variable "do_token" {
  description = "DigitalOcean API token"
  type        = string
  sensitive   = true
}

variable "vpc_id" {
  description = "VPC ID for the load balancer"
  type        = string
  default     = null
}
