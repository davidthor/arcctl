variable "name" {
  description = "Service name"
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

variable "target" {
  description = "Target Cloud Run service URL or endpoint"
  type        = string
}

variable "target_type" {
  description = "Type of the target (deployment or function)"
  type        = string
}

variable "port" {
  description = "Service port"
  type        = number
  default     = null
}
