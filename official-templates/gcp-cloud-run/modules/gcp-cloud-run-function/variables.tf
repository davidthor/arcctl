variable "name" {
  description = "Name for the Cloud Run function service"
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

variable "image" {
  description = "Container image to deploy"
  type        = string
}

variable "port" {
  description = "Container port"
  type        = number
  default     = null
}

variable "command" {
  description = "Container command override"
  type        = list(string)
  default     = null
}

variable "environment" {
  description = "Environment variables"
  type        = map(string)
  default     = {}
}

variable "cpu" {
  description = "CPU limit"
  type        = string
  default     = null
}

variable "memory" {
  description = "Memory limit"
  type        = string
  default     = null
}

variable "timeout" {
  description = "Request timeout in seconds"
  type        = number
  default     = null
}

variable "max_scale" {
  description = "Maximum number of instances"
  type        = number
  default     = null
}

variable "vpc_connector" {
  description = "VPC connector ID for private network access"
  type        = string
}
