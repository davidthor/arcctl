variable "name" {
  description = "Name for the backend resources"
  type        = string
}

variable "project" {
  description = "GCP project ID"
  type        = string
}

variable "region" {
  description = "GCP region"
  type        = string
}

variable "load_balancer" {
  description = "URL map ID of the load balancer"
  type        = string
}

variable "domain" {
  description = "Domain for this backend"
  type        = string
}

variable "target" {
  description = "Cloud Run service name to route to"
  type        = string
}

variable "target_type" {
  description = "Type of the target resource"
  type        = string
  default     = "deployment"
}

variable "port" {
  description = "Port of the target service"
  type        = number
  default     = 443
}
