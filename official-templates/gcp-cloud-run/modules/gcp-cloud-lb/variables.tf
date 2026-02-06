variable "name" {
  description = "Name for the load balancer resources"
  type        = string
}

variable "project" {
  description = "GCP project ID"
  type        = string
}

variable "domain" {
  description = "Base domain for SSL certificates"
  type        = string
}
