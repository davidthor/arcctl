variable "name" {
  description = "Monitoring stack name"
  type        = string
}

variable "region" {
  description = "DigitalOcean App Platform region"
  type        = string
}

variable "token" {
  description = "DigitalOcean API token"
  type        = string
  sensitive   = true
}
