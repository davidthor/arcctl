variable "name" {
  description = "VPC name"
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
