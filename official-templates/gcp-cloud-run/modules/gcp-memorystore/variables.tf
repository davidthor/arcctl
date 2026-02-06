variable "name" {
  description = "Name for the Memorystore Redis instance"
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

variable "version" {
  description = "Redis version (e.g., 7)"
  type        = string
  default     = "7"
}

variable "tier" {
  description = "Redis tier: BASIC or STANDARD_HA"
  type        = string
  default     = "BASIC"
}

variable "network" {
  description = "VPC network ID for private connectivity"
  type        = string
}
